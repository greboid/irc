package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/greboid/irc/rpc"
	"google.golang.org/grpc"
	"io/ioutil"
)

type mockIRCPluginClient struct{}

func (m *mockIRCPluginClient) SendChannelMessage(context.Context, *rpc.ChannelMessage, ...grpc.CallOption) (*rpc.Error, error) {
	return nil, nil
}

func (m *mockIRCPluginClient) SendRawMessage(context.Context, *rpc.RawMessage, ...grpc.CallOption) (*rpc.Error, error) {
	return nil, nil
}

func (m *mockIRCPluginClient) GetMessages(context.Context, *rpc.Channel, ...grpc.CallOption) (rpc.IRCPlugin_GetMessagesClient, error) {
	return nil, nil
}

func (m *mockIRCPluginClient) Ping(context.Context, *rpc.Empty, ...grpc.CallOption) (*rpc.Empty, error) {
	return nil, nil
}

func getTestData(filename string, output interface{}) error {
	data, err := ioutil.ReadFile(fmt.Sprintf("./testdata/%s", filename))
	if err != nil {
		return err
	}
	err = json.Unmarshal(data, &output)
	return err
}
