package auth

import (
	"github.com/rs/zerolog/log"
	"golang.org/x/crypto/bcrypt"
	"strings"
)

type passwordHasherCfg struct {
	hash       string
	workFactor int
	salt       string
}

func getDefaultHasherConfig() passwordHasherCfg {
	return passwordHasherCfg{
		hash:       "bcrypt",
		workFactor: 15,
		// TODO: use a better salt source
		salt: makeSessionCookieString(),
	}
}

func parseHasherConfigFromPwdField(field string) passwordHasherCfg {
	parts := strings.Split(field, ":")

	if len(parts) < 2 {
		log.Error().Msg("Error parsing password field")
		return passwordHasherCfg{}
	}

	if parts[0] == "bcrypt" && len(parts) == 3 {
		cfg := passwordHasherCfg{}
		cfg.hash = parts[0]
		cfg.salt = parts[1]
		bcrypt.Cost([]byte(parts[2]))
		return cfg
	}
	log.Error().Str("hasher", parts[0]).Msg("Unknown password hasher")

	return passwordHasherCfg{}
}

func makePwField(cfg passwordHasherCfg, hash []byte) string {
	if cfg.hash == "bcrypt" {
		return cfg.hash + ":" + cfg.salt + ":" + string(hash)
	}

	log.Error().Str("hasher", cfg.hash).Msg("Unknown password hasher")

	return ""
}

func hashPassword(pwd string, cfg passwordHasherCfg) string {
	working := pwd + cfg.salt
	result := []byte{}
	if cfg.hash == "bcrypt" {
		var err error
		result, err = bcrypt.GenerateFromPassword([]byte(working), cfg.workFactor)
		if err != nil {
			result = []byte{}
		}

	} else {
		log.Error().Str("hasher", cfg.hash).Msg("Unknown password hasher")
	}

	return makePwField(cfg, result)
}

func checkPassword(given string, expected string) bool {
	parts := strings.Split(expected, ":")

	if len(parts) < 2 {
		log.Error().Msg("Error parsing password field")
		return false
	}

	if parts[0] == "bcrypt" && len(parts) == 3 {
		// this returns nil if the hash matched
		return bcrypt.CompareHashAndPassword([]byte(parts[2]), []byte(given+parts[1])) == nil
	}

	return false
}
