package main

import (
	"net/http"
	"sync"
	"testing"
)

func BenchmarkMux(b *testing.B) {
	var group sync.WaitGroup

	for i := 0; i < b.N; i++ {
		group.Add(1)
		go func() {
			req, _ := http.NewRequest("GET", "/user/foobar/1000", nil)
			r.ServeHTTP(w, req)
			group.Done()
		}()
	}

	group.Wait()
}
