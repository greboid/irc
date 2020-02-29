package irc

import (
	"strings"
	"sync"
)

type capabilityHandler struct {
	available map[capabilityStruct]bool
	wanted    map[string]bool
	acked     map[string]bool
	listing   bool
	requested bool
	finished  bool
	mutex     *sync.Mutex
}

type capabilityStruct struct {
	name   string
	values string
}

func (h *capabilityHandler) install(c *Connection) {
	h.available = map[capabilityStruct]bool{}
	h.wanted = map[string]bool{"echo-message": true, "message-tags": true, "multi-prefix": true}
	h.acked = map[string]bool{}
	h.listing = false
	h.requested = false
	h.finished = false
	h.mutex = &sync.Mutex{}

	c.AddCallback("CAP", h.handleCaps)
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
	default:
		//TODO: Support NEW and DEL
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
	reqs := []string{}
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
	}
	h.mutex.Unlock()
	if len(h.acked) == len(h.wanted) {
		c.SendRaw("CAP END")
	}
}

func (h *capabilityHandler) handleNAK(tokenised []string) {
	delete(h.wanted, tokenised[0])
}
