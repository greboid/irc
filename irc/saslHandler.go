package irc

import (
	"encoding/base64"
	"fmt"
	"go.uber.org/zap"
	"strings"
)

type SaslHandler struct {
	SASLAuth    bool
	SASLUser    string
	SASLPass    string
	authed      bool
	readyToAuth bool
	authing     bool
	saslMethods []string
	logger      *zap.SugaredLogger
}

func NewSASLHandler(useSasl bool, saslUser string, saslPass string, logger *zap.SugaredLogger) *SaslHandler {
	return &SaslHandler{
		SASLAuth: useSasl,
		SASLUser: saslUser,
		SASLPass: saslPass,
		logger:   logger,
	}
}

func (h *SaslHandler) Install(em *EventManager, c *Connection) {
	em.SubscribeCapAdd(h.handleCapAdd)
	em.SubscribeCapDel(h.handleCapDel)
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

func (h *SaslHandler) handleCapAdd(c *Connection, cap *CapabilityStruct) {
	if cap.name != "sasl" || c.saslStarted {
		return
	}
	if !h.authed {
		c.saslStarted = true
	}
	if !h.checkSASLSupported(cap) {
		h.logger.Debugf("SASL Finished")
		c.saslFinishedChan <- true
		c.saslFinished = true
		return
	}
	h.readyToAuth = true
	c.SendRaw("AUTHENTICATE PLAIN")
}

func (h *SaslHandler) checkSASLSupported(cap *CapabilityStruct) bool {
	if h.SASLAuth && h.SASLUser != "" && h.SASLPass != "" {
		h.logger.Debugf("SASL configured")
		if !contains(strings.Split(cap.values, ","), "PLAIN") {
			h.logger.Debugf("No supported SASL methods")
			return false
		}
		return true
	}
	h.logger.Debugf("SASL not configured")
	return false
}

func (h *SaslHandler) handleCapDel(*Connection, *CapabilityStruct) {}

func (h *SaslHandler) handleLoggedinAs(*EventManager, *Connection, *Message) {}

func (h *SaslHandler) handleLoggedOut(*EventManager, *Connection, *Message) {}

func (h *SaslHandler) handleNickLocked(*EventManager, *Connection, *Message) {}

func (h *SaslHandler) handleAuthSuccess(_ *EventManager, c *Connection, _ *Message) {
	h.logger.Debugf("SASL Auth success")
	c.saslFinishedChan <- true
}

func (h *SaslHandler) handleAuthFail(_ *EventManager, c *Connection, _ *Message) {
	h.logger.Debugf("SASL Auth failed")
	c.saslFinishedChan <- true
}

func (h *SaslHandler) handleMessageTooLong(*EventManager, *Connection, *Message) {
}

func (h *SaslHandler) handleAborted(*EventManager, *Connection, *Message) {
}

func (h *SaslHandler) handleAlreadyAuthed(*EventManager, *Connection, *Message) {
}

func (h *SaslHandler) handleMechanisms(*EventManager, *Connection, *Message) {
}

func (h *SaslHandler) handleAuthenticate(em *EventManager, c *Connection, m *Message) {
	if h.readyToAuth {
		if len(m.Params) == 1 && m.Params[0] == "+" {
			h.authing = true
			h.doAuth(em, c)
		}
	}
}

func (h *SaslHandler) doAuth(_ *EventManager, c *Connection) {
	encoded := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s\x00%s\x00%s", h.SASLUser, h.SASLUser, h.SASLPass)))
	for i := 0; i < len(encoded); i += 400 {
		c.SendRawf("AUTHENTICATE %s", encoded[i:])
		if len(encoded[i:]) == 400 {
			c.SendRaw("AUTHENTICATE +")
		}
	}
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
