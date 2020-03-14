package main

import (
	"context"
	"crypto/tls"
	"fmt"
	"github.com/greboid/irc/rpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"io"
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
	_, err = client.Ping(rpc.CtxWithToken(context.Background(), "bearer", conf.RPCToken), &rpc.Empty{})
	if err != nil {
		log.Fatalf("Error getting messages: %s", err.Error())
	}
	_, err = client.SendChannelMessage(rpc.CtxWithToken(context.Background(), "bearer", conf.RPCToken), &rpc.ChannelMessage{
		Channel: conf.Channel,
		Message: "RPC",
	})
	if err != nil {
		log.Printf("Error sending: %s", err.Error())
	} else {
		log.Print("Sent message, exiting.")
	}
	handler, err := client.GetMessages(rpc.CtxWithToken(context.Background(), "bearer", conf.RPCToken), &rpc.Channel{Name: conf.Channel})
	if err != nil {
		log.Fatalf("Error getting messages: %s", err.Error())
	}
	for {
		msg, err := handler.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("Err receiving message: %v", err)
		}
		log.Printf("%s: %s - %s", msg.Channel, msg.Source, msg.Message)
	}

}
