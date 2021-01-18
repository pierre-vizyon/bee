package feeds

import (
	"context"
	"encoding/binary"
	"hash"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethersphere/bee/pkg/storage"
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

type Id struct {
	topic [32]byte
	index [9]byte
}

func NewId(topic []byte, epoch uint64, level uint8) (*Id, error) {
	hasher := hashPool.Get().(hash.Hash)
	defer func() {
		hasher.Reset()
		hashPool.Put(hasher)
	}()

	_, err := hasher.Write(topic)
	if err != nil {
		return nil, err
	}
	sum := hasher.Sum(nil)
	i := &Id{
		index: newIndex(epoch, level),
	}
	copy(i.topic[:], sum)
	return i, nil
}

func (i *Id) Identifier() []byte {
	b := make([]byte, 41)
	copy(b, i.topic[:])
	copy(b[32:], i.index[:])
	hasher := hashPool.Get().(hash.Hash)
	defer func() {
		hasher.Reset()
		hashPool.Put(hasher)
	}()

	_, err := hasher.Write(b)
	if err != nil {
		//return nil, err
		panic(err)
	}

	return hasher.Sum(nil)
}

func (i *Id) Bytes() []byte {
	return []byte{}
}

// address returns the chunk address for this update
//func (i *Id) Address() []byte {
//b := make([]byte, 32+20)
//id := i.Identifier()
//copy(b, id)
//copy(b[32:], i.owner)
//hasher := hashPool.Get().(hash.Hash)
//defer func() {
//hasher.Reset()
//hashPool.Put(hasher)
//}()

//_, err := hasher.Write(b)
//if err != nil {
//panic(err)
//}

//sum, err := hasher.Sum()
//if err != nil {
//panic(err)
//}
//return sum
//}

// newIndexReturns a new index based on a unix epoch in uint64 representation and a level
func newIndex(t uint64, l uint8) [9]byte {
	var b [9]byte
	binary.LittleEndian.PutUint64(b[:], t)
	b[8] = l
	return b
}

// Lookup retrieves the latest feed update
func SimpleLookup(ctx context.Context, getter storage.Getter, user common.Address, topic []byte) ([]byte, error) {
	//t := time.Now()

	// find first bit of our time
	//for {

	//}
	return nil, nil
}
