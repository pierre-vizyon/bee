// Copyright 2021 The Swarm Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package feeds_test

import (
	"bytes"
	"context"
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
		desc       string
		updates    []update // the last update we'd like to find must be the last element in this list
		lookupTime uint64
	}{
		{
			desc:    "one update at root",
			updates: []update{updateAt(0, 0)},
		},
		{
			desc:    "one update at root, one at the next level (left)",
			updates: []update{updateAt(0, 0), updateAt(1, 0)},
		},
		{
			desc:    "one update at root, one at the next level (left), and another one left down",
			updates: []update{updateAt(0, 0), updateAt(1, 0), updateAt(2, 0)},
		},
		{
			desc:    "one update at root, down left, left and right",
			updates: []update{updateAt(0, 0), updateAt(1, 0), updateAt(2, 0), updateAt(2, (1 << 30))},
		},
		{
			desc:       "one update at root, down left, left and right",
			updates:    []update{updateAt(0, 0), updateAt(1, 0), updateAt(2, 0), updateAt(2, (1 << 30)), updateAt(1, (1 << 31))},
			lookupTime: 2147483648 + 10,
		},
	} {
		storer := mock.NewStorer()
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

			_, err = storer.Put(context.Background(), storage.ModePutUpload, ch)
			if err != nil {
				t.Fatal(err)
			}
		}
		lookupTime := uint64(time.Now().Unix())
		if tc.lookupTime != 0 {
			lookupTime = tc.lookupTime
		}

		result, err := feeds.SimpleLookupAt(context.Background(), storer, ethAddr, topic, lookupTime)
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
