package snapshot

import (
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
		Cache:   cache.NewSnapshotCache(false, Hash{}, nil),
	}
}

func (s *SnapShot) DefaultCache(nodes []*core.Node, defaultListener *v2.Listener) error {
	for _, node := range nodes {
		snapshotCache := cache.NewSnapshot("1", nil, nil, nil, []types.Resource{defaultListener}, nil)
		if err := s.Cache.SetSnapshot(node.Cluster+"/"+node.Id, snapshotCache); err != nil {
			return err
		}
	}
	return nil
}

func (s *SnapShot) UpdateListener(listener *v2.Listener, nodes []*core.Node, version string) error {
	for _, node := range nodes {
		shapshotCache := cache.NewSnapshot(version, nil, nil, nil, []types.Resource{listener}, nil)
		if err := s.Cache.SetSnapshot(node.Cluster+"/"+node.Id, shapshotCache); err != nil {
			return err
		}
	}

	return nil
}
