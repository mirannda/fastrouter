package main

import (
	"expvar"
	"fastrouter"
	"log"
	"net/http"
	_ "net/http/pprof"
	"time"
	"sync"
	"runtime"
)

var r *fastrouter.Router
var w http.ResponseWriter
var count expvar.Int
var handler = func(w http.ResponseWriter, req *http.Request) {
	count.Add(1)
}

func init() {
	r = fastrouter.NewRouter(nil)
	r.GET("/user/{name}/{age}", handler)
	r.GET("/usar/{name}/{age}", handler)
	r.GET("/usbr/{name}/{age}", handler)
	r.GET("/usdr/{name}/{age}", handler)
	r.GET("/usfr/{name}/{age}", handler)
	r.GET("/usgr/{name}/{age}", handler)
	r.GET("/ushr/{name}/{age}", handler)
	r.GET("/usir/{name}/{age}", handler)
	r.GET("/usjr/{name}/{age}", handler)
	r.GET("/uskr/{name}/{age}", handler)
	r.GET("/uslr/{name}/{age}", handler)
	r.GET("/usmr/{name}/{age}", handler)
	r.GET("/usnr/{name}/{age}", handler)
	r.GET("/usor/{name}/{age}", handler)
	r.GET("/uspr/{name}/{age}", handler)
	r.GET("/usqr/{name}/{age}", handler)
	r.GET("/usrr/{name}/{age}", handler)
	r.GET("/ussr/{name}/{age}", handler)
	r.GET("/ustr/{name}/{age}", handler)
	r.GET("/usur/{name}/{age}", handler)
	r.GET("/usvr/{name}/{age}", handler)
	r.GET("/uswr/{name}/{age}", handler)
	r.GET("/usxr/{name}/{age}", handler)
	r.GET("/usyr/{name}/{age}", handler)
	r.GET("/uszr/{name}/{age}", handler)
	r.GET("/use/{name}/{age}", handler)
	r.GET("/us/{name}/{age}", handler)
	r.GET("/u/{name}/{age}", handler)
	r.GET("/hello", handler)
}

func sequence() {
	begin := time.Now()
	for i := 0; i < 10000000; i++ {
		req, _ := http.NewRequest("GET", "/uszr/foobar/1000", nil)
		r.ServeHTTP(w, req)
	}
	log.Println(time.Now().Sub(begin), count)
}

func goroutines() {
	var wg sync.WaitGroup
	runtime.GOMAXPROCS(4)

	begin := time.Now()
	for i := 0; i < 10000000; i++ {
		wg.Add(1)
		go func() {
			req, _ := http.NewRequest("GET", "/uszr/foobar/1000", nil)
			r.ServeHTTP(w, req)
			wg.Done()
		}()
	}
	wg.Wait()
	log.Println(time.Now().Sub(begin), count)
}

func main() {
	go func() {
		http.ListenAndServe(":6060", nil)
	}()

	sequence()
	goroutines()
}
