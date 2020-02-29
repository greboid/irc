package irc

import (
	"strings"
)

var (
	defaultCallbacks = map[string]func(c *Connection, m *Message){
		"PING":  pong,
		"ERROR": quitOnError,
	}
)

func (irc *Connection) AddCallbacks(callbacks map[string]func(c *Connection, m *Message)) {
	for verb, callback := range callbacks {
		irc.AddCallback(verb, callback)
	}
}

func (irc *Connection) AddCallback(s string, f func(c *Connection, m *Message)) {
	if !irc.initialised {
		irc.Init()
	}
	s = strings.ToUpper(s)
	irc.callbacks[s] = append(irc.callbacks[s], f)
}

func (irc *Connection) runCallbacks(m *Message) {
	callbacks := irc.callbacks[m.Verb]
	for _, callback := range callbacks {
		go callback(irc, m)
	}
}

func pong(c *Connection, m *Message) {
	c.SendRawf("PONG :%v", m.Params)
}

func quitOnError(c *Connection, _ *Message) {
	c.Finished <- true
}