// Copyright 2024 Gravitational, Inc
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.2
// 	protoc        (unknown)
// source: teleport/access_graph/v1/private_key.proto

package accessgraphv1

import (
	v1 "github.com/gravitational/teleport/api/gen/proto/go/teleport/header/v1"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

// PublicKeyMode is the mode of the public key.
// The public key can be derived from the private key, stored in a separate file, or the private key was password protected
// and we could not extract the public key from it or from the file.
type PublicKeyMode int32

const (
	// PUBLIC_KEY_MODE_UNSPECIFIED is an invalid state.
	PublicKeyMode_PUBLIC_KEY_MODE_UNSPECIFIED PublicKeyMode = 0
	// PUBLIC_KEY_MODE_DERIVED is the state where the public key is derived from the private key.
	PublicKeyMode_PUBLIC_KEY_MODE_DERIVED PublicKeyMode = 1
	// PUBLIC_KEY_MODE_PUB_FILE is a state where the public key is stored in a separate file from the private key.
	// The private key is password protected and we could not extract the public key from it.
	// This mode is used when the private key is password protected and there is a <key>.pub file next to the private key
	// that contains the public key.
	PublicKeyMode_PUBLIC_KEY_MODE_PUB_FILE PublicKeyMode = 2
	// PUBLIC_KEY_MODE_PROTECTED is a state where the private key is password protected and we could not extract the public key from it
	// or from the .pub file.
	PublicKeyMode_PUBLIC_KEY_MODE_PROTECTED PublicKeyMode = 3
)

// Enum value maps for PublicKeyMode.
var (
	PublicKeyMode_name = map[int32]string{
		0: "PUBLIC_KEY_MODE_UNSPECIFIED",
		1: "PUBLIC_KEY_MODE_DERIVED",
		2: "PUBLIC_KEY_MODE_PUB_FILE",
		3: "PUBLIC_KEY_MODE_PROTECTED",
	}
	PublicKeyMode_value = map[string]int32{
		"PUBLIC_KEY_MODE_UNSPECIFIED": 0,
		"PUBLIC_KEY_MODE_DERIVED":     1,
		"PUBLIC_KEY_MODE_PUB_FILE":    2,
		"PUBLIC_KEY_MODE_PROTECTED":   3,
	}
)

func (x PublicKeyMode) Enum() *PublicKeyMode {
	p := new(PublicKeyMode)
	*p = x
	return p
}

func (x PublicKeyMode) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (PublicKeyMode) Descriptor() protoreflect.EnumDescriptor {
	return file_teleport_access_graph_v1_private_key_proto_enumTypes[0].Descriptor()
}

func (PublicKeyMode) Type() protoreflect.EnumType {
	return &file_teleport_access_graph_v1_private_key_proto_enumTypes[0]
}

func (x PublicKeyMode) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use PublicKeyMode.Descriptor instead.
func (PublicKeyMode) EnumDescriptor() ([]byte, []int) {
	return file_teleport_access_graph_v1_private_key_proto_rawDescGZIP(), []int{0}
}

// The `PrivateKey` message represents a private key entry for a specific local user.
// It serves as a reference to a private key located on a user's laptop. Note that it *NEVER* contains the private key itself.
// Instead, it stores metadata related to the key, including the fingerprint of the public key, the device trust identifier, and the public key mode.
// The Teleport Access Graph uses this metadata to assess whether a particular private key is authorized to access a user on the node without using Teleport.
type PrivateKey struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	// metadata is the PrivateKey's metadata.
	Metadata *v1.Metadata `protobuf:"bytes,1,opt,name=metadata,proto3" json:"metadata,omitempty"`
	// kind is a resource kind.
	Kind string `protobuf:"bytes,2,opt,name=kind,proto3" json:"kind,omitempty"`
	// sub_kind is an optional resource sub kind, used in some resources.
	SubKind string `protobuf:"bytes,3,opt,name=sub_kind,json=subKind,proto3" json:"sub_kind,omitempty"`
	// version is version.
	Version string `protobuf:"bytes,4,opt,name=version,proto3" json:"version,omitempty"`
	// Spec is a PrivateKey specification.
	Spec          *PrivateKeySpec `protobuf:"bytes,5,opt,name=spec,proto3" json:"spec,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *PrivateKey) Reset() {
	*x = PrivateKey{}
	mi := &file_teleport_access_graph_v1_private_key_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *PrivateKey) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PrivateKey) ProtoMessage() {}

func (x *PrivateKey) ProtoReflect() protoreflect.Message {
	mi := &file_teleport_access_graph_v1_private_key_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PrivateKey.ProtoReflect.Descriptor instead.
func (*PrivateKey) Descriptor() ([]byte, []int) {
	return file_teleport_access_graph_v1_private_key_proto_rawDescGZIP(), []int{0}
}

func (x *PrivateKey) GetMetadata() *v1.Metadata {
	if x != nil {
		return x.Metadata
	}
	return nil
}

func (x *PrivateKey) GetKind() string {
	if x != nil {
		return x.Kind
	}
	return ""
}

func (x *PrivateKey) GetSubKind() string {
	if x != nil {
		return x.SubKind
	}
	return ""
}

func (x *PrivateKey) GetVersion() string {
	if x != nil {
		return x.Version
	}
	return ""
}

func (x *PrivateKey) GetSpec() *PrivateKeySpec {
	if x != nil {
		return x.Spec
	}
	return nil
}

// PrivateKeySpec is the private key spec.
type PrivateKeySpec struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	// device_id is the device trust identifier of the device that owns the key.
	DeviceId string `protobuf:"bytes,1,opt,name=device_id,json=deviceId,proto3" json:"device_id,omitempty"`
	// public_key_fingerprint is the SHA256 of the SSH public key corresponding to
	// the private key.
	PublicKeyFingerprint string `protobuf:"bytes,2,opt,name=public_key_fingerprint,json=publicKeyFingerprint,proto3" json:"public_key_fingerprint,omitempty"`
	// public_key_mode is the public key mode.
	PublicKeyMode PublicKeyMode `protobuf:"varint,3,opt,name=public_key_mode,json=publicKeyMode,proto3,enum=teleport.access_graph.v1.PublicKeyMode" json:"public_key_mode,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *PrivateKeySpec) Reset() {
	*x = PrivateKeySpec{}
	mi := &file_teleport_access_graph_v1_private_key_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *PrivateKeySpec) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PrivateKeySpec) ProtoMessage() {}

