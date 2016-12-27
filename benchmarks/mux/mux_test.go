package main

import (
	"net/http"
	"sync"
	"testing"

	"github.com/gorilla/mux"
)

var router *mux.Router
var w http.ResponseWriter
var handler = func(w http.ResponseWriter, req *http.Request) {
}

func init() {
	router = mux.NewRouter()
	router.HandleFunc("/usar/{name}/{age}", handler)
	router.HandleFunc("/usbr/{name}/{age}", handler)
	router.HandleFunc("/uscr/{name}/{age}", handler)
	router.HandleFunc("/usdr/{name}/{age}", handler)
	router.HandleFunc("/user/{name}/{age}", handler)
	router.HandleFunc("/usfr/{name}/{age}", handler)
	router.HandleFunc("/usgr/{name}/{age}", handler)
	router.HandleFunc("/ushr/{name}/{age}", handler)
	router.HandleFunc("/usir/{name}/{age}", handler)
	router.HandleFunc("/usjr/{name}/{age}", handler)
	router.HandleFunc("/uskr/{name}/{age}", handler)
	router.HandleFunc("/uslr/{name}/{age}", handler)
	router.HandleFunc("/usmr/{name}/{age}", handler)
	router.HandleFunc("/usnr/{name}/{age}", handler)
	router.HandleFunc("/usor/{name}/{age}", handler)
	router.HandleFunc("/uspr/{name}/{age}", handler)
	router.HandleFunc("/usqr/{name}/{age}", handler)
	router.HandleFunc("/usrr/{name}/{age}", handler)
	router.HandleFunc("/ussr/{name}/{age}", handler)
	router.HandleFunc("/ustr/{name}/{age}", handler)
	router.HandleFunc("/usur/{name}/{age}", handler)
	router.HandleFunc("/usvr/{name}/{age}", handler)
	router.HandleFunc("/uswr/{name}/{age}", handler)
	router.HandleFunc("/usxr/{name}/{age}", handler)
	router.HandleFunc("/usyr/{name}/{age}", handler)
	router.HandleFunc("/uszr/{name}/{age}", handler)

	router.HandleFunc("/params/{a}/{b}/{c}/{d}/{e}/{f}/{g}/{h}/{i}/{j}/{k}/{l}/{m}/{n}/{o}/{p}/{q}/{r}/{s}/{t}/{u}/{v}/{w}/{x}/{y}/{z}", handler)
	router.HandleFunc("/noparams/noparams/noparams", handler)
}

func BenchmarkFastRouter_10000000Goroutines_2Params(b *testing.B) {
	var group sync.WaitGroup

	b.N = 1000000
	for i := 0; i < b.N; i++ {
		group.Add(1)
		go func() {
			req, _ := http.NewRequest("GET", "/usor/foobar/1000", nil)
			router.ServeHTTP(w, req)
			group.Done()
		}()
	}

	group.Wait()
}

func BenchmarkFastRouter_8Goroutines_2Params(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			req, _ := http.NewRequest("GET", "/usor/foobar/1000", nil)
			router.ServeHTTP(w, req)
		}
	})
}

func BenchmarkFastRouter_1Goroutine_2Params(b *testing.B) {
	for i := 0; i < b.N; i++ {
		req, _ := http.NewRequest("GET", "/usor/foobar/1000", nil)
		router.ServeHTTP(w, req)
	}
}

func BenchmarkFastRouter_10000000Goroutines_26Params(b *testing.B) {
	var group sync.WaitGroup

	b.N = 10000000
	for i := 0; i < b.N; i++ {
		group.Add(1)
		go func() {
			req, _ := http.NewRequest("GET", "/params/a/b/c/d/e/f/g/h/i/j/k/l/m/n/o/p/q/r/s/t/u/v/w/x/y/z", nil)
			router.ServeHTTP(w, req)
			group.Done()
		}()
	}

	group.Wait()
}


func BenchmarkFastRouter_10000000Goroutines_NoParams(b *testing.B) {
	var group sync.WaitGroup

	b.N = 10000000
	for i := 0; i < b.N; i++ {
		group.Add(1)
		go func() {
			req, _ := http.NewRequest("GET", "/noparams/noparams/noparams", nil)
			router.ServeHTTP(w, req)
			group.Done()
		}()
	}

	group.Wait()
}
