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
)

func TestSimpleLookup_RootEpoch(t *testing.T) {
	storer := mock.NewStorer()
	topic := []byte("testtopic")
	//updateData := []byte("updateData")
	level := uint8(32)
	pk, _ := crypto.GenerateSecp256k1Key()
	signer := crypto.NewDefaultSigner(pk)
	i, err := feeds.NewId(topic, 0, level)
	if err != nil {
		t.Fatal(err)
	}

	ethAddr, err := signer.EthereumAddress()
	if err != nil {
		t.Fatal(err)
	}

	mockChunk := testingc.GenerateTestRandomChunk()
	ch, err := soc.NewChunk(i.Bytes(), mockChunk, signer)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("test chunk addr ", ch.Address().String())
	_, err = storer.Put(context.Background(), storage.ModePutUpload, ch)
	if err != nil {
		t.Fatal(err)
	}
	now := uint64(time.Now().Unix())

	result, err := feeds.SimpleLookupAt(context.Background(), storer, ethAddr, topic, now)
	if err != nil {
		t.Fatal(err)
	}

	if !bytes.Equal(result, mockChunk.Data()) {
		t.Fatalf("result mismatch. want %v got %v", mockChunk.Data(), result)
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
