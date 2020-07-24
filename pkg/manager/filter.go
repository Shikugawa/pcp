package manager

import (
	"fmt"
	"log"

	core "github.com/envoyproxy/go-control-plane/envoy/api/v2/core"

	"github.com/Shikugawa/pcp/pkg/factory"
	"github.com/Shikugawa/pcp/pkg/filter"
	"github.com/Shikugawa/pcp/pkg/nodes"
	"github.com/Shikugawa/pcp/pkg/snapshot"
)

var (
	wasmRuntime            = "envoy.wasm.runtime.v8"
	wasmFilterChainFactory = factory.NewHttpWasmFilterChainFactory(wasmRuntime)
	httpFilterChainFactory = factory.NewHttpFilterChainFactory(wasmFilterChainFactory)
	httpConnManagerFactory = factory.NewHttpConnectionManagerFactory(httpFilterChainFactory)
	listenerFilterFactory  = factory.NewListenerFilterFactory(httpConnManagerFactory)
	listenerFactory        = factory.NewListenerFactory(listenerFilterFactory)
)

type EnvoyFilterManager struct {
	registeredFilterSpecifiers []filter.FilterSpecifier
	SnapShot                   *snapshot.SnapShot
	Storage                    *filter.FilterStorage
}

func NewEnvoyFilterManager(runtime string) *EnvoyFilterManager {
	wasmRuntime = runtime

	snap := snapshot.InitSnapShot()
	snap.DefaultCache(nodes.ManagedNodes.GetAll(), listenerFactory.Create())

	manager := &EnvoyFilterManager{
		registeredFilterSpecifiers: []filter.FilterSpecifier{},
		SnapShot:                   &snap,
		Storage:                    filter.NewFilterStorage(),
	}

	return manager
}

func (h *EnvoyFilterManager) Append(filterType string, filterName string, targetNodes []*core.Node) error {
	specifier := filter.FilterSpecifier{
		FilterType: filterType,
		FilterName: filterName,
	}
	if h.existFilter(specifier) {
		log.Println(fmt.Sprintf("%s is already registered", specifier.String()))
		return nil
	}

	if !h.Storage.ExistFilter(specifier) {
		log.Println(fmt.Sprintf("%s is not already uploaded", specifier.String()))
		return nil
	}

	nextVersion := h.SnapShot.Version + 1
	h.SnapShot.Version = nextVersion
	h.addRegisteredFilter(specifier)

	actualNodes := []*core.Node{}
	for _, targetNode := range targetNodes {
		if !nodes.ManagedNodes.Exists(targetNode) {
			continue
		}
		actualNodes = append(actualNodes, targetNode)
		log.Println("Update " + nodes.NodeToString(targetNode))
		log.Println(h.registeredFilterSpecifiers)
	}

	wasmFilterChainFactory.Filters = h.registeredFilterSpecifiers
	listener := listenerFactory.Create()

	h.SnapShot.UpdateListener(listener, actualNodes, string(h.SnapShot.Version))
	return nil
}

func (h *EnvoyFilterManager) RemoveFilter(filterType string, filterName string, targetNodes []*core.Node) error {
	specifier := filter.FilterSpecifier{
		FilterType: filterType,
		FilterName: filterName,
	}
	if !h.existFilter(specifier) {
		log.Println(fmt.Sprintf("%s isn't registered", specifier.String()))
		return nil
	}

	nextVersion := h.SnapShot.Version + 1
	h.SnapShot.Version = nextVersion
	h.removeRegisteredFilter(specifier)

	actualNodes := []*core.Node{}
	for _, targetNode := range targetNodes {
		if !nodes.ManagedNodes.Exists(targetNode) {
			continue
		}
		actualNodes = append(actualNodes, targetNode)
		log.Println("Update " + nodes.NodeToString(targetNode))
		log.Println(h.registeredFilterSpecifiers)
	}

	wasmFilterChainFactory.Filters = h.registeredFilterSpecifiers
	listener := listenerFactory.Create()

	h.SnapShot.UpdateListener(listener, actualNodes, string(h.SnapShot.Version))

	return nil
}

func (h *EnvoyFilterManager) addRegisteredFilter(specifier filter.FilterSpecifier) {
	h.registeredFilterSpecifiers = append(h.registeredFilterSpecifiers, specifier)
}

func (h *EnvoyFilterManager) removeRegisteredFilter(specifier filter.FilterSpecifier) {
	var updatedFilters []filter.FilterSpecifier
	for _, registeredSpecifier := range h.registeredFilterSpecifiers {
		if specifier.FilterName == registeredSpecifier.FilterName && specifier.FilterType == registeredSpecifier.FilterType {
			continue
		}
		updatedFilters = append(updatedFilters, registeredSpecifier)
	}
	h.registeredFilterSpecifiers = updatedFilters
}

func (h *EnvoyFilterManager) existFilter(filter filter.FilterSpecifier) bool {
	for _, specifier := range h.registeredFilterSpecifiers {
		if filter.FilterName == specifier.FilterName && filter.FilterType == specifier.FilterType {
			return true
		}
	}
	return false
}
