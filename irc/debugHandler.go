package irc

type debugHandler struct {
	debug  bool
	logger Logger
}

func NewDebugHandler(logger Logger) *debugHandler {
	return &debugHandler{
		logger: logger,
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
	h.logger.Debugf(" IN | %s", m.message)
}

func (h *debugHandler) handleOutboundMessage(_ *Connection, m RawMessage) {
	h.logger.Debugf("OUT | %s", m.message)
}
