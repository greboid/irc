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
	c.AddRawHandler(h.handleRawMessage)
}

func (h *debugHandler) handleRawMessage(c *Connection, m RawMessage) {
	if m.out {
		h.handleOutboundMessage(c, m)
	} else {
		h.handleMessage(c, m)
	}
}

func (h *debugHandler) handleMessage(_ *Connection, m RawMessage) {
	if h.debug {
		log.Printf("In : %s", m.message)
	}
}

func (h *debugHandler) handleOutboundMessage(_ *Connection, m RawMessage) {
	if h.debug {
		log.Printf("Out: %s", m.message)
	}
}
