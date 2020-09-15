package main

import (
	"github.com/greboid/irc/irc"
	"strings"
)

func addBotCallbacks(c *irc.Connection) {
	c.AddInboundHandler("001", joinChannels)
	c.AddInboundHandler("PRIVMSG", publishMessages)
}

func joinChannels(_ *irc.EventManager, c *irc.Connection, _ *irc.Message) {
	channels := strings.Split(*Channel, ",")
	for index := range channels {
		c.SendRawf("JOIN :%s", channels[index])
	}
}

func publishMessages(em *irc.EventManager, _ *irc.Connection, m *irc.Message) {
	em.PublishChannelMessage(*m)
}
