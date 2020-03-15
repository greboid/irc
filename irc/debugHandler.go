package irc

import (
	"log"
	"strings"
)

type debugHandler struct {
	debug bool
}

func NewDebugHandler(debug bool) *debugHandler {
	return &debugHandler{
		debug: debug,
	}
}

func (h *debugHandler) install(c *Connection) {
	c.AddInboundHandler("*", h.handleMessage)
	c.AddOutboundHandler(h.handleOutboundMessage)
}

func (h *debugHandler) handleMessage(_ *Connection, m *Message) {
	if h.debug {
		log.Printf("In : %s %s %s", m.Source, m.Verb, strings.Join(m.Params, " "))
	}
}

func (h *debugHandler) handleOutboundMessage(_ *Connection, m string) {
	if h.debug {
		log.Printf("Out: %s", m)
	}
}
