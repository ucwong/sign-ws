package main

import (
	"fmt"
	"log"
	"net/http"
	//"strconv"
	"strings"
	"time"

	"github.com/ucwong/golang-kv"
	"github.com/ucwong/sign/util"
)

var db kv.Bucket

func main() {
	db = kv.Badger(".badger")
	mux := http.NewServeMux()
	mux.HandleFunc("/", handler)
	http.ListenAndServe("127.0.0.1:8080", mux)
}

func handler(w http.ResponseWriter, r *http.Request) {

	log.Printf("%v %v %v", r.Method, r.URL, time.Now().Unix())
	res := "OK"

	uri := strings.ToLower(r.URL.Path)
	u := strings.Split(uri, "/")
	if len(u) < 2 {
		fmt.Fprintf(w, "Invalid URL")
		return
	}

	addr, _ := u[len(u)-1], u[len(u)-2]
	//if !common.IsHexAddress(addr) {
	//	fmt.Fprintf(w, "Invalid infohash format")
	//	return
	//}
	q := r.URL.Query()
	switch r.Method {
	case "GET":
		res = Get(uri)
	case "POST":
		/*ts, err := strconv.ParseInt(q.Get("ts"), 10, 64)
		if err != nil {
			res = "Invalid timestamp"
			break
		}*/
		if !util.Verify(q.Get("msg"), addr, q.Get("sig")) {
			res = "Invalid signature"
			break
		}
		log.Printf("suc\n")
	default:
		res = "method not found"
	}
	fmt.Fprintf(w, res)
}

func Get(k string) string {
	return get(k)
}

func Set(k, v string) error {
	return set(k, v)
}

func get(k string) (v string) {
	if len(k) == 0 {
		return
	}
	v = string(db.Get([]byte(k)))
	return
}

func set(k, v string) (err error) {
	if len(k) == 0 || len(v) == 0 {
		return
	}

	err = db.Set([]byte(k), []byte(v))

	return
}
