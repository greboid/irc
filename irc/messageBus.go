package irc

import (
	externalbus "github.com/vardius/message-bus"
)

func newExternalBus(size int) externalbus.MessageBus {
	mb := externalbus.New(size)
	return mb
}

func (c *Connection) SubscribeCapAdd(receiver func(*Connection, *CapabilityStruct)) {
	_ = c.bus.Subscribe("+cap", receiver)
}

func (c *Connection) UnsubscribeCapAdd(receiver func(*Connection, *CapabilityStruct)) {
	_ = c.bus.Unsubscribe("+cap", receiver)
}

func (c *Connection) PublishCapAdd(conn *Connection, capability *CapabilityStruct) {
	c.bus.Publish("+cap", conn, capability)
}

func (c *Connection) SubscribeCapDel(receiver func(*Connection, *CapabilityStruct)) {
	_ = c.bus.Subscribe("-cap", receiver)
}

func (c *Connection) UnsubscribeCapDel(receiver func(*Connection, *CapabilityStruct)) {
	_ = c.bus.Unsubscribe("-cap", receiver)
}

func (c *Connection) PublishCapDel(conn *Connection, capability *CapabilityStruct) {
	c.bus.Publish("-cap", conn, capability)
}

func (c *Connection) SubscribeChannelPart(receiver func(Channel)) {
	_ = c.bus.Subscribe("ChannelPart", receiver)
}

func (c *Connection) UnsubscribeChannelPart(receiver func(Channel)) {
	_ = c.bus.Unsubscribe("ChannelPart", receiver)
}

func (c *Connection) PublishChannelPart(channel Channel) {
	c.bus.Publish("ChannelPart", channel)
}

func (c *Connection) SubscribeChannelMessage(receiver func(Message)) {
	_ = c.bus.Subscribe("ChannelMessage", receiver)
}

func (c *Connection) UnsubscribeChannelMessage(receiver func(Message)) {
	_ = c.bus.Unsubscribe("ChannelMessage", receiver)
}

func (c *Connection) PublishChannelMessage(channel Message) {
	c.bus.Publish("ChannelMessage", channel)
}
