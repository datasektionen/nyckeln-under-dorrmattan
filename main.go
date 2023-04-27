package main

import (
	"flag"
	"net/http"

	"github.com/datasektionen/nyckeln-under-dorrmattan/login"
	"github.com/datasektionen/nyckeln-under-dorrmattan/pls"
)

func main() {
	flag.Parse()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("true"))
	})

	go login.Listen()
	pls.Listen()
}
