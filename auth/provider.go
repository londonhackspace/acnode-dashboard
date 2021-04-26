package auth

type Provider interface {
	GetName() string
	GetWritable() bool
	LoginUser(username string, password string) (User,error)

	AddUser(username string, password string, usertype int) error
}