// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.0
// 	protoc        v3.12.3
// source: ws_protobuf.proto

//import "google/protobuf/any.proto";

package protobuf

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

type CheckHealthReq struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Ts int64 `protobuf:"varint,1,opt,name=ts,proto3" json:"ts,omitempty"`
}

func (x *CheckHealthReq) Reset() {
	*x = CheckHealthReq{}
	if protoimpl.UnsafeEnabled {
		mi := &file_ws_protobuf_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CheckHealthReq) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CheckHealthReq) ProtoMessage() {}

func (x *CheckHealthReq) ProtoReflect() protoreflect.Message {
	mi := &file_ws_protobuf_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CheckHealthReq.ProtoReflect.Descriptor instead.
func (*CheckHealthReq) Descriptor() ([]byte, []int) {
	return file_ws_protobuf_proto_rawDescGZIP(), []int{0}
}

func (x *CheckHealthReq) GetTs() int64 {
	if x != nil {
		return x.Ts
	}
	return 0
}

type CheckHealthRsp struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Ts int64 `protobuf:"varint,1,opt,name=ts,proto3" json:"ts,omitempty"`
}

func (x *CheckHealthRsp) Reset() {
	*x = CheckHealthRsp{}
	if protoimpl.UnsafeEnabled {
		mi := &file_ws_protobuf_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CheckHealthRsp) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CheckHealthRsp) ProtoMessage() {}

func (x *CheckHealthRsp) ProtoReflect() protoreflect.Message {
	mi := &file_ws_protobuf_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CheckHealthRsp.ProtoReflect.Descriptor instead.
func (*CheckHealthRsp) Descriptor() ([]byte, []int) {
	return file_ws_protobuf_proto_rawDescGZIP(), []int{1}
}

func (x *CheckHealthRsp) GetTs() int64 {
	if x != nil {
		return x.Ts
	}
	return 0
}

type OkRsp struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Errcode uint32 `protobuf:"varint,1,opt,name=errcode,proto3" json:"errcode,omitempty"`
	ErrMsg  string `protobuf:"bytes,2,opt,name=errMsg,proto3" json:"errMsg,omitempty"`
}

func (x *OkRsp) Reset() {
	*x = OkRsp{}
	if protoimpl.UnsafeEnabled {
		mi := &file_ws_protobuf_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *OkRsp) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*OkRsp) ProtoMessage() {}

func (x *OkRsp) ProtoReflect() protoreflect.Message {
	mi := &file_ws_protobuf_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use OkRsp.ProtoReflect.Descriptor instead.
func (*OkRsp) Descriptor() ([]byte, []int) {
	return file_ws_protobuf_proto_rawDescGZIP(), []int{2}
}

func (x *OkRsp) GetErrcode() uint32 {
	if x != nil {
		return x.Errcode
	}
	return 0
}

func (x *OkRsp) GetErrMsg() string {
	if x != nil {
		return x.ErrMsg
	}
	return ""
}

type SendMsgReq struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Clientid   string `protobuf:"bytes,1,opt,name=clientid,proto3" json:"clientid,omitempty"`                        //
	ProtocolId int64  `protobuf:"varint,2,opt,name=protocol_id,json=protocolId,proto3" json:"protocol_id,omitempty"` // protocol_id
	BodyType   int32  `protobuf:"varint,3,opt,name=body_type,json=bodyType,proto3" json:"body_type,omitempty"`       //
	Queue      int32  `protobuf:"varint,4,opt,name=queue,proto3" json:"queue,omitempty"`                             //
	Broadcast  int32  `protobuf:"varint,5,opt,name=broadcast,proto3" json:"broadcast,omitempty"`                     //全局广播
	//google.protobuf.Any body = 6; //
	Body []byte `protobuf:"bytes,6,opt,name=body,proto3" json:"body,omitempty"`
}

func (x *SendMsgReq) Reset() {
	*x = SendMsgReq{}
	if protoimpl.UnsafeEnabled {
		mi := &file_ws_protobuf_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SendMsgReq) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SendMsgReq) ProtoMessage() {}

func (x *SendMsgReq) ProtoReflect() protoreflect.Message {
	mi := &file_ws_protobuf_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SendMsgReq.ProtoReflect.Descriptor instead.
func (*SendMsgReq) Descriptor() ([]byte, []int) {
	return file_ws_protobuf_proto_rawDescGZIP(), []int{3}
}

func (x *SendMsgReq) GetClientid() string {
	if x != nil {
		return x.Clientid
	}
	return ""
}

func (x *SendMsgReq) GetProtocolId() int64 {
	if x != nil {
		return x.ProtocolId
	}
	return 0
}

func (x *SendMsgReq) GetBodyType() int32 {
	if x != nil {
		return x.BodyType
	}
	return 0
}

func (x *SendMsgReq) GetQueue() int32 {
	if x != nil {
		return x.Queue
	}
	return 0
}

func (x *SendMsgReq) GetBroadcast() int32 {
	if x != nil {
		return x.Broadcast
	}
	return 0
}

func (x *SendMsgReq) GetBody() []byte {
	if x != nil {
		return x.Body
	}
	return nil
}

type ForceDisconnectReq struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Clientid string `protobuf:"bytes,1,opt,name=clientid,proto3" json:"clientid,omitempty"`
}

