package main

import (
	"context"
	"encoding/json"
	"github.com/greboid/irc/rpc"
	"go.uber.org/zap"
)

type githubWebhookHandler struct {
	client rpc.IRCPluginClient
	log    *zap.SugaredLogger
}

func (g *githubWebhookHandler) handleWebhook(eventType string, bodyBytes []byte) error {
	switch eventType {
	case "ping":
		data := pinghook{}
		err := json.Unmarshal(bodyBytes, &data)
		if err == nil {
			go func() {
				err := g.sendMessage([]string{"Ping received."})
				if len(err) > 0 {
					for index := range err {
						g.log.Errorf("Error handling push: %s", err[index].Error())
					}
				}
			}()
		} else {
			g.log.Errorf("Error handling push: %s", err.Error())
			return err
		}
	case "push":
		data := pushhook{}
		handler := githubPushHandler{}
		err := json.Unmarshal(bodyBytes, &data)
		if err == nil {
			go func() {
				err := g.sendMessage(handler.handlePushEvent(data))
				if len(err) > 0 {
					for index := range err {
						g.log.Errorf("Error handling push: %s", err[index].Error())
					}
				}
			}()
		} else {
			g.log.Errorf("Error handling push: %s", err.Error())
			return err
		}
	case "pull_request":
		data := prhook{}
		handler := githubPRHandler{}
		err := json.Unmarshal(bodyBytes, &data)
		if err == nil {
			go func() {
				err := g.sendMessage(handler.handlePREvent(data))
				if len(err) > 0 {
					for index := range err {
						g.log.Errorf("Error handling push: %s", err[index].Error())
					}
				}
			}()
		} else {
			g.log.Errorf("Error handling PR: %s", err.Error())
			return err
		}
	case "issues":
		data := issuehook{}
		handler := githubissuehandler{}
		err := json.Unmarshal(bodyBytes, &data)
		if err == nil {
			go func() {
				err := g.sendMessage(handler.handleIssueEvent(data))
				if len(err) > 0 {
					for index := range err {
						g.log.Errorf("Error handling push: %s", err[index].Error())
					}
				}
			}()
		} else {
			g.log.Errorf("Error handling PR: %s", err.Error())
			return err
		}
	case "issue_comment":
		data := issuehook{}
		handler := githubIssueCommenthandler{}
		err := json.Unmarshal(bodyBytes, &data)
		if err == nil {
			go func() {
				err := g.sendMessage(handler.handleIssueCommentEvent(data))
				if len(err) > 0 {
					for index := range err {
						g.log.Errorf("Error handling push: %s", err[index].Error())
					}
				}
			}()
		} else {
			g.log.Errorf("Error handling PR: %s", err.Error())
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

func (g *githubWebhookHandler) sendMessage(messages []string) []error {
	errors := make([]error, 0)
	for index := range messages {
		_, err := g.client.SendChannelMessage(rpc.CtxWithToken(context.Background(), "bearer", *RPCToken), &rpc.ChannelMessage{
			Channel: *Channel,
			Message: messages[index],
		})
		if err != nil {
			errors = append(errors, err)
		}
	}
	return errors
}
