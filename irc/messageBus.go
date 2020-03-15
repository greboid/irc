package irc

import (
	"reflect"
)

type eventListeners struct {
	capadds        map[reflect.Value]func(*Connection, *CapabilityStruct)
	capdels        map[reflect.Value]func(*Connection, *CapabilityStruct)
	channelParts   map[reflect.Value]func(Channel)
	channelMessage map[reflect.Value]func(Message)
}

func newEventListeners() eventListeners {
	return eventListeners{
		capadds:        make(map[reflect.Value]func(*Connection, *CapabilityStruct)),
		capdels:        make(map[reflect.Value]func(*Connection, *CapabilityStruct)),
		channelParts:   make(map[reflect.Value]func(Channel)),
		channelMessage: make(map[reflect.Value]func(Message)),
	}
}

func (irc *Connection) SubscribeCapAdd(receiver func(*Connection, *CapabilityStruct)) {
	irc.listeners.capadds[reflect.ValueOf(receiver)] = receiver
}

func (irc *Connection) UnsubscribeCapAdd(receiver func(*Connection, *CapabilityStruct)) {
	delete(irc.listeners.capadds, reflect.ValueOf(receiver))
}

func (irc *Connection) PublishCapAdd(conn *Connection, capability *CapabilityStruct) {
	for _, value := range irc.listeners.capadds {
		value(conn, capability)
	}
}

func (irc *Connection) SubscribeCapDel(receiver func(*Connection, *CapabilityStruct)) {
	irc.listeners.capdels[reflect.ValueOf(receiver)] = receiver
}

func (irc *Connection) UnsubscribeCapDel(receiver func(*Connection, *CapabilityStruct)) {
	delete(irc.listeners.capdels, reflect.ValueOf(receiver))
}

func (irc *Connection) PublishCapDel(conn *Connection, capability *CapabilityStruct) {
	for _, value := range irc.listeners.capdels {
		value(conn, capability)
	}
}

func (irc *Connection) SubscribeChannelPart(receiver func(Channel)) {
	irc.listeners.channelParts[reflect.ValueOf(receiver)] = receiver
}

func (irc *Connection) UnsubscribeChannelPart(receiver func(Channel)) {
	delete(irc.listeners.channelParts, reflect.ValueOf(receiver))
}

func (irc *Connection) PublishChannelPart(channel Channel) {
	for _, value := range irc.listeners.channelParts {
		value(channel)
	}
}

func (irc *Connection) SubscribeChannelMessage(receiver func(Message)) {
	irc.listeners.channelMessage[reflect.ValueOf(receiver)] = receiver
}

func (irc *Connection) UnsubscribeChannelMessage(receiver func(Message)) {
	delete(irc.listeners.channelMessage, reflect.ValueOf(receiver))
}

func (irc *Connection) PublishChannelMessage(message Message) {
	for _, value := range irc.listeners.channelMessage {
		value(message)
	}
}
