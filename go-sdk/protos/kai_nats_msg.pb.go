// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        v4.23.4
// source: kai_nats_msg.proto

package kai

import (
	reflect "reflect"
	sync "sync"

	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	anypb "google.golang.org/protobuf/types/known/anypb"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type MessageType int32

const (
	MessageType_UNDEFINED MessageType = 0
	MessageType_OK        MessageType = 1
	MessageType_ERROR     MessageType = 2
)

// Enum value maps for MessageType.
var (
	MessageType_name = map[int32]string{
		0: "UNDEFINED",
		1: "OK",
		2: "ERROR",
	}
	MessageType_value = map[string]int32{
		"UNDEFINED": 0,
		"OK":        1,
		"ERROR":     2,
	}
)

func (x MessageType) Enum() *MessageType {
	p := new(MessageType)
	*p = x
	return p
}

func (x MessageType) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (MessageType) Descriptor() protoreflect.EnumDescriptor {
	return file_kai_nats_msg_proto_enumTypes[0].Descriptor()
}

func (MessageType) Type() protoreflect.EnumType {
	return &file_kai_nats_msg_proto_enumTypes[0]
}

func (x MessageType) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use MessageType.Descriptor instead.
func (MessageType) EnumDescriptor() ([]byte, []int) {
	return file_kai_nats_msg_proto_rawDescGZIP(), []int{0}
}

type KaiNatsMessage struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	RequestId   string      `protobuf:"bytes,1,opt,name=request_id,json=requestId,proto3" json:"request_id,omitempty"`
	Payload     *anypb.Any  `protobuf:"bytes,2,opt,name=payload,proto3" json:"payload,omitempty"`
	Error       string      `protobuf:"bytes,3,opt,name=error,proto3" json:"error,omitempty"`
	FromNode    string      `protobuf:"bytes,4,opt,name=from_node,json=fromNode,proto3" json:"from_node,omitempty"`
	MessageType MessageType `protobuf:"varint,5,opt,name=message_type,json=messageType,proto3,enum=MessageType" json:"message_type,omitempty"`
}

func (x *KaiNatsMessage) Reset() {
	*x = KaiNatsMessage{}
	if protoimpl.UnsafeEnabled {
		mi := &file_kai_nats_msg_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *KaiNatsMessage) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*KaiNatsMessage) ProtoMessage() {}

