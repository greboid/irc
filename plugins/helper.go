package plugins

import (
	"context"
	"crypto/tls"
	"fmt"
	"io"

	"github.com/greboid/irc/v2/rpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

type PluginHelper interface {
	GetHttpClient() error
	GetIRCClient() (rpc.IRCPluginClient, error)
	RegisterWebhook(path string, handler func(request *rpc.HttpRequest) *rpc.HttpResponse) error
	SendIRCMessage(channel string, messages []string) []error
	Ping() error
	SendRawMessage(messages []string) []error
	RegisterChannelMessageHandler(channel string, handler func(message *rpc.ChannelMessage)) error
}

type helper struct {
	rPCHost       string
	rPCPort       uint16
	rPCToken      string
	rpcConnection *grpc.ClientConn
	httpClient    rpc.HTTPPluginClient
	ircClient     rpc.IRCPluginClient
}

//NewHelper returns a PluginHelper that simplifies writing plugins by managing grpc connections and exposing a simple
//interface.
//It returns a PluginHelper or any errors encountered whilst creating
func NewHelper(rpchost string, rpcport uint16, rpctoken string) (PluginHelper, error) {
	if len(rpchost) == 0 {
		return nil, fmt.Errorf("rpchost must be set")
	}
	if len(rpctoken) == 0 {
		return nil, fmt.Errorf("rpctoken must be set")
	}
	return &helper{
		rPCHost:  rpchost,
		rPCPort:  rpcport,
		rPCToken: rpctoken,
	}, nil
}

func (h *helper) connectToRPC() error {
	creds := credentials.NewTLS(&tls.Config{InsecureSkipVerify: true})
	conn, err := grpc.Dial(fmt.Sprintf("%s:%d", h.rPCHost, h.rPCPort), grpc.WithTransportCredentials(creds))
	defer func() { _ = conn.Close() }()
	if err != nil {
		return err
	}
	h.rpcConnection = conn
	return nil
}

func (h *helper) GetHttpClient() error {
	if h.rpcConnection == nil {
		err := h.connectToRPC()
		if err != nil {
			return err
		}
	}
	err := h.connectHTTPClient()
	if err != nil {
		return err
	}
	return nil
}

func (h *helper) connectHTTPClient() error {
	client := rpc.NewHTTPPluginClient(h.rpcConnection)
	h.httpClient = client
	return nil
}

func (h *helper) RegisterWebhook(path string, handler func(request *rpc.HttpRequest) *rpc.HttpResponse) error {
	if h.rpcConnection == nil {
		err := h.connectToRPC()
		if err != nil {
			return err
		}
	}
	if h.httpClient == nil {
		err := h.connectHTTPClient()
		if err != nil {
			return err
		}
	}
	stream, err := h.httpClient.GetRequest(rpc.CtxWithTokenAndPath(context.Background(), "bearer", h.rPCToken, path))
	if err != nil {
		return err
	}
	for {
		request, err := stream.Recv()
		if err == io.EOF {
			return err
		}
		if err != nil {
			return err
		}
		response := handler(request)
		if err = stream.Send(response); err != nil {
			return err
		}
	}
}

func (h *helper) connectIRCClient() error {
	if h.rpcConnection == nil {
		err := h.connectToRPC()
		if err != nil {
			return err
		}
	}
	client := rpc.NewIRCPluginClient(h.rpcConnection)
	_, err := client.Ping(rpc.CtxWithToken(context.Background(), "bearer", h.rPCToken), &rpc.Empty{})
	if err != nil {
		return nil
	}
	h.ircClient = client
	return nil
}

func (h *helper) GetIRCClient() (rpc.IRCPluginClient, error) {
	if h.rpcConnection == nil {
		err := h.connectToRPC()
		if err != nil {
			return nil, err
		}
	}
	if h.ircClient == nil {
		err := h.connectIRCClient()
		if err != nil {
			return nil, err
		}
	}
	return h.ircClient, nil
}

func (h *helper) Ping() error {
	if h.rpcConnection == nil {
		err := h.connectToRPC()
		if err != nil {
			return err
		}
	}
	if h.ircClient == nil {
		err := h.connectIRCClient()
		if err != nil {
			return err
		}
	}
	_, err := h.ircClient.Ping(rpc.CtxWithToken(context.Background(), "bearer", h.rPCToken), &rpc.Empty{})
	return err
}

func (h *helper) SendIRCMessage(channel string, messages []string) []error {
	if h.rpcConnection == nil {
		err := h.connectToRPC()
		if err != nil {
			return []error{err}
		}
	}
	if h.ircClient == nil {
		err := h.connectIRCClient()
		if err != nil {
			return []error{err}
		}
	}
	errors := make([]error, 0)
	for index := range messages {
		_, err := h.ircClient.SendChannelMessage(rpc.CtxWithToken(context.Background(), "bearer", h.rPCToken), &rpc.ChannelMessage{
			Channel: channel,
			Message: messages[index],
		})
		if err != nil {
			errors = append(errors, err)
		}
	}
	return errors
}

func (h *helper) SendRawMessage(messages []string) []error {
	if h.rpcConnection == nil {
		err := h.connectToRPC()
		if err != nil {
			return []error{err}
		}
	}
	if h.ircClient == nil {
		err := h.connectIRCClient()
		if err != nil {
			return []error{err}
		}
	}
	errors := make([]error, 0)
	for index := range messages {
		_, err := h.ircClient.SendRawMessage(rpc.CtxWithToken(context.Background(), "bearer", h.rPCToken), &rpc.RawMessage{
			Message: messages[index],
		})
		if err != nil {
			errors = append(errors, err)
		}
	}
	return errors
}

func (h *helper) RegisterChannelMessageHandler(channel string, handler func(message *rpc.ChannelMessage)) error {
	if h.rpcConnection == nil {
		err := h.connectToRPC()
		if err != nil {
			return err
		}
	}
	if h.ircClient == nil {
		err := h.connectIRCClient()
		if err != nil {
			return err
		}
	}
	stream, err := h.ircClient.GetMessages(
		rpc.CtxWithToken(context.Background(), "bearer", h.rPCToken),
		&rpc.Channel{Name: channel},
	)
	if err != nil {
		return err
	}
	for {
		message, err := stream.Recv()
		if err == io.EOF {
			return err
		}
		if err != nil {
			return err
		}
		handler(message)
	}
}
