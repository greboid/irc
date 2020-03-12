package main

import (
	"github.com/greboid/irc/database"
	"github.com/greboid/irc/irc"
	"github.com/greboid/irc/rpc"
	"github.com/greboid/irc/web"
	"log"
)

//go:generate protoc -I ../../rpc plugin.proto --go_out=plugins=grpc:../../rpc

func main() {
	conf, err := GetConfig()
	if err != nil {
		log.Panicf("Unable to load config: %s", err.Error())
	}
	db, err := database.New(conf.DBPath)
	if err != nil {
		log.Panicf("Unable to load config: %s", err.Error())
	}
	connection := irc.NewIRC(conf.Server, conf.Password, conf.Nickname, conf.TLS, conf.Debug)
	rpcServer := rpc.GrpcServer{Conn: connection, DB: db}
	go web.NewWeb(conf.WebPort, conf.Channel, conf.AdminKey, connection, db).StartWeb()
	log.Print("Adding callbacks")
	connection.AddInboundHandler("001", func(c *irc.Connection, m *irc.Message) {
		c.SendRawf("JOIN :%s", conf.Channel)
	})
	go rpcServer.StartGRPC()
	err = connection.ConnectAndWait()
	if err != nil {
		log.Fatal(err)
	}
	log.Print("Exiting")
}
