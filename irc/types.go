package irc

import (
	"github.com/greboid/irc/config"
	"github.com/vardius/message-bus"
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
	conf              *config.Config
	ConnConfig        ConnectionConfig
	ClientConfig      ClientConfig
	socket            net.Conn
	lastMessage       time.Time
	saslFinished      chan bool
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
	Bus               messagebus.MessageBus
	saslStarted       bool
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
