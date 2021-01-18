package feeds_test

import (
	"bytes"
	"context"
	"testing"

	"github.com/ethersphere/bee/pkg/feeds"
	"github.com/ethersphere/bee/pkg/storage/mock"
)

func TestSimpleLookup_RootEpoch(t *testing.T) {

	storer := mock.NewStorer()
	ethAddr := newMockEthAddress()
	topic := []byte("testtopic")
	updateData := []byte("updateData")
	mockUpdate := newMockUpdate()
	_ = storer.Put(mockUpdate)

	result, err := feeds.SimpleLookup(context.Background(), storer, ethAddr, topic)
	if err != nil {
		t.Fatal(err)
	}

	if !bytes.Equal(result, updateData) {
		t.Fatalf("result mismatch. want %v got %v", updateData, result)
	}
}
