package config

import (
	"encoding/json"
	"github.com/rs/zerolog/log"
	"io/ioutil"
	"os"
	"strings"
)

type Config struct {
	MqttServer string `json:"mqtt_server",omitempty`
	MqttClientId string `json:"mqtt_clientid",omitempty`

	LdapEnable bool `json:"ldap_enable",omitempty`
	LdapServer string `json:"ldap_server",omitempty`
	LdapBindDN string `json:"ldap_binddn",omitempty`
	LdapBindPW string `json:"ldap_bindpw",omitempty`
	LdapBaseDN string `json:"ldap_basedn",omitempty`
	LdapUserOU string `json:"ldap_userou",omitempty`
	LdapGroupOU string `json:"ldap_groupou",omitempty`
	LdapSkipTLSVerify bool `json:"ldap_skipverify",omitempty`
}

func GetConfigurationFromEnvironment() Config {
	return Config{
		MqttServer: os.Getenv("MQTT_SERVER"),
		MqttClientId: os.Getenv("MQTT_CLIENTID"),
		LdapEnable: strings.ToLower(os.Getenv("LDAP_ENABLE")) == "true",
		LdapServer: os.Getenv("LDAP_SERVER"),
		LdapBindDN: os.Getenv("LDAP_BINDDN"),
		LdapBindPW: os.Getenv("LDAP_BINDPW"),
		LdapBaseDN: os.Getenv("LDAP_BASEDN"),
		LdapUserOU: os.Getenv("LDAP_USEROU"),
		LdapGroupOU: os.Getenv("LDAP_GROUPOU"),
		LdapSkipTLSVerify: strings.ToLower(os.Getenv("LDAP_SKIPTLSVERIFY")) == "true",
	}
}

func GetConfigurationFromFile(filename string) Config {
	f, err := os.Open(filename)
	if err != nil {
		log.Info().
			Str("filename", filename).
			Msg("Unable to read configuration from file")
		return Config{}
	}
	defer f.Close()

	data,err := ioutil.ReadAll(f)
	if err != nil {
		return Config{}
	}

	c := Config{}

	json.Unmarshal(data, &c)

	return c
}

func GetCombinedConfig(filename string) Config {
	combined := Config{}

	envvar := GetConfigurationFromEnvironment()
	fileconf := GetConfigurationFromFile(filename)

	if len(envvar.MqttServer) != 0 {
		combined.MqttServer = envvar.MqttServer
	} else {
		combined.MqttServer = fileconf.MqttServer
	}

	if os.Getenv("MQTT_CLIENTID") != "" {
		combined.MqttClientId = envvar.MqttClientId
	} else {
		combined.MqttClientId = fileconf.MqttClientId
	}


	if len(envvar.LdapServer) != 0 {
		combined.LdapServer = envvar.LdapServer
	} else {
		combined.LdapServer = fileconf.LdapServer
	}

	if len(envvar.LdapBaseDN) != 0 {
		combined.LdapBaseDN = envvar.LdapBaseDN
	} else {
		combined.LdapBaseDN = fileconf.LdapBaseDN
	}

	if len(envvar.LdapBindDN) != 0 {
		combined.LdapBindDN = envvar.LdapBindDN
	} else {
		combined.LdapBindDN = fileconf.LdapBindDN
	}

	if len(envvar.LdapBindPW) != 0 {
		combined.LdapBindPW = envvar.LdapBindPW
	} else {
		combined.LdapBindPW = fileconf.LdapBindPW
	}

	if os.Getenv("LDAP_SKIPTLSVERIFY") != "" {
		combined.LdapSkipTLSVerify = envvar.LdapSkipTLSVerify
	} else {
		combined.LdapSkipTLSVerify = fileconf.LdapSkipTLSVerify
	}

	if os.Getenv("LDAP_ENABLE") != "" {
		combined.LdapEnable = envvar.LdapEnable
	} else {
		combined.LdapEnable = fileconf.LdapEnable
	}

	if os.Getenv("LDAP_USEROU") != "" {
		combined.LdapUserOU = envvar.LdapUserOU
	} else {
		combined.LdapUserOU = fileconf.LdapUserOU
	}

	if os.Getenv("LDAP_GROUPOU") != "" {
		combined.LdapGroupOU = envvar.LdapGroupOU
	} else {
		combined.LdapGroupOU = fileconf.LdapGroupOU
	}

	// set sensible defaults where we can
	if combined.MqttClientId == "" {
		combined.MqttClientId = "ACNodeDash"
	}

	if combined.LdapUserOU == "" {
		combined.LdapUserOU = "ou=Users"
	}

	if combined.LdapGroupOU == "" {
		combined.LdapGroupOU = "ou=Groups"
	}

	return combined
}

func (c *Config) Validate()  bool {
	if c.MqttServer == "" {
		return false
	}

	if c.LdapEnable {
		if c.LdapServer == "" {
			log.Error().Msg("Empty LDAP server")
			return false
		}

		if c.LdapBaseDN == "" {
			log.Error().Msg("Empty LDAP BaseDN")
			return false
		}

		if c.LdapBindDN == "" {
			log.Error().Msg("Empty LDAP BindDN")
			return false
		}

		if c.LdapBindPW == "" {
			log.Error().Msg("Empty LDAP BindPW")
			return false
		}
	}

	return true
}