func (x *PrivateKeySpec) ProtoReflect() protoreflect.Message {
	mi := &file_teleport_access_graph_v1_private_key_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PrivateKeySpec.ProtoReflect.Descriptor instead.
func (*PrivateKeySpec) Descriptor() ([]byte, []int) {
	return file_teleport_access_graph_v1_private_key_proto_rawDescGZIP(), []int{1}
}

func (x *PrivateKeySpec) GetDeviceId() string {
	if x != nil {
		return x.DeviceId
	}
	return ""
}

func (x *PrivateKeySpec) GetPublicKeyFingerprint() string {
	if x != nil {
		return x.PublicKeyFingerprint
	}
	return ""
}

func (x *PrivateKeySpec) GetPublicKeyMode() PublicKeyMode {
	if x != nil {
		return x.PublicKeyMode
	}
	return PublicKeyMode_PUBLIC_KEY_MODE_UNSPECIFIED
}

var File_teleport_access_graph_v1_private_key_proto protoreflect.FileDescriptor

var file_teleport_access_graph_v1_private_key_proto_rawDesc = []byte{
	0x0a, 0x2a, 0x74, 0x65, 0x6c, 0x65, 0x70, 0x6f, 0x72, 0x74, 0x2f, 0x61, 0x63, 0x63, 0x65, 0x73,
	0x73, 0x5f, 0x67, 0x72, 0x61, 0x70, 0x68, 0x2f, 0x76, 0x31, 0x2f, 0x70, 0x72, 0x69, 0x76, 0x61,
	0x74, 0x65, 0x5f, 0x6b, 0x65, 0x79, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x18, 0x74, 0x65,
	0x6c, 0x65, 0x70, 0x6f, 0x72, 0x74, 0x2e, 0x61, 0x63, 0x63, 0x65, 0x73, 0x73, 0x5f, 0x67, 0x72,
	0x61, 0x70, 0x68, 0x2e, 0x76, 0x31, 0x1a, 0x21, 0x74, 0x65, 0x6c, 0x65, 0x70, 0x6f, 0x72, 0x74,
	0x2f, 0x68, 0x65, 0x61, 0x64, 0x65, 0x72, 0x2f, 0x76, 0x31, 0x2f, 0x6d, 0x65, 0x74, 0x61, 0x64,
	0x61, 0x74, 0x61, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0xcd, 0x01, 0x0a, 0x0a, 0x50, 0x72,
	0x69, 0x76, 0x61, 0x74, 0x65, 0x4b, 0x65, 0x79, 0x12, 0x38, 0x0a, 0x08, 0x6d, 0x65, 0x74, 0x61,
	0x64, 0x61, 0x74, 0x61, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1c, 0x2e, 0x74, 0x65, 0x6c,
	0x65, 0x70, 0x6f, 0x72, 0x74, 0x2e, 0x68, 0x65, 0x61, 0x64, 0x65, 0x72, 0x2e, 0x76, 0x31, 0x2e,
	0x4d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0x52, 0x08, 0x6d, 0x65, 0x74, 0x61, 0x64, 0x61,
	0x74, 0x61, 0x12, 0x12, 0x0a, 0x04, 0x6b, 0x69, 0x6e, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x04, 0x6b, 0x69, 0x6e, 0x64, 0x12, 0x19, 0x0a, 0x08, 0x73, 0x75, 0x62, 0x5f, 0x6b, 0x69,
	0x6e, 0x64, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x73, 0x75, 0x62, 0x4b, 0x69, 0x6e,
	0x64, 0x12, 0x18, 0x0a, 0x07, 0x76, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x18, 0x04, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x07, 0x76, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x12, 0x3c, 0x0a, 0x04, 0x73,
	0x70, 0x65, 0x63, 0x18, 0x05, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x28, 0x2e, 0x74, 0x65, 0x6c, 0x65,
	0x70, 0x6f, 0x72, 0x74, 0x2e, 0x61, 0x63, 0x63, 0x65, 0x73, 0x73, 0x5f, 0x67, 0x72, 0x61, 0x70,
	0x68, 0x2e, 0x76, 0x31, 0x2e, 0x50, 0x72, 0x69, 0x76, 0x61, 0x74, 0x65, 0x4b, 0x65, 0x79, 0x53,
	0x70, 0x65, 0x63, 0x52, 0x04, 0x73, 0x70, 0x65, 0x63, 0x22, 0xb4, 0x01, 0x0a, 0x0e, 0x50, 0x72,
	0x69, 0x76, 0x61, 0x74, 0x65, 0x4b, 0x65, 0x79, 0x53, 0x70, 0x65, 0x63, 0x12, 0x1b, 0x0a, 0x09,
	0x64, 0x65, 0x76, 0x69, 0x63, 0x65, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x08, 0x64, 0x65, 0x76, 0x69, 0x63, 0x65, 0x49, 0x64, 0x12, 0x34, 0x0a, 0x16, 0x70, 0x75, 0x62,
	0x6c, 0x69, 0x63, 0x5f, 0x6b, 0x65, 0x79, 0x5f, 0x66, 0x69, 0x6e, 0x67, 0x65, 0x72, 0x70, 0x72,
	0x69, 0x6e, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x14, 0x70, 0x75, 0x62, 0x6c, 0x69,
	0x63, 0x4b, 0x65, 0x79, 0x46, 0x69, 0x6e, 0x67, 0x65, 0x72, 0x70, 0x72, 0x69, 0x6e, 0x74, 0x12,
	0x4f, 0x0a, 0x0f, 0x70, 0x75, 0x62, 0x6c, 0x69, 0x63, 0x5f, 0x6b, 0x65, 0x79, 0x5f, 0x6d, 0x6f,
	0x64, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x27, 0x2e, 0x74, 0x65, 0x6c, 0x65, 0x70,
	0x6f, 0x72, 0x74, 0x2e, 0x61, 0x63, 0x63, 0x65, 0x73, 0x73, 0x5f, 0x67, 0x72, 0x61, 0x70, 0x68,
	0x2e, 0x76, 0x31, 0x2e, 0x50, 0x75, 0x62, 0x6c, 0x69, 0x63, 0x4b, 0x65, 0x79, 0x4d, 0x6f, 0x64,
	0x65, 0x52, 0x0d, 0x70, 0x75, 0x62, 0x6c, 0x69, 0x63, 0x4b, 0x65, 0x79, 0x4d, 0x6f, 0x64, 0x65,
	0x2a, 0x8a, 0x01, 0x0a, 0x0d, 0x50, 0x75, 0x62, 0x6c, 0x69, 0x63, 0x4b, 0x65, 0x79, 0x4d, 0x6f,
	0x64, 0x65, 0x12, 0x1f, 0x0a, 0x1b, 0x50, 0x55, 0x42, 0x4c, 0x49, 0x43, 0x5f, 0x4b, 0x45, 0x59,
	0x5f, 0x4d, 0x4f, 0x44, 0x45, 0x5f, 0x55, 0x4e, 0x53, 0x50, 0x45, 0x43, 0x49, 0x46, 0x49, 0x45,
	0x44, 0x10, 0x00, 0x12, 0x1b, 0x0a, 0x17, 0x50, 0x55, 0x42, 0x4c, 0x49, 0x43, 0x5f, 0x4b, 0x45,
	0x59, 0x5f, 0x4d, 0x4f, 0x44, 0x45, 0x5f, 0x44, 0x45, 0x52, 0x49, 0x56, 0x45, 0x44, 0x10, 0x01,
	0x12, 0x1c, 0x0a, 0x18, 0x50, 0x55, 0x42, 0x4c, 0x49, 0x43, 0x5f, 0x4b, 0x45, 0x59, 0x5f, 0x4d,
	0x4f, 0x44, 0x45, 0x5f, 0x50, 0x55, 0x42, 0x5f, 0x46, 0x49, 0x4c, 0x45, 0x10, 0x02, 0x12, 0x1d,
	0x0a, 0x19, 0x50, 0x55, 0x42, 0x4c, 0x49, 0x43, 0x5f, 0x4b, 0x45, 0x59, 0x5f, 0x4d, 0x4f, 0x44,
	0x45, 0x5f, 0x50, 0x52, 0x4f, 0x54, 0x45, 0x43, 0x54, 0x45, 0x44, 0x10, 0x03, 0x42, 0x5a, 0x5a,
	0x58, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x67, 0x72, 0x61, 0x76,
	0x69, 0x74, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x61, 0x6c, 0x2f, 0x74, 0x65, 0x6c, 0x65, 0x70, 0x6f,
	0x72, 0x74, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x67, 0x65, 0x6e, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x2f, 0x67, 0x6f, 0x2f, 0x74, 0x65, 0x6c, 0x65, 0x70, 0x6f, 0x72, 0x74, 0x2f, 0x61, 0x63, 0x63,
	0x65, 0x73, 0x73, 0x67, 0x72, 0x61, 0x70, 0x68, 0x2f, 0x76, 0x31, 0x3b, 0x61, 0x63, 0x63, 0x65,
	0x73, 0x73, 0x67, 0x72, 0x61, 0x70, 0x68, 0x76, 0x31, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x33,
}

var (
	file_teleport_access_graph_v1_private_key_proto_rawDescOnce sync.Once
	file_teleport_access_graph_v1_private_key_proto_rawDescData = file_teleport_access_graph_v1_private_key_proto_rawDesc
)

func file_teleport_access_graph_v1_private_key_proto_rawDescGZIP() []byte {
	file_teleport_access_graph_v1_private_key_proto_rawDescOnce.Do(func() {
		file_teleport_access_graph_v1_private_key_proto_rawDescData = protoimpl.X.CompressGZIP(file_teleport_access_graph_v1_private_key_proto_rawDescData)
	})
	return file_teleport_access_graph_v1_private_key_proto_rawDescData
}

var file_teleport_access_graph_v1_private_key_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_teleport_access_graph_v1_private_key_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_teleport_access_graph_v1_private_key_proto_goTypes = []any{
	(PublicKeyMode)(0),     // 0: teleport.access_graph.v1.PublicKeyMode
	(*PrivateKey)(nil),     // 1: teleport.access_graph.v1.PrivateKey
	(*PrivateKeySpec)(nil), // 2: teleport.access_graph.v1.PrivateKeySpec
	(*v1.Metadata)(nil),    // 3: teleport.header.v1.Metadata
}
var file_teleport_access_graph_v1_private_key_proto_depIdxs = []int32{
	3, // 0: teleport.access_graph.v1.PrivateKey.metadata:type_name -> teleport.header.v1.Metadata
	2, // 1: teleport.access_graph.v1.PrivateKey.spec:type_name -> teleport.access_graph.v1.PrivateKeySpec
	0, // 2: teleport.access_graph.v1.PrivateKeySpec.public_key_mode:type_name -> teleport.access_graph.v1.PublicKeyMode
	3, // [3:3] is the sub-list for method output_type
	3, // [3:3] is the sub-list for method input_type
	3, // [3:3] is the sub-list for extension type_name
	3, // [3:3] is the sub-list for extension extendee
	0, // [0:3] is the sub-list for field type_name
}

func init() { file_teleport_access_graph_v1_private_key_proto_init() }
func file_teleport_access_graph_v1_private_key_proto_init() {
	if File_teleport_access_graph_v1_private_key_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_teleport_access_graph_v1_private_key_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_teleport_access_graph_v1_private_key_proto_goTypes,
		DependencyIndexes: file_teleport_access_graph_v1_private_key_proto_depIdxs,
		EnumInfos:         file_teleport_access_graph_v1_private_key_proto_enumTypes,
		MessageInfos:      file_teleport_access_graph_v1_private_key_proto_msgTypes,
	}.Build()
	File_teleport_access_graph_v1_private_key_proto = out.File
	file_teleport_access_graph_v1_private_key_proto_rawDesc = nil
	file_teleport_access_graph_v1_private_key_proto_goTypes = nil
	file_teleport_access_graph_v1_private_key_proto_depIdxs = nil
}
