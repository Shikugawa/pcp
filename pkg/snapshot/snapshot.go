package snapshot

import (
	"github.com/Shikugawa/pcp/pkg/nodes"
	v2 "github.com/envoyproxy/go-control-plane/envoy/api/v2"

	core "github.com/envoyproxy/go-control-plane/envoy/api/v2/core"
	"github.com/envoyproxy/go-control-plane/pkg/cache/types"
	cache "github.com/envoyproxy/go-control-plane/pkg/cache/v2"
)

type SnapShot struct {
	Version int
	Cache   cache.SnapshotCache
}

func InitSnapShot() SnapShot {
	return SnapShot{
		Version: 0,
		Cache:   cache.NewSnapshotCache(false, CoreNodePtrHash{}, nil),
	}
}

func (s *SnapShot) DefaultCache(envoyNodes []*core.Node, defaultListener *v2.Listener) error {
	for _, node := range envoyNodes {
		snapshotCache := cache.NewSnapshot("1", nil, nil, nil, []types.Resource{defaultListener}, nil)
		if err := s.Cache.SetSnapshot(nodes.NodeToString(node), snapshotCache); err != nil {
			return err
		}
	}
	return nil
}

func (s *SnapShot) UpdateListener(listener *v2.Listener, envoyNode *core.Node, version string) error {
	shapshotCache := cache.NewSnapshot(version, nil, nil, nil, []types.Resource{listener}, nil)
	if err := s.Cache.SetSnapshot(nodes.NodeToString(envoyNode), shapshotCache); err != nil {
		return err
	}

	return nil
}
