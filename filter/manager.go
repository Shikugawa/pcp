package filter

import (
	"errors"
	"fmt"

	wasm_http "github.com/Shikugawa/pcp/envoy/extensions/filters/http/wasm/v3"
	wasm "github.com/Shikugawa/pcp/envoy/extensions/wasm/v3"
	v3 "github.com/envoyproxy/go-control-plane/envoy/config/core/v3"
	hcm "github.com/envoyproxy/go-control-plane/envoy/config/filter/network/http_connection_manager/v2"
	"github.com/golang/protobuf/ptypes"
)

type WasmFilterStorageDriver struct {
	Storage        *FilterStorage
	defaultRuntime string
}

func NewWasmFilterStorageDriver(runtime string) *WasmFilterStorageDriver {
	return &WasmFilterStorageDriver{
		Storage:        NewFilterStorage(),
		defaultRuntime: runtime,
	}
}

func (f *WasmFilterStorageDriver) EnvoyFilterConfig(filter FilterSpecifier) (*hcm.HttpFilter, error) {
	if f.Storage.ExistFilter(filter) {
		return nil, errors.New(fmt.Sprintln("unregistered wasm filter %s.&%s", filter.FilterType, filter.FilterName))
	}

	plugin_config := &wasm.PluginConfig{
		RootId: "my_root_id",
		VmConfig: &wasm.VmConfig{
			Runtime: f.defaultRuntime,
			Code: &v3.AsyncDataSource{
				Specifier: &v3.AsyncDataSource_Local{
					Local: &v3.DataSource{
						Specifier: &v3.DataSource_Filename{
							Filename: wasmCodePath(filter),
						},
					},
				},
			},
		},
	}
	wasm_config, _ := ptypes.MarshalAny(&wasm_http.Wasm{
		Config: plugin_config,
	})

	return &hcm.HttpFilter{
		Name: "envoy.filters.http.wasm",
		ConfigType: &hcm.HttpFilter_TypedConfig{
			TypedConfig: wasm_config,
		},
	}, nil
}
