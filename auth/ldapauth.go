package auth

import (
	"crypto/tls"
	"errors"
	"fmt"
	ldapv3 "github.com/go-ldap/ldap/v3"
	"github.com/londonhackspace/acnode-dashboard/config"
	"github.com/rs/zerolog/log"
)

type LDAPAuthenticator struct {
	conf *config.Config
}

func GetLDAPAuthenticator(conf *config.Config) LDAPAuthenticator {
	return LDAPAuthenticator{
		conf: conf,
	}
}

func (ldap *LDAPAuthenticator) GetName() string {
	return "LDAP"
}

func (ldap *LDAPAuthenticator) LoginUser(username string, password string) (User,error) {
	cfg := tls.Config{InsecureSkipVerify: ldap.conf.LdapSkipTLSVerify}
	c, err := ldapv3.DialURL(ldap.conf.LdapServer,
						ldapv3.DialWithTLSConfig(&cfg))
	if err != nil {
		log.Err(err).Msg("Error connecting to LDAP server")
		return User{}, errors.New("Error contacting LDAP server")
	}
	defer c.Close()

	err = c.Bind(ldap.conf.LdapBindDN, ldap.conf.LdapBindPW)
	if err != nil {
		log.Err(err).Msg("Error binding with LDAP server")
		return User{}, errors.New("Error contacting LDAP server")
	}

	log.Debug().
		Str("Server", ldap.conf.LdapServer).
		Msg("Successfully Bound to LDAP server")

	// search for the user
	query := ldapv3.NewSearchRequest(ldap.conf.LdapUserOU + "," + ldap.conf.LdapBaseDN,
		ldapv3.ScopeSingleLevel, ldapv3.NeverDerefAliases,
		0, 0, false,
		fmt.Sprintf("(&(objectClass=organizationalPerson)(uid=%s))", ldapv3.EscapeFilter(username)),
		[]string{"dn","givenName"}, nil)

	res, err := c.Search(query)
	if err != nil {
		log.Err(err).Msg("Error searching LDAP server")
		return User{}, errors.New("Error searching LDAP server")
	}

	if len(res.Entries) != 1 {
		log.Err(err).
			Str("username", username).
			Msg("User Not Found")
		return User{}, errors.New("User not found")
	}

	err = c.Bind(res.Entries[0].DN, password)
	if err != nil {
		log.Err(err).
			Str("username", username).
			Msg("Invalid Password")
		return User{}, errors.New("Invalid Password")
	}

	return User{
		Name: res.Entries[0].GetAttributeValue("givenName"),
		UserType: UserType_User,
	}, nil
}