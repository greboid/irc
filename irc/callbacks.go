package irc

import (
	"strings"
)

var (
	defaultCallbacks = map[string]func(c *Connection, m *Message){
		"PING":  pong,
		"ERROR": quitOnError,
		"CAP":   handleCaps,
		"001":   handle001,
	}
)

func (irc *Connection) AddCallbacks(callbacks map[string]func(c *Connection, m *Message)) {
	for verb, callback := range callbacks {
		irc.AddCallback(verb, callback)
	}
}

func (irc *Connection) AddCallback(s string, f func(c *Connection, m *Message)) {
	if !irc.initialised {
		irc.Init()
	}
	s = strings.ToUpper(s)
	irc.callbacks[s] = append(irc.callbacks[s], f)
}

func (irc *Connection) runCallbacks(m *Message) {
	callbacks := irc.callbacks[m.Verb]
	for _, callback := range callbacks {
		go callback(irc, m)
	}
}

func pong(c *Connection, m *Message) {
	c.SendRawf("PONG :%v", m.Params)
}

func quitOnError(c *Connection, _ *Message) {
	c.Finished <- true
}

func handleCaps(c *Connection, m *Message) {
	capList := strings.Split(m.Params, " ")
	if capList[2] == "*" {

	}
	if capList[1] == "LS" {
		if capList[2] == "*" {
			capList = capList[3:]
		} else {
			c.capListingEnded = true
			capList = capList[2:]
		}
		c.capsAdvertised = append(c.capsAdvertised, capList...)
	} else if capList[1] == "ACK" {
		c.capsGained = append(c.capsGained, capList[2])
		if len(c.capsGained) == len(c.capsWanted) {
			c.capsAcknowledged = true
		}
	} else if capList[1] == "NAK" {
		c.capsWanted = remove(c.capsWanted, capList[2])
	}
	if c.capListingEnded && !c.capRequestingEnded {
		var wantedAndAvailble []string
		for _, available := range c.capsAdvertised {
			for _, wanted := range c.capsWanted {
				if wanted == available {
					wantedAndAvailble = append(wantedAndAvailble, wanted)
				}
			}
		}
		c.capsWanted = wantedAndAvailble
		for _, capability := range c.capsWanted {
			c.SendRawf("CAP REQ :%s", capability)
		}
		c.capRequestingEnded = true
	}
	if c.capListingEnded && c.capRequestingEnded && c.capsAcknowledged {
		c.SendRaw("CAP END")
	}
}

func handle001(c *Connection, _ *Message) {
	c.registered = true
}

func remove(items []string, item string) []string {
	var newitems []string
	for _, i := range items {
		if i != item {
			newitems = append(newitems, i)
		}
	}
	return newitems
}
