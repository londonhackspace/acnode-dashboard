package auth

import "github.com/londonhackspace/acnode-dashboard/config"

const (
	UserType_User = iota
	UserType_Machine = iota
)

func UserTypeName(ut int) string {
	switch ut {
	case UserType_User:
		return "User"
	case UserType_Machine:
		return "Machine"
	}

	return "Unknown"
}

type User struct {
	UserType int
	Name string `json:"name"`
	UserName string `json:"username"`
	Groups []string `json:"groups"`

	source string `json:"source"`
}

func (u *User) IsAdmin(config *config.Config) bool {
	for _,ga := range u.Groups {
		for _,gb := range config.AdminGroups {
			if ga == gb {
				return true
			}
		}
	}
	return false
}

