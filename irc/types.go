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
	Server           string
	Password         string
	Nickname         string
	UseTLS           bool
	Debug            bool
	SASLAuth         bool
	SASLUser         string
	SASLPass         string
	ConnConfig       ConnectionConfig
	ClientConfig     ClientConfig
	socket           net.Conn
	lastMessage      time.Time
	saslFinished     chan bool
	quitting         chan bool
	Finished         chan bool
	writeChan        chan string
	errorChannel     chan error
	rawHandlers      []func(*Connection, RawMessage)
	inboundHandlers  map[string][]func(*Connection, *Message)
	outboundHandlers []func(*Connection, string)
	signals          chan os.Signal
	initialised      bool
	registered       bool
	listeners        eventListeners
	saslStarted      bool
}

type RawMessage struct {
	message string
	out     bool
}

type Message struct {
	Raw    string
	Tags   string
	Source string
	Verb   string
	Params []string
}

type InboundHandler struct {
	Verb    string
	Handler func(*Connection, *Message)
}

type Channel struct {
	Name string
}
