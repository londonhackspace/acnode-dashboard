package auth

type Provider interface {
	GetName() string
	LoginUser(username string, password string) (User,error)
}