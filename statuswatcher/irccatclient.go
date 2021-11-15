package main

import (
	"github.com/rs/zerolog/log"
	"net"
)

type IRCCatClient struct {
	server string
	channel string
}

func MakeIRCCatClient(server string, channel string) *IRCCatClient {
	return &IRCCatClient{
		server:  server,
		channel: channel,
	}
}

func (irccat *IRCCatClient) SendMessage(message string) {
	conn, err := net.Dial("tcp", irccat.server)
	if err != nil {
		log.Err(err).Msg("Error connecting to IRCCat")
		return
	}
	defer conn.Close()

	conn.Write([]byte(irccat.channel + " " + message))
}