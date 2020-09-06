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
)

func NewGrpcServer(conn *irc.Connection, eventManager *irc.EventManager, rpcPort int, plugins []Plugin, webPort int) GrpcServer {
	return GrpcServer{
		conn:         conn,
		eventManager: eventManager,
		rpcPort:      rpcPort,
		plugins:      plugins,
		webPort:      webPort,
	}
}

type GrpcServer struct {
	conn         *irc.Connection
	eventManager *irc.EventManager
	rpcPort      int
	plugins      []Plugin
	webPort      int
}

func (s *GrpcServer) StartGRPC() {
	log.Print("Generating certificate")
	certificate, err := generateSelfSignedCert()
	if err != nil {
		log.Fatalf("failed to generate certifcate: %s", err.Error())
	}
	log.Printf("Starting RPC: %d", s.rpcPort)
	lis, err := tls.Listen("tcp", fmt.Sprintf(":%d", s.rpcPort), &tls.Config{Certificates: []tls.Certificate{*certificate}})
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer(
		grpc.StreamInterceptor(grpcmiddleware.ChainStreamServer(grpcauth.StreamServerInterceptor(s.authPlugin))),
		grpc.UnaryInterceptor(grpcmiddleware.ChainUnaryServer(grpcauth.UnaryServerInterceptor(s.authPlugin))),
	)
	httpsServer := NewHttpServer(s.webPort, s.plugins)
	RegisterIRCPluginServer(grpcServer, &pluginServer{s.conn, s.eventManager})
	RegisterHTTPPluginServer(grpcServer, httpsServer)
	log.Printf("Starting webserver")
	httpsServer.Start()
	log.Printf("Starting RPC server")
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
	for _, plugin := range s.plugins {
		if plugin.Token == token {
			return true
		}
	}
	return false
}
