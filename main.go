package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"strings"
)

func main() {
	var port int

	flag.IntVar(&port, "port", 10917, "port to listen on")
	flag.Parse()

	http.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello Login!!!!"))
	})
	http.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		callback := r.URL.Query().Get("callback")
		callback = strings.TrimSuffix(callback, "/")
		http.Redirect(w, r, fmt.Sprintf("%s/%s", callback, "lol-this-is-token"), http.StatusSeeOther)
	})
	http.HandleFunc("/verify/", func(w http.ResponseWriter, r *http.Request) {
		code := strings.TrimPrefix(r.URL.String(), "/verify/")
		if code != "lol-this-is-token" {
			fmt.Println("wrong code, but who cares 🤷")
		}
		json.NewEncoder(w).Encode(map[string]string{
			"first_name": "Ture",
			"last_name":  "Teknolog",
			"user":       "turetek",
			"emails":     "turetek@kth.se",
			"ugkthid":    "u1e9kghi",
		})
	})

	fmt.Printf("Listening on http://localhost:%d\n", port)
	if err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil); err != nil {
		panic(err)
	}
}
