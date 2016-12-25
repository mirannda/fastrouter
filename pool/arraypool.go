package pool

import (
	"bytes"
	"math/rand"
	"sync"
)

type arrayPool struct {
	minSize int
	queue   []*bytes.Buffer
	mutex   sync.RWMutex
}

func newArrayPool(minSize int) *arrayPool {
	pool := &arrayPool{
		minSize: minSize,
		queue:   make([]*bytes.Buffer, minSize),
	}

	for i := 0; i < minSize; i++ {
		pool.queue[i] = &bytes.Buffer{}
	}

	return pool
}

func (pool *arrayPool) get() *bytes.Buffer {
	pool.mutex.Lock()
	defer pool.mutex.Unlock()

	n := len(pool.queue)
	if n == 0 {
		return &bytes.Buffer{}
	}

	buffer := pool.queue[0]
	pool.queue = pool.queue[1:]

	return buffer
}

func (pool *arrayPool) put(buffer *bytes.Buffer) {
	pool.mutex.Lock()
	pool.queue = append(pool.queue, buffer)
	pool.mutex.Unlock()
}

type arrayPools struct {
	size  int
	pools []*arrayPool
}

func newArrayPools(size int) *arrayPools {
	pools := &arrayPools{
		size:  size,
		pools: make([]*arrayPool, size),
	}

	for i := 0; i < size; i++ {
		pools.pools[i] = newArrayPool(256)
	}

	return pools
}

func (pools *arrayPools) get() (*bytes.Buffer, int) {
	i := rand.Intn(pools.size)
	return pools.pools[i].get(), i
}

func (pools *arrayPools) put(buffer *bytes.Buffer, i int) {
	pools.pools[i].put(buffer)
}
