package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/ucwong/sign/util"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/sign", sign)
	err := http.ListenAndServe("127.0.0.1:8080", mux)
	if err != nil {
		log.Fatal(err)
	}
}

func sign(w http.ResponseWriter, r *http.Request) {
	log.Printf("%v %v\n", r.Method, r.URL)
	res := "OK"

	q := r.URL.Query()
	switch r.Method {
	case "GET":
		// TODO
	case "POST":
		if !util.Verify(q.Get("msg"), q.Get("addr"), q.Get("sig")) {
			res = "Invalid signature"
			break
		}
		log.Printf("suc\n")
	default:
		res = "method not found"
	}
	fmt.Fprintf(w, res)
}
