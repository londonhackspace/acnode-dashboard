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

func (ldap *LDAPAuthenticator) GetWritable() bool {
	return false
}

func (ldap *LDAPAuthenticator) AddUser(username string, password string, usertype int) error {
	return errors.New("not implemented")
}

func (ldap *LDAPAuthenticator) AddUserToGroup(username string, group string) error {
	return errors.New("not implemented")
}

func (ldap *LDAPAuthenticator) getConnection() (*ldapv3.Conn,error) {
	cfg := tls.Config{InsecureSkipVerify: ldap.conf.LdapSkipTLSVerify}
	return ldapv3.DialURL(ldap.conf.LdapServer, ldapv3.DialWithTLSConfig(&cfg))
}

func (ldap *LDAPAuthenticator) makeLDAPQuery(req *ldapv3.SearchRequest, c *ldapv3.Conn) (*ldapv3.SearchResult, error) {
	var ourCon *ldapv3.Conn
	if c == nil {
		var err error
		ourCon, err = ldap.getConnection()
		if err != nil {
			log.Err(err).Msg("Error connecting to LDAP server")
			return nil, err
		}
		defer ourCon.Close()
	} else {
		ourCon = c
	}

	err := ourCon.Bind(ldap.conf.LdapBindDN, ldap.conf.LdapBindPW)
	if err != nil {
		log.Err(err).Msg("Error binding with LDAP server")
		return nil, err
	}

	log.Debug().
		Str("Server", ldap.conf.LdapServer).
		Msg("Successfully Bound to LDAP server")

	return ourCon.Search(req)
}

func (ldap *LDAPAuthenticator) getGroups(username string, c *ldapv3.Conn) []string {
	query := ldapv3.NewSearchRequest(ldap.conf.LdapGroupOU + ","+ ldap.conf.LdapBaseDN,
		ldapv3.ScopeWholeSubtree, ldapv3.NeverDerefAliases,
		0, 0, false,
		fmt.Sprintf("(&(objectClass=posixGroup)(memberUid=%s))", ldapv3.EscapeFilter(username)),
		[]string{"cn", "gidNumber"}, nil)
	res, err := ldap.makeLDAPQuery(query, c)
	if err != nil {
		log.Err(err).Msg("Error getting groups for user")
		return []string{}
	}

	var groups []string
	for _,entry := range(res.Entries) {
		groups = append(groups, entry.GetAttributeValue("cn"))
	}
	return groups
}

func (ldap *LDAPAuthenticator) LoginUser(username string, password string) (User,error) {
	c, err := ldap.getConnection()
	if err != nil {
		return User{}, errors.New("Error contacting LDAP server")
	}
	// search for the user
	query := ldapv3.NewSearchRequest(ldap.conf.LdapUserOU + "," + ldap.conf.LdapBaseDN,
		ldapv3.ScopeSingleLevel, ldapv3.NeverDerefAliases,
		0, 0, false,
		fmt.Sprintf("(&(objectClass=organizationalPerson)(uid=%s))", ldapv3.EscapeFilter(username)),
		[]string{"dn","givenName"}, nil)

	res, err := ldap.makeLDAPQuery(query, c)
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
		UserName: username,
		UserType: UserType_User,
		Groups: ldap.getGroups(username, c),
	}, nil
}

func (ldap *LDAPAuthenticator) GetUser(username string) (User, error)  {
	c, err := ldap.getConnection()
	if err != nil {
		return User{}, errors.New("Error contacting LDAP server")
	}
	// search for the user
	query := ldapv3.NewSearchRequest(ldap.conf.LdapUserOU + "," + ldap.conf.LdapBaseDN,
		ldapv3.ScopeSingleLevel, ldapv3.NeverDerefAliases,
		0, 0, false,
		fmt.Sprintf("(&(objectClass=organizationalPerson)(uid=%s))", ldapv3.EscapeFilter(username)),
		[]string{"dn","givenName"}, nil)

	res, err := ldap.makeLDAPQuery(query, c)
	if err != nil {
		log.Err(err).Msg("Error searching LDAP server")
		return User{}, err
	}

	if len(res.Entries) != 1 {
		log.Err(err).
			Str("username", username).
			Msg("User Not Found")
		return User{}, err
	}

	return User{
		Name: res.Entries[0].GetAttributeValue("givenName"),
		UserName: username,
		UserType: UserType_User,
		Groups: ldap.getGroups(username, c),
	}, nil
}