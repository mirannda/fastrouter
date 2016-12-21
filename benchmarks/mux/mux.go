package main

import (
	"github.com/gorilla/mux"
	"log"
	"net/http"
	_ "net/http/pprof"
	"sync"
	"time"
)

var r *mux.Router
var w http.ResponseWriter
var req *http.Request
var handler = func(w http.ResponseWriter, req *http.Request) {
	//log.Println(mux.Vars(req))
}

func init() {
	r = mux.NewRouter()
	r.HandleFunc("/user/{name}/{age}", handler)
	r.HandleFunc("/usar/{name}/{age}", handler)
	r.HandleFunc("/usbr/{name}/{age}", handler)
	r.HandleFunc("/usdr/{name}/{age}", handler)
	r.HandleFunc("/usfr/{name}/{age}", handler)
	r.HandleFunc("/usgr/{name}/{age}", handler)
	r.HandleFunc("/ushr/{name}/{age}", handler)
	r.HandleFunc("/usir/{name}/{age}", handler)
	r.HandleFunc("/usjr/{name}/{age}", handler)
	r.HandleFunc("/uskr/{name}/{age}", handler)
	r.HandleFunc("/uslr/{name}/{age}", handler)
	r.HandleFunc("/usmr/{name}/{age}", handler)
	r.HandleFunc("/usnr/{name}/{age}", handler)
	r.HandleFunc("/usor/{name}/{age}", handler)
	r.HandleFunc("/uspr/{name}/{age}", handler)
	r.HandleFunc("/usqr/{name}/{age}", handler)
	r.HandleFunc("/usrr/{name}/{age}", handler)
	r.HandleFunc("/ussr/{name}/{age}", handler)
	r.HandleFunc("/ustr/{name}/{age}", handler)
	r.HandleFunc("/usur/{name}/{age}", handler)
	r.HandleFunc("/usvr/{name}/{age}", handler)
	r.HandleFunc("/uswr/{name}/{age}", handler)
	r.HandleFunc("/usxr/{name}/{age}", handler)
	r.HandleFunc("/usyr/{name}/{age}", handler)
	r.HandleFunc("/uszr/{name}/{age}", handler)
	r.HandleFunc("/use/{name}/{age}", handler)
	r.HandleFunc("/us/{name}/{age}", handler)
	r.HandleFunc("/u/{name}/{age}", handler)
	r.HandleFunc("/hello", handler)
}

func main() {
	go func() {
		http.ListenAndServe(":6060", nil)
	}()

	var group sync.WaitGroup

	begin := time.Now()
	for i := 0; i < 10000000; i++ {
		group.Add(1)
		go func() {
			req, _ := http.NewRequest("GET", "/user/foobar/1000", nil)
			r.ServeHTTP(w, req)
			group.Done()
		}()
	}
	log.Println(time.Now().Sub(begin))
}
