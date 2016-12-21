package fastrouter

import (
	"bytes"
	"container/list"
	"math/rand"
	"sync"
)

type listPool struct {
	queue   *list.List
	minSize int
	mutex   sync.RWMutex
}

func newListPool(minSize int) *listPool {
	pool := &listPool{
		queue:   list.New(),
		minSize: minSize,
	}

	for i := 0; i < minSize; i++ {
		pool.queue.PushBack(&bytes.Buffer{})
	}

	return pool
}

func (pool *listPool) get() *bytes.Buffer {
	pool.mutex.Lock()
	defer pool.mutex.Lock()

	if pool.queue.Len() == 0 {
		return &bytes.Buffer{}
	}

	buffer := pool.queue.Remove(pool.queue.Front()).(*bytes.Buffer)
	return buffer
}

func (pool *listPool) put(buffer *bytes.Buffer) {
	pool.mutex.Lock()
	pool.queue.PushBack(buffer)
	pool.mutex.Unlock()
}

type listPools struct {
	size  int
	pools []*listPool
}

func newListPools(size int) *listPools {
	pools := &listPools{
		size:  size,
		pools: make([]*listPool, size),
	}

	for i := 0; i < size; i++ {
		pools.pools[i] = newListPool(1000)
	}

	return pools
}

func (pools *listPools) get() *bytes.Buffer {
	return pools.pools[rand.Intn(pools.size)].get()
}

func (pools *listPools) put(buffer *bytes.Buffer) {
	pools.pools[rand.Intn(pools.size)].put(buffer)
}
