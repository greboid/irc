package main

import (
	"flag"
	"github.com/greboid/irc/database"
	"github.com/greboid/irc/irc"
	"github.com/greboid/irc/rpc"
	"github.com/greboid/irc/web"
	"github.com/kouhin/envflag"
	"log"
)

//go:generate protoc -I ../../rpc plugin.proto --go_out=plugins=grpc:../../rpc

func main() {
	var (
		Server        = flag.String("server", "", "")
		Password      = flag.String("password", "", "")
		TLS           = flag.Bool("tls", true, "")
		Nickname      = flag.String("nick", "", "")
		WebPort       = flag.Int("web-port", 8000, "")
		Channel       = flag.String("channel", "", "")
		DBPath        = flag.String("db-path", "/data/db", "")
		AdminKey      = flag.String("admin-key", "", "")
		Debug         = flag.Bool("debug", false, "")
		SASLAuth      = flag.Bool("sasl-auth", false, "")
		SASLUser      = flag.String("sasl-user", "", "")
		SASLPass      = flag.String("sasl-pass", "", "")
		RPCPort       = flag.Int("rpc-port", 8001, "")
		PluginsString = flag.String("plugins", "", "")
	)
	if err := envflag.Parse(); err != nil {
		log.Fatalf("Unable to load config: %s", err.Error())
	}
	Plugins, err := rpc.ParsePluginString(*PluginsString)
	if err != nil {
		log.Fatalf("Unable to load config: %s", err.Error())
	}
	db, err := database.New(*DBPath)
	if err != nil {
		log.Panicf("Unable to load config: %s", err.Error())
	}
	connection := irc.NewIRC(*Server, *Password, *Nickname, *TLS, *SASLAuth, *SASLUser, *SASLPass, *Debug)
	rpcServer := rpc.GrpcServer{Conn: connection, DB: db, RPCPort: *RPCPort, Plugins: Plugins}
	go web.NewWeb(*WebPort, *Channel, *AdminKey, connection, db).StartWeb()
	log.Print("Adding callbacks")
	connection.AddInboundHandler("001", func(c *irc.Connection, m *irc.Message) {
		c.SendRawf("JOIN :%s", *Channel)
	})
	connection.AddInboundHandler("PRIVMSG", func(c *irc.Connection, m *irc.Message) {
		c.PublishChannelMessage(*m)
	})
	go rpcServer.StartGRPC()
	err = connection.ConnectAndWait()
	if err != nil {
		log.Fatal(err)
	}
	log.Print("Exiting")
}
