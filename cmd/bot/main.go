package main

import (
	"flag"
	"github.com/greboid/irc/irc"
	"github.com/greboid/irc/rpc"
	"github.com/kouhin/envflag"
	"log"
)

//go:generate protoc -I ../../rpc plugin.proto --go_out=plugins=grpc:../../rpc

var (
	Server        = flag.String("server", "", "Which IRC server to connect to")
	Password      = flag.String("password", "", "The server password, if required")
	TLS           = flag.Bool("tls", true, "Connect with TLS?")
	Nickname      = flag.String("nick", "", "Nickname to use")
	Realname      = flag.String("realname", "", "'Real name' to use")
	Channel       = flag.String("channel", "", "Channel to join on connect")
	Debug         = flag.Bool("debug", false, "Enable IRC debug output")
	SASLAuth      = flag.Bool("sasl-auth", false, "Authenticate via SASL?")
	SASLUser      = flag.String("sasl-user", "", "SASL username")
	SASLPass      = flag.String("sasl-pass", "", "SASL password")
	RPCPort       = flag.Int("rpc-port", 8001, "gRPC server port")
	PluginsString = flag.String("plugins", "", "Comma separated list of plugins, name=token")
	FloodProfile  = flag.String("flood-profile", "restrictive", "Flood profile: restrictive, unlimited")
)

func main() {
	if err := envflag.Parse(); err != nil {
		log.Fatalf("Unable to load config: %s", err.Error())
	}
	Plugins, err := rpc.ParsePluginString(*PluginsString)
	if err != nil {
		log.Fatalf("Unable to load config: %s", err.Error())
	}
	connection := irc.NewIRC(*Server, *Password, *Nickname, *Realname, *TLS, *SASLAuth, *SASLUser, *SASLPass, *Debug,
		*FloodProfile, irc.NewEventManager())
	rpcServer := rpc.GrpcServer{Conn: connection, RPCPort: *RPCPort, Plugins: Plugins}
	log.Print("Adding callbacks")
	connection.AddInboundHandler("001", func(_ *irc.EventManager, c *irc.Connection, m *irc.Message) {
		c.SendRawf("JOIN :%s", *Channel)
	})
	connection.AddInboundHandler("PRIVMSG", func(em *irc.EventManager, _ *irc.Connection, m *irc.Message) {
		em.PublishChannelMessage(*m)
	})
	go rpcServer.StartGRPC()
	err = connection.ConnectAndWaitWithRetry(5)
	if err != nil {
		log.Fatal(err)
	}
	log.Print("Exiting")
}
