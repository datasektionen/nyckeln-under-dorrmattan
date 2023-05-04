package login

import (
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"strings"
)

var port int

func init() {
	flag.IntVar(&port, "login-port", 7002, "port to listen on for login")
}

func Listen() {
	h := http.NewServeMux()

	h.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello Login!!!!"))
	})
	h.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		callback := r.URL.Query().Get("callback")
		callback = strings.TrimSuffix(callback, "/")
		http.Redirect(w, r, fmt.Sprintf("%s/%s", callback, "lol-this-is-token"), http.StatusSeeOther)
	})
	h.HandleFunc("/verify/", func(w http.ResponseWriter, r *http.Request) {
		code := strings.TrimPrefix(r.URL.String(), "/verify/")
		if code != "lol-this-is-token" {
			fmt.Println("wrong code, but who cares 🤷")
		}
		json.NewEncoder(w).Encode(map[string]string{
			"first_name": "Ture",
			"last_name":  "Teknokrat",
			"user":       "turetek",
			"emails":     "turetek@kth.se",
			"ugkthid":    "u1jwkms6",
		})
	})

	fmt.Printf("login listening on http://localhost:%d\n", port)
	if err := http.ListenAndServe(fmt.Sprintf(":%d", port), h); err != nil {
		panic(err)
	}
}
