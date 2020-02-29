package irc

import (
	"net"
	"os"
	"time"
)

var (
	DefaultConnectionConfig = ConnectionConfig{
		KeepAlive: 4 * time.Minute,
		Timeout:   1 * time.Minute,
	}
)

type ClientConfig struct {
	Server   string
	Nick     string
	User     string
	Password string
	UseTLS   bool
	Realname string
}

type ConnectionConfig struct {
	KeepAlive time.Duration
	Timeout   time.Duration
}

type Connection struct {
	ConnConfig        ConnectionConfig
	ClientConfig      ClientConfig
	socket            net.Conn
	lastMessage       time.Time
	quitting          chan bool
	Finished          chan bool
	writeChan         chan string
	errorChannel      chan error
	inboundHandlers   map[string][]func(*Connection, *Message)
	outboundHandlers  []func(*Connection, *string)
	signals           chan os.Signal
	initialised       bool
	registered        bool
	capabilityHandler capabilityHandler
	nickHandler       nickHandler
	debugHandler      debugHandler
}

type Message struct {
	Raw         string
	Tags        string
	Source      string
	Verb        string
	ParamsArray []string
	Params      string
}

type InboundHandler struct {
	Verb    string
	Handler func(*Connection, *Message)
}
