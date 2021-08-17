package irc

import (
	"reflect"

	"github.com/ergochat/irc-go/ircmsg"
)

type EventManager struct {
	capadds        map[reflect.Value]func(*Connection, *CapabilityStruct)
	capdels        map[reflect.Value]func(*Connection, *CapabilityStruct)
	channelParts   map[reflect.Value]func(Channel)
	channelMessage map[reflect.Value]func(*ircmsg.Message)
}

func NewEventManager() *EventManager {
	return &EventManager{
		capadds:        make(map[reflect.Value]func(*Connection, *CapabilityStruct)),
		capdels:        make(map[reflect.Value]func(*Connection, *CapabilityStruct)),
		channelParts:   make(map[reflect.Value]func(Channel)),
		channelMessage: make(map[reflect.Value]func(*ircmsg.Message)),
	}
}

func (irc *EventManager) SubscribeCapAdd(receiver func(*Connection, *CapabilityStruct)) {
	irc.capadds[reflect.ValueOf(receiver)] = receiver
}

func (irc *EventManager) UnsubscribeCapAdd(receiver func(*Connection, *CapabilityStruct)) {
	delete(irc.capadds, reflect.ValueOf(receiver))
}

func (irc *EventManager) PublishCapAdd(conn *Connection, capability *CapabilityStruct) {
	for _, value := range irc.capadds {
		value(conn, capability)
	}
}

func (irc *EventManager) SubscribeCapDel(receiver func(*Connection, *CapabilityStruct)) {
	irc.capdels[reflect.ValueOf(receiver)] = receiver
}

func (irc *EventManager) UnsubscribeCapDel(receiver func(*Connection, *CapabilityStruct)) {
	delete(irc.capdels, reflect.ValueOf(receiver))
}

func (irc *EventManager) PublishCapDel(conn *Connection, capability *CapabilityStruct) {
	for _, value := range irc.capdels {
		value(conn, capability)
	}
}

func (irc EventManager) SubscribeChannelPart(receiver func(Channel)) {
	irc.channelParts[reflect.ValueOf(receiver)] = receiver
}

func (irc *EventManager) UnsubscribeChannelPart(receiver func(Channel)) {
	delete(irc.channelParts, reflect.ValueOf(receiver))
}

func (irc *EventManager) PublishChannelPart(channel Channel) {
	for _, value := range irc.channelParts {
		value(channel)
	}
}

func (irc *EventManager) SubscribeChannelMessage(receiver func(*ircmsg.Message)) {
	irc.channelMessage[reflect.ValueOf(receiver)] = receiver
}

func (irc *EventManager) UnsubscribeChannelMessage(receiver func(*ircmsg.Message)) {
	delete(irc.channelMessage, reflect.ValueOf(receiver))
}

func (irc *EventManager) PublishChannelMessage(message *ircmsg.Message) {
	for _, value := range irc.channelMessage {
		value(message)
	}
}
