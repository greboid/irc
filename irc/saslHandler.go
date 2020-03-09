package irc

import (
	"encoding/base64"
	"fmt"
	"github.com/greboid/irc/config"
	"log"
	"strings"
)

type saslHandler struct {
	conf        *config.Config
	authed      bool
	readyToAuth bool
	authing     bool
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
	if cap != "sasl" || c.saslStarted {
		return
	}
	if !h.authed {
		c.saslStarted = true
	}
	if !checkSASLSupported(c) {
		log.Printf("SASL Finished")
		c.saslFinished <- true
		return
	}
	h.readyToAuth = true
	c.SendRaw("AUTHENTICATE PLAIN")
}

func checkSASLSupported(c *Connection) bool {
	if c.conf.SASLAuth && c.conf.SASLUser != "" && c.conf.SASLPass != "" {
		log.Print("SASL configured")
		var saslmethods []string
		for key, _ := range c.capabilityHandler.available {
			if key.name == "sasl" {
				saslmethods = strings.Split(key.values, ",")
			}
		}
		_, ok := contains(saslmethods, "PLAIN")
		if !ok {
			log.Printf("No supported SASL methods")
		}
		return ok
	}
	log.Print("SASL not configured")
	return false
}

func (h *saslHandler) handleCapDel(_ string) {}

func (h *saslHandler) handleLoggedinAs(*Connection, *Message) {}

func (h *saslHandler) handleLoggedOut(*Connection, *Message) {}

func (h *saslHandler) handleNickLocked(*Connection, *Message) {}

func (h *saslHandler) handleAuthSuccess(c *Connection, _ *Message) {
	log.Print("SASL Auth success")
	c.saslFinished <- true
}

func (h *saslHandler) handleAuthFail(c *Connection, _ *Message) {
	log.Print("SASL Auth failed")
	c.saslFinished <- true
}

func (h *saslHandler) handleMessageTooLong(*Connection, *Message) {
}

func (h *saslHandler) handleAborted(*Connection, *Message) {
}

func (h *saslHandler) handleAlreadyAuthed(*Connection, *Message) {
}

func (h *saslHandler) handleMechanisms(*Connection, *Message) {
}

func (h *saslHandler) handleAuthenticate(c *Connection, m *Message) {
	if h.readyToAuth {
		if m.Params == "+" {
			h.authing = true
			h.doAuth(c)
		}
	}
}

func (h *saslHandler) doAuth(c *Connection) {
	encoded := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s\x00%s\x00%s", c.conf.SASLUser, c.conf.SASLUser, c.conf.SASLPass)))
	for i := 0; i < len(encoded); i += 400 {
		c.SendRawf("AUTHENTICATE %s", encoded[i:])
		if len(encoded[i:]) == 400 {
			c.SendRaw("AUTHENTICATE +")
		}
	}
}

func contains(s []string, e string) (*string, bool) {
	for _, a := range s {
		if a == e {
			return &e, true
		}
	}
	return nil, false
}
