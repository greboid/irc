package main

import (
	"context"
	"crypto/tls"
	"fmt"
	"github.com/greboid/irc/rpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"log"
	"time"
)

func main() {
	conf, err := GetConfig()
	if err != nil {
		log.Fatalf("Unable to load config: %s", err.Error())
	}
	time.Sleep(5 * time.Second)
	creds := credentials.NewTLS(&tls.Config{InsecureSkipVerify: true})
	conn, err := grpc.Dial(fmt.Sprintf("%s:%d", conf.RPCHost, conf.RPCPort), grpc.WithTransportCredentials(creds))
	if err != nil {
		log.Fatalf("failed to fake plugin: %v", err)
	}
	defer conn.Close()
	client := rpc.NewIRCPluginClient(conn)
	context.Background()
	_, err = client.SendChannelMessage(rpc.CtxWithToken(context.Background(), "bearer", conf.RPCToken), &rpc.ChannelMessage{
		Channel: conf.Channel,
		Message: "RPC",
	})
	if err != nil {
		log.Printf("Error sending: %s", err.Error())
	} else {
		log.Print("Sent message, exiting.")
	}
}
