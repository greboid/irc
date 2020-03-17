package main

import (
	"context"
	"crypto/tls"
	"flag"
	"fmt"
	"github.com/greboid/irc/rpc"
	"github.com/kouhin/envflag"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"log"
)

var (
	RPCHost  = flag.String("rpc-host", "localhost", "gRPC server to connect to")
	RPCPort  = flag.Int("rpc-port", 8001, "gRPC server port")
	RPCToken = flag.String("rpc-token", "", "gRPC authentication token")
	Channel  = flag.String("channel", "", "Channel to send messages to")

	WebPort  = flag.Int("web-port", 8000, "Port for the web server to listen")
	DBPath   = flag.String("db-path", "/data/db", "Path to user/plugin database")
	AdminKey = flag.String("admin-key", "", "Admin key for API")
)

func main() {
	if err := envflag.Parse(); err != nil {
		log.Fatalf("Unable to load config: %s", err.Error())
	}
	db, err := New(*DBPath)
	if err != nil {
		log.Panicf("Unable to load config: %s", err.Error())
	}
	creds := credentials.NewTLS(&tls.Config{InsecureSkipVerify: true})
	conn, err := grpc.Dial(fmt.Sprintf("%s:%d", *RPCHost, *RPCPort), grpc.WithTransportCredentials(creds))
	defer func() { _ = conn.Close() }()
	client := rpc.NewIRCPluginClient(conn)
	_, err = client.Ping(rpc.CtxWithToken(context.Background(), "bearer", *RPCToken), &rpc.Empty{})
	if err != nil {
		log.Fatalf("Error with connection: %s", err.Error())
	}
	NewWeb(*WebPort, *Channel, *AdminKey, client, db).StartWeb()
}
