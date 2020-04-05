package main

import (
	"github.com/greboid/irc/rpc"
	"testing"
)

func Test_githubWebhookHandler_handleWebhook(t *testing.T) {
	type fields struct {
		client rpc.IRCPluginClient
	}
	type args struct {
		eventType string
		bodyBytes []byte
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := &githubWebhookHandler{
				client: tt.fields.client,
			}
			if err := g.handleWebhook(tt.args.eventType, tt.args.bodyBytes); (err != nil) != tt.wantErr {
				t.Errorf("handleWebhook() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_githubWebhookHandler_sendMessage(t *testing.T) {
	type fields struct {
		client rpc.IRCPluginClient
	}
	type args struct {
		messages []string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := &githubWebhookHandler{
				client: tt.fields.client,
			}
		})
	}
}