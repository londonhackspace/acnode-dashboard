package auth

import (
	"github.com/rs/zerolog/log"
	"net/http"
)

// Authenticator Right now at least, this is a singleton
type Authenticator struct {
	providers []Provider
	sessionstore SessionStore
}

var authenticator *Authenticator = nil
func GetAuthenticator() *Authenticator {
	if authenticator == nil {
		sessionStore := CreateMemorySessionStore()
		newAuth := Authenticator{
			make([]Provider, 0),
			sessionStore,
		}
		authenticator = &newAuth
	}

	return authenticator
}

func SetSessionStore(ss SessionStore) {
	GetAuthenticator().sessionstore = ss
}

func AddProvider(p Provider) {
	auth := GetAuthenticator()
	auth.providers =  append(auth.providers, p)
}

// CheckAuthAPI Check for auth - either for a (human) user, or an API user
func CheckAuthAPI(w http.ResponseWriter, r* http.Request) (bool, *User) {
	// first see if there is a valid user cookie
	success, u := CheckAuthUser(w, r)
	if success {
		return true, u
	}

	return false,nil
}

func unsetAuthCookie(w http.ResponseWriter) {
	unsetter := http.Cookie{Name: "ACNodeDashboardSession", MaxAge: -1}
	http.SetCookie(w, &unsetter)
}

// CheckAuthUser Check for auth - endpoint only for (human) user use
func CheckAuthUser(w http.ResponseWriter, r *http.Request) (bool, *User) {
	c, err := r.Cookie("ACNodeDashboardSession")

	if err == nil {
		u := GetAuthenticator().sessionstore.GetUser(c.Value)
		if u != nil {
			return true, u
		} else {
			// remove junk cookie
			unsetAuthCookie(w)
		}
	}

	return false,nil
}

func Logout(w http.ResponseWriter, r *http.Request) {
	c,err := r.Cookie("ACNodeDashboardSession")

	if err == nil {
		GetAuthenticator().sessionstore.RemoveUser(c.Value)
		unsetAuthCookie(w)
	}
}

func AuthenticateUser(w http.ResponseWriter, username string, password string) bool {
	log.Info().
		Str("Username", username).
		Msg("Attempting authentication")

	authenticator := GetAuthenticator()
	for _, a := range authenticator.providers {
		u, err := a.LoginUser(username, password)
		if err == nil {
			cookieString := authenticator.sessionstore.AddUser(&u)
			cookie := http.Cookie{
				Name: "ACNodeDashboardSession",
				Path: "/",
				Value: cookieString,
				MaxAge: 0,
			}
			http.SetCookie(w, &cookie)
			log.Info().Str("Provider", a.GetName()).Msg("User authenticated successfully")
			return true
		}
	}
	log.Info().Msg("No providers authenticated user")
	return false
}