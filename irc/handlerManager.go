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

func (irc *Connection) addOutboundHandlers(handlers []func(c *Connection, m *string)) {
	for _, handler := range handlers {
		irc.AddOutboundHandler(handler)
	}
}

func (irc *Connection) AddOutboundHandler(f func(c *Connection, m *string)) {
	if !irc.initialised {
		irc.Init()
	}
	irc.outboundHandlers = append(irc.outboundHandlers, f)
}

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
	irc.inboundHandlers[s] = append(irc.inboundHandlers[s], f)
}

func (irc *Connection) runInboundHandlers(m *Message) {
	handlers := irc.inboundHandlers[m.Verb]
	handlers = append(handlers, irc.inboundHandlers["*"]...)
	for _, handler := range handlers {
		go handler(irc, m)
	}
}

func (irc *Connection) runOutboundHandlers(m *string) {
	for _, handler := range irc.outboundHandlers {
		go handler(irc, m)
	}
}

func pong(c *Connection, m *Message) {
	c.SendRawf("PONG :%v", m.ParamsArray[0])
}

func quitOnError(c *Connection, _ *Message) {
	c.Finished <- true
}
