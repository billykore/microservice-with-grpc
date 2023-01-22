// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        v3.21.12
// source: pb/payment/v1/payment.proto

package v1

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

type QrisRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	MerchantId         string `protobuf:"bytes,1,opt,name=merchantId,proto3" json:"merchantId,omitempty"`
	TrxNumber          string `protobuf:"bytes,2,opt,name=trxNumber,proto3" json:"trxNumber,omitempty"`
	AccountSource      string `protobuf:"bytes,3,opt,name=accountSource,proto3" json:"accountSource,omitempty"`
	AccountDestination string `protobuf:"bytes,4,opt,name=accountDestination,proto3" json:"accountDestination,omitempty"`
	Amount             string `protobuf:"bytes,5,opt,name=amount,proto3" json:"amount,omitempty"`
}

func (x *QrisRequest) Reset() {
	*x = QrisRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pb_payment_v1_payment_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *QrisRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*QrisRequest) ProtoMessage() {}

func (x *QrisRequest) ProtoReflect() protoreflect.Message {
	mi := &file_pb_payment_v1_payment_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use QrisRequest.ProtoReflect.Descriptor instead.
func (*QrisRequest) Descriptor() ([]byte, []int) {
	return file_pb_payment_v1_payment_proto_rawDescGZIP(), []int{0}
}

func (x *QrisRequest) GetMerchantId() string {
	if x != nil {
		return x.MerchantId
	}
	return ""
}

func (x *QrisRequest) GetTrxNumber() string {
	if x != nil {
		return x.TrxNumber
	}
	return ""
}

func (x *QrisRequest) GetAccountSource() string {
	if x != nil {
		return x.AccountSource
	}
	return ""
}

func (x *QrisRequest) GetAccountDestination() string {
	if x != nil {
		return x.AccountDestination
	}
	return ""
}

func (x *QrisRequest) GetAmount() string {
	if x != nil {
		return x.Amount
	}
	return ""
}

type QrisResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Success bool `protobuf:"varint,1,opt,name=success,proto3" json:"success,omitempty"`
}

func (x *QrisResponse) Reset() {
	*x = QrisResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pb_payment_v1_payment_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *QrisResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*QrisResponse) ProtoMessage() {}

