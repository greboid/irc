package irc

import (
	"log"
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
		log.Printf("In : %s", m.Raw)
	}
}

func (h *debugHandler) handleOutboundMessage(_ *Connection, m string) {
	if h.debug {
		log.Printf("Out: %s", m)
	}
}
