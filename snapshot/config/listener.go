package config

import (
	"github.com/Shikugawa/pcp/util"

	core "github.com/envoyproxy/go-control-plane/envoy/config/core/v3"
	listener "github.com/envoyproxy/go-control-plane/envoy/config/listener/v3"
	route "github.com/envoyproxy/go-control-plane/envoy/config/route/v3"
	hcm "github.com/envoyproxy/go-control-plane/envoy/extensions/filters/network/http_connection_manager/v3"
	"github.com/envoyproxy/go-control-plane/pkg/conversion"
)

func HCM(filters []*hcm.HttpFilter) *hcm.HttpConnectionManager {
	routerExists := func() bool {
		for _, filter := range filters {
			if filter.Name == "envoy.router" {
				return true
			}
		}
		return false
	}

	if !routerExists() {
		filters = append(filters, &hcm.HttpFilter{
			Name: "envoy.router",
		})
	}

	return &hcm.HttpConnectionManager{
		StatPrefix: "ingress_http",
		RouteSpecifier: &hcm.HttpConnectionManager_RouteConfig{
			RouteConfig: &route.RouteConfiguration{
				Name: "local_route",
				VirtualHosts: []*route.VirtualHost{
					&route.VirtualHost{
						Name:    "local_service",
						Domains: []string{"*"},
						Routes: []*route.Route{
							&route.Route{
								Match: &route.RouteMatch{
									PathSpecifier: &route.RouteMatch_Prefix{
										Prefix: "/",
									},
								},
								Action: &route.Route_Route{
									Route: &route.RouteAction{
										ClusterSpecifier: &route.RouteAction_Cluster{
											Cluster: "service",
										},
									},
								},
							},
						},
					},
				},
			},
		},
		HttpFilters: filters,
	}
}

func GetListener(h *hcm.HttpConnectionManager) (*listener.Listener, error) {
	if h == nil {
		h = HCM(nil)
	}
	config, _ := conversion.MessageToStruct(h)

	filterChains := []*listener.FilterChain{
		{
			Filters: []*listener.Filter{{
				Name: "envoy.http_connection_manager",
				ConfigType: &listener.Filter_TypedConfig{
					TypedConfig: util.MarshalAny(
						"type.googleapis.com/envoy.extensions.filters.network.http_connection_manager.v3.HttpConnectionManagerr", config),
				},
			}},
		},
	}
	return &listener.Listener{
		Name: "default_listener",
		Address: &core.Address{
			Address: &core.Address_SocketAddress{
				SocketAddress: &core.SocketAddress{
					Address: "0.0.0.0",
					PortSpecifier: &core.SocketAddress_PortValue{
						PortValue: 5000,
					},
				},
			},
		},
		FilterChains: filterChains,
	}, nil
}
