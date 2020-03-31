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
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"time"
)

var (
	RPCHost      = flag.String("rpc-host", "localhost", "gRPC server to connect to")
	RPCPort      = flag.Int("rpc-port", 8001, "gRPC server port")
	RPCToken     = flag.String("rpc-token", "", "gRPC authentication token")
	Channel      = flag.String("channel", "", "Channel to send messages to")
	GithubSecret = flag.String("github-secret", "", "Github secret for validating webhooks")
	WebPort      = flag.Int("web-port", 8000, "Web port for receiving github webhooks")
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
	github.doWeb()

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

func (g *github) doWeb() {
	mux := http.NewServeMux()
	mux.HandleFunc("/github", g.handleGithub)
	server := http.Server{
		Addr:    fmt.Sprintf(":%d", *WebPort),
		Handler: mux,
	}
	go func() {
		log.Print(server.ListenAndServe())
	}()
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, os.Kill)
	log.Printf("Waiting for stop")
	<-stop
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Unable to shutdown: %s", err.Error())
	}
}

func (g *github) handleGithub(writer http.ResponseWriter, request *http.Request) {
	bodyBytes, err := ioutil.ReadAll(request.Body)
	defer func() { _ = request.Body.Close() }()
	if err != nil {
		log.Printf("Error reading body: %s", err.Error())
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	eventType := request.Header.Get("X-GitHub-Event")
	header := strings.SplitN(request.Header.Get("X-Hub-Signature"), "=", 2)
	if header[0] != "sha1" {
		log.Printf("Error: %s", "Bad header")
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	if !CheckGithubSecret(bodyBytes, header[1], *GithubSecret) {
		log.Printf("Error: %s", "Bad hash")
		writer.WriteHeader(http.StatusBadRequest)
	}
	_, _ = writer.Write([]byte("Delivered."))
	webhookHandler := githubWebhookHandler{
		client: g.client,
	}
	if err := webhookHandler.handleWebhook(eventType, bodyBytes); err != nil {
		_, _ = writer.Write([]byte("Error."))
	}
}

func CheckGithubSecret(bodyBytes []byte, headerSecret string, githubSecret string) bool {
	h := hmac.New(sha1.New, []byte(githubSecret))
	h.Write(bodyBytes)
	expected := fmt.Sprintf("%s", hex.EncodeToString(h.Sum(nil)))
	return len(expected) == len(headerSecret) && subtle.ConstantTimeCompare([]byte(expected), []byte(headerSecret)) == 1
}
