package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/ucwong/golang-kv"
	"github.com/ucwong/sign/util"
)

var db kv.Bucket

type Body struct {
	Timestamp int64  `json:"ts"`
	Addr      string `json:"addr"`
}

func main() {
	db = kv.Badger(".badger")
	mux := http.NewServeMux()
	mux.HandleFunc("/", handler)
	http.ListenAndServe("127.0.0.1:8080", mux)
}

func handler(w http.ResponseWriter, r *http.Request) {

	log.Printf("%v %v", r.Method, r.URL)
	res := "OK"

	uri := strings.ToLower(r.URL.Path)
	u := strings.Split(uri, "/")
	if len(u) < 2 {
		fmt.Fprintf(w, "Invalid URL")
		return
	}

	addr := u[len(u)-1] //, u[len(u)-2]
	//if !common.IsHexAddress(addr) {
	//	fmt.Fprintf(w, "Invalid infohash format")
	//	return
	//}
	q := r.URL.Query()
	switch r.Method {
	case "GET":
		res = Get(uri)
	case "POST":
		if reqBody, err := ioutil.ReadAll(r.Body); err == nil {
			//if err := Set(uri, string(reqBody)); err != nil {
			//	res = "ERROR" //fmt.Sprintf("%v", err)
			//
			var body Body
			//var to string
			if len(reqBody) > 0 {
				if err := json.Unmarshal(reqBody, &body); err != nil {
					log.Printf("%v", err)
					res = "Invalid json"
					break
				}
			}

			//to = strings.ToLower(body.Addr)
			//if len(to) > 0 && !common.IsHexAddress(to) {
			//	res = "Invalid addr format"
			//	break
			//}
			log.Println(string(reqBody))
			if !util.Verify(string(reqBody), addr, q.Get("sig"), body.Timestamp) {
				//if !Verify(string(reqBody), addr, q.Get("sig"), 1) {
				res = "Invalid signature"
				break
			}
		}
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
