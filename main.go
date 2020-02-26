package main

import (
	"github.com/greboid/irc/irc"
	"log"
)

func main() {
	//Fill in the values for irc.ClientConfig{}
	connection := irc.IRCConnection{
		ClientConfig: irc.ClientConfig{},
		ConnConfig: irc.DefaultConnectionConfig,
	}
	connection.Init()
	//Add some callbacks
	err := connection.ConnectAndWait()
	if err != nil {
		log.Fatal(err)
	}
}
