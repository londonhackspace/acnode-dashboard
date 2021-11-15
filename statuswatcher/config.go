package main

import "os"

type Config struct {
	ApiKey string
	HTTPListen string
	ACNodeDashAPI string
	SlackToken string
	SlackChannel string
}

func GetConfig() Config {
	c:= Config{
		ApiKey: os.Getenv("API_KEY"),
		HTTPListen: os.Getenv("HTTP_LISTEN"),
		ACNodeDashAPI: os.Getenv("ACNODE_DASH"),
		SlackToken: os.Getenv("SLACK_TOKEN"),
		SlackChannel: os.Getenv("SLACK_CHANNEL"),
	}

	if len(c.HTTPListen) == 0 {
		c.HTTPListen = ":8081"
	}

	if len(c.ACNodeDashAPI) == 0 {
		c.ACNodeDashAPI = "https://acnode-dash.london.hackspace.org.uk/api/"
	}

	if len(c.SlackChannel) == 0 {
		c.SlackChannel = "bot-test"
	}

	return c
}
