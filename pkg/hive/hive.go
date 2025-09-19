package hive

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/datasektionen/nyckeln-under-dorrmattan/pkg/config"
	"github.com/datasektionen/nyckeln-under-dorrmattan/pkg/dao"
)

func Listen(cfg *config.Config, dao *dao.Dao) {
	h := http.NewServeMux()

	// User endpoints

	h.HandleFunc("GET /api/v1/user/{id}/permissions", func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if !strings.Contains(authHeader, "Bearer") {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		id := r.PathValue("id")
		permissions := dao.GetHivePermissions(id)

		w.Header().Add("Content-Type", "application/json;charset=utf-8")
		json.NewEncoder(w).Encode(permissions)
	})

	h.HandleFunc("GET /api/v1/user/{id}/permission/{perm_id}", func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if !strings.Contains(authHeader, "Bearer") {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		id := r.PathValue("id")
		perm_id := r.PathValue("perm_id")
		permissions := dao.GetHivePermissions(id)

		hasPerm := false

		for _, perm := range permissions {
			if perm.Id == perm_id {
				hasPerm = true
			}
		}

		w.Header().Add("Content-Type", "application/json;charset=utf-8")
		json.NewEncoder(w).Encode(hasPerm)
	})

	h.HandleFunc("GET /api/v1/user/{id}/permission/{perm_id}/scopes", func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if !strings.Contains(authHeader, "Bearer") {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		id := r.PathValue("id")
		permId := r.PathValue("perm_id")
		permissions := dao.GetHivePermissions(id)

		var scopes []string

		for _, perm := range permissions {
			if perm.Id == permId {
				scopes = append(scopes, perm.Scope)
			}
		}

		w.Header().Add("Content-Type", "application/json;charset=utf-8")
		json.NewEncoder(w).Encode(scopes)
	})

	h.HandleFunc("GET /api/v1/user/{id}/permission/{perm_id}/scope/{scope}", func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if !strings.Contains(authHeader, "Bearer") {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		id := r.PathValue("id")
		perm_id := r.PathValue("perm_id")
		scope := r.PathValue("scope")
		permissions := dao.GetHivePermissions(id)

		hasScope := false

		for _, perm := range permissions {
			if perm.Id == perm_id && perm.Scope == scope {
				hasScope = true
			}
		}

		w.Header().Add("Content-Type", "application/json;charset=utf-8")
		json.NewEncoder(w).Encode(hasScope)
	})

	// Token endpoints

	h.HandleFunc("GET /api/v1/token/{secret}/permissions", func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if !strings.Contains(authHeader, "Bearer") {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		secret := r.PathValue("secret")
		permissions := dao.GetHivePermissionsToken(secret)

		w.Header().Add("Content-Type", "application/json;charset=utf-8")
		json.NewEncoder(w).Encode(permissions)
	})

	h.HandleFunc("GET /api/v1/token/{secret}/permission/{perm_id}", func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if !strings.Contains(authHeader, "Bearer") {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		secret := r.PathValue("secret")
		perm_id := r.PathValue("perm_id")
		permissions := dao.GetHivePermissionsToken(secret)

		hasPerm := false

		for _, perm := range permissions {
			if perm.Id == perm_id {
				hasPerm = true
			}
		}

		w.Header().Add("Content-Type", "application/json;charset=utf-8")
		json.NewEncoder(w).Encode(hasPerm)
	})

	h.HandleFunc("GET /api/v1/token/{secret}/permission/{perm_id}/scopes", func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if !strings.Contains(authHeader, "Bearer") {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		secret := r.PathValue("secret")
		perm_id := r.PathValue("perm_id")
		permissions := dao.GetHivePermissionsToken(secret)

		scopes := make([]string, 0)

		for _, perm := range permissions {
			if perm.Id == perm_id {
				scopes = append(scopes, perm.Scope)
			}
		}

		w.Header().Add("Content-Type", "application/json;charset=utf-8")
		json.NewEncoder(w).Encode(scopes)
	})

	h.HandleFunc("GET /api/v1/token/{secret}/permission/{perm_id}/scope/{scope}", func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if !strings.Contains(authHeader, "Bearer") {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		secret := r.PathValue("secret")
		perm_id := r.PathValue("perm_id")
		scope := r.PathValue("scope")
		permissions := dao.GetHivePermissionsToken(secret)

		hasScope := false

		for _, perm := range permissions {
			if perm.Id == perm_id && perm.Scope == scope {
				hasScope = true
			}
		}

		w.Header().Add("Content-Type", "application/json;charset=utf-8")
		json.NewEncoder(w).Encode(hasScope)
	})

	// Tag endpoints

	h.HandleFunc("GET /api/v1/tagged/{tag_id}/groups", func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if !strings.Contains(authHeader, "Bearer") {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		tagId := r.PathValue("tag_id")
		groups := dao.GetHiveTagGroups(tagId)

		w.Header().Add("Content-Type", "application/json;charset=utf-8")
		json.NewEncoder(w).Encode(groups)
	})

	h.HandleFunc("GET /api/v1/tagged/{tag_id}/memberships/{username}", func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if !strings.Contains(authHeader, "Bearer") {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		tagId := r.PathValue("tag_id")
		username := r.PathValue("username")
		groups := dao.GetHiveTagGroupsUser(tagId, username)

		w.Header().Add("Content-Type", "application/json;charset=utf-8")
		json.NewEncoder(w).Encode(groups)
	})

	h.HandleFunc("GET /api/v1/tagged/{tag_id}/users", func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if !strings.Contains(authHeader, "Bearer") {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		tagId := r.PathValue("tag_id")
		users := dao.GetHiveUsersWithTag(tagId)

		w.Header().Add("Content-Type", "application/json;charset=utf-8")
		json.NewEncoder(w).Encode(users)
	})

	h.HandleFunc("GET /api/v1/group/{group_domain}/{group_id}/members", func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if !strings.Contains(authHeader, "Bearer") {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		groupDomain := r.PathValue("group_domain")
		groupId := r.PathValue("group_id")
		users := dao.GetHiveMembership(groupDomain, groupId)

		w.Header().Add("Content-Type", "application/json;charset=utf-8")
		json.NewEncoder(w).Encode(users)
	})

	fmt.Printf("hive listening on http://localhost:%s\n", cfg.HivePort)
	if err := http.ListenAndServe(fmt.Sprintf(":%s", cfg.HivePort), h); err != nil {
		panic(err)
	}
}
