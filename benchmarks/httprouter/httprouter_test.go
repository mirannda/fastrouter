package main

import (
	"net/http"
	"sync"
	"testing"

	"github.com/julienschmidt/httprouter"
)

var router *httprouter.Router
var w http.ResponseWriter
var handler = func(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
}

func init() {
	router = httprouter.New()
	router.GET("/usar/:name/:age", handler)
	router.GET("/usbr/:name/:age", handler)
	router.GET("/uscr/:name/:age", handler)
	router.GET("/usdr/:name/:age", handler)
	router.GET("/user/:name/:age", handler)
	router.GET("/usfr/:name/:age", handler)
	router.GET("/usgr/:name/:age", handler)
	router.GET("/ushr/:name/:age", handler)
	router.GET("/usir/:name/:age", handler)
	router.GET("/usjr/:name/:age", handler)
	router.GET("/uskr/:name/:age", handler)
	router.GET("/uslr/:name/:age", handler)
	router.GET("/usmr/:name/:age", handler)
	router.GET("/usnr/:name/:age", handler)
	router.GET("/usor/:name/:age", handler)
	router.GET("/uspr/:name/:age", handler)
	router.GET("/usqr/:name/:age", handler)
	router.GET("/usrr/:name/:age", handler)
	router.GET("/ussr/:name/:age", handler)
	router.GET("/ustr/:name/:age", handler)
	router.GET("/usur/:name/:age", handler)
	router.GET("/usvr/:name/:age", handler)
	router.GET("/uswr/:name/:age", handler)
	router.GET("/usxr/:name/:age", handler)
	router.GET("/usyr/:name/:age", handler)
	router.GET("/uszr/:name/:age", handler)

	router.GET("/params/:a/:b/:c/:d/:e/:f/:g/:h/:i/:j/:k/:l/:m/:n/:o/:p/:q/:r/:s/:t/:u/:v/:w/:x/:y/:z", handler)
	router.GET("/noparams/noparams/noparams", handler)
}

func BenchmarkHttpRouter_10000000Goroutines_2Params(b *testing.B) {
	var group sync.WaitGroup

	b.N = 10000000
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

func BenchmarkHttpRouter_8Goroutines_2Params(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			req, _ := http.NewRequest("GET", "/usor/foobar/1000", nil)
			router.ServeHTTP(w, req)
		}
	})
}

func BenchmarkHttpRouter_1Goroutine_2Params(b *testing.B) {
	for i := 0; i < b.N; i++ {
		req, _ := http.NewRequest("GET", "/usor/foobar/1000", nil)
		router.ServeHTTP(w, req)
	}
}

func BenchmarkHttpRouter_10000000Goroutines_26Params(b *testing.B) {
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


func BenchmarkHttpRouter_10000000Goroutines_NoParams(b *testing.B) {
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
