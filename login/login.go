package login

import (
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"strings"
)

var port int
var hodisURL string

func init() {
	flag.IntVar(&port, "login-port", 7002, "port to listen on for login")
	flag.StringVar(&hodisURL, "hodis-url", "https://hodis.datasektionen.se", "url to hodis")
}

type loginUser struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	User      string `json:"user"`
	Email     string `json:"emails"`
	UGKthid   string `json:"ugkthid"`
}

// Mocks https://login.datasektionen.se
// If a kthID is sent on the channel, it will be used to fetch the user from
// hodis and make subsequent login requests return that user.
func Listen(kthIDs <-chan string) {
	h := http.NewServeMux()

	user := loginUser{
		FirstName: "Ture",
		LastName:  "Teknokrat",
		User:      "turetek",
		Email:     "turetek@kth.se",
		UGKthid:   "u1jwkms6",
	}

	h.HandleFunc("/hello", func(w http.ResponseWriter, _ *http.Request) {
		w.Write([]byte("Hello Login!!!!"))
	})
	h.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		callback := r.URL.Query().Get("callback")
		callback = strings.TrimSuffix(callback, "/")
		http.Redirect(w, r, fmt.Sprintf("%s/%s", callback, "dummy-token"), http.StatusSeeOther)
	})
	h.HandleFunc("/verify/", func(w http.ResponseWriter, _ *http.Request) {
		json.NewEncoder(w).Encode(user)
	})

	go func() {
		fmt.Printf("login listening on http://localhost:%d\n", port)
		if err := http.ListenAndServe(fmt.Sprintf(":%d", port), h); err != nil {
			panic(err)
		}
	}()
	for kthID := range kthIDs {
		res, err := http.Get(hodisURL + "/users/" + kthID)
		if err != nil {
			fmt.Println(err)
			continue
		}
		var users []struct {
			UGKthid     string `json:"ugKthid"`
			UID         string `json:"uid"`
			CN          string `json:"cn"`
			Mail        string `json:"mail"`
			GivenName   string `json:"givenName"`
			DisplayName string `json:"displayName"`
			Year        int    `json:"year"`
			Tag         string `json:"tag"`
		}
		json.NewDecoder(res.Body).Decode(&users)
		if len(users) == 0 {
			fmt.Println("no users found")
			continue
		}
		hodis := users[0]
		user = loginUser{
			FirstName: hodis.GivenName,
			LastName: strings.TrimPrefix(
				hodis.DisplayName,
				hodis.GivenName+" ",
			),
			User:    hodis.UID,
			Email:   hodis.Mail,
			UGKthid: hodis.UGKthid,
		}
		fmt.Println("Now logging in as:", hodis.CN)
	}
}
