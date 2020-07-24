package nodes

import (
	"fmt"

	core "github.com/envoyproxy/go-control-plane/envoy/api/v2/core"
)

func NodeToString(node *core.Node) string {
	return fmt.Sprintf("%s/%s", node.Cluster, node.Id)
}
