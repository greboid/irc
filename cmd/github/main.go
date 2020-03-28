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
	if err := g.handleWebhook(eventType, bodyBytes); err != nil {
		_, _ = writer.Write([]byte("Error."))
	}
}

func (g *github) handleWebhook(eventType string, bodyBytes []byte) error {
	switch eventType {
	case "push":
		data := pushhook{}
		err := json.Unmarshal(bodyBytes, &data)
		if err == nil {
			go g.sendMessage(g.handlePushEvent(data))
		} else {
			log.Printf("Error handling push: %s", err.Error())
			return err
		}
	case "pull_request":
		data := prhook{}
		err := json.Unmarshal(bodyBytes, &data)
		if err == nil {
			go g.sendMessage(g.handlePREvent(data))
		} else {
			log.Printf("Error handling PR: %s", err.Error())
			return err
		}
	}
	return nil
}

func (g *github) handlePREvent(data prhook) (messages []string) {
	if data.Action == "opened" {
		return g.handlePROpen(data)
	} else if data.Action == "closed" {
		if data.PullRequest.Merged == "" {
			return g.handlePRClose(data)
		} else {
			return g.handlePRMerged(data)
		}
	}
	return
}

func (g *github) handlePRClose(data prhook) (messages []string) {
	messages = append(messages, fmt.Sprintf(
		"[%s] %s closed PR: %s -  %s",
		data.Repository.FullName,
		data.PullRequest.User.Login,
		data.PullRequest.Title,
		data.PullRequest.HtmlURL,
	))
	return
}

func (g *github) handlePRMerged(data prhook) (messages []string) {
	messages = append(messages, fmt.Sprintf(
		"[%s] %s merged PR from %s: %s -  %s",
		data.Repository.FullName,
		data.PullRequest.MergedBy.Login,
		data.PullRequest.User.Login,
		data.PullRequest.Title,
		data.PullRequest.HtmlURL,
	))
	return
}

func (g *github) handlePROpen(data prhook) (messages []string) {
	messages = append(messages, fmt.Sprintf(
		"[%s] %s submitted PR: %s -  %s",
		data.Repository.FullName,
		data.PullRequest.User.Login,
		data.PullRequest.Title,
		data.PullRequest.HtmlURL,
	))
	return
}

func (g *github) handlePushEvent(data pushhook) (messages []string) {
	if strings.HasPrefix(data.Refspec, "refs/heads/") {
		data.Refspec = fmt.Sprintf("branch %s", strings.TrimPrefix(data.Refspec, "refs/heads/"))
	} else if strings.HasPrefix(data.Refspec, "refs/tags/") {
		data.Refspec = fmt.Sprintf("tag %s", strings.TrimPrefix(data.Refspec, "refs/tags/"))
	}
	if data.Created {
		return g.handleCreate(data)
	} else if data.Deleted {
		return g.handleDelete(data)
	} else {
		return g.handleCommit(data)
	}
}

func (g *github) handleDelete(data pushhook) (messages []string) {
	messages = append(messages, fmt.Sprintf(
		"[%s] %s deleted %s",
		data.Repository.FullName,
		data.Pusher.Name,
		data.Refspec,
	))
	return
}

func (g *github) handleCreate(data pushhook) (messages []string) {
	if data.Baserefspec == "" {
		messages = append(messages, fmt.Sprintf(
			"[%s] %s created %s - %s",
			data.Repository.FullName,
			data.Pusher.Name,
			data.Refspec,
			data.CompareLink,
		))
	} else {
		messages = append(messages, fmt.Sprintf(
			"[%s] %s created %s from %s - %s",
			data.Repository.FullName,
			data.Pusher.Name,
			data.Refspec,
			data.Baserefspec,
			data.CompareLink,
		))
	}
	return
}

func (g *github) handleCommit(data pushhook) (messages []string) {
	messages = append(messages, fmt.Sprintf(
		"[%s] %s pushed %d commits to %s - %s",
		data.Repository.FullName,
		data.Pusher.Name,
		len(data.Commits),
		data.Refspec,
		data.CompareLink,
	))
	for _, commit := range data.Commits {
		messages = append(messages, fmt.Sprintf(
			"[%s] %s committed %s - %s",
			data.Repository.FullName,
			commit.Author.User,
			commit.ID[len(commit.ID)-6:],
			strings.SplitN(commit.Message, "\n", 2)[0],
		))
	}
	return
}

func (g *github) sendMessage(messages []string) {
	for index := range messages {
		_, err := g.client.SendChannelMessage(rpc.CtxWithToken(context.Background(), "bearer", *RPCToken), &rpc.ChannelMessage{
			Channel: *Channel,
			Message: messages[index],
		})
		if err != nil {
			log.Printf("Error sending to channel: %s", err.Error())
		}
	}
}
