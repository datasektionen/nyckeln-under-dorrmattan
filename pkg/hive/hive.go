package hive

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/datasektionen/nyckeln-under-dorrmattan/pkg/config"
	"github.com/datasektionen/nyckeln-under-dorrmattan/pkg/dao"
)

func Listen(cfg *config.Config, dao *dao.Dao) {
	h := http.NewServeMux()

	h.HandleFunc("GET /api/v1/user/{id}/permissions", func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader != "Bearer fake-hive-token" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		id := r.PathValue("id")
		permissions := dao.GetHivePermissions(id)

		w.Header().Add("Content-Type", "application/json;charset=utf-8")
		json.NewEncoder(w).Encode(permissions)
	})

	fmt.Printf("hive listening on http://localhost:%s\n", cfg.HivePort)
	if err := http.ListenAndServe(fmt.Sprintf(":%s", cfg.HivePort), h); err != nil {
		panic(err)
	}
}
