package xds

import (
	"context"
	"log"
	"net"

	"github.com/Shikugawa/pcp/manager"

	v2 "github.com/envoyproxy/go-control-plane/envoy/api/v2"
	srv "github.com/envoyproxy/go-control-plane/pkg/server/v2"
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
	v2.RegisterListenerDiscoveryServiceServer(grpcServer, server)

	go func() {
		log.Println("xDS control plane started...")

		listen, err := net.Listen("tcp", ":"+port)
		if err != nil {
			log.Fatalln("Failed to create listen socket")
		}
		if err := grpcServer.Serve(listen); err != nil {
			log.Fatalln("Server closed with error:", err)
		}
	}()

	return grpcServer
}
