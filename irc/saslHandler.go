package irc

type saslHandler struct {
	authed bool
}

func (h *saslHandler) install(c *Connection) {
	_ = c.Bus.Subscribe("+cap", h.handleCapAdd)
	_ = c.Bus.Subscribe("-cap", h.handleCapDel)
	c.AddInboundHandler("900", h.handleLoggedinAs)
	c.AddInboundHandler("901", h.handleLoggedOut)
	c.AddInboundHandler("902", h.handleNickLocked)
	c.AddInboundHandler("903", h.handleAuthSuccess)
	c.AddInboundHandler("904", h.handleAuthFail)
	c.AddInboundHandler("905", h.handleMessageTooLong)
	c.AddInboundHandler("906", h.handleAborted)
	c.AddInboundHandler("907", h.handleAlreadyAuthed)
	c.AddInboundHandler("908", h.handleMechanisms)
	c.AddInboundHandler("AUTHENTICATE", h.handleAuthenticate)
}

func (h *saslHandler) handleCapAdd(c *Connection, cap string) {
	if !h.authed {
		c.saslStarted = true
	}
}

func (h *saslHandler) handleCapDel(cap string) {

}

func (h *saslHandler) handleLoggedinAs(c *Connection, m *Message) {

}

func (h *saslHandler) handleLoggedOut(c *Connection, m *Message) {

}

func (h *saslHandler) handleNickLocked(c *Connection, m *Message) {

}

func (h *saslHandler) handleAuthSuccess(c *Connection, m *Message) {

}

func (h *saslHandler) handleAuthFail(c *Connection, m *Message) {

}

func (h *saslHandler) handleMessageTooLong(c *Connection, m *Message) {

}

func (h *saslHandler) handleAborted(c *Connection, m *Message) {

}

func (h *saslHandler) handleAlreadyAuthed(c *Connection, m *Message) {

}

func (h *saslHandler) handleMechanisms(c *Connection, m *Message) {

}

func (h *saslHandler) handleAuthenticate(c *Connection, m *Message) {

}
