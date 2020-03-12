package rpc

import (
	"context"
	"crypto/tls"
	"fmt"
	"github.com/greboid/irc/database"
	"github.com/greboid/irc/irc"
	grpcmiddleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpcauth "github.com/grpc-ecosystem/go-grpc-middleware/auth"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log"
)

type GrpcServer struct {
	Conn *irc.Connection
	DB   *database.DB
}

func (s *GrpcServer) StartGRPC() {
	certificate, err := generateSelfSignedCert()
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	lis, err := tls.Listen("tcp", fmt.Sprintf(":%d", 8081), &tls.Config{Certificates: []tls.Certificate{*certificate}})
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer(
		grpc.StreamInterceptor(grpcmiddleware.ChainStreamServer(grpcauth.StreamServerInterceptor(s.authPlugin))),
		grpc.UnaryInterceptor(grpcmiddleware.ChainUnaryServer(grpcauth.UnaryServerInterceptor(s.authPlugin))),
	)
	RegisterIRCPluginServer(grpcServer, &pluginServer{s.Conn})
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
	if !s.DB.CheckPlugin(token) {
		return nil, status.Errorf(codes.Unauthenticated, "access denied")
	}
	return ctx, nil
}

type pluginServer struct {
	conn *irc.Connection
}

func (ps *pluginServer) SendChannelMesssage(_ context.Context, req *ChannelMessage) (*Error, error) {
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
