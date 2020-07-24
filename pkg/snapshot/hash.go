package snapshot

import core "github.com/envoyproxy/go-control-plane/envoy/api/v2/core"
type Hash struct{}

func (Hash) ID(node *core.Node) string {
	if node == nil {
		return "unknown"
	}
	return node.Cluster + "/" + node.Id
}
