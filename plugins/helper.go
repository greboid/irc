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

type PluginHelper struct {
	RPCTarget     string
	RPCToken      string
	rpcConnection *grpc.ClientConn
	httpClient    rpc.HTTPPluginClient
	ircClient     rpc.IRCPluginClient
}

//NewHelper returns a PluginHelper that simplifies writing plugins by managing grpc connections and exposing a simple
//interface.
//It returns a PluginHelper or any errors encountered whilst creating
func NewHelper(target string, rpctoken string) (*PluginHelper, error) {
	if len(target) == 0 {
		return nil, fmt.Errorf("rpchost must be set")
	}
	if len(rpctoken) == 0 {
		return nil, fmt.Errorf("rpctoken must be set")
	}
	return &PluginHelper{
		RPCTarget: target,
		RPCToken:  rpctoken,
	}, nil
}

func (h *PluginHelper) connectToRPC() error {
	creds := credentials.NewTLS(&tls.Config{InsecureSkipVerify: true})
	conn, err := grpc.Dial(h.RPCToken, grpc.WithTransportCredentials(creds))
	if err != nil {
		return err
	}
	h.rpcConnection = conn
	return nil
}

func (h *PluginHelper) HTTPClient() error {
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

func (h *PluginHelper) connectHTTPClient() error {
	client := rpc.NewHTTPPluginClient(h.rpcConnection)
	h.httpClient = client
	return nil
}

func (h *PluginHelper) RegisterWebhook(path string, handler func(request *rpc.HttpRequest) *rpc.HttpResponse) error {
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
	stream, err := h.httpClient.GetRequest(rpc.CtxWithTokenAndPath(context.Background(), "bearer", h.RPCToken, path))
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

func (h *PluginHelper) connectIRCClient() error {
	if h.rpcConnection == nil {
		err := h.connectToRPC()
		if err != nil {
			return err
		}
	}
	client := rpc.NewIRCPluginClient(h.rpcConnection)
	_, err := client.Ping(rpc.CtxWithToken(context.Background(), "bearer", h.RPCToken), &rpc.Empty{})
	if err != nil {
		return nil
	}
	h.ircClient = client
	return nil
}

func (h *PluginHelper) IRCClient() (rpc.IRCPluginClient, error) {
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

func (h *PluginHelper) Ping() error {
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
	_, err := h.ircClient.Ping(rpc.CtxWithToken(context.Background(), "bearer", h.RPCToken), &rpc.Empty{})
	return err
}

func (h *PluginHelper) SendIRCMessage(channel string, messages []string) error {
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
	for index := range messages {
		_, err := h.ircClient.SendChannelMessage(rpc.CtxWithToken(context.Background(), "bearer", h.RPCToken), &rpc.ChannelMessage{
			Channel: channel,
			Message: messages[index],
		})
		if err != nil {
			return err
		}
	}
	return nil
}

func (h *PluginHelper) SendRawMessage(messages ...string) error {
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
	for index := range messages {
		_, err := h.ircClient.SendRawMessage(rpc.CtxWithToken(context.Background(), "bearer", h.RPCToken), &rpc.RawMessage{
			Message: messages[index],
		})
		if err != nil {
			return err
		}
	}
	return nil
}

func (h *PluginHelper) RegisterChannelMessageHandler(channel string, handler func(message *rpc.ChannelMessage)) error {
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
		rpc.CtxWithToken(context.Background(), "bearer", h.RPCToken),
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
