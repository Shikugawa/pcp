package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Shikugawa/pcp/pkg/director"
	"github.com/Shikugawa/pcp/pkg/manager"
	"github.com/Shikugawa/pcp/pkg/nodes"
	"github.com/Shikugawa/pcp/pkg/xds"
)

func main() {
	var listen string
	var admin string

	flag.StringVar(&listen, "listen", "20000", "listen port")
	flag.StringVar(&admin, "admin", "3000", "listen port")
	flag.Parse()

	manager := manager.NewEnvoyFilterManager(nodes.Nodes)

	director := director.NewServer(manager).Start(admin)
	xdsServer := xds.NewServer(manager).Start(listen)

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM, syscall.SIGINT, syscall.SIGABRT)
	<-sig

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if director != nil {
		if err := director.Shutdown(ctx); err != nil {
			log.Println("Failed to graceful shutdown director:", err)
		}
	}

	if xdsServer != nil {
		xdsServer.GracefulStop()
	}
}
