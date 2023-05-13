package pls

import (
	"flag"
	"fmt"
	"net/http"
)

var port int

func init() {
	flag.IntVar(&port, "pls-port", 7001, "port to listen on for pls")
}

func Listen() {
	h := http.NewServeMux()

	h.HandleFunc("/", func(w http.ResponseWriter, _ *http.Request) {
		w.Write([]byte("true"))
	})

	fmt.Printf("pls listening on http://localhost:%d\n", port)
	if err := http.ListenAndServe(fmt.Sprintf(":%d", port), h); err != nil {
		panic(err)
	}
}
