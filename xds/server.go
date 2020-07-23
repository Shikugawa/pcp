package xds

import (
	"context"
	"log"
	"net"

	"github.com/Shikugawa/pcp/manager"

	listener "github.com/envoyproxy/go-control-plane/envoy/service/listener/v3"
	srv "github.com/envoyproxy/go-control-plane/pkg/server/v3"
	"google.golang.org/grpc"
)

type Server struct {
	envoyFilterManager *manager.EnvoyFilterManager
}

func NewServer(envoyFilterManager *manager.EnvoyFilterManager) *Server {
	return &Server{
		envoyFilterManager: envoyFilterManager,
	}
}

func (s *Server) Start(port string) *grpc.Server {
	server := srv.NewServer(context.Background(), s.envoyFilterManager.SnapShot.Cache, nil)
	grpcServer := grpc.NewServer()
	listener.RegisterListenerDiscoveryServiceServer(grpcServer, server)

	listen, err := net.Listen("tcp", port)
	if err != nil {
		return nil
	}

	go func() {
		if err := grpcServer.Serve(listen); err != nil {
			log.Fatalln("Server closed with error:", err)
		}
	}()

	return grpcServer
}
