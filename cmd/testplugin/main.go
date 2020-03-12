package main

import (
	"context"
	"crypto/tls"
	"github.com/greboid/irc/config"
	"github.com/greboid/irc/protos"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"log"
)

func main() {
	conf := config.GetConfig()
	creds := credentials.NewTLS(&tls.Config{InsecureSkipVerify: true})
	conn, err := grpc.Dial("localhost:8081", grpc.WithTransportCredentials(creds))
	if err != nil {
		log.Fatalf("failed to fake plugin: %v", err)
	}
	defer conn.Close()
	client := protos.NewIRCPluginClient(conn)
	_, err = client.SendChannelMesssage(context.Background(), &protos.ChannelMessage{
		Channel: conf.Channel,
		Message: "RPC",
	})
	if err != nil {
		log.Printf("Error sending: %s", err.Error())
	}
}
