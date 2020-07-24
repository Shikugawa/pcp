package nodes

import core "github.com/envoyproxy/go-control-plane/envoy/api/v2/core"
var Nodes = []*core.Node{
	&core.Node{
		Cluster: "cluster.local",
		Id:      "node_0",
	},
	&core.Node{
		Cluster: "cluster.local",
		Id:      "node_1",
	},
}
