package sso

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"encoding/json"
	"fmt"
	"html/template"
	"log/slog"
	"net/http"
	"os"
	"strconv"
	"sync/atomic"

	"github.com/datasektionen/nyckeln-under-dorrmattan/pkg/config"
	"github.com/datasektionen/nyckeln-under-dorrmattan/pkg/dao"
	jose "github.com/go-jose/go-jose/v4"
	"github.com/google/uuid"
	"github.com/zitadel/oidc/v3/pkg/op"
	"golang.org/x/text/language"
)

var port int
var hodisURL string
var initKTHID string

var (
	loginTmpl, _ = template.New("login").Parse(`
	<!DOCTYPE html>
	<html>
		<head>
			<meta charset="UTF-8">
			<title>Mock SSO Login</title>
		</head>

		<body style="display: flex; align-items: center; justify-content: center;">
			<form method="POST" action="/login" style="display: flex; flex-direction: column; align-items: start; justify-content: start; gap: .5rem;">
			   <div>Mock SSO Login</div>
			   <input type="hidden" name="id" value="{{.ID}}">
			   <label for="kthid">KTH id to login as</label>
			   <input type="text" name="kthid">
			   <button type="submit">Login</button>
			   <p style="color:red;">{{.Error}}</p>
			</form>
		</body>
	</html>`)
	counter atomic.Int64
)

var SupportedScopes = []string{"openid", "profile", "email", "offline_access", "pls_*", "permissions", "year_tag"}

type auth interface {
	CheckLogin(kthid, id string) error
}

type ssoUser struct {
	Email      string `json:"email,omitempty"`
	FirstName  string `json:"firstName,omitempty"`
	FamilyName string `json:"familyName,omitempty"`
	YearTag    string `json:"yearTag,omitempty"`
}

