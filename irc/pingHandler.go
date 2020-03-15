package irc

type pingHandler struct {
}

func NewPingHandler() *pingHandler {
	return &pingHandler{}
}

func (h *pingHandler) install(c *Connection) {
	c.AddInboundHandler("PING", h.pong)
}

func (h *pingHandler) pong(c *Connection, m *Message) {
	c.SendRawf("PONG :%v", m.Params[0])
}
