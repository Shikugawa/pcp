package nodes

import (
	"sync"

	set "github.com/deckarep/golang-set"
	core "github.com/envoyproxy/go-control-plane/envoy/api/v2/core"
)

type ManagedNodes struct {
	nodes set.Set
	mux   sync.RWMutex
}

func NewManagedNodes() *ManagedNodes {
	return &ManagedNodes{nodes: set.NewSet(), mux: sync.RWMutex{}}
}

func (m *ManagedNodes) AddNode(clusterName string, nodeId string) {
	node := &core.Node{
		Cluster: clusterName,
		Id:      nodeId,
	}
	m.mux.Lock()
	defer m.mux.Unlock()
	m.nodes.Add(NodeToString(node))
}

func (m *ManagedNodes) Exists(clusterName string, nodeId string) bool {
	node := &core.Node{
		Cluster: clusterName,
		Id:      nodeId,
	}
	return m.nodes.Contains(NodeToString(node))
}

func (m *ManagedNodes) GetAll() set.Set {
	m.mux.Lock()
	defer m.mux.Unlock()
	return m.nodes
}
