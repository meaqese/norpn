// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.6
// 	protoc        v6.30.2
// source: internal/pb/norpn.proto

package pb

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
	unsafe "unsafe"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type Empty struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *Empty) Reset() {
	*x = Empty{}
	mi := &file_internal_pb_norpn_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Empty) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Empty) ProtoMessage() {}

func (x *Empty) ProtoReflect() protoreflect.Message {
	mi := &file_internal_pb_norpn_proto_msgTypes[0]
	if x != nil {
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
	return file_internal_pb_norpn_proto_rawDescGZIP(), []int{0}
}

type TaskResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	ID            string                 `protobuf:"bytes,1,opt,name=ID,proto3" json:"ID,omitempty"`
	Arg1          float32                `protobuf:"fixed32,2,opt,name=Arg1,proto3" json:"Arg1,omitempty"`
	Arg2          float32                `protobuf:"fixed32,3,opt,name=Arg2,proto3" json:"Arg2,omitempty"`
	Operation     int32                  `protobuf:"varint,4,opt,name=Operation,proto3" json:"Operation,omitempty"`
	OperationTime int32                  `protobuf:"varint,5,opt,name=OperationTime,proto3" json:"OperationTime,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *TaskResponse) Reset() {
	*x = TaskResponse{}
	mi := &file_internal_pb_norpn_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *TaskResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TaskResponse) ProtoMessage() {}

func (x *TaskResponse) ProtoReflect() protoreflect.Message {
	mi := &file_internal_pb_norpn_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TaskResponse.ProtoReflect.Descriptor instead.
func (*TaskResponse) Descriptor() ([]byte, []int) {
	return file_internal_pb_norpn_proto_rawDescGZIP(), []int{1}
}

func (x *TaskResponse) GetID() string {
	if x != nil {
		return x.ID
	}
	return ""
}

func (x *TaskResponse) GetArg1() float32 {
	if x != nil {
		return x.Arg1
	}
	return 0
}

func (x *TaskResponse) GetArg2() float32 {
	if x != nil {
		return x.Arg2
	}
	return 0
}

func (x *TaskResponse) GetOperation() int32 {
	if x != nil {
		return x.Operation
	}
	return 0
}

func (x *TaskResponse) GetOperationTime() int32 {
	if x != nil {
		return x.OperationTime
	}
	return 0
}

type TaskResult struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	ID            string                 `protobuf:"bytes,1,opt,name=ID,proto3" json:"ID,omitempty"`
	Result        float32                `protobuf:"fixed32,2,opt,name=Result,proto3" json:"Result,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *TaskResult) Reset() {
	*x = TaskResult{}
	mi := &file_internal_pb_norpn_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *TaskResult) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TaskResult) ProtoMessage() {}

func (x *TaskResult) ProtoReflect() protoreflect.Message {
	mi := &file_internal_pb_norpn_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TaskResult.ProtoReflect.Descriptor instead.
func (*TaskResult) Descriptor() ([]byte, []int) {
	return file_internal_pb_norpn_proto_rawDescGZIP(), []int{2}
}

func (x *TaskResult) GetID() string {
	if x != nil {
		return x.ID
	}
	return ""
}

func (x *TaskResult) GetResult() float32 {
	if x != nil {
		return x.Result
	}
	return 0
}

var File_internal_pb_norpn_proto protoreflect.FileDescriptor

const file_internal_pb_norpn_proto_rawDesc = "" +
	"\n" +
	"\x17internal/pb/norpn.proto\x12\x05norpn\"\a\n" +
	"\x05Empty\"\x8a\x01\n" +
	"\fTaskResponse\x12\x0e\n" +
	"\x02ID\x18\x01 \x01(\tR\x02ID\x12\x12\n" +
	"\x04Arg1\x18\x02 \x01(\x02R\x04Arg1\x12\x12\n" +
	"\x04Arg2\x18\x03 \x01(\x02R\x04Arg2\x12\x1c\n" +
	"\tOperation\x18\x04 \x01(\x05R\tOperation\x12$\n" +
	"\rOperationTime\x18\x05 \x01(\x05R\rOperationTime\"4\n" +
	"\n" +
	"TaskResult\x12\x0e\n" +
	"\x02ID\x18\x01 \x01(\tR\x02ID\x12\x16\n" +
	"\x06Result\x18\x02 \x01(\x02R\x06Result2n\n" +
	"\vOrchService\x12,\n" +
	"\aGetTask\x12\f.norpn.Empty\x1a\x13.norpn.TaskResponse\x121\n" +
	"\x0eSendTaskResult\x12\x11.norpn.TaskResult\x1a\f.norpn.EmptyB&Z$github.com/meaqese/norpn/internal/pbb\x06proto3"

var (
	file_internal_pb_norpn_proto_rawDescOnce sync.Once
	file_internal_pb_norpn_proto_rawDescData []byte
)

func file_internal_pb_norpn_proto_rawDescGZIP() []byte {
	file_internal_pb_norpn_proto_rawDescOnce.Do(func() {
		file_internal_pb_norpn_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_internal_pb_norpn_proto_rawDesc), len(file_internal_pb_norpn_proto_rawDesc)))
	})
	return file_internal_pb_norpn_proto_rawDescData
}

var file_internal_pb_norpn_proto_msgTypes = make([]protoimpl.MessageInfo, 3)
var file_internal_pb_norpn_proto_goTypes = []any{
	(*Empty)(nil),        // 0: norpn.Empty
	(*TaskResponse)(nil), // 1: norpn.TaskResponse
	(*TaskResult)(nil),   // 2: norpn.TaskResult
}
var file_internal_pb_norpn_proto_depIdxs = []int32{
	0, // 0: norpn.OrchService.GetTask:input_type -> norpn.Empty
	2, // 1: norpn.OrchService.SendTaskResult:input_type -> norpn.TaskResult
	1, // 2: norpn.OrchService.GetTask:output_type -> norpn.TaskResponse
	0, // 3: norpn.OrchService.SendTaskResult:output_type -> norpn.Empty
	2, // [2:4] is the sub-list for method output_type
	0, // [0:2] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_internal_pb_norpn_proto_init() }
func file_internal_pb_norpn_proto_init() {
	if File_internal_pb_norpn_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_internal_pb_norpn_proto_rawDesc), len(file_internal_pb_norpn_proto_rawDesc)),
			NumEnums:      0,
			NumMessages:   3,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_internal_pb_norpn_proto_goTypes,
		DependencyIndexes: file_internal_pb_norpn_proto_depIdxs,
		MessageInfos:      file_internal_pb_norpn_proto_msgTypes,
	}.Build()
	File_internal_pb_norpn_proto = out.File
	file_internal_pb_norpn_proto_goTypes = nil
	file_internal_pb_norpn_proto_depIdxs = nil
}