func Listen(cfg *config.Config, dao *dao.Dao) {
	logger := slog.New(
		slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
			AddSource: true,
			Level:     slog.LevelDebug,
		}),
	)

	key, _ := rsa.GenerateKey(rand.Reader, 2048)

	issuer := fmt.Sprintf("http://localhost:%s", cfg.SsoPort)

	storage := storage{
		signingKey: signingKey{
			id:        uuid.NewString(),
			algorithm: jose.RS256,
			key:       key,
		},
		authRequests: make(map[string]*authRequest),
		codes:        make(map[string]string),
		tokens:       make(map[string]*accessToken),
		dao:          dao,
	}
	var opts []op.Option

	opts = append(opts, op.WithAllowInsecure())

	provider, err := op.NewProvider(&op.Config{
		SupportedUILocales: []language.Tag{language.English},
		SupportedClaims: []string{
			"aud", "exp", "iat", "iss", "c_hash", "at_hash", "azp", // "scopes",
			"sub",
			"name", "family_name", "given_name",
			"email", "email_verified",
			"pls_*",
			"permissions",
			"year_tag",
		},
		SupportedScopes: SupportedScopes,
	}, &storage, op.StaticIssuer(issuer), opts...)
	if err != nil {
		logger.Error("failed to create provider", "error", err)
		os.Exit(1)
	}

	mux := http.NewServeMux()

	loginHandler := newLogin(&storage, op.AuthCallbackURL(provider), op.NewIssuerInterceptor(provider.IssuerFromRequest))

	mux.Handle("/login", loginHandler)
	mux.Handle("/", provider.Handler)

	mux.HandleFunc("GET /api/users", func(w http.ResponseWriter, r *http.Request) {

		query := r.URL.Query()

		format := query["format"]
		kthids := query["u"]

		switch format[0] {
		case "single":
			if len(kthids) != 1 {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			user, err := dao.GetUser(kthids[0])
			if err != nil {
				w.WriteHeader(http.StatusNotFound)
				return
			}
			resp := ssoUser{Email: user.Email, FirstName: user.FirstName, FamilyName: user.FamilyName, YearTag: user.YearTag}
			w.Header().Add("Content-Type", "application/json")
			json.NewEncoder(w).Encode(resp)
		case "array":
			users := []ssoUser{}
			for _, kthid := range kthids {
				user, err := dao.GetUser(kthid)
				if err != nil {
					users = append(users, ssoUser{})
				} else {
					users = append(users, ssoUser{Email: user.Email, FirstName: user.FirstName, FamilyName: user.FamilyName, YearTag: user.YearTag})
				}
			}
			w.Header().Add("Content-Type", "application/json")
			json.NewEncoder(w).Encode(users)
		case "map":
			users := make(map[string]ssoUser)
			for _, kthid := range kthids {
				user, err := dao.GetUser(kthid)
				if err == nil {
					users[kthid] = ssoUser{Email: user.Email, FirstName: user.FirstName, FamilyName: user.FamilyName, YearTag: user.YearTag}
				}
			}
			w.Header().Add("Content-Type", "application/json")
			json.NewEncoder(w).Encode(users)
		default:
			w.WriteHeader(http.StatusNotFound)
			return
		}
	})

	mux.HandleFunc("GET /api/search", func(w http.ResponseWriter, r *http.Request) {
		limitStr := r.FormValue("limit")
		if limitStr == "" {
			limitStr = "5"
		}
		i, err := strconv.ParseInt(limitStr, 10, 32)
		if err != nil {
			w.WriteHeader(http.StatusBadGateway)
			w.Write([]byte(err.Error()))
			return
		}
		limit := int32(i)

		offsetStr := r.FormValue("offset")
		if offsetStr == "" {
			offsetStr = "0"
		}
		i, err = strconv.ParseInt(offsetStr, 10, 32)
		if err != nil {
			w.WriteHeader(http.StatusBadGateway)
			w.Write([]byte(err.Error()))
			return
		}
		offset := int32(i)

		search := r.FormValue("query")
		year := r.FormValue("year")

		dbUsers := dao.ListUsers(search, year)

		limitedUsers := dbUsers[offset:(min(offset+limit, int32(len(dbUsers))))]

		type User struct {
			KTHID      string `json:"kthid"`
			Email      string `json:"email,omitempty"`
			FirstName  string `json:"firstName,omitempty"`
			FamilyName string `json:"familyName,omitempty"`
			YearTag    string `json:"yearTag,omitempty"`
		}

		users := make([]User, len(dbUsers))
		for i, user := range limitedUsers {
			users[i] = User{
				KTHID: user.KTHID,
				Email: user.Email,
				FirstName: user.FirstName,
				FamilyName: user.FamilyName,
				YearTag: user.YearTag,
			}
		}

		w.Header().Add("Content-Type", "application/json")
		json.NewEncoder(w).Encode(users)
	})

	handler := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			attr := slog.Int64("id", counter.Add(1))
			logger.With(attr).Debug("request", "method", r.Method, "url", r.URL.Path)

			next.ServeHTTP(w, r)
		})
	}(mux)

	server := &http.Server{
		Addr:    ":" + cfg.SsoPort,
		Handler: handler,
	}
	logger.Info("starting SSO server", "port", cfg.SsoPort)
	if server.ListenAndServe() != http.ErrServerClosed {
		logger.Error("server terminated")
		os.Exit(1)
	}
}

func newLogin(auth auth, callback func(context.Context, string) string, issuerInterceptor *op.IssuerInterceptor) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseForm()
		if err != nil {
			http.Error(w, fmt.Sprintf("cannot parse form:%s", err), http.StatusInternalServerError)
			return
		}
		renderLogin(w, r.FormValue("authRequestID"), nil)
	})

	mux.HandleFunc("POST /",
		issuerInterceptor.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				err := r.ParseForm()
				if err != nil {
					http.Error(w, fmt.Sprintf("cannot parse form:%s", err), http.StatusInternalServerError)
					return
				}
				id := r.FormValue("id")
				kthid := r.FormValue("kthid")
				err = auth.CheckLogin(kthid, id)
				if err != nil {
					renderLogin(w, id, err)
					return
				}
				http.Redirect(w, r, callback(r.Context(), id), http.StatusFound)
			}))

	return mux
}

func renderLogin(w http.ResponseWriter, id string, err error) {
	var errMsg string
	if err != nil {
		errMsg = err.Error()
	}

	data := &struct {
		ID    string
		Error string
	}{
		ID:    id,
		Error: errMsg,
	}
	err = loginTmpl.Execute(w, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
