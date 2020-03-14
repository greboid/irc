package irc

import (
	"log"
)

type debugHandler struct {
	debug bool
}

func (h *debugHandler) install(c *Connection) {
	c.AddInboundHandler("*", h.handleMessage)
	c.AddOutboundHandler(h.handleOutboundMessage)
}

func (h *debugHandler) handleMessage(_ *Connection, m *Message) {
	if h.debug {
		log.Printf("In : %s %s %s", m.Source, m.Verb, m.Params)
	}
}

func (h *debugHandler) handleOutboundMessage(_ *Connection, m string) {
	if h.debug {
		log.Printf("Out: %s", m)
	}
}