func (x *QrisResponse) ProtoReflect() protoreflect.Message {
	mi := &file_pb_payment_v1_payment_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use QrisResponse.ProtoReflect.Descriptor instead.
func (*QrisResponse) Descriptor() ([]byte, []int) {
	return file_pb_payment_v1_payment_proto_rawDescGZIP(), []int{1}
}

func (x *QrisResponse) GetSuccess() bool {
	if x != nil {
		return x.Success
	}
	return false
}

var File_pb_payment_v1_payment_proto protoreflect.FileDescriptor

var file_pb_payment_v1_payment_proto_rawDesc = []byte{
	0x0a, 0x1b, 0x70, 0x62, 0x2f, 0x70, 0x61, 0x79, 0x6d, 0x65, 0x6e, 0x74, 0x2f, 0x76, 0x31, 0x2f,
	0x70, 0x61, 0x79, 0x6d, 0x65, 0x6e, 0x74, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x07, 0x70,
	0x61, 0x79, 0x6d, 0x65, 0x6e, 0x74, 0x22, 0xb9, 0x01, 0x0a, 0x0b, 0x51, 0x72, 0x69, 0x73, 0x52,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x1e, 0x0a, 0x0a, 0x6d, 0x65, 0x72, 0x63, 0x68, 0x61,
	0x6e, 0x74, 0x49, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0a, 0x6d, 0x65, 0x72, 0x63,
	0x68, 0x61, 0x6e, 0x74, 0x49, 0x64, 0x12, 0x1c, 0x0a, 0x09, 0x74, 0x72, 0x78, 0x4e, 0x75, 0x6d,
	0x62, 0x65, 0x72, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x74, 0x72, 0x78, 0x4e, 0x75,
	0x6d, 0x62, 0x65, 0x72, 0x12, 0x24, 0x0a, 0x0d, 0x61, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x53,
	0x6f, 0x75, 0x72, 0x63, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0d, 0x61, 0x63, 0x63,
	0x6f, 0x75, 0x6e, 0x74, 0x53, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x12, 0x2e, 0x0a, 0x12, 0x61, 0x63,
	0x63, 0x6f, 0x75, 0x6e, 0x74, 0x44, 0x65, 0x73, 0x74, 0x69, 0x6e, 0x61, 0x74, 0x69, 0x6f, 0x6e,
	0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x12, 0x61, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x44,
	0x65, 0x73, 0x74, 0x69, 0x6e, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x16, 0x0a, 0x06, 0x61, 0x6d,
	0x6f, 0x75, 0x6e, 0x74, 0x18, 0x05, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x61, 0x6d, 0x6f, 0x75,
	0x6e, 0x74, 0x22, 0x28, 0x0a, 0x0c, 0x51, 0x72, 0x69, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e,
	0x73, 0x65, 0x12, 0x18, 0x0a, 0x07, 0x73, 0x75, 0x63, 0x63, 0x65, 0x73, 0x73, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x08, 0x52, 0x07, 0x73, 0x75, 0x63, 0x63, 0x65, 0x73, 0x73, 0x32, 0x3e, 0x0a, 0x07,
	0x50, 0x61, 0x79, 0x6d, 0x65, 0x6e, 0x74, 0x12, 0x33, 0x0a, 0x04, 0x51, 0x72, 0x69, 0x73, 0x12,
	0x14, 0x2e, 0x70, 0x61, 0x79, 0x6d, 0x65, 0x6e, 0x74, 0x2e, 0x51, 0x72, 0x69, 0x73, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x15, 0x2e, 0x70, 0x61, 0x79, 0x6d, 0x65, 0x6e, 0x74, 0x2e,
	0x51, 0x72, 0x69, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x42, 0x3a, 0x0a, 0x18,
	0x69, 0x6f, 0x2e, 0x67, 0x72, 0x70, 0x63, 0x2e, 0x65, 0x78, 0x61, 0x6d, 0x70, 0x6c, 0x65, 0x73,
	0x2e, 0x70, 0x61, 0x79, 0x6d, 0x65, 0x6e, 0x74, 0x42, 0x0c, 0x50, 0x61, 0x79, 0x6d, 0x65, 0x6e,
	0x74, 0x50, 0x72, 0x6f, 0x74, 0x6f, 0x50, 0x01, 0x5a, 0x0e, 0x67, 0x65, 0x6e, 0x2f, 0x70, 0x61,
	0x79, 0x6d, 0x65, 0x6e, 0x74, 0x2f, 0x76, 0x31, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_pb_payment_v1_payment_proto_rawDescOnce sync.Once
	file_pb_payment_v1_payment_proto_rawDescData = file_pb_payment_v1_payment_proto_rawDesc
)

func file_pb_payment_v1_payment_proto_rawDescGZIP() []byte {
	file_pb_payment_v1_payment_proto_rawDescOnce.Do(func() {
		file_pb_payment_v1_payment_proto_rawDescData = protoimpl.X.CompressGZIP(file_pb_payment_v1_payment_proto_rawDescData)
	})
	return file_pb_payment_v1_payment_proto_rawDescData
}

var file_pb_payment_v1_payment_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_pb_payment_v1_payment_proto_goTypes = []interface{}{
	(*QrisRequest)(nil),  // 0: payment.QrisRequest
	(*QrisResponse)(nil), // 1: payment.QrisResponse
}
var file_pb_payment_v1_payment_proto_depIdxs = []int32{
	0, // 0: payment.Payment.Qris:input_type -> payment.QrisRequest
	1, // 1: payment.Payment.Qris:output_type -> payment.QrisResponse
	1, // [1:2] is the sub-list for method output_type
	0, // [0:1] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_pb_payment_v1_payment_proto_init() }
func file_pb_payment_v1_payment_proto_init() {
	if File_pb_payment_v1_payment_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_pb_payment_v1_payment_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*QrisRequest); i {
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
		file_pb_payment_v1_payment_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*QrisResponse); i {
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
			RawDescriptor: file_pb_payment_v1_payment_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_pb_payment_v1_payment_proto_goTypes,
		DependencyIndexes: file_pb_payment_v1_payment_proto_depIdxs,
		MessageInfos:      file_pb_payment_v1_payment_proto_msgTypes,
	}.Build()
	File_pb_payment_v1_payment_proto = out.File
	file_pb_payment_v1_payment_proto_rawDesc = nil
	file_pb_payment_v1_payment_proto_goTypes = nil
	file_pb_payment_v1_payment_proto_depIdxs = nil
}
