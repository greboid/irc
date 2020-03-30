package irc

import (
	"strings"
)

func (irc *Connection) addOutboundHandlers(handlers []func(em *EventManager, c *Connection, m string)) {
	for _, handler := range handlers {
		irc.AddOutboundHandler(handler)
	}
}

func (irc *Connection) AddOutboundHandler(f func(em *EventManager, c *Connection, m string)) {
	if !irc.initialised {
		irc.Init()
	}
	irc.outboundHandlers = append(irc.outboundHandlers, f)
}

func (irc *Connection) AddRawHandlers(handlers []func(c *Connection, m RawMessage)) {
	for _, handler := range handlers {
		irc.AddRawHandler(handler)
	}
}

func (irc *Connection) AddRawHandler(f func(c *Connection, m RawMessage)) {
	if !irc.initialised {
		irc.Init()
	}
	irc.rawHandlers = append(irc.rawHandlers, f)
}

func (irc *Connection) AddInboundHandlers(handlers map[string]func(em *EventManager, c *Connection, m *Message)) {
	for verb, handler := range handlers {
		irc.AddInboundHandler(verb, handler)
	}
}

func (irc *Connection) AddInboundHandler(s string, f func(em *EventManager, c *Connection, m *Message)) {
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
		go handler(&irc.listeners, irc, m)
	}
}

func (irc *Connection) runOutboundHandlers(m string) {
	for _, handler := range irc.outboundHandlers {
		go handler(&irc.listeners, irc, m)
	}
}

func (irc *Connection) runRawHandlers(m RawMessage) {
	for _, handler := range irc.rawHandlers {
		go handler(irc, m)
	}
}
