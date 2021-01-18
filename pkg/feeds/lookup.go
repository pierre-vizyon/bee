package feeds

import (
	"context"
	"encoding/binary"
	"hash"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethersphere/bee/pkg/soc"
	"github.com/ethersphere/bee/pkg/storage"
	"golang.org/x/crypto/sha3"
)

var (
	retrieveTimeout = 10 * time.Second
	maxLevel        = 32
)

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
	hasher := hashPool.Get().(hash.Hash)
	defer func() {
		hasher.Reset()
		hashPool.Put(hasher)
	}()

	_, err := hasher.Write(i.topic[:])
	if err != nil {
		//return nil, err
		panic(err)
	}
	_, err = hasher.Write(i.index[:])
	if err != nil {
		//return nil, err
		panic(err)
	}

	return hasher.Sum(nil)
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
func SimpleLookupAt(ctx context.Context, getter storage.Getter, user common.Address, topic []byte, time uint64) ([]byte, error) {
	return simpleLookupAt(ctx, getter, user, topic, 0, time, 32, nil)
}

func simpleLookupAt(ctx context.Context, getter storage.Getter, user common.Address, topic []byte, current, time uint64, level uint8, data []byte) ([]byte, error) {
	id, _ := NewId(topic, current, level)
	owner, _ := soc.NewOwner(user[:])
	addr, err := soc.CreateAddress(id.Bytes(), owner)
	if err != nil {
		return nil, err
	}
	data1, err := getter.Get(ctx, storage.ModeGetRequest, addr)
	if err != nil {
		if data == nil {
			return nil, err
		}
		return data, nil
	}

	dd, _ := soc.FromChunk(data1)
	branch := time & (1 << level)
	if branch == 0 {
		return simpleLookupAt(ctx, getter, user, topic, current-1, time, level-1, dd.Chunk.Data()) // fetch right
	}
	current |= branch
	return simpleLookupAt(ctx, getter, user, topic, current, time, level-1, dd.Chunk.Data()) // fetch right
}
