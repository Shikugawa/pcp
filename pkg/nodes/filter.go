package nodes

import (
	"sync"

	"github.com/Shikugawa/pcp/pkg/filter"
	set "github.com/deckarep/golang-set"
	core "github.com/envoyproxy/go-control-plane/envoy/api/v2/core"
)

type NodeFilters struct {
	filters map[string]set.Set
	mux     sync.Mutex
}

func NewNodeFilters() *NodeFilters {
	return &NodeFilters{
		filters: make(map[string]set.Set, 0),
	}
}

func (n *NodeFilters) Add(node *core.Node, specifier filter.FilterSpecifier) {
	n.mux.Lock()
	defer n.mux.Unlock()
	nodeName := NodeToString(node)
	if _, ok := n.filters[nodeName]; !ok {
		n.filters[nodeName] = set.NewSet()
	}
	n.filters[nodeName].Add(specifier.String())

}

func (n *NodeFilters) Remove(node *core.Node, specifier filter.FilterSpecifier) {
	n.mux.Lock()
	defer n.mux.Unlock()
	nodeName := NodeToString(node)
	if _, ok := n.filters[nodeName]; !ok {
		return
	}
	n.filters[nodeName].Remove(specifier.String())
}

func (n *NodeFilters) Filters(node *core.Node) set.Set {
	if s, ok := n.filters[NodeToString(node)]; ok {
		return s
	}
	return set.NewSet()

}

func (n *NodeFilters) IsRegistered(node *core.Node, specifier filter.FilterSpecifier) bool {
	if _, ok := n.filters[NodeToString(node)]; !ok {
		return false
	}
	return n.filters[NodeToString(node)].Contains(specifier.String())

}
