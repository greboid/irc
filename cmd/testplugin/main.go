package main

import (
	"context"
	"crypto/tls"
	"fmt"
	"github.com/greboid/irc/rpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"log"
)

func main() {
	conf, err := GetConfig()
	if err != nil {
		log.Fatalf("Unable to load config: %s", err.Error())
	}
	creds := credentials.NewTLS(&tls.Config{InsecureSkipVerify: true})
	conn, err := grpc.Dial(fmt.Sprintf("%s:%d", conf.RPCHost, conf.RPCPort), grpc.WithTransportCredentials(creds))
	if err != nil {
		log.Fatalf("failed to fake plugin: %v", err)
	}
	defer conn.Close()
	client := rpc.NewIRCPluginClient(conn)
	context.Background()
	_, err = client.SendChannelMesssage(rpc.CtxWithToken(context.Background(), "bearer", "bad_token"), &rpc.ChannelMessage{
		Channel: conf.Channel,
		Message: "RPC",
	})
	if err != nil {
		log.Printf("Error sending: %s", err.Error())
	}
}
