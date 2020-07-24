package nodes

import (
	"sync"

	core "github.com/envoyproxy/go-control-plane/envoy/api/v2/core"
)

type ManagedNodesType struct {
	nodes []*core.Node
	mux   sync.RWMutex
}

func (m *ManagedNodesType) AddNode(clusterName string, nodeId string) {
	node := &core.Node{
		Cluster: clusterName,
		Id:      nodeId,
	}
	m.mux.Lock()
	defer m.mux.Unlock()
	m.nodes = append(m.nodes, node)
}

func (m *ManagedNodesType) Exists(checkNode *core.Node) bool {
	for _, node := range m.nodes {
		if node.Cluster == checkNode.Cluster && node.Id == checkNode.Id {
			return true
		}
	}
	return false
}

func (m *ManagedNodesType) GetAll() []*core.Node {
	m.mux.Lock()
	defer m.mux.Unlock()
	return m.nodes
}

var ManagedNodes = &ManagedNodesType{nodes: make([]*core.Node, 0), mux: sync.RWMutex{}}
