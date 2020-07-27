package nodes

import (
	"fmt"
	"strings"

	core "github.com/envoyproxy/go-control-plane/envoy/api/v2/core"
)

func NodeToString(node *core.Node) string {
	return fmt.Sprintf("%s/%s", node.Cluster, node.Id)
}

func StringToNode(s string) *core.Node {
	splitedNode := strings.Split(s, "/")
	node := &core.Node{
		Cluster: splitedNode[0],
		Id:      splitedNode[1],
	}
	return node
}
