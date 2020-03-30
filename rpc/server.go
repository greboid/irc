package rpc

import (
	"context"
	"crypto/tls"
	"fmt"
	"github.com/greboid/irc/irc"
	grpcmiddleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpcauth "github.com/grpc-ecosystem/go-grpc-middleware/auth"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log"
	"strings"
)

type GrpcServer struct {
	Conn    *irc.Connection
	EventManager irc.EventManager
	RPCPort int
	Plugins []Plugin
}

func (s *GrpcServer) StartGRPC() {
	log.Print("Generating certificate")
	certificate, err := generateSelfSignedCert()
	if err != nil {
		log.Fatalf("failed to generate certifcate: %s", err.Error())
	}
	log.Printf("Starting RPC: %d", s.RPCPort)
	lis, err := tls.Listen("tcp", fmt.Sprintf(":%d", s.RPCPort), &tls.Config{Certificates: []tls.Certificate{*certificate}})
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer(
		grpc.StreamInterceptor(grpcmiddleware.ChainStreamServer(grpcauth.StreamServerInterceptor(s.authPlugin))),
		grpc.UnaryInterceptor(grpcmiddleware.ChainUnaryServer(grpcauth.UnaryServerInterceptor(s.authPlugin))),
	)
	RegisterIRCPluginServer(grpcServer, &pluginServer{s.Conn, s.EventManager})
	err = grpcServer.Serve(lis)
	if err != nil {
		log.Printf("Error listening: %s", err.Error())
	}
}

func (s *GrpcServer) authPlugin(ctx context.Context) (context.Context, error) {
	token, err := grpcauth.AuthFromMD(ctx, "bearer")
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "invalid auth token: %s", err.Error())
	}
	if !s.checkPlugin(token) {
		return nil, status.Errorf(codes.Unauthenticated, "access denied")
	}
	return ctx, nil
}

func (s *GrpcServer) checkPlugin(token string) bool {
	for _, plugin := range s.Plugins {
		if plugin.Token == token {
			return true
		}
	}
	return false
}

type pluginServer struct {
	conn irc.Sender
	EventManager irc.EventManager
}

func (ps *pluginServer) SendChannelMessage(_ context.Context, req *ChannelMessage) (*Error, error) {
	ps.conn.SendRawf("PRIVMSG %s :%s", req.Channel, req.Message)
	return &Error{
		Message: "",
	}, nil
}
func (*pluginServer) SendRawMessage(_ context.Context, _ *RawMessage) (*Error, error) {
	return &Error{
		Message: "",
	}, nil
}

func (ps *pluginServer) GetMessages(channel *Channel, stream IRCPlugin_GetMessagesServer) error {
	exitLoop := make(chan bool, 1)
	chanMessage := make(chan *irc.Message, 1)
	channelName := channel.Name
	partHandler := func(channelPart irc.Channel) {
		if channelPart.Name == channelName {
			exitLoop <- true
		}
	}
	messageHandler := func(message irc.Message) {
		if channelName == "*" || strings.ToLower(message.Params[0]) == strings.ToLower(channelName) {
			chanMessage <- &message
		}
	}
	ps.EventManager.SubscribeChannelPart(partHandler)
	defer ps.EventManager.UnsubscribeChannelPart(partHandler)
	ps.EventManager.SubscribeChannelMessage(messageHandler)
	defer ps.EventManager.UnsubscribeChannelMessage(messageHandler)
	for {
		select {
		case <-exitLoop:
			return nil
		case msg := <-chanMessage:
			if err := stream.Send(&ChannelMessage{Channel: strings.ToLower(msg.Params[0]), Message: strings.Join(msg.Params[1:], " "), Source: msg.Source}); err != nil {
				return err
			}
		}
	}
}

func (ps *pluginServer) Ping(context.Context, *Empty) (*Empty, error) {
	return &Empty{}, nil
}
