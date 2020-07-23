// Code generated by protoc-gen-go. DO NOT EDIT.
// source: envoy/extensions/wasm/v3/wasm.proto

package envoy_extensions_wasm_v3

import (
	fmt "fmt"
	_ "github.com/cncf/udpa/go/udpa/annotations"
	v3 "github.com/envoyproxy/go-control-plane/envoy/config/core/v3"
	_ "github.com/envoyproxy/protoc-gen-validate/validate"
	proto "github.com/golang/protobuf/proto"
	math "math"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

type VmConfig struct {
	VmId                 string              `protobuf:"bytes,1,opt,name=vm_id,json=vmId,proto3" json:"vm_id,omitempty"`
	Runtime              string              `protobuf:"bytes,2,opt,name=runtime,proto3" json:"runtime,omitempty"`
	Code                 *v3.AsyncDataSource `protobuf:"bytes,3,opt,name=code,proto3" json:"code,omitempty"`
	Configuration        string              `protobuf:"bytes,4,opt,name=configuration,proto3" json:"configuration,omitempty"`
	AllowPrecompiled     bool                `protobuf:"varint,5,opt,name=allow_precompiled,json=allowPrecompiled,proto3" json:"allow_precompiled,omitempty"`
	XXX_NoUnkeyedLiteral struct{}            `json:"-"`
	XXX_unrecognized     []byte              `json:"-"`
	XXX_sizecache        int32               `json:"-"`
}

func (m *VmConfig) Reset()         { *m = VmConfig{} }
func (m *VmConfig) String() string { return proto.CompactTextString(m) }
func (*VmConfig) ProtoMessage()    {}
func (*VmConfig) Descriptor() ([]byte, []int) {
	return fileDescriptor_fb6accf9c32c7c2b, []int{0}
}

func (m *VmConfig) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_VmConfig.Unmarshal(m, b)
}
func (m *VmConfig) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_VmConfig.Marshal(b, m, deterministic)
}
func (m *VmConfig) XXX_Merge(src proto.Message) {
	xxx_messageInfo_VmConfig.Merge(m, src)
}
func (m *VmConfig) XXX_Size() int {
	return xxx_messageInfo_VmConfig.Size(m)
}
func (m *VmConfig) XXX_DiscardUnknown() {
	xxx_messageInfo_VmConfig.DiscardUnknown(m)
}

var xxx_messageInfo_VmConfig proto.InternalMessageInfo

func (m *VmConfig) GetVmId() string {
	if m != nil {
		return m.VmId
	}
	return ""
}

func (m *VmConfig) GetRuntime() string {
	if m != nil {
		return m.Runtime
	}
	return ""
}

func (m *VmConfig) GetCode() *v3.AsyncDataSource {
	if m != nil {
		return m.Code
	}
	return nil
}

func (m *VmConfig) GetConfiguration() string {
	if m != nil {
		return m.Configuration
	}
	return ""
}

func (m *VmConfig) GetAllowPrecompiled() bool {
	if m != nil {
		return m.AllowPrecompiled
	}
	return false
}

