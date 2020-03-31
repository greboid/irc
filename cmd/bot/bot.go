package main

import "github.com/greboid/irc/irc"

func addBotCallbacks(c *irc.Connection) {
	c.AddInboundHandler("001", joinChannels)
	c.AddInboundHandler("PRIVMSG", publishMessages)
}

func joinChannels(_ *irc.EventManager, c *irc.Connection, _ *irc.Message) {
	c.SendRawf("JOIN :%s", *Channel)
}

func publishMessages(em *irc.EventManager, _ *irc.Connection, m *irc.Message) {
	em.PublishChannelMessage(*m)
}