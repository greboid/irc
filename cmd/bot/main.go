package main

import (
	"github.com/greboid/irc/config"
	"github.com/greboid/irc/database"
	"github.com/greboid/irc/irc"
	"github.com/greboid/irc/web"
	"log"
)

func main() {
	conf := config.GetConfig()
	db := database.New(conf.DBPath)
	connection := irc.NewIRC(conf)
	go web.NewWeb(conf, connection, db).StartWeb()
	log.Print("Adding callbacks")
	connection.AddInboundHandler("001", func(c*irc.Connection, m *irc.Message) {
		c.SendRawf("JOIN :%s", conf.Channel)
	})
	err := connection.ConnectAndWait()
	if err != nil {
		log.Fatal(err)
	}
	log.Print("Exiting")
}
