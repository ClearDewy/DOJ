// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.32.0
// 	protoc        v4.25.2
// source: judge_server.proto

package jspb

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type JudgeItem struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Uid         string `protobuf:"bytes,1,opt,name=uid,proto3" json:"uid,omitempty"`
	Pid         int32  `protobuf:"varint,2,opt,name=pid,proto3" json:"pid,omitempty"`
	Jid         int32  `protobuf:"varint,3,opt,name=jid,proto3" json:"jid,omitempty"`
	Parallelism int32  `protobuf:"varint,4,opt,name=parallelism,proto3" json:"parallelism,omitempty"`
}

func (x *JudgeItem) Reset() {
	*x = JudgeItem{}
	if protoimpl.UnsafeEnabled {
		mi := &file_judge_server_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *JudgeItem) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*JudgeItem) ProtoMessage() {}

func (x *JudgeItem) ProtoReflect() protoreflect.Message {
	mi := &file_judge_server_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use JudgeItem.ProtoReflect.Descriptor instead.
func (*JudgeItem) Descriptor() ([]byte, []int) {
	return file_judge_server_proto_rawDescGZIP(), []int{0}
}

func (x *JudgeItem) GetUid() string {
	if x != nil {
		return x.Uid
	}
	return ""
}

func (x *JudgeItem) GetPid() int32 {
	if x != nil {
		return x.Pid
	}
	return 0
}

func (x *JudgeItem) GetJid() int32 {
	if x != nil {
		return x.Jid
	}
	return 0
}

func (x *JudgeItem) GetParallelism() int32 {
	if x != nil {
		return x.Parallelism
	}
	return 0
}

var File_judge_server_proto protoreflect.FileDescriptor

var file_judge_server_proto_rawDesc = []byte{
	0x0a, 0x12, 0x6a, 0x75, 0x64, 0x67, 0x65, 0x5f, 0x73, 0x65, 0x72, 0x76, 0x65, 0x72, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x12, 0x03, 0x6a, 0x70, 0x62, 0x1a, 0x1b, 0x67, 0x6f, 0x6f, 0x67, 0x6c,
	0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x65, 0x6d, 0x70, 0x74, 0x79,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x63, 0x0a, 0x09, 0x4a, 0x75, 0x64, 0x67, 0x65, 0x49,
	0x74, 0x65, 0x6d, 0x12, 0x10, 0x0a, 0x03, 0x75, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x03, 0x75, 0x69, 0x64, 0x12, 0x10, 0x0a, 0x03, 0x70, 0x69, 0x64, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x05, 0x52, 0x03, 0x70, 0x69, 0x64, 0x12, 0x10, 0x0a, 0x03, 0x6a, 0x69, 0x64, 0x18, 0x03,
	0x20, 0x01, 0x28, 0x05, 0x52, 0x03, 0x6a, 0x69, 0x64, 0x12, 0x20, 0x0a, 0x0b, 0x70, 0x61, 0x72,
	0x61, 0x6c, 0x6c, 0x65, 0x6c, 0x69, 0x73, 0x6d, 0x18, 0x04, 0x20, 0x01, 0x28, 0x05, 0x52, 0x0b,
	0x70, 0x61, 0x72, 0x61, 0x6c, 0x6c, 0x65, 0x6c, 0x69, 0x73, 0x6d, 0x32, 0x3e, 0x0a, 0x0b, 0x4a,
	0x75, 0x64, 0x67, 0x65, 0x53, 0x65, 0x72, 0x76, 0x65, 0x72, 0x12, 0x2f, 0x0a, 0x05, 0x4a, 0x75,
	0x64, 0x67, 0x65, 0x12, 0x0e, 0x2e, 0x6a, 0x70, 0x62, 0x2e, 0x4a, 0x75, 0x64, 0x67, 0x65, 0x49,
	0x74, 0x65, 0x6d, 0x1a, 0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x42, 0x0d, 0x5a, 0x0b, 0x64,
	0x6f, 0x6a, 0x2d, 0x67, 0x6f, 0x2f, 0x6a, 0x73, 0x70, 0x62, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x33,
}

var (
	file_judge_server_proto_rawDescOnce sync.Once
	file_judge_server_proto_rawDescData = file_judge_server_proto_rawDesc
)

func file_judge_server_proto_rawDescGZIP() []byte {
	file_judge_server_proto_rawDescOnce.Do(func() {
		file_judge_server_proto_rawDescData = protoimpl.X.CompressGZIP(file_judge_server_proto_rawDescData)
	})
	return file_judge_server_proto_rawDescData
}

var file_judge_server_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_judge_server_proto_goTypes = []interface{}{
	(*JudgeItem)(nil),     // 0: jpb.JudgeItem
	(*emptypb.Empty)(nil), // 1: google.protobuf.Empty
}
var file_judge_server_proto_depIdxs = []int32{
	0, // 0: jpb.JudgeServer.Judge:input_type -> jpb.JudgeItem
	1, // 1: jpb.JudgeServer.Judge:output_type -> google.protobuf.Empty
	1, // [1:2] is the sub-list for method output_type
	0, // [0:1] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_judge_server_proto_init() }
func file_judge_server_proto_init() {
	if File_judge_server_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_judge_server_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*JudgeItem); i {
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
			RawDescriptor: file_judge_server_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_judge_server_proto_goTypes,
		DependencyIndexes: file_judge_server_proto_depIdxs,
		MessageInfos:      file_judge_server_proto_msgTypes,
	}.Build()
	File_judge_server_proto = out.File
	file_judge_server_proto_rawDesc = nil
	file_judge_server_proto_goTypes = nil
	file_judge_server_proto_depIdxs = nil
}