type PluginConfig struct {
	Name                 string    `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	RootId               string    `protobuf:"bytes,2,opt,name=root_id,json=rootId,proto3" json:"root_id,omitempty"`
	VmConfig             *VmConfig `protobuf:"bytes,3,opt,name=vm_config,json=vmConfig,proto3" json:"vm_config,omitempty"`
	Configuration        string    `protobuf:"bytes,4,opt,name=configuration,proto3" json:"configuration,omitempty"`
	XXX_NoUnkeyedLiteral struct{}  `json:"-"`
	XXX_unrecognized     []byte    `json:"-"`
	XXX_sizecache        int32     `json:"-"`
}

func (m *PluginConfig) Reset()         { *m = PluginConfig{} }
func (m *PluginConfig) String() string { return proto.CompactTextString(m) }
func (*PluginConfig) ProtoMessage()    {}
func (*PluginConfig) Descriptor() ([]byte, []int) {
	return fileDescriptor_fb6accf9c32c7c2b, []int{1}
}

func (m *PluginConfig) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_PluginConfig.Unmarshal(m, b)
}
func (m *PluginConfig) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_PluginConfig.Marshal(b, m, deterministic)
}
func (m *PluginConfig) XXX_Merge(src proto.Message) {
	xxx_messageInfo_PluginConfig.Merge(m, src)
}
func (m *PluginConfig) XXX_Size() int {
	return xxx_messageInfo_PluginConfig.Size(m)
}
func (m *PluginConfig) XXX_DiscardUnknown() {
	xxx_messageInfo_PluginConfig.DiscardUnknown(m)
}

var xxx_messageInfo_PluginConfig proto.InternalMessageInfo

func (m *PluginConfig) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *PluginConfig) GetRootId() string {
	if m != nil {
		return m.RootId
	}
	return ""
}

func (m *PluginConfig) GetVmConfig() *VmConfig {
	if m != nil {
		return m.VmConfig
	}
	return nil
}

func (m *PluginConfig) GetConfiguration() string {
	if m != nil {
		return m.Configuration
	}
	return ""
}

type WasmService struct {
	Config               *PluginConfig `protobuf:"bytes,1,opt,name=config,proto3" json:"config,omitempty"`
	Singleton            bool          `protobuf:"varint,2,opt,name=singleton,proto3" json:"singleton,omitempty"`
	StatPrefix           string        `protobuf:"bytes,3,opt,name=stat_prefix,json=statPrefix,proto3" json:"stat_prefix,omitempty"`
	XXX_NoUnkeyedLiteral struct{}      `json:"-"`
	XXX_unrecognized     []byte        `json:"-"`
	XXX_sizecache        int32         `json:"-"`
}

func (m *WasmService) Reset()         { *m = WasmService{} }
func (m *WasmService) String() string { return proto.CompactTextString(m) }
func (*WasmService) ProtoMessage()    {}
func (*WasmService) Descriptor() ([]byte, []int) {
	return fileDescriptor_fb6accf9c32c7c2b, []int{2}
}

func (m *WasmService) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_WasmService.Unmarshal(m, b)
}
func (m *WasmService) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_WasmService.Marshal(b, m, deterministic)
}
func (m *WasmService) XXX_Merge(src proto.Message) {
	xxx_messageInfo_WasmService.Merge(m, src)
}
func (m *WasmService) XXX_Size() int {
	return xxx_messageInfo_WasmService.Size(m)
}
func (m *WasmService) XXX_DiscardUnknown() {
	xxx_messageInfo_WasmService.DiscardUnknown(m)
}

var xxx_messageInfo_WasmService proto.InternalMessageInfo

func (m *WasmService) GetConfig() *PluginConfig {
	if m != nil {
		return m.Config
	}
	return nil
}

func (m *WasmService) GetSingleton() bool {
	if m != nil {
		return m.Singleton
	}
	return false
}

func (m *WasmService) GetStatPrefix() string {
	if m != nil {
		return m.StatPrefix
	}
	return ""
}

func init() {
	proto.RegisterType((*VmConfig)(nil), "envoy.extensions.wasm.v3.VmConfig")
	proto.RegisterType((*PluginConfig)(nil), "envoy.extensions.wasm.v3.PluginConfig")
	proto.RegisterType((*WasmService)(nil), "envoy.extensions.wasm.v3.WasmService")
}

func init() {
	proto.RegisterFile("envoy/extensions/wasm/v3/wasm.proto", fileDescriptor_fb6accf9c32c7c2b)
}

var fileDescriptor_fb6accf9c32c7c2b = []byte{
	// 466 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x8c, 0x92, 0xd1, 0x6e, 0xd3, 0x3e,
	0x14, 0xc6, 0xe5, 0xff, 0xbf, 0xeb, 0x9a, 0x53, 0x90, 0xc0, 0x5c, 0x2c, 0x9a, 0x80, 0x76, 0xd9,
	0x18, 0x95, 0x90, 0x1c, 0x69, 0x15, 0x17, 0xf4, 0x02, 0xc4, 0xe0, 0x66, 0x77, 0x51, 0x26, 0xc1,
	0x65, 0xe5, 0x25, 0x5e, 0x65, 0x29, 0xb6, 0x23, 0xc7, 0xf1, 0xda, 0x37, 0xe0, 0x19, 0x78, 0x0f,
	0x1e, 0x81, 0x3b, 0x5e, 0x85, 0x77, 0x40, 0xb6, 0x13, 0xda, 0x49, 0xab, 0xc4, 0x55, 0xec, 0xe3,
	0xf3, 0x7d, 0xf9, 0x7d, 0xc7, 0x86, 0x53, 0x26, 0xad, 0xda, 0xa4, 0x6c, 0x6d, 0x98, 0x6c, 0xb8,
	0x92, 0x4d, 0x7a, 0x47, 0x1b, 0x91, 0xda, 0xb9, 0xff, 0x92, 0x5a, 0x2b, 0xa3, 0x70, 0xec, 0x9b,
	0xc8, 0xb6, 0x89, 0xf8, 0x43, 0x3b, 0x3f, 0x9e, 0x04, 0x79, 0xa1, 0xe4, 0x2d, 0x5f, 0xa5, 0x85,
	0xd2, 0xcc, 0x49, 0x6f, 0x68, 0xc3, 0x82, 0xf4, 0xf8, 0xa4, 0x2d, 0x6b, 0x9a, 0x52, 0x29, 0x95,
	0xa1, 0xc6, 0xfb, 0x5b, 0xa6, 0x9d, 0x07, 0x97, 0xab, 0xae, 0xe5, 0xc8, 0xd2, 0x8a, 0x97, 0xd4,
	0xb0, 0xb4, 0x5f, 0x84, 0x83, 0xe4, 0x37, 0x82, 0xd1, 0x17, 0xf1, 0xc9, 0x7b, 0xe3, 0x67, 0x70,
	0x60, 0xc5, 0x92, 0x97, 0x31, 0x9a, 0xa2, 0x59, 0x94, 0x0f, 0xac, 0xb8, 0x2a, 0x71, 0x0c, 0x87,
	0xba, 0x95, 0x86, 0x0b, 0x16, 0xff, 0xe7, 0xcb, 0xfd, 0x16, 0xbf, 0x83, 0x41, 0xa1, 0x4a, 0x16,
	0xff, 0x3f, 0x45, 0xb3, 0xf1, 0xc5, 0x2b, 0x12, 0x12, 0x04, 0x4e, 0xe2, 0x38, 0x89, 0x9d, 0x93,
	0x8f, 0xcd, 0x46, 0x16, 0x9f, 0xa9, 0xa1, 0xd7, 0xaa, 0xd5, 0x05, 0xcb, 0xbd, 0x04, 0x9f, 0xc1,
	0xe3, 0xd0, 0xd7, 0x6a, 0x8f, 0x1c, 0x0f, 0xbc, 0xf5, 0xfd, 0x22, 0x7e, 0x03, 0x4f, 0x69, 0x55,
	0xa9, 0xbb, 0x65, 0xad, 0x59, 0xa1, 0x44, 0xcd, 0x2b, 0x56, 0xc6, 0x07, 0x53, 0x34, 0x1b, 0xe5,
	0x4f, 0xfc, 0x41, 0xb6, 0xad, 0x2f, 0xce, 0xbe, 0xff, 0xfc, 0xf6, 0x72, 0x02, 0x2f, 0xee, 0x51,
	0x84, 0x19, 0x5e, 0x90, 0x3e, 0x62, 0xf2, 0x0b, 0xc1, 0xa3, 0xac, 0x6a, 0x57, 0x5c, 0x76, 0x99,
	0x31, 0x0c, 0x24, 0x15, 0xac, 0x8f, 0xec, 0xd6, 0xf8, 0x08, 0x0e, 0xb5, 0x52, 0xc6, 0x4d, 0x22,
	0x44, 0x1e, 0xba, 0xed, 0x55, 0x89, 0x3f, 0x40, 0x64, 0xc5, 0x32, 0x78, 0x77, 0xb1, 0x13, 0xb2,
	0xef, 0xe2, 0xfe, 0xfe, 0x34, 0x1f, 0xd9, 0x7e, 0xc2, 0xff, 0x94, 0x7b, 0x31, 0x73, 0x51, 0x4e,
	0xe1, 0xe4, 0xc1, 0x28, 0xbb, 0xf4, 0xc9, 0x0f, 0x04, 0xe3, 0xaf, 0xb4, 0x11, 0xd7, 0x4c, 0x5b,
	0x5e, 0x30, 0xfc, 0x1e, 0x86, 0x1d, 0x1d, 0xf2, 0x74, 0xe7, 0xfb, 0xe9, 0x76, 0x7d, 0xf2, 0x4e,
	0x85, 0x9f, 0x43, 0xd4, 0x70, 0xb9, 0xaa, 0x98, 0x51, 0xd2, 0x67, 0x1f, 0xe5, 0xdb, 0x02, 0x9e,
	0xc0, 0xb8, 0x31, 0xd4, 0xb8, 0xeb, 0xb8, 0xe5, 0x6b, 0x3f, 0x80, 0x28, 0x07, 0x57, 0xca, 0x7c,
	0x65, 0xf1, 0xda, 0x81, 0x27, 0x30, 0x7d, 0x10, 0x7c, 0x87, 0xf3, 0xf2, 0x2d, 0x9c, 0x73, 0x15,
	0xd8, 0x6a, 0xad, 0xd6, 0x9b, 0xbd, 0x98, 0x97, 0x91, 0x93, 0x65, 0xee, 0xad, 0x66, 0xe8, 0x66,
	0xe8, 0x1f, 0xed, 0xfc, 0x4f, 0x00, 0x00, 0x00, 0xff, 0xff, 0x2c, 0x34, 0x95, 0x54, 0x52, 0x03,
	0x00, 0x00,
}
