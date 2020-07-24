package main

import (
	"context"
	"flag"
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Shikugawa/pcp/pkg/config"
	"github.com/Shikugawa/pcp/pkg/director"
	"github.com/Shikugawa/pcp/pkg/manager"
	"github.com/Shikugawa/pcp/pkg/nodes"
	"github.com/Shikugawa/pcp/pkg/xds"
)

func readConfig(configPath string) ([]byte, error) {
	res, err := ioutil.ReadFile(configPath)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func loadNodes(attachedNodes []config.Node) {
	for _, node := range attachedNodes {
		nodes.ManagedNodes.AddNode(node.Cluster, node.Id)
	}
}

func main() {
	var listen string
	var admin string
	var configPath string

	flag.StringVar(&listen, "listen", "20000", "listen port")
	flag.StringVar(&admin, "admin", "3000", "listen port")
	flag.StringVar(&configPath, "config", "", "config path")
	flag.Parse()

	configData, err := readConfig(configPath)
	if err != nil {
		log.Fatalln(err)
	}

	parsedConfig, err := config.ParseConfigRoot(string(configData))
	if err != nil {
		log.Fatalln(err)
	}

	loadNodes(parsedConfig.Nodes)
	manager := manager.NewEnvoyFilterManager(parsedConfig.Runtime, parsedConfig.StoragePath)

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
