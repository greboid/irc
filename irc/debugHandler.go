package irc

import (
	"github.com/greboid/irc/config"
	"log"
)

type debugHandler struct {
	conf config.Config
}

func (h *debugHandler) install(c *Connection) {
	c.AddInboundHandler("*", h.handleMessage)
	c.AddOutboundHandler(h.handleOutboundMessage)
}

func (h *debugHandler) handleMessage(c *Connection, m *Message) {
	if h.conf.Debug {
		log.Printf("In : %s %s %s", m.Source, m.Verb, m.Params)
	}
}

func (h *debugHandler) handleOutboundMessage(_ *Connection, m *string) {
	log.Printf("Out: %s", *m)
}
