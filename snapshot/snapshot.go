package snapshot

import (
	"github.com/Shikugawa/pcp/snapshot/config"
	core "github.com/envoyproxy/go-control-plane/envoy/config/core/v3"
	listener "github.com/envoyproxy/go-control-plane/envoy/config/listener/v3"
	"github.com/envoyproxy/go-control-plane/pkg/cache/types"
	cache "github.com/envoyproxy/go-control-plane/pkg/cache/v3"
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

func (s *SnapShot) DefaultCache(nodes []*core.Node) error {
	defaultListener, err := config.GetListener(nil)
	if err != nil {
		return err
	}

	for _, node := range nodes {
		snapshotCache := cache.NewSnapshot("1", nil, nil, nil, []types.Resource{defaultListener}, nil)
		if err := s.Cache.SetSnapshot(node.Cluster+"/"+node.Id, snapshotCache); err != nil {
			return err
		}
	}
	return nil
}

func (s *SnapShot) UpdateListener(listener *listener.Listener, nodes []*core.Node, version string) error {
	for _, node := range nodes {
		shapshotCache := cache.NewSnapshot(version, nil, nil, nil, []types.Resource{listener}, nil)
		if err := s.Cache.SetSnapshot(node.Cluster+"/"+node.Id, shapshotCache); err != nil {
			return err
		}
	}

	return nil
}
