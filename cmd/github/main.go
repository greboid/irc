package main

import (
	"context"
	"crypto/tls"
	"encoding/json"
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
	switch eventType {
	case "push":
		data := pushhook{}
		err = json.Unmarshal(bodyBytes, &data)
		if err == nil {
			if strings.HasPrefix(data.Refspec, "refs/heads/") {
				data.Refspec = fmt.Sprintf("branch %s", strings.TrimPrefix(data.Refspec, "refs/heads/"))
			} else if strings.HasPrefix(data.Refspec, "refs/tags/") {
				data.Refspec = fmt.Sprintf("tag %s", strings.TrimPrefix(data.Refspec, "refs/tags/"))
			}
			if data.Created {
				go g.handleCreate(data)
			} else if data.Deleted {
				go g.handleDelete(data)
			} else {
				go g.handleCommit(data)
			}
		} else {
			log.Printf("Error handling PR: %s", err.Error())
			_, _ = writer.Write([]byte("Error."))
		}
	case "pull_request":
		data := prhook{}
		err = json.Unmarshal(bodyBytes, &data)
		if err == nil {
			if data.Action == "opened" {
				go g.handlePROpen(data)
			} else if data.Action == "closed" {
				if data.PullRequest.Merged == "" {
					go g.handlePRClose(data)
				} else {
					go g.handlePRMerged(data)
				}
			}
		} else {
			log.Printf("Error handling PR: %s", err.Error())
			_, _ = writer.Write([]byte("Error."))
		}
	}
}

func (g *github) handlePRClose(data prhook) {
	g.sendMessage(fmt.Sprintf("[%s] %s closed PR: %s -  %s",
		data.Repository.FullName,
		data.PullRequest.User.Login,
		data.PullRequest.Title,
		data.PullRequest.HtmlURL,
	))
}

func (g *github) handlePRMerged(data prhook) {
	g.sendMessage(fmt.Sprintf("[%s] %s merged PR from %s: %s -  %s",
		data.Repository.FullName,
		data.PullRequest.MergedBy.Login,
		data.PullRequest.User.Login,
		data.PullRequest.Title,
		data.PullRequest.HtmlURL,
	))
}

func (g *github) handlePROpen(data prhook) {
	g.sendMessage(fmt.Sprintf("[%s] %s submitted PR: %s -  %s",
		data.Repository.FullName,
		data.PullRequest.User.Login,
		data.PullRequest.Title,
		data.PullRequest.HtmlURL,
	))
}

func (g *github) handleDelete(data pushhook) {
	g.sendMessage(fmt.Sprintf("[%s] %s deleted %s",
		data.Repository.FullName,
		data.Pusher.Name,
		data.Refspec,
	))
}

func (g *github) handleCreate(data pushhook) {
	if data.Baserefspec == "" {
		g.sendMessage(fmt.Sprintf("[%s] %s created %s - %s",
			data.Repository.FullName,
			data.Pusher.Name,
			data.Refspec,
			data.CompareLink,
		))
	} else {
		g.sendMessage(fmt.Sprintf("[%s] %s created %s from %s - %s",
			data.Repository.FullName,
			data.Pusher.Name,
			data.Refspec,
			data.Baserefspec,
			data.CompareLink,
		))
	}
}

func (g *github) handleCommit(data pushhook) {
	g.sendMessage(fmt.Sprintf("[%s] %s pushed %d commits to %s - %s",
		data.Repository.FullName,
		data.Pusher.Name,
		len(data.Commits),
		data.Refspec,
		data.CompareLink))
	for _, commit := range data.Commits {
		g.sendMessage(fmt.Sprintf("[%s] %s committed %s - %s",
			data.Repository.FullName,
			commit.Committer.User,
			commit.ID[len(commit.ID)-6:],
			commit.Message))
	}
}

func (g *github) sendMessage(message string) {
	_, err := g.client.SendChannelMessage(rpc.CtxWithToken(context.Background(), "bearer", *RPCToken), &rpc.ChannelMessage{
		Channel: *Channel,
		Message: message,
	})
	if err != nil {
		log.Printf("Error sending to channel: %s", err.Error())
		return
	}
}
