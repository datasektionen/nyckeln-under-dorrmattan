package dao

import (
	"fmt"
	"os"
	"time"

	"maps"

	"github.com/datasektionen/nyckeln-under-dorrmattan/pkg/config"
	"golang.org/x/exp/slices"
	"gopkg.in/yaml.v3"
)

type Dao struct {
	cfg *config.Config
	db  db
}

type db struct {
	Clients     []Client `yaml:"clients"`
	Users       []User   `yaml:"users"`
	permissions map[string]map[string][]string
}

type Client struct {
	Id           string   `yaml:"id"`
	Secret       string   `yaml:"secret"`
	RedirectURIs []string `yaml:"redirect_uris"`
}

type User struct {
	KTHID                   string              `yaml:"kth_id"`
	UGKTHID                 string              `yaml:"ug_kth_id"`
	Email                   string              `yaml:"email"`
	FirstName               string              `yaml:"first_name"`
	FamilyName              string              `yaml:"family_name"`
	YearTag                 string              `yaml:"year_tag"`
	MemberTo                time.Time           `yaml:"member_to"`
	WebAuthnID              []byte              `yaml:"web_authn_id"`
	FirstNameChangeRequest  string              `yaml:"first_name_change_request"`
	FamilyNameChangeRequest string              `yaml:"family_name_change_request"`
	PlsPermissions          map[string][]string `yaml:"pls_permissions"`
}

func New(cfg *config.Config) *Dao {
	file, err := os.ReadFile(cfg.ConfigFile)
	if err != nil {
		panic(err)
	}
	var db db
	err = yaml.Unmarshal(file, &db)
	if err != nil {
		panic(err)
	}

	db.permissions = make(map[string]map[string][]string)
	for _, user := range db.Users {
		db.permissions[user.KTHID] = make(map[string][]string)
		maps.Copy(db.permissions[user.KTHID], user.PlsPermissions)
	}

	return &Dao{
		db: db,
	}
}

func (d *Dao) GetClient(id string) (*Client, error) {
	for _, client := range d.db.Clients {
		if client.Id == id {
			return &client, nil
		}
	}
	return nil, fmt.Errorf("client not found")
}

func (d *Dao) GetUser(kthid string) (*User, error) {
	for _, user := range d.db.Users {
		if user.KTHID == kthid {
			return &user, nil
		}
	}
	return nil, fmt.Errorf("user not found")
}

func (d *Dao) GetUserPermissionsForGroup(kthid string, group string) []string {
	if permissions, ok := d.db.permissions[kthid][group]; ok {
		return permissions
	}
	return []string{}
}

func (d *Dao) GetUserGroups(kthid string) map[string][]string {
	if groups, ok := d.db.permissions[kthid]; ok {
		return groups
	}
	return map[string][]string{}
}

func (d *Dao) HasPermission(kthid string, group string, permission string) bool {
	if groups, ok := d.db.permissions[kthid][group]; ok {
		return slices.Contains(groups, permission)
	}
	return false
}
