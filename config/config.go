package config

import (
	"encoding/json"
	"github.com/rs/zerolog/log"
	"io/ioutil"
	"os"
	"strings"
)

type Config struct {
	MqttServer   string `json:"mqtt_server,omitempty"`
	MqttClientId string `json:"mqtt_clientid,omitempty"`

	AcserverUrl    string `json:"acserver_url,omitempty"`
	AcserverApiKey string `json:"acserver_key,omitempty"`

	LdapEnable        bool   `json:"ldap_enable,omitempty"`
	LdapServer        string `json:"ldap_server,omitempty"`
	LdapBindDN        string `json:"ldap_binddn,omitempty"`
	LdapBindPW        string `json:"ldap_bindpw,omitempty"`
	LdapBaseDN        string `json:"ldap_basedn,omitempty"`
	LdapUserOU        string `json:"ldap_userou,omitempty"`
	LdapGroupOU       string `json:"ldap_groupou,omitempty"`
	LdapSkipTLSVerify bool   `json:"ldap_skipverify,omitempty"`

	RedisEnable bool   `json:"redis_enable,omitempty"`
	RedisServer string `json:"redis_server,omitempty"`

	LogJSON bool `json:"log_json"`

	AdminGroups []string `json:"admin_groups"`
}

func GetConfigurationFromEnvironment() Config {
	return Config{
		MqttServer:        os.Getenv("MQTT_SERVER"),
		MqttClientId:      os.Getenv("MQTT_CLIENTID"),
		AcserverUrl:       os.Getenv("ACSERVER_URL"),
		AcserverApiKey:    os.Getenv("ACSERVER_APIKEY"),
		LdapEnable:        strings.ToLower(os.Getenv("LDAP_ENABLE")) == "true",
		LdapServer:        os.Getenv("LDAP_SERVER"),
		LdapBindDN:        os.Getenv("LDAP_BINDDN"),
		LdapBindPW:        os.Getenv("LDAP_BINDPW"),
		LdapBaseDN:        os.Getenv("LDAP_BASEDN"),
		LdapUserOU:        os.Getenv("LDAP_USEROU"),
		LdapGroupOU:       os.Getenv("LDAP_GROUPOU"),
		LdapSkipTLSVerify: strings.ToLower(os.Getenv("LDAP_SKIPTLSVERIFY")) == "true",
		RedisEnable:       strings.ToLower(os.Getenv("REDIS_ENABLE")) == "true",
		RedisServer:       os.Getenv("REDIS_SERVER"),
		LogJSON:           strings.ToLower(os.Getenv("LOG_FMT_JSON")) == "true",
		AdminGroups:       strings.Split(os.Getenv("ADMIN_GROUPS"), ","),
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

	data, err := ioutil.ReadAll(f)
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

	if len(envvar.AcserverUrl) != 0 {
		combined.AcserverUrl = envvar.AcserverUrl
	} else {
		combined.AcserverUrl = fileconf.AcserverUrl
	}

	if len(envvar.AcserverApiKey) != 0 {
		combined.AcserverApiKey = envvar.AcserverApiKey
	} else {
		combined.AcserverApiKey = fileconf.AcserverApiKey
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

	combined.AdminGroups = envvar.AdminGroups
	combined.AdminGroups = append(combined.AdminGroups, fileconf.AdminGroups...)

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

	if os.Getenv("REDIS_ENABLE") != "" {
		combined.RedisEnable = envvar.RedisEnable
	} else {
		combined.RedisEnable = fileconf.RedisEnable
	}

	if envvar.RedisServer != "" {
		combined.RedisServer = envvar.RedisServer
	} else {
		combined.RedisServer = fileconf.RedisServer
	}

	// set sensible defaults where we can
	if combined.MqttClientId == "" {
		combined.MqttClientId = "ACNodeDash"
	}

	if combined.AcserverUrl == "" {
		combined.AcserverUrl = "https://acserver.london.hackspace.org.uk"
	}

	if combined.LdapUserOU == "" {
		combined.LdapUserOU = "ou=Users"
	}

	if combined.LdapGroupOU == "" {
		combined.LdapGroupOU = "ou=Groups"
	}

	if os.Getenv("LOG_FMT_JSON") != "" {
		combined.LogJSON = envvar.LogJSON
	} else {
		combined.LogJSON = fileconf.LogJSON
	}

	return combined
}

func (c *Config) Validate() bool {
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

	if c.RedisEnable {
		if c.RedisServer == "" {
			log.Error().Msg("Empty Redis Server")
			return false
		}
	}

	return true
}
