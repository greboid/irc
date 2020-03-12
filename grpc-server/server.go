package grpc_server

import (
	"context"
	"crypto/tls"
	"fmt"
	"github.com/greboid/irc/irc"
	"github.com/greboid/irc/protos"
	grpcmiddleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpcauth "github.com/grpc-ecosystem/go-grpc-middleware/auth"
	grpcctxtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log"
)

func StartGRPC(conn *irc.Connection) {
	certificate, err := generateSelfSignedCert()
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	lis, err := tls.Listen("tcp", fmt.Sprintf(":%d", 8081), &tls.Config{Certificates: []tls.Certificate{*certificate}})
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer(
		grpc.StreamInterceptor(grpcmiddleware.ChainStreamServer(grpcauth.StreamServerInterceptor(myAuthFunction))),
		grpc.UnaryInterceptor(grpcmiddleware.ChainUnaryServer(grpcauth.UnaryServerInterceptor(myAuthFunction))),
	)
	protos.RegisterIRCPluginServer(grpcServer, &pluginServer{conn})
	err = grpcServer.Serve(lis)
	if err != nil {
		log.Printf("Error listening: %s", err.Error())
	}
}

func myAuthFunction(ctx context.Context) (context.Context, error) {
	token, err := grpcauth.AuthFromMD(ctx, "bearer")
	log.Printf("Auth: token: %v", token)
	if err != nil {
		return nil, err
	}
	tokenInfo, err := parseToken()
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "invalid auth token: %v", err)
	}
	grpcctxtags.Extract(ctx).Set("auth.sub", userClaimFromToken())
	newCtx := context.WithValue(ctx, "tokenInfo", tokenInfo)
	return newCtx, nil
}

func parseToken() (struct{}, error) {
	return struct{}{}, nil
}

func userClaimFromToken() string {
	return "foobar"
}

type pluginServer struct {
	conn *irc.Connection
}

func (ps *pluginServer) SendChannelMesssage(_ context.Context, req *protos.ChannelMessage) (*protos.Error, error) {
	ps.conn.SendRawf("PRIVMSG %s :%s", req.Channel, req.Message)
	return &protos.Error{
		Message: "",
	}, nil
}
func (*pluginServer) SendRawMessage(_ context.Context, _ *protos.RawMessage) (*protos.Error, error) {
	return &protos.Error{
		Message: "",
	}, nil
}
