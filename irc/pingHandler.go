package irc

type pingHandler struct {
}

func NewPingHandler() *pingHandler {
	return &pingHandler{}
}

func (h *pingHandler) install(_ *EventManager, c *Connection) {
	c.AddInboundHandler("PING", h.pong)
}

func (h *pingHandler) pong(_ *EventManager, c *Connection, m *Message) {
	c.SendRawf("PONG :%v", m.Params[0])
}
