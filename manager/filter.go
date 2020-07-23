package manager

import (
	"fmt"
	"log"

	core "github.com/envoyproxy/go-control-plane/envoy/config/core/v3"
	hcm "github.com/envoyproxy/go-control-plane/envoy/extensions/filters/network/http_connection_manager/v3"

	"github.com/Shikugawa/pcp/filter"
	"github.com/Shikugawa/pcp/snapshot"
	"github.com/Shikugawa/pcp/snapshot/config"
)

type EnvoyFilterManager struct {
	registeredFilterSpecifiers []filter.FilterSpecifier
	SnapShot                   *snapshot.SnapShot
	StorageDriver              *filter.WasmFilterStorageDriver
}

func NewEnvoyFilterManager(defaultNodes []*core.Node) *EnvoyFilterManager {
	snap := snapshot.InitSnapShot()
	snap.DefaultCache(defaultNodes)

	manager := &EnvoyFilterManager{
		registeredFilterSpecifiers: []filter.FilterSpecifier{},
		SnapShot:                   &snap,
		StorageDriver:              filter.NewWasmFilterStorageDriver("envoy.wasm.runtime.v8"),
	}

	return manager
}

func (h *EnvoyFilterManager) Append(filterType string, filterName string, nodes []*core.Node) error {
	specifier := filter.FilterSpecifier{
		FilterType: filterType,
		FilterName: filterName,
	}
	if h.existFilter(specifier) {
		log.Println(fmt.Sprintf("%s is already registered", specifier.String()))
		return nil
	}

	if !h.StorageDriver.Storage.ExistFilter(specifier) {
		log.Println(fmt.Sprintf("%s is not already uploaded", specifier.String()))
		return nil
	}

	nextVersion := h.SnapShot.Version + 1
	h.SnapShot.Version = nextVersion
	h.registeredFilterSpecifiers = append(h.registeredFilterSpecifiers, specifier)

	filters := h.currentFilters()

	for _, node := range nodes {
		log.Println("Update " + node.Cluster + "/" + node.Id)
		log.Println(h.registeredFilterSpecifiers)
	}

	listener, err := config.GetListener(config.HCM(filters))
	if err != nil {
		return err
	}

	h.SnapShot.UpdateListener(listener, nodes, string(h.SnapShot.Version))
	return nil
}

func (h *EnvoyFilterManager) RemoveFilter(filterType string, filterName string, nodes []*core.Node) error {
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
	h.registeredFilterSpecifiers = func() []filter.FilterSpecifier {
		var updatedFilters []filter.FilterSpecifier
		for _, registeredSpecifier := range h.registeredFilterSpecifiers {
			if specifier.FilterName == registeredSpecifier.FilterName && specifier.FilterType == registeredSpecifier.FilterType {
				continue
			}
			updatedFilters = append(updatedFilters, registeredSpecifier)
		}
		return updatedFilters
	}()

	filters := h.currentFilters()

	for _, node := range nodes {
		log.Println("Update " + node.Cluster + "/" + node.Id)
		log.Println(h.registeredFilterSpecifiers)
	}

	l, err := config.GetListener(config.HCM(filters))
	if err != nil {
		return err
	}

	h.SnapShot.UpdateListener(l, nodes, string(h.SnapShot.Version))

	return nil
}

func (h *EnvoyFilterManager) currentFilters() []*hcm.HttpFilter {
	var filters []*hcm.HttpFilter
	for _, specifier := range h.registeredFilterSpecifiers {
		reg, _ := h.StorageDriver.EnvoyFilterConfig(specifier)
		filters = append(filters, reg)
	}
	return filters
}

func (h *EnvoyFilterManager) existFilter(filter filter.FilterSpecifier) bool {
	for _, specifier := range h.registeredFilterSpecifiers {
		if filter.FilterName == specifier.FilterName && filter.FilterType == specifier.FilterType {
			return true
		}
	}
	return false
}
