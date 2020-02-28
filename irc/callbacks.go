package irc

import (
	"strings"
)

var (
	pong = func(c *IRCConnection, m *Message) {
		c.SendRawf("PONG :%v", m.Params)
	}
	quitOnError = func(c *IRCConnection, m *Message) {
		c.Finished <- true
	}
	defaultCallbacks = map[string]func(c *IRCConnection, m *Message){
		"PING":  pong,
		"ERROR": quitOnError,
	}
)

func (irc *IRCConnection) AddCallbacks(callbacks map[string]func(c *IRCConnection, m *Message)) {
	for verb, callback := range callbacks {
		irc.AddCallback(verb, callback)
	}
}

func (irc *IRCConnection) AddCallback(s string, f func(c *IRCConnection, m *Message)) {
	if !irc.initialised {
		irc.Init()
	}
	s = strings.ToUpper(s)
	irc.callbacks[s] = append(irc.callbacks[s], f)
}

func (irc *IRCConnection) runCallbacks(m *Message) {
	callbacks := irc.callbacks[m.Verb]
	for _, callback := range callbacks {
		go callback(irc, m)
	}
}