func (x *ForceDisconnectReq) Reset() {
	*x = ForceDisconnectReq{}
	if protoimpl.UnsafeEnabled {
		mi := &file_ws_protobuf_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ForceDisconnectReq) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ForceDisconnectReq) ProtoMessage() {}

func (x *ForceDisconnectReq) ProtoReflect() protoreflect.Message {
	mi := &file_ws_protobuf_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ForceDisconnectReq.ProtoReflect.Descriptor instead.
func (*ForceDisconnectReq) Descriptor() ([]byte, []int) {
	return file_ws_protobuf_proto_rawDescGZIP(), []int{4}
}

func (x *ForceDisconnectReq) GetClientid() string {
	if x != nil {
		return x.Clientid
	}
	return ""
}

var File_ws_protobuf_proto protoreflect.FileDescriptor

var file_ws_protobuf_proto_rawDesc = []byte{
	0x0a, 0x11, 0x77, 0x73, 0x5f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x12, 0x08, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x22, 0x20, 0x0a,
	0x0e, 0x43, 0x68, 0x65, 0x63, 0x6b, 0x48, 0x65, 0x61, 0x6c, 0x74, 0x68, 0x52, 0x65, 0x71, 0x12,
	0x0e, 0x0a, 0x02, 0x74, 0x73, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x02, 0x74, 0x73, 0x22,
	0x20, 0x0a, 0x0e, 0x43, 0x68, 0x65, 0x63, 0x6b, 0x48, 0x65, 0x61, 0x6c, 0x74, 0x68, 0x52, 0x73,
	0x70, 0x12, 0x0e, 0x0a, 0x02, 0x74, 0x73, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x02, 0x74,
	0x73, 0x22, 0x39, 0x0a, 0x05, 0x4f, 0x6b, 0x52, 0x73, 0x70, 0x12, 0x18, 0x0a, 0x07, 0x65, 0x72,
	0x72, 0x63, 0x6f, 0x64, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x07, 0x65, 0x72, 0x72,
	0x63, 0x6f, 0x64, 0x65, 0x12, 0x16, 0x0a, 0x06, 0x65, 0x72, 0x72, 0x4d, 0x73, 0x67, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x65, 0x72, 0x72, 0x4d, 0x73, 0x67, 0x22, 0xae, 0x01, 0x0a,
	0x0a, 0x53, 0x65, 0x6e, 0x64, 0x4d, 0x73, 0x67, 0x52, 0x65, 0x71, 0x12, 0x1a, 0x0a, 0x08, 0x63,
	0x6c, 0x69, 0x65, 0x6e, 0x74, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x63,
	0x6c, 0x69, 0x65, 0x6e, 0x74, 0x69, 0x64, 0x12, 0x1f, 0x0a, 0x0b, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x63, 0x6f, 0x6c, 0x5f, 0x69, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x03, 0x52, 0x0a, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x63, 0x6f, 0x6c, 0x49, 0x64, 0x12, 0x1b, 0x0a, 0x09, 0x62, 0x6f, 0x64, 0x79,
	0x5f, 0x74, 0x79, 0x70, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x05, 0x52, 0x08, 0x62, 0x6f, 0x64,
	0x79, 0x54, 0x79, 0x70, 0x65, 0x12, 0x14, 0x0a, 0x05, 0x71, 0x75, 0x65, 0x75, 0x65, 0x18, 0x04,
	0x20, 0x01, 0x28, 0x05, 0x52, 0x05, 0x71, 0x75, 0x65, 0x75, 0x65, 0x12, 0x1c, 0x0a, 0x09, 0x62,
	0x72, 0x6f, 0x61, 0x64, 0x63, 0x61, 0x73, 0x74, 0x18, 0x05, 0x20, 0x01, 0x28, 0x05, 0x52, 0x09,
	0x62, 0x72, 0x6f, 0x61, 0x64, 0x63, 0x61, 0x73, 0x74, 0x12, 0x12, 0x0a, 0x04, 0x62, 0x6f, 0x64,
	0x79, 0x18, 0x06, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x04, 0x62, 0x6f, 0x64, 0x79, 0x22, 0x30, 0x0a,
	0x12, 0x46, 0x6f, 0x72, 0x63, 0x65, 0x44, 0x69, 0x73, 0x63, 0x6f, 0x6e, 0x6e, 0x65, 0x63, 0x74,
	0x52, 0x65, 0x71, 0x12, 0x1a, 0x0a, 0x08, 0x63, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x69, 0x64, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x63, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x69, 0x64, 0x32,
	0xc7, 0x01, 0x0a, 0x08, 0x57, 0x53, 0x53, 0x65, 0x72, 0x76, 0x65, 0x72, 0x12, 0x43, 0x0a, 0x0b,
	0x43, 0x68, 0x65, 0x63, 0x6b, 0x48, 0x65, 0x61, 0x6c, 0x74, 0x68, 0x12, 0x18, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x43, 0x68, 0x65, 0x63, 0x6b, 0x48, 0x65, 0x61, 0x6c,
	0x74, 0x68, 0x52, 0x65, 0x71, 0x1a, 0x18, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66,
	0x2e, 0x43, 0x68, 0x65, 0x63, 0x6b, 0x48, 0x65, 0x61, 0x6c, 0x74, 0x68, 0x52, 0x73, 0x70, 0x22,
	0x00, 0x12, 0x32, 0x0a, 0x07, 0x53, 0x65, 0x6e, 0x64, 0x4d, 0x73, 0x67, 0x12, 0x14, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x53, 0x65, 0x6e, 0x64, 0x4d, 0x73, 0x67, 0x52,
	0x65, 0x71, 0x1a, 0x0f, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x4f, 0x6b,
	0x52, 0x73, 0x70, 0x22, 0x00, 0x12, 0x42, 0x0a, 0x0f, 0x46, 0x6f, 0x72, 0x63, 0x65, 0x44, 0x69,
	0x73, 0x63, 0x6f, 0x6e, 0x6e, 0x65, 0x63, 0x74, 0x12, 0x1c, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x62, 0x75, 0x66, 0x2e, 0x46, 0x6f, 0x72, 0x63, 0x65, 0x44, 0x69, 0x73, 0x63, 0x6f, 0x6e, 0x6e,
	0x65, 0x63, 0x74, 0x52, 0x65, 0x71, 0x1a, 0x0f, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75,
	0x66, 0x2e, 0x4f, 0x6b, 0x52, 0x73, 0x70, 0x22, 0x00, 0x42, 0x39, 0x0a, 0x19, 0x69, 0x6f, 0x2e,
	0x67, 0x72, 0x70, 0x63, 0x2e, 0x65, 0x78, 0x61, 0x6d, 0x70, 0x6c, 0x65, 0x73, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x42, 0x0d, 0x50, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66,
	0x50, 0x72, 0x6f, 0x74, 0x6f, 0x50, 0x01, 0x5a, 0x0b, 0x2e, 0x2e, 0x2f, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x62, 0x75, 0x66, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_ws_protobuf_proto_rawDescOnce sync.Once
	file_ws_protobuf_proto_rawDescData = file_ws_protobuf_proto_rawDesc
)

func file_ws_protobuf_proto_rawDescGZIP() []byte {
	file_ws_protobuf_proto_rawDescOnce.Do(func() {
		file_ws_protobuf_proto_rawDescData = protoimpl.X.CompressGZIP(file_ws_protobuf_proto_rawDescData)
	})
	return file_ws_protobuf_proto_rawDescData
}

var file_ws_protobuf_proto_msgTypes = make([]protoimpl.MessageInfo, 5)
var file_ws_protobuf_proto_goTypes = []interface{}{
	(*CheckHealthReq)(nil),     // 0: protobuf.CheckHealthReq
	(*CheckHealthRsp)(nil),     // 1: protobuf.CheckHealthRsp
	(*OkRsp)(nil),              // 2: protobuf.OkRsp
	(*SendMsgReq)(nil),         // 3: protobuf.SendMsgReq
	(*ForceDisconnectReq)(nil), // 4: protobuf.ForceDisconnectReq
}
var file_ws_protobuf_proto_depIdxs = []int32{
	0, // 0: protobuf.WSServer.CheckHealth:input_type -> protobuf.CheckHealthReq
	3, // 1: protobuf.WSServer.SendMsg:input_type -> protobuf.SendMsgReq
	4, // 2: protobuf.WSServer.ForceDisconnect:input_type -> protobuf.ForceDisconnectReq
	1, // 3: protobuf.WSServer.CheckHealth:output_type -> protobuf.CheckHealthRsp
	2, // 4: protobuf.WSServer.SendMsg:output_type -> protobuf.OkRsp
	2, // 5: protobuf.WSServer.ForceDisconnect:output_type -> protobuf.OkRsp
	3, // [3:6] is the sub-list for method output_type
	0, // [0:3] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_ws_protobuf_proto_init() }
func file_ws_protobuf_proto_init() {
	if File_ws_protobuf_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_ws_protobuf_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CheckHealthReq); i {
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
		file_ws_protobuf_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CheckHealthRsp); i {
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
		file_ws_protobuf_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*OkRsp); i {
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
		file_ws_protobuf_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SendMsgReq); i {
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
		file_ws_protobuf_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ForceDisconnectReq); i {
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
			RawDescriptor: file_ws_protobuf_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   5,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_ws_protobuf_proto_goTypes,
		DependencyIndexes: file_ws_protobuf_proto_depIdxs,
		MessageInfos:      file_ws_protobuf_proto_msgTypes,
	}.Build()
	File_ws_protobuf_proto = out.File
	file_ws_protobuf_proto_rawDesc = nil
	file_ws_protobuf_proto_goTypes = nil
	file_ws_protobuf_proto_depIdxs = nil
}
