package factory

import (
	wasm_http "github.com/Shikugawa/pcp/envoy/extensions/filters/http/wasm/v3"
	wasm "github.com/Shikugawa/pcp/envoy/extensions/wasm/v3"
	v2 "github.com/envoyproxy/go-control-plane/envoy/api/v2"
	core "github.com/envoyproxy/go-control-plane/envoy/api/v2/core"
	listener "github.com/envoyproxy/go-control-plane/envoy/api/v2/listener"
	route "github.com/envoyproxy/go-control-plane/envoy/api/v2/route"
	v3 "github.com/envoyproxy/go-control-plane/envoy/config/core/v3"
	hcm "github.com/envoyproxy/go-control-plane/envoy/config/filter/network/http_connection_manager/v2"

	"github.com/Shikugawa/pcp/filter"
	"github.com/envoyproxy/go-control-plane/pkg/conversion"
	"github.com/golang/protobuf/ptypes"
)

type HttpWasmFilterChainFactory struct {
	DefaultRuntime string
	Filters        []filter.FilterSpecifier
}

func NewHttpWasmFilterChainFactory(defaultRuntime string) *HttpWasmFilterChainFactory {
	return &HttpWasmFilterChainFactory{
		DefaultRuntime: defaultRuntime,
		Filters:        make([]filter.FilterSpecifier, 0),
	}
}

func (h *HttpWasmFilterChainFactory) Create() []*hcm.HttpFilter {
	var httpFilterChainConfig []*hcm.HttpFilter
	for _, filt := range h.Filters {
		pluginConfig := &wasm.PluginConfig{
			RootId: "my_root_id",
			VmConfig: &wasm.VmConfig{
				Runtime: h.DefaultRuntime,
				Code: &v3.AsyncDataSource{
					Specifier: &v3.AsyncDataSource_Local{
						Local: &v3.DataSource{
							Specifier: &v3.DataSource_Filename{
								Filename: filter.WasmCodePath(filt),
							},
						},
					},
				},
			},
		}
		wasmConfig, _ := ptypes.MarshalAny(&wasm_http.Wasm{
			Config: pluginConfig,
		})
		httpFilterChainConfig = append(httpFilterChainConfig, &hcm.HttpFilter{
			Name: "envoy.filters.http.wasm",
			ConfigType: &hcm.HttpFilter_TypedConfig{
				TypedConfig: wasmConfig,
			},
		})
	}

	return httpFilterChainConfig
}

type HttpFilterChainFactory struct {
	wasmFilterChainFactory *HttpWasmFilterChainFactory
}

func NewHttpFilterChainFactory(wasmFilterChainFactory *HttpWasmFilterChainFactory) *HttpFilterChainFactory {
	return &HttpFilterChainFactory{
		wasmFilterChainFactory: wasmFilterChainFactory,
	}
}

func (h *HttpFilterChainFactory) Create() []*hcm.HttpFilter {
	wasmFilterChain := h.wasmFilterChainFactory.Create()
	wasmFilterChain = append(wasmFilterChain, &hcm.HttpFilter{
		Name: "envoy.router",
	})
	return wasmFilterChain
}

type HttpConnectionManagerFactory struct {
	httpFilterChainFactory *HttpFilterChainFactory
}

func NewHttpConnectionManagerFactory(httpFilterChainFactory *HttpFilterChainFactory) *HttpConnectionManagerFactory {
	return &HttpConnectionManagerFactory{
		httpFilterChainFactory: httpFilterChainFactory,
	}
}

func (h *HttpConnectionManagerFactory) Create() *hcm.HttpConnectionManager {
	httpFilterChain := h.httpFilterChainFactory.Create()

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
		HttpFilters: httpFilterChain,
	}
}

type ListenerFilterFactory struct {
	httpConnManagerFactory *HttpConnectionManagerFactory
}

func NewListenerFilterFactory(httpConnManagerFactory *HttpConnectionManagerFactory) *ListenerFilterFactory {
	return &ListenerFilterFactory{
		httpConnManagerFactory: httpConnManagerFactory,
	}
}

func (l *ListenerFilterFactory) Create() *listener.Filter {
	connManager := l.httpConnManagerFactory.Create()
	config, _ := conversion.MessageToStruct(connManager)

	return &listener.Filter{
		Name: "envoy.http_connection_manager",
		ConfigType: &listener.Filter_Config{
			Config: config,
		},
	}
}

type ListenerFactory struct {
	listenerFilterFactory *ListenerFilterFactory
}

func NewListenerFactory(listenerFilterFactory *ListenerFilterFactory) *ListenerFactory {
	return &ListenerFactory{
		listenerFilterFactory: listenerFilterFactory,
	}
}

func (l *ListenerFactory) Create() *v2.Listener {
	listenerFilterChains := []*listener.FilterChain{
		{
			Filters: []*listener.Filter{l.listenerFilterFactory.Create()},
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
		FilterChains: listenerFilterChains,
	}
}
