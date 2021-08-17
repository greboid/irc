package irc

import (
	"github.com/ergochat/irc-go/ircmsg"
)

type errorHandler struct {
}

func NewErrorHandler() *errorHandler {
	return &errorHandler{}
}

func (h *errorHandler) install(c *Connection) {
	c.AddInboundHandler("ERROR", h.quitOnError)
}

func (h *errorHandler) quitOnError(_ *EventManager, c *Connection, _ *ircmsg.Message) {
	c.Finished <- true
}
