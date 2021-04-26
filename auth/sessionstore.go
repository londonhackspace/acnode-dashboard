package auth

import (
	"crypto/rand"
	"math/big"
)

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
		rn, err := rand.Int(rand.Reader, big.NewInt(int64(len(charBytes))))
		if err != nil {
			// We rarely panic but this is an awkward case to recover from
			panic(err)
		}
		ret[i] = charBytes[rn.Int64()]

	}

	return string(ret)
}