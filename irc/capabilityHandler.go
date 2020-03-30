package irc

import (
	"github.com/imdario/mergo"
	"log"
	"strings"
	"sync"
	"time"
)

type capabilityHandler struct {
	List      map[string]*CapabilityStruct
	wanted    map[string]bool
	listing   bool
	requested bool
	finished  bool
	mutex     *sync.Mutex
}

type CapabilityStruct struct {
	name         string
	values       string
	acked        bool
	waitingonAck bool
}

func NewCapabilityHandler() *capabilityHandler {
	return &capabilityHandler{
		List:      map[string]*CapabilityStruct{},
		wanted:    map[string]bool{"echo-message": true, "message-tags": true, "multi-prefix": true, "sasl": true},
		listing:   false,
		requested: false,
		finished:  false,
		mutex:     &sync.Mutex{},
	}
}

func (h *capabilityHandler) install(_ *EventManager, c *Connection) {
	c.AddInboundHandler("CAP", h.handleCaps)
	c.AddInboundHandler("001", h.handleRegistered)
	h.Negotiate(c)
}

func (h *capabilityHandler) handleRegistered(*EventManager, *Connection, *Message) {
	h.finished = true
	h.listing = false
	h.requested = false
}

func (h *capabilityHandler) Negotiate(irc Sender) {
	irc.SendRaw("CAP LS 302")
}

func (h *capabilityHandler) handleCaps(eventManager *EventManager, c *Connection, m *Message) {
	switch m.Params[1] {
	case "LS":
		h.handleLS(strings.Split(m.Params[2], " "))
		if !h.listing {
			h.capReq(c)
		}
		break
	case "ACK":
		h.handleACK(eventManager, c, strings.Split(m.Params[2], " "))
		break
	case "NAK":
		h.handleNAK(strings.Split(m.Params[2], " "))
		break
	case "NEW":
		h.handleLS(strings.Split(m.Params[2], " "))
		if !h.listing {
			h.capReq(c)
		}
		break
	case "DEL":
		h.handleDel(eventManager, c, strings.Split(m.Params[2], " "))
		break
	}
}

func (h *capabilityHandler) handleLS(tokenised []string) {
	if len(tokenised) == 0 {
		return
	}
	if tokenised[0] == "*" {
		tokenised = tokenised[1:]
	} else {
		h.listing = false
	}
	_ = mergo.Merge(&h.List, h.parseCapabilities(tokenised), mergo.WithOverride)
}

func (_ *capabilityHandler) parseCapabilities(tokenised []string) map[string]*CapabilityStruct {
	capabilities := map[string]*CapabilityStruct{}
	for _, token := range tokenised {
		capability := &CapabilityStruct{}
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

func (h *capabilityHandler) capReq(c Sender) {
	var reqs []string
	for name, capability := range h.List {
		_, ok := h.wanted[name]
		if ok {
			reqs = append(reqs, name)
			capability.waitingonAck = true
		}
	}
	if len(reqs) > 0 {
		c.SendRawf("CAP REQ :%s", strings.Join(reqs, " "))
	}
}

func (h *capabilityHandler) handleACK(m *EventManager, c *Connection, tokenised []string) {
	h.mutex.Lock()
	for _, token := range tokenised {
		capability, ok := h.List[token]
		if ok {
			capability.acked = true
			capability.waitingonAck = false
			m.PublishCapAdd(c, capability)
		}
	}
	h.mutex.Unlock()
	if countAcked(h.List) == len(h.wanted) {
		h.finished = true
		if _, ok := h.List["sasl"]; ok && !c.saslFinished {
			log.Print("Waiting for SASL")
			h.waitonSasl(c)
		}
		c.SendRaw("CAP END")
	}
}

func (h *capabilityHandler) waitonSasl(c *Connection) {
	select {
	case <-c.saslFinishedChan:
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

func (h *capabilityHandler) handleDel(m *EventManager, c *Connection, tokenised []string) {
	toRemove := h.parseCapabilities(tokenised)
	for _, capability := range toRemove {
		capability.acked = false
		capability.waitingonAck = false
		m.PublishCapDel(c, capability)
	}
}

func countAcked(list map[string]*CapabilityStruct) int {
	acked := 0
	for _, capability := range list {
		if !capability.waitingonAck && capability.acked {
			acked++
		}
	}
	return acked
}
