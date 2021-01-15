package feeds

import (
	"context"
	"encoding/binary"
	"hash"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethersphere/bee/pkg/file"
	"golang.org/x/crypto/sha3"
)

var retrieveTimeout = 10 * time.Second

// hashPool contains a pool of ready hashers
var hashPool sync.Pool

// init initializes the package and hashPool
func init() {
	hashPool = sync.Pool{
		New: func() interface{} {
			return sha3.NewLegacyKeccak256()
		},
	}
}

type id struct {
	owner [20]byte
	topic [32]byte
	index [9]byte
}

func newId(owner, topic []byte, epoch uint64, level uint8) (*id, error) {
	hasher := hashPool.Get().(hash.Hash)
	defer func() {
		hasher.Reset()
		hashPool.Put(hasher)
	}()

	_, err := hasher.Write(topic)
	if err != nil {
		return nil, err
	}

	return &id{
		owner: owner,
		topic: hasher.Sum(),
		index: newIndex(epoch, level),
	}
}

func (i *id) identifier() []byte {
	b := make([]byte, 41)
	copy(b, i.topic)
	copy(b[32:], i.index)
	hasher := hashPool.Get().(hash.Hash)
	defer func() {
		hasher.Reset()
		hashPool.Put(hasher)
	}()

	_, err := hasher.Write(b)
	if err != nil {
		return nil, err
	}

	return hasher.Sum()
}

// address returns the chunk address for this update
func (i *id) address(owner []byte) []byte {
	b := make([]byte, 32+20)
	id := i.identifier()
	copy(b, id)
	copy(b[32:], i.owner)
	hasher := hashPool.Get().(hash.Hash)
	defer func() {
		hasher.Reset()
		hashPool.Put(hasher)
	}()

	_, err := hasher.Write(b)
	if err != nil {
		return nil, err
	}

	return hasher.Sum()
}

// newIndexReturns a new index based on a unix epoch in uint64 representation and a level
func newIndex(t uint64, l uint8) [9]byte {
	b := make([]byte, 9)
	binary.LittleEndian.PutUint64(b, t)
	b[8] = l
	return b
}

// Lookup retrieves the latest feed update
func Lookup(ctx context.Context, ls file.LoadSaver, user common.Address, topic []byte) (interface{}, error) {
	t := time.Now()
	a

	mask := 1 << (64 - 1)

	// find first bit of our time
	for {

	}

}
