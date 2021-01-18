package feeds_test

import (
	"bytes"
	"context"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethersphere/bee/pkg/crypto"
	"github.com/ethersphere/bee/pkg/feeds"
	"github.com/ethersphere/bee/pkg/soc"
	"github.com/ethersphere/bee/pkg/storage"
	"github.com/ethersphere/bee/pkg/storage/mock"
	testingc "github.com/ethersphere/bee/pkg/storage/testing"
)

func TestSimpleLookup_RootEpoch(t *testing.T) {
	storer := mock.NewStorer()
	addr := []byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19}
	ethAddr := common.BytesToAddress(addr)
	topic := []byte("testtopic")
	updateData := []byte("updateData")
	epoch := uint64(time.Now().Unix())
	level := uint8(32)
	pk, _ := crypto.GenerateSecp256k1Key()
	signer := crypto.NewDefaultSigner(pk)
	i, err := feeds.NewId(topic, epoch, level)
	if err != nil {
		t.Fatal(err)
	}
	mockChunk := testingc.GenerateTestRandomChunk()
	ch, err := soc.NewChunk(i.Bytes(), mockChunk, signer)
	if err != nil {
		t.Fatal(err)
	}
	_, err = storer.Put(context.Background(), storage.ModePutUpload, ch)
	if err != nil {
		t.Fatal(err)
	}

	result, err := feeds.SimpleLookup(context.Background(), storer, ethAddr, topic)
	if err != nil {
		t.Fatal(err)
	}

	if !bytes.Equal(result, updateData) {
		t.Fatalf("result mismatch. want %v got %v", updateData, result)
	}
}

//func newChunk(content []byte) swarm.Chunk {
//s
//hasher := bmtpool.Get()
//defer bmtpool.Put(hasher)

//// execute hash, compare and return result
//err := hasher.SetSpanBytes(span)
//if err != nil {
//return false
//}
//_, err = hasher.Write(content)
//if err != nil {
//return false
//}
//s := hasher.Sum(nil)

//}
