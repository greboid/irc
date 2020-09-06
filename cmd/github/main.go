package main

import (
	"context"
	"crypto/hmac"
	"crypto/sha1"
	"crypto/subtle"
	"crypto/tls"
	"encoding/hex"
	"flag"
	"fmt"
	"github.com/greboid/irc/rpc"
	"github.com/kouhin/envflag"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"io"
	"log"
	"net/http"
	"strings"
)

var (
	RPCHost      = flag.String("rpc-host", "localhost", "gRPC server to connect to")
	RPCPort      = flag.Int("rpc-port", 8001, "gRPC server port")
	RPCToken     = flag.String("rpc-token", "", "gRPC authentication token")
	Channel      = flag.String("channel", "", "Channel to send messages to")
	GithubSecret = flag.String("github-secret", "", "Github secret for validating webhooks")
)

type github struct {
	client rpc.IRCPluginClient
}

func main() {
	if err := envflag.Parse(); err != nil {
		log.Fatalf("Unable to load config: %s", err.Error())
	}
	github := github{}
	log.Printf("Creating Github RPC Client")
	client, err := github.doRPC()
	if err != nil {
		log.Fatalf("Unable to create RPC Client: %s", err.Error())
	}
	github.client = client
	log.Printf("Starting github web server")
	err = github.doWeb()
	if err != nil {
		log.Panicf("Error handling web: %s", err.Error())
	}
	log.Printf("exiting")
}

func (g *github) doRPC() (rpc.IRCPluginClient, error) {
	creds := credentials.NewTLS(&tls.Config{InsecureSkipVerify: true})
	conn, err := grpc.Dial(fmt.Sprintf("%s:%d", *RPCHost, *RPCPort), grpc.WithTransportCredentials(creds))
	client := rpc.NewIRCPluginClient(conn)
	_, err = client.Ping(rpc.CtxWithToken(context.Background(), "bearer", *RPCToken), &rpc.Empty{})
	if err != nil {
		return nil, err
	}
	return client, nil
}

func (g *github) doWeb() error {
	creds := credentials.NewTLS(&tls.Config{InsecureSkipVerify: true})
	conn, err := grpc.Dial(fmt.Sprintf("%s:%d", *RPCHost, *RPCPort), grpc.WithTransportCredentials(creds))
	if err != nil {
		return err
	}
	client := rpc.NewHTTPPluginClient(conn)
	_, err = client.RegisterRoute(rpc.CtxWithToken(context.Background(), "bearer", *RPCToken), &rpc.Route{Prefix: "github"})
	if err != nil {
		return err
	}
	stream, err := client.GetRequest(rpc.CtxWithToken(context.Background(), "bearer", *RPCToken))
	for {
		request, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return nil
		}
		response := g.handleGithub(request)
		err = stream.Send(response)
		if err != nil {
			return err
		}
	}
}

func (g *github) handleGithub(request *rpc.HttpRequest) *rpc.HttpResponse {
	headers := rpc.ConvertFromRPCHeaders(request.Header)
	eventType := headers.Get("X-GitHub-Event")
	header := strings.SplitN(headers.Get("X-Hub-Signature"), "=", 2)
	if header[0] != "sha1" {
		log.Printf("Error: %s", "Bad header")
		return &rpc.HttpResponse{
			Header:               nil,
			Body:                 []byte("Bad headers"),
			Status:               http.StatusInternalServerError,
		}
	}
	if !CheckGithubSecret(request.Body, header[1], *GithubSecret) {
		log.Printf("Error: %s", "Bad hash")
		return &rpc.HttpResponse{
			Header:               nil,
			Body:                 []byte("Bad hash"),
			Status:               http.StatusBadRequest,
		}
	}
	go func() {
		webhookHandler := githubWebhookHandler{
			client: g.client,
		}
		_ = webhookHandler.handleWebhook(eventType, request.Body)
	}()
	return &rpc.HttpResponse{
		Header:               nil,
		Body:                 []byte("Delivered"),
		Status:               http.StatusOK,
	}
}

func CheckGithubSecret(bodyBytes []byte, headerSecret string, githubSecret string) bool {
	h := hmac.New(sha1.New, []byte(githubSecret))
	h.Write(bodyBytes)
	expected := fmt.Sprintf("%s", hex.EncodeToString(h.Sum(nil)))
	return len(expected) == len(headerSecret) && subtle.ConstantTimeCompare([]byte(expected), []byte(headerSecret)) == 1
}
