package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/julienschmidt/httprouter"
)

var (
	authURL string
	listen  string
)

func handle(w http.ResponseWriter, r *http.Request, rp httprouter.Params) {
	log.Println("handling callback")

	nonce := r.URL.Query().Get("nonce")
	if nonce == "" {
		log.Println("received callback")
		return
	}

	res, err := http.Get(authURL + "/sign?nonce=" + nonce)
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte(err.Error()))
		return
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		w.WriteHeader(500)
		w.Write([]byte(res.Status))
		return
	}
	io.Copy(w, res.Body)
}

func main() {
	flag.StringVar(&authURL, "auth-url", "http://127.0.0.1:3000", "Auth api address")
	flag.StringVar(&listen, "listen", "127.0.0.1:6129", "Listen address")
	flag.Parse()

	r := httprouter.New()
	r.GET("/", handle)

	if err := http.ListenAndServe(listen, r); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
