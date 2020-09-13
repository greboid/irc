package irc

import (
	"strings"
)

type supportParser struct {
	values []supportedValue
	conn   *Connection
}

type supportedValue struct {
	name  string
	value string
}

func NewSupportHandler() *supportParser {
	return &supportParser{
		values: []supportedValue{},
	}
}

func (h *supportParser) install(c *Connection) {
	h.conn = c
	c.AddInboundHandler("005", h.handleSupport)
}

func (h *supportParser) handleSupport(em *EventManager, c *Connection, m *Message) {
	var tokenised []supportedValue
	if m == nil || m.Params == nil || len(m.Params) == 0 {
		tokenised = []supportedValue{}
	} else {
		tokenised = h.tokenise(m.Params[1 : len(m.Params)-1])
	}
	for index := range tokenised {
		value := tokenised[index]
		if strings.HasPrefix(value.name, "-") {
			h.values = h.remove(h.values, index)
		} else {
			h.values = append(h.values, value)
		}
	}
}

func (h *supportParser) tokenise(input []string) []supportedValue {
	values := make([]supportedValue, 0)
	for index := range input {
		value := strings.Split(input[index], "=")
		if len(value) == 2 {
			values = append(values, supportedValue{
				name:  value[0],
				value: value[1],
			})
		}
	}
	return values
}

func (h *supportParser) remove(s []supportedValue, i int) []supportedValue {
	if i >= len(s) || i < 0 {
		return s
	}
	s[len(s)-1], s[i] = s[i], s[len(s)-1]
	return s[:len(s)-1]
}
