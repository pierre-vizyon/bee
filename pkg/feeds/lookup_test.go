package feeds_test

import (
	"bytes"
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/ethersphere/bee/pkg/crypto"
	"github.com/ethersphere/bee/pkg/feeds"
	"github.com/ethersphere/bee/pkg/soc"
	"github.com/ethersphere/bee/pkg/storage"
	"github.com/ethersphere/bee/pkg/storage/mock"
	testingc "github.com/ethersphere/bee/pkg/storage/testing"
	"github.com/ethersphere/bee/pkg/swarm"
)

var (
	topic      = []byte("testtopic")
	mockChunk  = testingc.GenerateTestRandomChunk()
	lastUpdate = testingc.GenerateTestRandomChunk()
)

func TestSimpleLookup(t *testing.T) {
	pk, _ := crypto.GenerateSecp256k1Key()
	signer := crypto.NewDefaultSigner(pk)
	ethAddr, err := signer.EthereumAddress()
	if err != nil {
		t.Fatal(err)
	}
	for _, tc := range []struct {
		desc    string
		updates []update // the last update we'd like to find must be the last element in this list
	}{
		//{
		//desc:    "one update at root",
		//updates: []update{updateAt(32, 0)},
		//},
		//{
		//desc:    "one update at root, one at the next level (left)",
		//updates: []update{updateAt(32, 0), updateAt(31, 0)},
		//},
		//{
		//desc:    "one update at root, one at the next level (left), and another one left down",
		//updates: []update{updateAt(32, 0), updateAt(31, 0), updateAt(30, 0)},
		//},
		{
			desc:    "one update at root, down left, left and right",
			updates: []update{updateAt(0, 0), updateAt(1, 0), updateAt(2, 0), updateAt(2, (1 << 30))},
		},
	} {
		storer := mock.NewStorer()
		fmt.Println("update at ", 1<<29)
		for i, v := range tc.updates {
			id, err := feeds.NewId(topic, v.epoch, v.level)
			if err != nil {
				t.Fatal(err)
			}
			var ch swarm.Chunk
			if i == len(tc.updates)-1 {
				// create the last update from the different mock chunk
				ch, err = soc.NewChunk(id.Bytes(), lastUpdate, signer)
			} else {
				ch, err = soc.NewChunk(id.Bytes(), mockChunk, signer)
			}
			if err != nil {
				t.Fatal(err)
			}

			fmt.Println("test chunk addr ", ch.Address().String())
			_, err = storer.Put(context.Background(), storage.ModePutUpload, ch)
			if err != nil {
				t.Fatal(err)
			}
		}
		now := uint64(time.Now().Unix())

		result, err := feeds.SimpleLookupAt(context.Background(), storer, ethAddr, topic, now)
		if err != nil {
			t.Fatal(err)
		}

		if !bytes.Equal(result, lastUpdate.Data()) {
			t.Fatalf("wrong result")
		}
	}
}

func updateAt(l uint8, e uint64) update {
	return update{l, e}
}

type update struct {
	level uint8
	epoch uint64
}
