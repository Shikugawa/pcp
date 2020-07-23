package nodes

import core "github.com/envoyproxy/go-control-plane/envoy/config/core/v3"

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
