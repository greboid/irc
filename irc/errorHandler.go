package irc

type errorHandler struct {
}

func NewErrorHandler() *errorHandler {
	return &errorHandler{
	}
}

func (h *errorHandler) install(c *Connection) {
	c.AddInboundHandler("ERROR", h.quitOnError)
}

func (h *errorHandler) quitOnError(c *Connection, _ *Message) {
	c.Finished <- true
}