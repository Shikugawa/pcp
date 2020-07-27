package snapshot

import (
	"github.com/Shikugawa/pcp/pkg/nodes"
	core "github.com/envoyproxy/go-control-plane/envoy/api/v2/core"
)

type CoreNodePtrHash struct{}

func (CoreNodePtrHash) ID(node *core.Node) string {
	if node == nil {
		return "unknown"
	}
	return nodes.NodeToString(node)
}
