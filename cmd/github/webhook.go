package main

import (
	"context"
	"encoding/json"
	"github.com/greboid/irc/rpc"
	"log"
)

type githubWebhookHandler struct {
	client rpc.IRCPluginClient
}

func (g *githubWebhookHandler) handleWebhook(eventType string, bodyBytes []byte) error {
	switch eventType {
	case "push":
		data := pushhook{}
		handler := githubPushHandler{}
		err := json.Unmarshal(bodyBytes, &data)
		if err == nil {
			go g.sendMessage(handler.handlePushEvent(data))
		} else {
			log.Printf("Error handling push: %s", err.Error())
			return err
		}
	case "pull_request":
		data := prhook{}
		handler := githubPRHandler{}
		err := json.Unmarshal(bodyBytes, &data)
		if err == nil {
			go g.sendMessage(handler.handlePREvent(data))
		} else {
			log.Printf("Error handling PR: %s", err.Error())
			return err
		}
	case "issues":
		data := issuehook{}
		handler := githubissuehandler{}
		err := json.Unmarshal(bodyBytes, &data)
		if err == nil {
			go g.sendMessage(handler.handleIssueEvent(data))
		} else {
			log.Printf("Error handling PR: %s", err.Error())
			return err
		}
	case "issue_comment":
		data := issuehook{}
		handler := githubIssueCommenthandler{}
		err := json.Unmarshal(bodyBytes, &data)
		if err == nil {
			go g.sendMessage(handler.handleIssueCommentEvent(data))
		} else {
			log.Printf("Error handling PR: %s", err.Error())
			return err
		}
	case "check_run":
		// TODO: Handle
		return nil
	case "release":
		// TODO: Handle
		return nil
	case "create":
		// TODO: Handle
		return nil
	case "check_suite":
		// TODO: Handle
		return nil
	}
	return nil
}

func (g *githubWebhookHandler) sendMessage(messages []string) {
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
