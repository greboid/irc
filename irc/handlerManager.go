package irc

import (
	"strings"
)

var (
	defaultInboundHandlers = map[string]func(c *Connection, m *Message){
		"PING":  pong,
		"ERROR": quitOnError,
	}
)

func (irc *Connection) AddInboundHandlers(handlers map[string]func(c *Connection, m *Message)) {
	for verb, handler := range handlers {
		irc.AddInboundHandler(verb, handler)
	}
}

func (irc *Connection) AddInboundHandler(s string, f func(c *Connection, m *Message)) {
	if !irc.initialised {
		irc.Init()
	}
	s = strings.ToUpper(s)
	irc.handlers[s] = append(irc.handlers[s], f)
}

func (irc *Connection) runCallbacks(m *Message) {
	handlers := irc.handlers[m.Verb]
	handlers = append(handlers, irc.handlers["*"]...)
	for _, handler := range handlers {
		go handler(irc, m)
	}
}

func pong(c *Connection, m *Message) {
	c.SendRawf("PONG :%v", m.ParamsArray[0])
}

func quitOnError(c *Connection, _ *Message) {
	c.Finished <- true
}
