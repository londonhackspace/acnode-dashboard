package auth

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

	source string `json:"source"`
}

func CreateUser(usertype int, name string) User {
	return User{
		UserType: usertype,
		Name: name,
	}
}
