package kthldap

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/datasektionen/nyckeln-under-dorrmattan/pkg/config"
	"github.com/datasektionen/nyckeln-under-dorrmattan/pkg/dao"
)

func Listen(cfg *config.Config, dao *dao.Dao) {
	h := http.NewServeMux()

	h.HandleFunc("GET /user", func(w http.ResponseWriter, r *http.Request) {
		id := r.FormValue("kthid")
		user, err := dao.GetLdapUser(id)

		if err != nil {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		}

		w.Header().Add("Content-Type", "application/json;charset=utf-8")
		json.NewEncoder(w).Encode(user)
	})

	fmt.Printf("ldap listening on http://localhost:%s\n", cfg.LdapPort)
	if err := http.ListenAndServe(fmt.Sprintf(":%s", cfg.LdapPort), h); err != nil {
		panic(err)
	}
}
