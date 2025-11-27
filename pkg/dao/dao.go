package dao

import (
	"fmt"
	"os"
	"strings"
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
	Clients     []Client   `yaml:"clients"`
	Users       []User     `yaml:"users"`
	Hive        Hive       `yaml:"hive"`
	Ldap        []LdapUser `yaml:"ldap"`
	permissions map[string]map[string][]string
}

type Client struct {
	Id           string   `yaml:"id"`
	Secret       string   `yaml:"secret"`
	AllowsGuests  bool     `yaml:"allows_guests"`
	RedirectURIs []string `yaml:"redirect_uris"`
}

type Hive struct {
	Tokens []HiveToken `yaml:"tokens"`
	Groups []HiveGroup `yaml:"groups"`
}

type HiveToken struct {
	Secret      string           `yaml:"secret"`
	Permissions []HivePermission `yaml:"permissions"`
}

type HiveGroup struct {
	Name        string           `yaml:"name"`
	Id          string           `yaml:"id"`
	Domain      string           `yaml:"domain"`
	Members     []string         `yaml:"members"`
	Tags        []HiveTag        `yaml:"tags"`
	Permissions []HivePermission `yaml:"permissions"`
}

type HiveTag struct {
	Id      string `yaml:"id"`
	Content string `yaml:"content"`
}

type HivePermission struct {
	Id    string `yaml:"id" json:"id"`
	Scope string `yaml:"scope" json:"scope"`
}

type HiveTagGroup struct {
	GroupName   string `json:"group_name"`
	GroupId     string `json:"group_id"`
	GroupDomain string `json:"group_domain"`
	TagContent  string `json:"tag_content"`
}

type HiveTagUser struct {
	Username   string `json:"username"`
	TagContent string `json:"tag_content"`
}

type User struct {
	KTHID                   string              `yaml:"kth_id"`
	UGKTHID                 string              `yaml:"ug_kth_id"`
	Email                   string              `yaml:"email"`
	FirstName               string              `yaml:"first_name"`
	FamilyName              string              `yaml:"family_name"`
	Picture                 string              `yaml:"picture"`
	Thumbnail               string              `yaml:"thumbnail"`
	YearTag                 string              `yaml:"year_tag"`
	MemberTo                time.Time           `yaml:"member_to"`
	WebAuthnID              []byte              `yaml:"web_authn_id"`
	FirstNameChangeRequest  string              `yaml:"first_name_change_request"`
	FamilyNameChangeRequest string              `yaml:"family_name_change_request"`
	PlsPermissions          map[string][]string `yaml:"pls_permissions"`
	HiveTags                []HiveTag           `yaml:"hive_tags"`
}

type LdapUser struct {
	KTHID      string `yaml:"kth_id" json:"kthid"`
	UGKTHID    string `yaml:"ug_kth_id" json:"ug_kthid"`
	FirstName  string `yaml:"first_name" json:"first_name"`
	FamilyName string `yaml:"family_name" json:"last_name"`
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

func (d *Dao) GetLdapUser(kthid string) (*LdapUser, error) {
	for _, user := range d.db.Ldap {
		if user.KTHID == kthid {
			return &user, nil
		}
	}
	return nil, fmt.Errorf("ldap user not found")
}

func (d *Dao) ListUsers(query string, year string) []User {
	users := []User{}
	for _, user := range d.db.Users {
		if (user.KTHID == query ||
			strings.Contains(user.FirstName, query) ||
			strings.Contains(user.FamilyName, query) ||
			strings.Contains(user.FirstName+" "+user.FamilyName, query)) &&
			(year == user.YearTag || year == "") {

			users = append(users, user)
		}
	}
	return users
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

func (d *Dao) GetHivePermissions(kthid string) []HivePermission {
	permissions := []HivePermission{}
	for _, group := range d.db.Hive.Groups {
		for _, user := range group.Members {
			if user == kthid {
				permissions = append(permissions, group.Permissions...)
			}
		}
	}
	return permissions
}

func (d *Dao) GetHivePermissionsToken(tokenId string) []HivePermission {
	permissions := []HivePermission{}

	for _, token := range d.db.Hive.Tokens {
		if token.Secret == tokenId {
			for _, permission := range token.Permissions {
				permissions = append(permissions, permission)
			}
		}
	}

	return permissions
}

func (d *Dao) GetHiveTagGroups(tagId string) []HiveTagGroup {
	tagGroups := []HiveTagGroup{}
	for _, group := range d.db.Hive.Groups {
		for _, tag := range group.Tags {
			if tag.Id == tagId {
				tagGroups = append(tagGroups, HiveTagGroup{GroupName: group.Name, GroupId: group.Id, GroupDomain: group.Domain, TagContent: tag.Content})
			}
		}
	}

	return tagGroups
}

func (d *Dao) GetHiveTagGroupsUser(tagId string, kthid string) []HiveTagGroup {
	tagGroups := []HiveTagGroup{}
	for _, group := range d.db.Hive.Groups {
		for _, user := range group.Members {
			if user == kthid {
				for _, tag := range group.Tags {
					if tag.Id == tagId {
						tagGroups = append(tagGroups, HiveTagGroup{GroupName: group.Name, GroupId: group.Id, GroupDomain: group.Domain, TagContent: tag.Content})
					}
				}
			}
		}
	}

	return tagGroups
}

func (d *Dao) GetHiveUsersWithTag(tagId string) []HiveTagUser {
	tagUsers := []HiveTagUser{}

	for _, user := range d.db.Users {
		for _, tag := range user.HiveTags {
			if tag.Id == tagId {
				tagUsers = append(tagUsers, HiveTagUser{Username: user.KTHID, TagContent: tag.Content})
			}
		}
	}

	return tagUsers
}

func (d *Dao) GetHiveMembership(groupDomain string, groupId string) []string {
	members := make([]string, 0)

	for _, group := range d.db.Hive.Groups {
		if group.Domain == groupDomain && group.Id == groupId {
			for _, user := range group.Members {
				members = append(members, user)
			}
		}
	}

	return members
}
