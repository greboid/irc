package irc

import (
	"log"
	"strings"
	"sync"
	"time"
)

type capabilityHandler struct {
	available   map[capabilityStruct]bool
	wanted      map[string]bool
	acked       map[string]bool
	listing     bool
	requested   bool
	finished    bool
	mutex       *sync.Mutex
	saslHandler *saslHandler
}

type capabilityStruct struct {
	name   string
	values string
}

func (h *capabilityHandler) install(c *Connection) {
	h.available = map[capabilityStruct]bool{}
	h.wanted = map[string]bool{"echo-message": true, "message-tags": true, "multi-prefix": true, "sasl": true}
	h.acked = map[string]bool{}
	h.listing = false
	h.requested = false
	h.finished = false
	h.mutex = &sync.Mutex{}
	h.saslHandler = &saslHandler{
		SASLAuth: c.SASLAuth,
		SASLUser: c.SASLUser,
		SASLPass: c.SASLPass,
	}

	c.AddInboundHandler("CAP", h.handleCaps)
	c.AddInboundHandler("001", h.handleRegistered)
	c.AddOutboundHandler(passHandler(h))
	h.saslHandler.install(c)
}

func passHandler(h *capabilityHandler) func(c *Connection, m string) {
	return func(c *Connection, m string) {
		if strings.HasPrefix(m, "PASS") {
			h.Negotiate(c)
		}
	}
}

func (h *capabilityHandler) handleRegistered(*Connection, *Message) {
	h.finished = true
	h.listing = false
	h.requested = false
}

func (h *capabilityHandler) Negotiate(irc *Connection) {
	irc.SendRaw("CAP LS 302")
}

func (h *capabilityHandler) handleCaps(c *Connection, m *Message) {
	tokenised := strings.Split(m.Params, " ")[1:]
	switch tokenised[0] {
	case "LS":
		h.handleLS(c, tokenised[1:])
		break
	case "ACK":
		h.handleACK(c, tokenised[1:])
		break
	case "NAK":
		h.handleNAK(tokenised[1:])
		break
	case "NEW":
		h.handleLS(c, tokenised[1:])
		break
	case "DEL":
		h.handleDel(c, tokenised[1:])
		break
	}
}

func (h *capabilityHandler) handleLS(c *Connection, tokenised []string) {
	if tokenised[0] == "*" {
		tokenised = tokenised[1:]
	} else {
		h.listing = false
	}
	h.available = h.parseCapabilities(tokenised)
	if !h.listing {
		h.capReq(c)
	}
}

func (_ *capabilityHandler) parseCapabilities(tokenised []string) map[capabilityStruct]bool {
	capabilities := map[capabilityStruct]bool{}
	for _, token := range tokenised {
		capability := capabilityStruct{}
		if strings.Contains(token, "=") {
			values := strings.SplitN(token, "=", 2)
			capability.name = values[0]
			capability.values = values[1]
		} else {
			capability.name = token
			capability.values = ""
		}
		capabilities[capability] = true
	}
	return capabilities
}

func (h *capabilityHandler) capReq(c *Connection) {
	var reqs []string
	for capability := range h.available {
		_, ok := h.wanted[capability.name]
		if ok {
			reqs = append(reqs, capability.name)
		}
	}
	c.SendRawf("CAP REQ :%s", strings.Join(reqs, " "))
}

func (h *capabilityHandler) handleACK(c *Connection, tokenised []string) {
	h.mutex.Lock()
	for _, token := range tokenised {
		h.acked[token] = true
		c.Bus.Publish("+cap", c, token)
	}
	h.mutex.Unlock()
	if len(h.acked) == len(h.wanted) {
		h.finished = true
		if _, ok := h.acked["sasl"]; ok {
			log.Print("Waiting for SASL")
			h.waitonSasl(c)
		}
		c.SendRaw("CAP END")
	}
}

func (h *capabilityHandler) waitonSasl(c *Connection) {
	select {
	case <-c.saslFinished:
		return
	case <-time.After(5 * time.Second):
		return
	}
}

func (h *capabilityHandler) handleNAK(tokenised []string) {
	delete(h.wanted, tokenised[0])
}

func (h *capabilityHandler) handleDel(c *Connection, tokenised []string) {
	toRemove := h.parseCapabilities(tokenised)
	for remove := range toRemove {
		delete(h.wanted, remove.name)
		delete(h.available, remove)
		c.Bus.Publish("-cap", remove.name)
	}
}
