package manager

import (
	"errors"
	"fmt"
	"log"

	core "github.com/envoyproxy/go-control-plane/envoy/api/v2/core"

	"github.com/Shikugawa/pcp/pkg/config"
	"github.com/Shikugawa/pcp/pkg/factory"
	"github.com/Shikugawa/pcp/pkg/filter"
	"github.com/Shikugawa/pcp/pkg/nodes"
	"github.com/Shikugawa/pcp/pkg/snapshot"
)

var (
	listenerName                         = "default_listener"
	listenerHost                         = "0.0.0.0"
	listenerPort                         = 5000
	wasmRuntime                          = "envoy.wasm.runtime.v8"
	wasmFilterChainFactory               = factory.NewHttpWasmFilterChainFactory(wasmRuntime)
	httpFilterChainFactory               = factory.NewHttpFilterChainFactory(wasmFilterChainFactory)
	httpConnManagerFactory               = factory.NewHttpConnectionManagerFactory(httpFilterChainFactory)
	httpConnManagerListenerFilterFactory = factory.NewLHttpConnManagerListenerFilterFactory(httpConnManagerFactory)
	listenerFactory                      = factory.NewListenerFactory(httpConnManagerListenerFilterFactory)
)

type EnvoyFilterManager struct {
	NodeFilters  *nodes.NodeFilters
	ManagedNodes *nodes.ManagedNodes
	SnapShot     *snapshot.SnapShot
	Storage      *filter.FilterStorage
}

func NewEnvoyFilterManager(runtime string, wasmStoragePath string, defaultNodes []config.Node) *EnvoyFilterManager {
	wasmRuntime = runtime
	managedNodes := nodes.NewManagedNodes()
	for _, node := range defaultNodes {
		managedNodes.AddNode(node.Cluster, node.Id)
	}

	snap := snapshot.InitSnapShot()

	var nodesSlice []*core.Node
	managedNodes.GetAll().Each(func(n interface{}) bool {
		nodesSlice = append(nodesSlice, nodes.StringToNode(n.(string)))
		return false
	})

	snap.DefaultCache(nodesSlice, listenerFactory.Create(listenerName, listenerHost, uint32(listenerPort)))

	manager := &EnvoyFilterManager{
		NodeFilters:  nodes.NewNodeFilters(),
		ManagedNodes: managedNodes,
		SnapShot:     &snap,
		Storage:      filter.NewFilterStorage(wasmStoragePath),
	}

	return manager
}

func (h *EnvoyFilterManager) Append(filterType string, filterName string, targetNodes []core.Node) error {
	specifier := filter.FilterSpecifier{
		FilterType: filterType,
		FilterName: filterName,
	}

	if !h.Storage.ExistFilter(specifier) {
		return errors.New(fmt.Sprintf("%s haven't uploaded yet", specifier.String()))
	}

	nextVersion := h.SnapShot.Version + 1
	h.SnapShot.Version = nextVersion

	for _, targetNode := range targetNodes {
		if !h.ManagedNodes.Exists(targetNode.Cluster, targetNode.Id) {
			log.Println(nodes.NodeToString(&targetNode) + " not found")
			continue
		}

		if h.NodeFilters.IsRegistered(&targetNode, specifier) {
			log.Println(specifier.String() + " had been already registered to " + nodes.NodeToString(&targetNode))
			continue
		}

		h.NodeFilters.Add(&targetNode, specifier)
		h.NodeFilters.Filters(&targetNode).Each(func(f interface{}) bool {
			stringSpecifier := f.(string)
			wasmFilterChainFactory.Filters = append(wasmFilterChainFactory.Filters, filter.StringToSpecifier(stringSpecifier))
			return false
		})

		listener := listenerFactory.Create(listenerName, listenerHost, uint32(listenerPort))
		wasmFilterChainFactory.Filters = wasmFilterChainFactory.Filters[:0]

		if err := h.SnapShot.UpdateListener(listener, &targetNode, string(h.SnapShot.Version)); err != nil {
			log.Println("Failed to update " + nodes.NodeToString(&targetNode))
			continue
		}

		log.Println("Update " + nodes.NodeToString(&targetNode))
	}

	return nil
}

func (h *EnvoyFilterManager) RemoveFilter(filterType string, filterName string, targetNodes []core.Node) error {
	specifier := filter.FilterSpecifier{
		FilterType: filterType,
		FilterName: filterName,
	}

	nextVersion := h.SnapShot.Version + 1
	h.SnapShot.Version = nextVersion

	for _, targetNode := range targetNodes {
		if !h.ManagedNodes.Exists(targetNode.Cluster, targetNode.Id) {
			log.Println(nodes.NodeToString(&targetNode) + " not found")
			continue
		}

		if !h.NodeFilters.IsRegistered(&targetNode, specifier) {
			log.Println(specifier.String() + " had not been registered to " + nodes.NodeToString(&targetNode) + " yet")
			continue
		}
		fmt.Println(h.NodeFilters.Filters(&targetNode).Cardinality())

		h.NodeFilters.Remove(&targetNode, specifier)
		fmt.Println(h.NodeFilters.Filters(&targetNode).Cardinality())
		h.NodeFilters.Filters(&targetNode).Each(func(f interface{}) bool {
			stringSpecifier := f.(string)
			wasmFilterChainFactory.Filters = append(wasmFilterChainFactory.Filters, filter.StringToSpecifier(stringSpecifier))
			return false
		})

		listener := listenerFactory.Create(listenerName, listenerHost, uint32(listenerPort))
		wasmFilterChainFactory.Filters = wasmFilterChainFactory.Filters[:0]

		if err := h.SnapShot.UpdateListener(listener, &targetNode, string(h.SnapShot.Version)); err != nil {
			log.Println("Failed to update " + nodes.NodeToString(&targetNode))
			continue
		}

		log.Println("Update " + nodes.NodeToString(&targetNode))
	}

	return nil
}
