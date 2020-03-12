package main

import (
	"context"
	"crypto/tls"
	"fmt"
	"github.com/greboid/irc/config"
	"github.com/greboid/irc/rpc"
	"github.com/grpc-ecosystem/go-grpc-middleware/util/metautils"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"
	"log"
)

func ctxWithToken(ctx context.Context, scheme string, token string) context.Context {
	md := metadata.Pairs("authorization", fmt.Sprintf("%s %v", scheme, token))
	nCtx := metautils.NiceMD(md).ToOutgoing(ctx)
	return nCtx
}

func main() {
	conf := config.GetConfig()
	creds := credentials.NewTLS(&tls.Config{InsecureSkipVerify: true})
	conn, err := grpc.Dial("localhost:8081", grpc.WithTransportCredentials(creds))
	if err != nil {
		log.Fatalf("failed to fake plugin: %v", err)
	}
	defer conn.Close()
	client := rpc.NewIRCPluginClient(conn)
	context.Background()
	_, err = client.SendChannelMesssage(ctxWithToken(context.Background(), "bearer", "bad_token"), &rpc.ChannelMessage{
		Channel: conf.Channel,
		Message: "RPC",
	})
	if err != nil {
		log.Printf("Error sending: %s", err.Error())
	}
}
