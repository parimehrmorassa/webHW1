// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.27.1
// 	protoc        v3.6.1
// source: get_users.proto

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

type Userr struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id        int32  `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	Name      string `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
	Family    string `protobuf:"bytes,3,opt,name=family,proto3" json:"family,omitempty"`
	Age       int32  `protobuf:"varint,4,opt,name=age,proto3" json:"age,omitempty"`
	Sex       string `protobuf:"bytes,5,opt,name=sex,proto3" json:"sex,omitempty"`
	CreatedAt string `protobuf:"bytes,6,opt,name=createdAt,proto3" json:"createdAt,omitempty"`
}

func (x *Userr) Reset() {
	*x = Userr{}
	if protoimpl.UnsafeEnabled {
		mi := &file_get_users_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Userr) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Userr) ProtoMessage() {}

func (x *Userr) ProtoReflect() protoreflect.Message {
	mi := &file_get_users_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Userr.ProtoReflect.Descriptor instead.
func (*Userr) Descriptor() ([]byte, []int) {
	return file_get_users_proto_rawDescGZIP(), []int{0}
}

func (x *Userr) GetId() int32 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *Userr) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *Userr) GetFamily() string {
	if x != nil {
		return x.Family
	}
	return ""
}

func (x *Userr) GetAge() int32 {
	if x != nil {
		return x.Age
	}
	return 0
}

func (x *Userr) GetSex() string {
	if x != nil {
		return x.Sex
	}
	return ""
}

func (x *Userr) GetCreatedAt() string {
	if x != nil {
		return x.CreatedAt
	}
	return ""
}

type GetDataRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	UserId    int32  `protobuf:"varint,1,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"`
	AuthKey   []byte `protobuf:"bytes,2,opt,name=auth_key,json=authKey,proto3" json:"auth_key,omitempty"`
	MessageId int32  `protobuf:"varint,3,opt,name=message_id,json=messageId,proto3" json:"message_id,omitempty"`
	RedisKey  string `protobuf:"bytes,4,opt,name=redis_key,json=redisKey,proto3" json:"redis_key,omitempty"`
}

func (x *GetDataRequest) Reset() {
	*x = GetDataRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_get_users_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetDataRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetDataRequest) ProtoMessage() {}

func (x *GetDataRequest) ProtoReflect() protoreflect.Message {
	mi := &file_get_users_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetDataRequest.ProtoReflect.Descriptor instead.
func (*GetDataRequest) Descriptor() ([]byte, []int) {
	return file_get_users_proto_rawDescGZIP(), []int{1}
}

func (x *GetDataRequest) GetUserId() int32 {
	if x != nil {
		return x.UserId
	}
	return 0
}

func (x *GetDataRequest) GetAuthKey() []byte {
	if x != nil {
		return x.AuthKey
	}
	return nil
}

func (x *GetDataRequest) GetMessageId() int32 {
	if x != nil {
		return x.MessageId
	}
	return 0
}

func (x *GetDataRequest) GetRedisKey() string {
	if x != nil {
		return x.RedisKey
	}
	return ""
}

type GetDataResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ReturnUsers []*Userr `protobuf:"bytes,1,rep,name=return_users,json=returnUsers,proto3" json:"return_users,omitempty"`
	MessageId   int32    `protobuf:"varint,2,opt,name=message_id,json=messageId,proto3" json:"message_id,omitempty"`
}

func (x *GetDataResponse) Reset() {
	*x = GetDataResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_get_users_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetDataResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetDataResponse) ProtoMessage() {}

