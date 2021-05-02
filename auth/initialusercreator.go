package auth

import "github.com/rs/zerolog/log"

func CreateInitialUser(p Provider) {
	if !p.GetWritable() {
		log.Error().Msg("Cannot create initial user - Provider not writable")
		return
	}
	pw := makeSessionCookieString()
	err := p.AddUser("admin", pw, UserType_User)
	if err != nil {
		log.Err(err).Msg("")
	} else {
		err = p.AddUserToGroup("admin", "Admins")
		if err != nil {
			log.Err(err).Msg("")
		} else {
			log.Info().Str("username", "admin").Str("password", pw).Msg("Created initial user")
		}
	}
}
