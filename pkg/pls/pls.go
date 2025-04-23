package pls

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/datasektionen/nyckeln-under-dorrmattan/pkg/config"
	"github.com/datasektionen/nyckeln-under-dorrmattan/pkg/dao"
)

func Listen(cfg *config.Config, dao *dao.Dao) {
	h := http.NewServeMux()

	h.HandleFunc("GET /api/user/{id}", func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		groups := dao.GetUserGroups(id)

		w.Header().Add("Content-Type", "application/json;charset=utf-8")
		json.NewEncoder(w).Encode(groups)
	})

	h.HandleFunc("GET /api/user/{id}/{group}", func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		group := r.PathValue("group")
		permissions := dao.GetUserPermissionsForGroup(id, group)

		w.Header().Add("Content-Type", "application/json;charset=utf-8")
		json.NewEncoder(w).Encode(permissions)
	})

	h.HandleFunc("GET /api/user/{id}/{group}/{permission}", func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		group := r.PathValue("group")
		permisison := r.PathValue("permission")
		has_permission := dao.HasPermission(id, group, permisison)
		w.Header().Add("Content-Type", "application/json;charset=utf-8")
		json.NewEncoder(w).Encode(has_permission)
	})

	fmt.Printf("pls listening on http://localhost:%s\n", cfg.PlsPort)
	if err := http.ListenAndServe(fmt.Sprintf(":%s", cfg.PlsPort), h); err != nil {
		panic(err)
	}
}
