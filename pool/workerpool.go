package pool

import (
	"bytes"
	"math/rand"
	"net/http"
	"sync"
)

type packet struct {
	w   *http.ResponseWriter
	req *http.Request
}

type workerPool struct {
	mutex sync.Mutex
	tree  *radixTree
	chans []chan *packet
	size  int
}

func newWorkerPool(size int) *workerPool {
	pool := &workerPool{
		chans: make([]chan *packet, size),
		size:  size,
	}
	return pool
}

func (pool *workerPool) serve() {
	for i := 0; i < pool.size; i++ {
		pool.chans[i] = make(chan *packet, 1000)
		go func(ch chan *packet) {
			buffer := new(bytes.Buffer)
			for p := range ch {
				pool.tree.searchWorker(p, buffer)
			}
		}(pool.chans[i])
	}
}

func (pool *workerPool) newJob(p *packet) {
	pool.chans[rand.Intn(pool.size)] <- p
}
