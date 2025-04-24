package login

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/datasektionen/nyckeln-under-dorrmattan/pkg/config"
)

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
func Listen(cfg *config.Config, kthIDs <-chan string) {
	h := http.NewServeMux()

	// Fun fact: this is an acual user that actually exists in kth's systems
	user := loginUser{
		FirstName: "Ture",
		LastName:  "Teknokrat",
		User:      "turetek",
		Email:     "turetek@kth.se",
		UGKthid:   "u1jwkms6",
	}

	if cfg.InitKTHID != "" {
		if u, err := getUserFromHodis(cfg.HodisURL, cfg.InitKTHID); err != nil {
			fmt.Println(err)
		} else {
			user = u
		}
	}
	fmt.Println("Now logging in as:", user.FirstName, user.LastName)

	var tokenBytes [9]byte
	if _, err := rand.Read(tokenBytes[:]); err != nil {
		panic(err)
	}
	randomToken := base64.URLEncoding.EncodeToString(tokenBytes[:])

	h.HandleFunc("/hello", func(w http.ResponseWriter, _ *http.Request) {
		w.Write([]byte("Hello Login!!!!"))
	})
	h.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		callback := r.URL.Query().Get("callback")
		http.Redirect(w, r, callback+randomToken, http.StatusSeeOther)
	})
	h.HandleFunc("/verify/{token}", func(w http.ResponseWriter, r *http.Request) {
		token := r.PathValue("token")
		token = strings.TrimSuffix(token, ".json")
		if token != randomToken {
			fmt.Printf("Invalid token received. Expected %s, got %s\n", randomToken, token)
			http.Error(w, "Invalid code", http.StatusNotFound)
			return
		}
		w.Header().Add("Content-Type", "application/json;charset=utf-8")
		json.NewEncoder(w).Encode(user)
	})

	go func() {
		fmt.Printf("login listening on http://localhost:%s\n", cfg.LoginPort)
		if err := http.ListenAndServe(fmt.Sprintf(":%s", cfg.LoginPort), h); err != nil {
			panic(err)
		}
	}()
	for kthID := range kthIDs {
		var err error
		user, err = getUserFromHodis(cfg.HodisURL, kthID)
		if err != nil {
			fmt.Println(err)
			fmt.Println("Still logging in as:", user.FirstName, user.LastName)
			continue
		}
		fmt.Println("Now logging in as:", user.FirstName, user.LastName)
	}
}

func getUserFromHodis(hodisURL, kthID string) (loginUser, error) {
	res, err := http.Get(hodisURL + "/users/" + kthID)
	if err != nil {
		return loginUser{}, nil
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
		return loginUser{}, fmt.Errorf("No user with the kth id '%s' found", kthID)
	}
	hodis := users[0]
	return loginUser{
		FirstName: hodis.GivenName,
		LastName: strings.TrimPrefix(
			hodis.DisplayName,
			hodis.GivenName+" ",
		),
		User:    hodis.UID,
		Email:   hodis.Mail,
		UGKthid: hodis.UGKthid,
	}, nil
}
