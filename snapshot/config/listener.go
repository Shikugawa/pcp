package config

import (
	v2 "github.com/envoyproxy/go-control-plane/envoy/api/v2"
	core "github.com/envoyproxy/go-control-plane/envoy/api/v2/core"
	listener "github.com/envoyproxy/go-control-plane/envoy/api/v2/listener"
	route "github.com/envoyproxy/go-control-plane/envoy/api/v2/route"
	hcm "github.com/envoyproxy/go-control-plane/envoy/config/filter/network/http_connection_manager/v2"
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
			RouteConfig: &v2.RouteConfiguration{
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

func GetListener(h *hcm.HttpConnectionManager) (*v2.Listener, error) {
	if h == nil {
		h = HCM(nil)
	}
	config, _ := conversion.MessageToStruct(h)

	filterChains := []*listener.FilterChain{
		{
			Filters: []*listener.Filter{{
				Name: "envoy.http_connection_manager",
				ConfigType: &listener.Filter_Config{
					Config: config,
				},
			}},
		},
	}
	return &v2.Listener{
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