func (x *GetDataResponse) ProtoReflect() protoreflect.Message {
	mi := &file_get_users_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetDataResponse.ProtoReflect.Descriptor instead.
func (*GetDataResponse) Descriptor() ([]byte, []int) {
	return file_get_users_proto_rawDescGZIP(), []int{2}
}

func (x *GetDataResponse) GetReturnUsers() []*Userr {
	if x != nil {
		return x.ReturnUsers
	}
	return nil
}

func (x *GetDataResponse) GetMessageId() int32 {
	if x != nil {
		return x.MessageId
	}
	return 0
}

var File_get_users_proto protoreflect.FileDescriptor

var file_get_users_proto_rawDesc = []byte{
	0x0a, 0x0f, 0x67, 0x65, 0x74, 0x5f, 0x75, 0x73, 0x65, 0x72, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x22, 0x85, 0x01, 0x0a, 0x05, 0x55, 0x73, 0x65, 0x72, 0x72, 0x12, 0x0e, 0x0a, 0x02, 0x69,
	0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x05, 0x52, 0x02, 0x69, 0x64, 0x12, 0x12, 0x0a, 0x04, 0x6e,
	0x61, 0x6d, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12,
	0x16, 0x0a, 0x06, 0x66, 0x61, 0x6d, 0x69, 0x6c, 0x79, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x06, 0x66, 0x61, 0x6d, 0x69, 0x6c, 0x79, 0x12, 0x10, 0x0a, 0x03, 0x61, 0x67, 0x65, 0x18, 0x04,
	0x20, 0x01, 0x28, 0x05, 0x52, 0x03, 0x61, 0x67, 0x65, 0x12, 0x10, 0x0a, 0x03, 0x73, 0x65, 0x78,
	0x18, 0x05, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x73, 0x65, 0x78, 0x12, 0x1c, 0x0a, 0x09, 0x63,
	0x72, 0x65, 0x61, 0x74, 0x65, 0x64, 0x41, 0x74, 0x18, 0x06, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09,
	0x63, 0x72, 0x65, 0x61, 0x74, 0x65, 0x64, 0x41, 0x74, 0x22, 0x80, 0x01, 0x0a, 0x0e, 0x47, 0x65,
	0x74, 0x44, 0x61, 0x74, 0x61, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x17, 0x0a, 0x07,
	0x75, 0x73, 0x65, 0x72, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x05, 0x52, 0x06, 0x75,
	0x73, 0x65, 0x72, 0x49, 0x64, 0x12, 0x19, 0x0a, 0x08, 0x61, 0x75, 0x74, 0x68, 0x5f, 0x6b, 0x65,
	0x79, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x07, 0x61, 0x75, 0x74, 0x68, 0x4b, 0x65, 0x79,
	0x12, 0x1d, 0x0a, 0x0a, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x5f, 0x69, 0x64, 0x18, 0x03,
	0x20, 0x01, 0x28, 0x05, 0x52, 0x09, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x49, 0x64, 0x12,
	0x1b, 0x0a, 0x09, 0x72, 0x65, 0x64, 0x69, 0x73, 0x5f, 0x6b, 0x65, 0x79, 0x18, 0x04, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x08, 0x72, 0x65, 0x64, 0x69, 0x73, 0x4b, 0x65, 0x79, 0x22, 0x5b, 0x0a, 0x0f,
	0x47, 0x65, 0x74, 0x44, 0x61, 0x74, 0x61, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12,
	0x29, 0x0a, 0x0c, 0x72, 0x65, 0x74, 0x75, 0x72, 0x6e, 0x5f, 0x75, 0x73, 0x65, 0x72, 0x73, 0x18,
	0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x06, 0x2e, 0x55, 0x73, 0x65, 0x72, 0x72, 0x52, 0x0b, 0x72,
	0x65, 0x74, 0x75, 0x72, 0x6e, 0x55, 0x73, 0x65, 0x72, 0x73, 0x12, 0x1d, 0x0a, 0x0a, 0x6d, 0x65,
	0x73, 0x73, 0x61, 0x67, 0x65, 0x5f, 0x69, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x05, 0x52, 0x09,
	0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x49, 0x64, 0x32, 0x3b, 0x0a, 0x09, 0x67, 0x65, 0x74,
	0x5f, 0x75, 0x73, 0x65, 0x72, 0x73, 0x12, 0x2e, 0x0a, 0x07, 0x47, 0x65, 0x74, 0x44, 0x61, 0x74,
	0x61, 0x12, 0x0f, 0x2e, 0x47, 0x65, 0x74, 0x44, 0x61, 0x74, 0x61, 0x52, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x1a, 0x10, 0x2e, 0x47, 0x65, 0x74, 0x44, 0x61, 0x74, 0x61, 0x52, 0x65, 0x73, 0x70,
	0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x42, 0x0d, 0x5a, 0x0b, 0x67, 0x65, 0x74, 0x5f, 0x75, 0x73,
	0x65, 0x72, 0x2f, 0x70, 0x62, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_get_users_proto_rawDescOnce sync.Once
	file_get_users_proto_rawDescData = file_get_users_proto_rawDesc
)

func file_get_users_proto_rawDescGZIP() []byte {
	file_get_users_proto_rawDescOnce.Do(func() {
		file_get_users_proto_rawDescData = protoimpl.X.CompressGZIP(file_get_users_proto_rawDescData)
	})
	return file_get_users_proto_rawDescData
}

var file_get_users_proto_msgTypes = make([]protoimpl.MessageInfo, 3)
var file_get_users_proto_goTypes = []interface{}{
	(*Userr)(nil),           // 0: Userr
	(*GetDataRequest)(nil),  // 1: GetDataRequest
	(*GetDataResponse)(nil), // 2: GetDataResponse
}
var file_get_users_proto_depIdxs = []int32{
	0, // 0: GetDataResponse.return_users:type_name -> Userr
	1, // 1: get_users.GetData:input_type -> GetDataRequest
	2, // 2: get_users.GetData:output_type -> GetDataResponse
	2, // [2:3] is the sub-list for method output_type
	1, // [1:2] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_get_users_proto_init() }
func file_get_users_proto_init() {
	if File_get_users_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_get_users_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Userr); i {
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
		file_get_users_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetDataRequest); i {
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
		file_get_users_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetDataResponse); i {
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
			RawDescriptor: file_get_users_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   3,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_get_users_proto_goTypes,
		DependencyIndexes: file_get_users_proto_depIdxs,
		MessageInfos:      file_get_users_proto_msgTypes,
	}.Build()
	File_get_users_proto = out.File
	file_get_users_proto_rawDesc = nil
	file_get_users_proto_goTypes = nil
	file_get_users_proto_depIdxs = nil
}
