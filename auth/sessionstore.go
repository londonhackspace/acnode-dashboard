package auth

import "math/rand"

type SessionStore interface {
	AddUser(u *User) string
	RemoveUser(cookie string)
	GetUser(cookie string) *User
}

const charBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ01234567890"

func makeSessionCookieString() string {
	const length = 24
	ret := make([]byte, length)

	for i := range ret {
		ret[i] = charBytes[rand.Intn(len(charBytes))]
	}

	return string(ret)
}