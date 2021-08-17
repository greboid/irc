package irc

import (
	"github.com/ergochat/irc-go/ircmsg"
)

type pingHandler struct {
}

func NewPingHandler() *pingHandler {
	return &pingHandler{}
}

func (h *pingHandler) install(_ *EventManager, c *Connection) {
	c.AddInboundHandler("PING", h.pong)
}

func (h *pingHandler) pong(_ *EventManager, c *Connection, m *ircmsg.Message) {
	c.SendRawf("PONG :%v", m.Params[0])
}