func (x *KaiNatsMessage) ProtoReflect() protoreflect.Message {
	mi := &file_kai_nats_msg_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use KaiNatsMessage.ProtoReflect.Descriptor instead.
func (*KaiNatsMessage) Descriptor() ([]byte, []int) {
	return file_kai_nats_msg_proto_rawDescGZIP(), []int{0}
}

func (x *KaiNatsMessage) GetRequestId() string {
	if x != nil {
		return x.RequestId
	}
	return ""
}

func (x *KaiNatsMessage) GetPayload() *anypb.Any {
	if x != nil {
		return x.Payload
	}
	return nil
}

func (x *KaiNatsMessage) GetError() string {
	if x != nil {
		return x.Error
	}
	return ""
}

func (x *KaiNatsMessage) GetFromNode() string {
	if x != nil {
		return x.FromNode
	}
	return ""
}

func (x *KaiNatsMessage) GetMessageType() MessageType {
	if x != nil {
		return x.MessageType
	}
	return MessageType_UNDEFINED
}

var File_kai_nats_msg_proto protoreflect.FileDescriptor

var file_kai_nats_msg_proto_rawDesc = []byte{
	0x0a, 0x12, 0x6b, 0x61, 0x69, 0x5f, 0x6e, 0x61, 0x74, 0x73, 0x5f, 0x6d, 0x73, 0x67, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x19, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x61, 0x6e, 0x79, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22,
	0xc3, 0x01, 0x0a, 0x0e, 0x4b, 0x61, 0x69, 0x4e, 0x61, 0x74, 0x73, 0x4d, 0x65, 0x73, 0x73, 0x61,
	0x67, 0x65, 0x12, 0x1d, 0x0a, 0x0a, 0x72, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x5f, 0x69, 0x64,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x72, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x49,
	0x64, 0x12, 0x2e, 0x0a, 0x07, 0x70, 0x61, 0x79, 0x6c, 0x6f, 0x61, 0x64, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x0b, 0x32, 0x14, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x62, 0x75, 0x66, 0x2e, 0x41, 0x6e, 0x79, 0x52, 0x07, 0x70, 0x61, 0x79, 0x6c, 0x6f, 0x61,
	0x64, 0x12, 0x14, 0x0a, 0x05, 0x65, 0x72, 0x72, 0x6f, 0x72, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x05, 0x65, 0x72, 0x72, 0x6f, 0x72, 0x12, 0x1b, 0x0a, 0x09, 0x66, 0x72, 0x6f, 0x6d, 0x5f,
	0x6e, 0x6f, 0x64, 0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x66, 0x72, 0x6f, 0x6d,
	0x4e, 0x6f, 0x64, 0x65, 0x12, 0x2f, 0x0a, 0x0c, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x5f,
	0x74, 0x79, 0x70, 0x65, 0x18, 0x05, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x0c, 0x2e, 0x4d, 0x65, 0x73,
	0x73, 0x61, 0x67, 0x65, 0x54, 0x79, 0x70, 0x65, 0x52, 0x0b, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67,
	0x65, 0x54, 0x79, 0x70, 0x65, 0x2a, 0x2f, 0x0a, 0x0b, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65,
	0x54, 0x79, 0x70, 0x65, 0x12, 0x0d, 0x0a, 0x09, 0x55, 0x4e, 0x44, 0x45, 0x46, 0x49, 0x4e, 0x45,
	0x44, 0x10, 0x00, 0x12, 0x06, 0x0a, 0x02, 0x4f, 0x4b, 0x10, 0x01, 0x12, 0x09, 0x0a, 0x05, 0x45,
	0x52, 0x52, 0x4f, 0x52, 0x10, 0x02, 0x42, 0x07, 0x5a, 0x05, 0x2e, 0x2f, 0x6b, 0x61, 0x69, 0x62,
	0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_kai_nats_msg_proto_rawDescOnce sync.Once
	file_kai_nats_msg_proto_rawDescData = file_kai_nats_msg_proto_rawDesc
)

func file_kai_nats_msg_proto_rawDescGZIP() []byte {
	file_kai_nats_msg_proto_rawDescOnce.Do(func() {
		file_kai_nats_msg_proto_rawDescData = protoimpl.X.CompressGZIP(file_kai_nats_msg_proto_rawDescData)
	})
	return file_kai_nats_msg_proto_rawDescData
}

var file_kai_nats_msg_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_kai_nats_msg_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_kai_nats_msg_proto_goTypes = []interface{}{
	(MessageType)(0),       // 0: MessageType
	(*KaiNatsMessage)(nil), // 1: KaiNatsMessage
	(*anypb.Any)(nil),      // 2: google.protobuf.Any
}
var file_kai_nats_msg_proto_depIdxs = []int32{
	2, // 0: KaiNatsMessage.payload:type_name -> google.protobuf.Any
	0, // 1: KaiNatsMessage.message_type:type_name -> MessageType
	2, // [2:2] is the sub-list for method output_type
	2, // [2:2] is the sub-list for method input_type
	2, // [2:2] is the sub-list for extension type_name
	2, // [2:2] is the sub-list for extension extendee
	0, // [0:2] is the sub-list for field type_name
}

func init() { file_kai_nats_msg_proto_init() }
func file_kai_nats_msg_proto_init() {
	if File_kai_nats_msg_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_kai_nats_msg_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*KaiNatsMessage); i {
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
			RawDescriptor: file_kai_nats_msg_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_kai_nats_msg_proto_goTypes,
		DependencyIndexes: file_kai_nats_msg_proto_depIdxs,
		EnumInfos:         file_kai_nats_msg_proto_enumTypes,
		MessageInfos:      file_kai_nats_msg_proto_msgTypes,
	}.Build()
	File_kai_nats_msg_proto = out.File
	file_kai_nats_msg_proto_rawDesc = nil
	file_kai_nats_msg_proto_goTypes = nil
	file_kai_nats_msg_proto_depIdxs = nil
}
