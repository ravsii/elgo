// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.30.0
// 	protoc        v3.21.12
// source: grpc/pb/pool.proto

package pb

import (
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

type Empty struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *Empty) Reset() {
	*x = Empty{}
	if protoimpl.UnsafeEnabled {
		mi := &file_grpc_pb_pool_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Empty) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Empty) ProtoMessage() {}

func (x *Empty) ProtoReflect() protoreflect.Message {
	mi := &file_grpc_pb_pool_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Empty.ProtoReflect.Descriptor instead.
func (*Empty) Descriptor() ([]byte, []int) {
	return file_grpc_pb_pool_proto_rawDescGZIP(), []int{0}
}

type SizeResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Size int32 `protobuf:"varint,1,opt,name=size,proto3" json:"size,omitempty"`
}

func (x *SizeResponse) Reset() {
	*x = SizeResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_grpc_pb_pool_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SizeResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SizeResponse) ProtoMessage() {}

func (x *SizeResponse) ProtoReflect() protoreflect.Message {
	mi := &file_grpc_pb_pool_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SizeResponse.ProtoReflect.Descriptor instead.
func (*SizeResponse) Descriptor() ([]byte, []int) {
	return file_grpc_pb_pool_proto_rawDescGZIP(), []int{1}
}

func (x *SizeResponse) GetSize() int32 {
	if x != nil {
		return x.Size
	}
	return 0
}

type Player struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id  string  `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Elo float64 `protobuf:"fixed64,2,opt,name=elo,proto3" json:"elo,omitempty"`
}

func (x *Player) Reset() {
	*x = Player{}
	if protoimpl.UnsafeEnabled {
		mi := &file_grpc_pb_pool_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Player) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Player) ProtoMessage() {}

func (x *Player) ProtoReflect() protoreflect.Message {
	mi := &file_grpc_pb_pool_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Player.ProtoReflect.Descriptor instead.
func (*Player) Descriptor() ([]byte, []int) {
	return file_grpc_pb_pool_proto_rawDescGZIP(), []int{2}
}

func (x *Player) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *Player) GetElo() float64 {
	if x != nil {
		return x.Elo
	}
	return 0
}

type PlayerMatch struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	P1 *Player `protobuf:"bytes,1,opt,name=p1,proto3" json:"p1,omitempty"`
	P2 *Player `protobuf:"bytes,2,opt,name=p2,proto3" json:"p2,omitempty"`
}

func (x *PlayerMatch) Reset() {
	*x = PlayerMatch{}
	if protoimpl.UnsafeEnabled {
		mi := &file_grpc_pb_pool_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PlayerMatch) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PlayerMatch) ProtoMessage() {}

func (x *PlayerMatch) ProtoReflect() protoreflect.Message {
	mi := &file_grpc_pb_pool_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PlayerMatch.ProtoReflect.Descriptor instead.
func (*PlayerMatch) Descriptor() ([]byte, []int) {
	return file_grpc_pb_pool_proto_rawDescGZIP(), []int{3}
}

func (x *PlayerMatch) GetP1() *Player {
	if x != nil {
		return x.P1
	}
	return nil
}

func (x *PlayerMatch) GetP2() *Player {
	if x != nil {
		return x.P2
	}
	return nil
}

var File_grpc_pb_pool_proto protoreflect.FileDescriptor

var file_grpc_pb_pool_proto_rawDesc = []byte{
	0x0a, 0x12, 0x67, 0x72, 0x70, 0x63, 0x2f, 0x70, 0x62, 0x2f, 0x70, 0x6f, 0x6f, 0x6c, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x12, 0x04, 0x70, 0x6f, 0x6f, 0x6c, 0x22, 0x07, 0x0a, 0x05, 0x45, 0x6d,
	0x70, 0x74, 0x79, 0x22, 0x22, 0x0a, 0x0c, 0x53, 0x69, 0x7a, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f,
	0x6e, 0x73, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x73, 0x69, 0x7a, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x05, 0x52, 0x04, 0x73, 0x69, 0x7a, 0x65, 0x22, 0x2a, 0x0a, 0x06, 0x50, 0x6c, 0x61, 0x79, 0x65,
	0x72, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x69,
	0x64, 0x12, 0x10, 0x0a, 0x03, 0x65, 0x6c, 0x6f, 0x18, 0x02, 0x20, 0x01, 0x28, 0x01, 0x52, 0x03,
	0x65, 0x6c, 0x6f, 0x22, 0x49, 0x0a, 0x0b, 0x50, 0x6c, 0x61, 0x79, 0x65, 0x72, 0x4d, 0x61, 0x74,
	0x63, 0x68, 0x12, 0x1c, 0x0a, 0x02, 0x70, 0x31, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0c,
	0x2e, 0x70, 0x6f, 0x6f, 0x6c, 0x2e, 0x50, 0x6c, 0x61, 0x79, 0x65, 0x72, 0x52, 0x02, 0x70, 0x31,
	0x12, 0x1c, 0x0a, 0x02, 0x70, 0x32, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0c, 0x2e, 0x70,
	0x6f, 0x6f, 0x6c, 0x2e, 0x50, 0x6c, 0x61, 0x79, 0x65, 0x72, 0x52, 0x02, 0x70, 0x32, 0x32, 0xa1,
	0x01, 0x0a, 0x04, 0x50, 0x6f, 0x6f, 0x6c, 0x12, 0x20, 0x0a, 0x03, 0x41, 0x64, 0x64, 0x12, 0x0c,
	0x2e, 0x70, 0x6f, 0x6f, 0x6c, 0x2e, 0x50, 0x6c, 0x61, 0x79, 0x65, 0x72, 0x1a, 0x0b, 0x2e, 0x70,
	0x6f, 0x6f, 0x6c, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x12, 0x29, 0x0a, 0x05, 0x4d, 0x61, 0x74,
	0x63, 0x68, 0x12, 0x0b, 0x2e, 0x70, 0x6f, 0x6f, 0x6c, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x1a,
	0x11, 0x2e, 0x70, 0x6f, 0x6f, 0x6c, 0x2e, 0x50, 0x6c, 0x61, 0x79, 0x65, 0x72, 0x4d, 0x61, 0x74,
	0x63, 0x68, 0x30, 0x01, 0x12, 0x23, 0x0a, 0x06, 0x52, 0x65, 0x6d, 0x6f, 0x76, 0x65, 0x12, 0x0c,
	0x2e, 0x70, 0x6f, 0x6f, 0x6c, 0x2e, 0x50, 0x6c, 0x61, 0x79, 0x65, 0x72, 0x1a, 0x0b, 0x2e, 0x70,
	0x6f, 0x6f, 0x6c, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x12, 0x27, 0x0a, 0x04, 0x53, 0x69, 0x7a,
	0x65, 0x12, 0x0b, 0x2e, 0x70, 0x6f, 0x6f, 0x6c, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x1a, 0x12,
	0x2e, 0x70, 0x6f, 0x6f, 0x6c, 0x2e, 0x53, 0x69, 0x7a, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e,
	0x73, 0x65, 0x42, 0x20, 0x5a, 0x1e, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d,
	0x2f, 0x72, 0x61, 0x76, 0x73, 0x69, 0x69, 0x2f, 0x65, 0x6c, 0x67, 0x6f, 0x2f, 0x67, 0x72, 0x70,
	0x63, 0x2f, 0x70, 0x62, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_grpc_pb_pool_proto_rawDescOnce sync.Once
	file_grpc_pb_pool_proto_rawDescData = file_grpc_pb_pool_proto_rawDesc
)

func file_grpc_pb_pool_proto_rawDescGZIP() []byte {
	file_grpc_pb_pool_proto_rawDescOnce.Do(func() {
		file_grpc_pb_pool_proto_rawDescData = protoimpl.X.CompressGZIP(file_grpc_pb_pool_proto_rawDescData)
	})
	return file_grpc_pb_pool_proto_rawDescData
}

var file_grpc_pb_pool_proto_msgTypes = make([]protoimpl.MessageInfo, 4)
var file_grpc_pb_pool_proto_goTypes = []interface{}{
	(*Empty)(nil),        // 0: pool.Empty
	(*SizeResponse)(nil), // 1: pool.SizeResponse
	(*Player)(nil),       // 2: pool.Player
	(*PlayerMatch)(nil),  // 3: pool.PlayerMatch
}
var file_grpc_pb_pool_proto_depIdxs = []int32{
	2, // 0: pool.PlayerMatch.p1:type_name -> pool.Player
	2, // 1: pool.PlayerMatch.p2:type_name -> pool.Player
	2, // 2: pool.Pool.Add:input_type -> pool.Player
	0, // 3: pool.Pool.Match:input_type -> pool.Empty
	2, // 4: pool.Pool.Remove:input_type -> pool.Player
	0, // 5: pool.Pool.Size:input_type -> pool.Empty
	0, // 6: pool.Pool.Add:output_type -> pool.Empty
	3, // 7: pool.Pool.Match:output_type -> pool.PlayerMatch
	0, // 8: pool.Pool.Remove:output_type -> pool.Empty
	1, // 9: pool.Pool.Size:output_type -> pool.SizeResponse
	6, // [6:10] is the sub-list for method output_type
	2, // [2:6] is the sub-list for method input_type
	2, // [2:2] is the sub-list for extension type_name
	2, // [2:2] is the sub-list for extension extendee
	0, // [0:2] is the sub-list for field type_name
}

func init() { file_grpc_pb_pool_proto_init() }
func file_grpc_pb_pool_proto_init() {
	if File_grpc_pb_pool_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_grpc_pb_pool_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Empty); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_grpc_pb_pool_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SizeResponse); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_grpc_pb_pool_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Player); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_grpc_pb_pool_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PlayerMatch); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_grpc_pb_pool_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   4,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_grpc_pb_pool_proto_goTypes,
		DependencyIndexes: file_grpc_pb_pool_proto_depIdxs,
		MessageInfos:      file_grpc_pb_pool_proto_msgTypes,
	}.Build()
	File_grpc_pb_pool_proto = out.File
	file_grpc_pb_pool_proto_rawDesc = nil
	file_grpc_pb_pool_proto_goTypes = nil
	file_grpc_pb_pool_proto_depIdxs = nil
}