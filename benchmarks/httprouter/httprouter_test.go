package main

import (
	"net/http"
	"sync"
	"testing"
)

func BenchmarkHttpRouter(b *testing.B) {
	var group sync.WaitGroup

	b.N = 10000000
	for i := 0; i < b.N; i++ {
		group.Add(1)
		go func() {
			req, _ := http.NewRequest("GET", "/uszr/foobar/1000", nil)
			r.ServeHTTP(w, req)
			group.Done()
		}()
	}

	group.Wait()
}

func BenchmarkHttpRouterP(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			req, _ := http.NewRequest("GET", "/uszr/foobar/1000", nil)
			r.ServeHTTP(w, req)
		}
	})
}

func BenchmarkHttpFor(b *testing.B) {
	b.N = 10000000

	for i := 0; i < b.N; i++ {
		req, _ := http.NewRequest("GET", "/uszr/foobar/1000", nil)
		r.ServeHTTP(w, req)
	}
}
