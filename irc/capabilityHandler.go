package irc

import (
	"log"
	"strings"
	"sync"
	"time"
)

type capabilityHandler struct {
	List        map[string]*capabilityStruct
	wanted      map[string]bool
	listing     bool
	requested   bool
	finished    bool
	mutex       *sync.Mutex
}

type capabilityStruct struct {
	name         string
	values       string
	acked        bool
	waitingonAck bool
}

func (h *capabilityHandler) install(c *Connection) {
	h.List = map[string]*capabilityStruct{}
	h.wanted = map[string]bool{"echo-message": true, "message-tags": true, "multi-prefix": true, "sasl": true}
	h.listing = false
	h.requested = false
	h.finished = false
	h.mutex = &sync.Mutex{}

	c.AddInboundHandler("CAP", h.handleCaps)
	c.AddInboundHandler("001", h.handleRegistered)
	c.AddOutboundHandler(h.passHandler)
}

func (h *capabilityHandler) passHandler(c *Connection, m string) {
	if strings.HasPrefix(m, "PASS") {
		h.Negotiate(c)
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
	h.List = h.parseCapabilities(tokenised)
	if !h.listing {
		h.capReq(c)
	}
}

func (_ *capabilityHandler) parseCapabilities(tokenised []string) map[string]*capabilityStruct {
	capabilities := map[string]*capabilityStruct{}
	for _, token := range tokenised {
		capability := &capabilityStruct{}
		if strings.Contains(token, "=") {
			values := strings.SplitN(token, "=", 2)
			capability.name = values[0]
			capability.values = values[1]
		} else {
			capability.name = token
		}
		capabilities[capability.name] = capability
	}
	return capabilities
}

func (h *capabilityHandler) capReq(c *Connection) {
	var reqs []string
	for name, capability := range h.List {
		_, ok := h.wanted[name]
		if ok {
			reqs = append(reqs, name)
			capability.waitingonAck = true
		}
	}
	c.SendRawf("CAP REQ :%s", strings.Join(reqs, " "))
}

func (h *capabilityHandler) handleACK(c *Connection, tokenised []string) {
	h.mutex.Lock()
	for _, token := range tokenised {
		capability, ok := h.List[token]
		if ok {
			capability.acked = true
			capability.waitingonAck = false
			c.Bus.Publish("+cap", c, capability)
		}
	}
	h.mutex.Unlock()
	if countAcked(h.List) == len(h.wanted) {
		h.finished = true
		if _, ok := h.List["sasl"]; ok {
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
	for _, token := range tokenised {
		capability, ok := h.List[token]
		if ok {
			capability.waitingonAck = false
		}
	}
}

func (h *capabilityHandler) handleDel(c *Connection, tokenised []string) {
	toRemove := h.parseCapabilities(tokenised)
	for remove, capability := range toRemove {
		capability.acked = false
		capability.waitingonAck = false
		c.Bus.Publish("-cap", remove)
	}
}

func countAcked(list map[string]*capabilityStruct) int {
	acked := 0
	for _, capability := range list {
		if !capability.waitingonAck && capability.acked {
			acked++
		}
	}
	return acked
}
