package main

import (
	"github.com/julienschmidt/httprouter"
	"log"
	"net/http"
	"sync"
	"time"
)

var r *httprouter.Router
var w http.ResponseWriter
var handler = func(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	//log.Println(ps)
}
var req *http.Request

func init() {
	r = httprouter.New()
	r.GET("/user/:name/:age", handler)
	r.GET("/usar/:name/:age", handler)
	r.GET("/usbr/:name/:age", handler)
	r.GET("/usdr/:name/:age", handler)
	r.GET("/usfr/:name/:age", handler)
	r.GET("/usgr/:name/:age", handler)
	r.GET("/ushr/:name/:age", handler)
	r.GET("/usir/:name/:age", handler)
	r.GET("/usjr/:name/:age", handler)
	r.GET("/uskr/:name/:age", handler)
	r.GET("/uslr/:name/:age", handler)
	r.GET("/usmr/:name/:age", handler)
	r.GET("/usnr/:name/:age", handler)
	r.GET("/usor/:name/:age", handler)
	r.GET("/uspr/:name/:age", handler)
	r.GET("/usqr/:name/:age", handler)
	r.GET("/usrr/:name/:age", handler)
	r.GET("/ussr/:name/:age", handler)
	r.GET("/ustr/:name/:age", handler)
	r.GET("/usur/:name/:age", handler)
	r.GET("/usvr/:name/:age", handler)
	r.GET("/uswr/:name/:age", handler)
	r.GET("/usxr/:name/:age", handler)
	r.GET("/usyr/:name/:age", handler)
	r.GET("/uszr/:name/:age", handler)
	r.GET("/use/:name/:age", handler)
	r.GET("/us/:name/:age", handler)
	r.GET("/u/:name/:age", handler)
	r.GET("/hello", handler)
}

func main() {
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
	group.Wait()
	log.Println(time.Now().Sub(begin))
}