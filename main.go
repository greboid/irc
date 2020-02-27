package main

import (
	"github.com/greboid/irc/irc"
	"log"
	"strings"
)

func main() {
	config := getConfig()
	connection := irc.IRCConnection{
		ClientConfig: irc.ClientConfig{
			Server:   config.server,
			Password: config.password,
			Nick:     config.nickname,
			User:     config.nickname,
			UseTLS:   true,
		},
		ConnConfig: irc.DefaultConnectionConfig,
	}
	web := Web{config.channel, config.WebPort, &connection}
	go web.StartWeb()
	connection.Init()
	//Add some callbacks
	connection.AddCallback("001", func(c *irc.IRCConnection, m *irc.Message) {
		c.SendRawf("JOIN %s", config.channel)
	})
	connection.AddCallback("PRIVMSG", func(c *irc.IRCConnection, m *irc.Message) {
		if strings.ToLower(m.Params) == config.channel+" hi" {
			c.SendRawf("PRIVMSG %s Hey.", config.channel)
		}
	})
	connection.AddCallback("PRIVMSG", func(c *irc.IRCConnection, m *irc.Message) {
		if strings.ToLower(m.Params) == config.channel+" bye" {
			c.Quit()
		}
	})
	err := connection.ConnectAndWait()
	if err != nil {
		log.Fatal(err)
	}
}
