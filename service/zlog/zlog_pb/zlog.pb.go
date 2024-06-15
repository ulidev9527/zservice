// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.32.0
// 	protoc        v4.25.1
// source: zlog.proto

package zlog_pb

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

type Default_REQ struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *Default_REQ) Reset() {
	*x = Default_REQ{}
	if protoimpl.UnsafeEnabled {
		mi := &file_zlog_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Default_REQ) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Default_REQ) ProtoMessage() {}

func (x *Default_REQ) ProtoReflect() protoreflect.Message {
	mi := &file_zlog_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Default_REQ.ProtoReflect.Descriptor instead.
func (*Default_REQ) Descriptor() ([]byte, []int) {
	return file_zlog_proto_rawDescGZIP(), []int{0}
}

type Default_RES struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Code uint32 `protobuf:"varint,1,opt,name=code,proto3" json:"code,omitempty"` // 服务状态码 约定: 1 表示成功, 其它数字根据业务进行返回
}

func (x *Default_RES) Reset() {
	*x = Default_RES{}
	if protoimpl.UnsafeEnabled {
		mi := &file_zlog_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Default_RES) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Default_RES) ProtoMessage() {}

func (x *Default_RES) ProtoReflect() protoreflect.Message {
	mi := &file_zlog_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Default_RES.ProtoReflect.Descriptor instead.
func (*Default_RES) Descriptor() ([]byte, []int) {
	return file_zlog_proto_rawDescGZIP(), []int{1}
}

func (x *Default_RES) GetCode() uint32 {
	if x != nil {
		return x.Code
	}
	return 0
}

type LogKV_REQ struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Uid      uint32 `protobuf:"varint,1,opt,name=uid,proto3" json:"uid,omitempty"`
	Key      string `protobuf:"bytes,2,opt,name=key,proto3" json:"key,omitempty"`
	Value    string `protobuf:"bytes,3,opt,name=value,proto3" json:"value,omitempty"`
	SaveTime int64  `protobuf:"varint,4,opt,name=saveTime,proto3" json:"saveTime,omitempty"`
}

func (x *LogKV_REQ) Reset() {
	*x = LogKV_REQ{}
	if protoimpl.UnsafeEnabled {
		mi := &file_zlog_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *LogKV_REQ) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*LogKV_REQ) ProtoMessage() {}

func (x *LogKV_REQ) ProtoReflect() protoreflect.Message {
	mi := &file_zlog_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use LogKV_REQ.ProtoReflect.Descriptor instead.
func (*LogKV_REQ) Descriptor() ([]byte, []int) {
	return file_zlog_proto_rawDescGZIP(), []int{2}
}

func (x *LogKV_REQ) GetUid() uint32 {
	if x != nil {
		return x.Uid
	}
	return 0
}

func (x *LogKV_REQ) GetKey() string {
	if x != nil {
		return x.Key
	}
	return ""
}

func (x *LogKV_REQ) GetValue() string {
	if x != nil {
		return x.Value
	}
	return ""
}

func (x *LogKV_REQ) GetSaveTime() int64 {
	if x != nil {
		return x.SaveTime
	}
	return 0
}

var File_zlog_proto protoreflect.FileDescriptor

var file_zlog_proto_rawDesc = []byte{
	0x0a, 0x0a, 0x7a, 0x6c, 0x6f, 0x67, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x07, 0x7a, 0x6c,
	0x6f, 0x67, 0x5f, 0x70, 0x62, 0x22, 0x0d, 0x0a, 0x0b, 0x44, 0x65, 0x66, 0x61, 0x75, 0x6c, 0x74,
	0x5f, 0x52, 0x45, 0x51, 0x22, 0x21, 0x0a, 0x0b, 0x44, 0x65, 0x66, 0x61, 0x75, 0x6c, 0x74, 0x5f,
	0x52, 0x45, 0x53, 0x12, 0x12, 0x0a, 0x04, 0x63, 0x6f, 0x64, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x0d, 0x52, 0x04, 0x63, 0x6f, 0x64, 0x65, 0x22, 0x61, 0x0a, 0x09, 0x4c, 0x6f, 0x67, 0x4b, 0x56,
	0x5f, 0x52, 0x45, 0x51, 0x12, 0x10, 0x0a, 0x03, 0x75, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x0d, 0x52, 0x03, 0x75, 0x69, 0x64, 0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x12, 0x14, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75,
	0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x12, 0x1a,
	0x0a, 0x08, 0x73, 0x61, 0x76, 0x65, 0x54, 0x69, 0x6d, 0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x03,
	0x52, 0x08, 0x73, 0x61, 0x76, 0x65, 0x54, 0x69, 0x6d, 0x65, 0x32, 0x3c, 0x0a, 0x04, 0x7a, 0x6c,
	0x6f, 0x67, 0x12, 0x34, 0x0a, 0x08, 0x41, 0x64, 0x64, 0x4c, 0x6f, 0x67, 0x4b, 0x56, 0x12, 0x12,
	0x2e, 0x7a, 0x6c, 0x6f, 0x67, 0x5f, 0x70, 0x62, 0x2e, 0x4c, 0x6f, 0x67, 0x4b, 0x56, 0x5f, 0x52,
	0x45, 0x51, 0x1a, 0x14, 0x2e, 0x7a, 0x6c, 0x6f, 0x67, 0x5f, 0x70, 0x62, 0x2e, 0x44, 0x65, 0x66,
	0x61, 0x75, 0x6c, 0x74, 0x5f, 0x52, 0x45, 0x53, 0x42, 0x0b, 0x5a, 0x09, 0x2e, 0x2f, 0x7a, 0x6c,
	0x6f, 0x67, 0x5f, 0x70, 0x62, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_zlog_proto_rawDescOnce sync.Once
	file_zlog_proto_rawDescData = file_zlog_proto_rawDesc
)

func file_zlog_proto_rawDescGZIP() []byte {
	file_zlog_proto_rawDescOnce.Do(func() {
		file_zlog_proto_rawDescData = protoimpl.X.CompressGZIP(file_zlog_proto_rawDescData)
	})
	return file_zlog_proto_rawDescData
}

var file_zlog_proto_msgTypes = make([]protoimpl.MessageInfo, 3)
var file_zlog_proto_goTypes = []interface{}{
	(*Default_REQ)(nil), // 0: zlog_pb.Default_REQ
	(*Default_RES)(nil), // 1: zlog_pb.Default_RES
	(*LogKV_REQ)(nil),   // 2: zlog_pb.LogKV_REQ
}
var file_zlog_proto_depIdxs = []int32{
	2, // 0: zlog_pb.zlog.AddLogKV:input_type -> zlog_pb.LogKV_REQ
	1, // 1: zlog_pb.zlog.AddLogKV:output_type -> zlog_pb.Default_RES
	1, // [1:2] is the sub-list for method output_type
	0, // [0:1] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_zlog_proto_init() }
func file_zlog_proto_init() {
	if File_zlog_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_zlog_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Default_REQ); i {
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
		file_zlog_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Default_RES); i {
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
		file_zlog_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*LogKV_REQ); i {
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
			RawDescriptor: file_zlog_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   3,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_zlog_proto_goTypes,
		DependencyIndexes: file_zlog_proto_depIdxs,
		MessageInfos:      file_zlog_proto_msgTypes,
	}.Build()
	File_zlog_proto = out.File
	file_zlog_proto_rawDesc = nil
	file_zlog_proto_goTypes = nil
	file_zlog_proto_depIdxs = nil
}