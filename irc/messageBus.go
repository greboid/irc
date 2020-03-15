package irc

import (
	externalbus "github.com/vardius/message-bus"
)

func newExternalBus(size int) externalbus.MessageBus {
	mb := externalbus.New(size)
	return mb
}

func (b *Connection) SubscribeCapAdd(receiver func(*Connection, *CapabilityStruct)) {
	_ = b.bus.Subscribe("+cap", receiver)
}

func (b *Connection) UnsubscribeCapAdd(receiver func(*Connection, *CapabilityStruct)) {
	_ = b.bus.Unsubscribe("+cap", receiver)
}

func (b *Connection) PublishCapAdd(conn *Connection, capability *CapabilityStruct) {
	b.bus.Publish("+cap", conn, capability)
}

func (b *Connection) SubscribeCapDel(receiver func(*Connection, *CapabilityStruct)) {
	_ = b.bus.Subscribe("-cap", receiver)
}

func (b *Connection) UnsubscribeCapDel(receiver func(*Connection, *CapabilityStruct)) {
	_ = b.bus.Unsubscribe("-cap", receiver)
}

func (b *Connection) PublishCapDel(conn *Connection, capability *CapabilityStruct) {
	b.bus.Publish("-cap", conn, capability)
}

func (b *Connection) SubscribeChannelPart(receiver func(Channel)) {
	_ = b.bus.Subscribe("ChannelPart", receiver)
}

func (b *Connection) UnsubscribeChannelPart(receiver func(Channel)) {
	_ = b.bus.Unsubscribe("ChannelPart", receiver)
}

func (b *Connection) PublishChannelPart(channel Channel) {
	b.bus.Publish("ChannelPart", channel)
}

func (b *Connection) SubscribeChannelMessage(receiver func(Message)) {
	_ = b.bus.Subscribe("ChannelMessage", receiver)
}

func (b *Connection) UnsubscribeChannelMessage(receiver func(Message)) {
	_ = b.bus.Unsubscribe("ChannelMessage", receiver)
}

func (b *Connection) PublishChannelMessage(channel Message) {
	b.bus.Publish("ChannelMessage", channel)
}
