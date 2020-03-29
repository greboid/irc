package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/greboid/irc/rpc"
	"github.com/sebdah/goldie/v2"
	"google.golang.org/grpc"
	"io/ioutil"
	"reflect"
	"strings"
	"testing"
)

type mockIRCPluginClient struct {
}

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

func Test_github_tidyPushRefspecs(t *testing.T) {
	type fields struct {
		client rpc.IRCPluginClient
	}
	tests := []struct {
		name   string
		fields fields
		args   *pushhook
		want   *pushhook
	}{
		{
			name:   "refspec: master branch",
			fields: fields{},
			args:   &pushhook{Refspec: "refs/heads/master"},
			want:   &pushhook{Refspec: "branch master"},
		},
		{
			name:   "refspec: tag v1",
			fields: fields{},
			args:   &pushhook{Refspec: "refs/tags/v1.0.0"},
			want:   &pushhook{Refspec: "tag v1.0.0"},
		},
		{
			name:   "baserefspec: master branch",
			fields: fields{},
			args:   &pushhook{Baserefspec: "refs/heads/master"},
			want:   &pushhook{Baserefspec: "branch master"},
		},
		{
			name:   "baserefspec: tag v1",
			fields: fields{},
			args:   &pushhook{Baserefspec: "refs/tags/v1.0.0"},
			want:   &pushhook{Baserefspec: "tag v1.0.0"},
		},
		{
			name:   "refspec: non master",
			fields: fields{},
			args:   &pushhook{Baserefspec: "refs/heads/testing"},
			want:   &pushhook{Baserefspec: "branch testing"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := &github{
				client: &mockIRCPluginClient{},
			}
			if g.tidyPushRefspecs(tt.args); !reflect.DeepEqual(tt.args.Refspec, tt.want.Refspec) {
				t.Errorf("%v != %v", tt.args.Refspec, tt.want.Refspec)
			}
		})
	}
}

func Test_github_handlePushEvent(t *testing.T) {
	tests := []string{"push/basic.json", "push/tag.json"}
	gold := goldie.New(t)
	for index := range tests {
		t.Run(tests[index], func(t *testing.T) {
			g := github{
				client: &mockIRCPluginClient{},
			}
			hook := pushhook{}
			err := getTestData(tests[index], &hook)
			if err != nil {
				t.Fatal("Unable to parse example data")
			}
			got := []byte(strings.Join(g.handlePushEvent(hook), "\n"))
			gold.Assert(t, tests[index], got)
		})
	}
}

func Test_github_handleCommit(t *testing.T) {
	tests := []string{"push/commit/basic.json", "push/commit/multiline_commit_message.json", "push/commit/external_merge.json"}
	gold := goldie.New(t)
	for index := range tests {
		t.Run(tests[index], func(t *testing.T) {
			g := github{
				client: &mockIRCPluginClient{},
			}
			hook := pushhook{}
			err := getTestData(tests[index], &hook)
			if err != nil {
				t.Fatal("Unable to parse example data")
			}
			got := []byte(strings.Join(g.handleCommit(hook), "\n"))
			gold.Assert(t, tests[index], got)
		})

	}
}
