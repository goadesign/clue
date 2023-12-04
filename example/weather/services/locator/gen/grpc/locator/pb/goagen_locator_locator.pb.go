// Code generated with goa v3.14.0, DO NOT EDIT.
//
// locator protocol buffer definition
//
// Command:
// $ goa gen goa.design/clue/example/weather/services/locator/design -o
// services/locator

// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.27.1
// 	protoc        v4.25.1
// source: goagen_locator_locator.proto

package locatorpb

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

type GetLocationRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Field string `protobuf:"bytes,1,opt,name=field,proto3" json:"field,omitempty"`
}

func (x *GetLocationRequest) Reset() {
	*x = GetLocationRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_goagen_locator_locator_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetLocationRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetLocationRequest) ProtoMessage() {}

func (x *GetLocationRequest) ProtoReflect() protoreflect.Message {
	mi := &file_goagen_locator_locator_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetLocationRequest.ProtoReflect.Descriptor instead.
func (*GetLocationRequest) Descriptor() ([]byte, []int) {
	return file_goagen_locator_locator_proto_rawDescGZIP(), []int{0}
}

func (x *GetLocationRequest) GetField() string {
	if x != nil {
		return x.Field
	}
	return ""
}

type GetLocationResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Latitude
	Lat float64 `protobuf:"fixed64,1,opt,name=lat,proto3" json:"lat,omitempty"`
	// Longitude
	Long float64 `protobuf:"fixed64,2,opt,name=long,proto3" json:"long,omitempty"`
	// City
	City string `protobuf:"bytes,3,opt,name=city,proto3" json:"city,omitempty"`
	// State, region etc.
	Region string `protobuf:"bytes,4,opt,name=region,proto3" json:"region,omitempty"`
	// Country
	Country string `protobuf:"bytes,5,opt,name=country,proto3" json:"country,omitempty"`
}

func (x *GetLocationResponse) Reset() {
	*x = GetLocationResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_goagen_locator_locator_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetLocationResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetLocationResponse) ProtoMessage() {}

func (x *GetLocationResponse) ProtoReflect() protoreflect.Message {
	mi := &file_goagen_locator_locator_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetLocationResponse.ProtoReflect.Descriptor instead.
func (*GetLocationResponse) Descriptor() ([]byte, []int) {
	return file_goagen_locator_locator_proto_rawDescGZIP(), []int{1}
}

func (x *GetLocationResponse) GetLat() float64 {
	if x != nil {
		return x.Lat
	}
	return 0
}

func (x *GetLocationResponse) GetLong() float64 {
	if x != nil {
		return x.Long
	}
	return 0
}

func (x *GetLocationResponse) GetCity() string {
	if x != nil {
		return x.City
	}
	return ""
}

func (x *GetLocationResponse) GetRegion() string {
	if x != nil {
		return x.Region
	}
	return ""
}

func (x *GetLocationResponse) GetCountry() string {
	if x != nil {
		return x.Country
	}
	return ""
}

var File_goagen_locator_locator_proto protoreflect.FileDescriptor

var file_goagen_locator_locator_proto_rawDesc = []byte{
	0x0a, 0x1c, 0x67, 0x6f, 0x61, 0x67, 0x65, 0x6e, 0x5f, 0x6c, 0x6f, 0x63, 0x61, 0x74, 0x6f, 0x72,
	0x5f, 0x6c, 0x6f, 0x63, 0x61, 0x74, 0x6f, 0x72, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x07,
	0x6c, 0x6f, 0x63, 0x61, 0x74, 0x6f, 0x72, 0x22, 0x2a, 0x0a, 0x12, 0x47, 0x65, 0x74, 0x4c, 0x6f,
	0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x14, 0x0a,
	0x05, 0x66, 0x69, 0x65, 0x6c, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x66, 0x69,
	0x65, 0x6c, 0x64, 0x22, 0x81, 0x01, 0x0a, 0x13, 0x47, 0x65, 0x74, 0x4c, 0x6f, 0x63, 0x61, 0x74,
	0x69, 0x6f, 0x6e, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x10, 0x0a, 0x03, 0x6c,
	0x61, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x01, 0x52, 0x03, 0x6c, 0x61, 0x74, 0x12, 0x12, 0x0a,
	0x04, 0x6c, 0x6f, 0x6e, 0x67, 0x18, 0x02, 0x20, 0x01, 0x28, 0x01, 0x52, 0x04, 0x6c, 0x6f, 0x6e,
	0x67, 0x12, 0x12, 0x0a, 0x04, 0x63, 0x69, 0x74, 0x79, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x04, 0x63, 0x69, 0x74, 0x79, 0x12, 0x16, 0x0a, 0x06, 0x72, 0x65, 0x67, 0x69, 0x6f, 0x6e, 0x18,
	0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x72, 0x65, 0x67, 0x69, 0x6f, 0x6e, 0x12, 0x18, 0x0a,
	0x07, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x72, 0x79, 0x18, 0x05, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07,
	0x63, 0x6f, 0x75, 0x6e, 0x74, 0x72, 0x79, 0x32, 0x53, 0x0a, 0x07, 0x4c, 0x6f, 0x63, 0x61, 0x74,
	0x6f, 0x72, 0x12, 0x48, 0x0a, 0x0b, 0x47, 0x65, 0x74, 0x4c, 0x6f, 0x63, 0x61, 0x74, 0x69, 0x6f,
	0x6e, 0x12, 0x1b, 0x2e, 0x6c, 0x6f, 0x63, 0x61, 0x74, 0x6f, 0x72, 0x2e, 0x47, 0x65, 0x74, 0x4c,
	0x6f, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1c,
	0x2e, 0x6c, 0x6f, 0x63, 0x61, 0x74, 0x6f, 0x72, 0x2e, 0x47, 0x65, 0x74, 0x4c, 0x6f, 0x63, 0x61,
	0x74, 0x69, 0x6f, 0x6e, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x42, 0x0c, 0x5a, 0x0a,
	0x2f, 0x6c, 0x6f, 0x63, 0x61, 0x74, 0x6f, 0x72, 0x70, 0x62, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x33,
}

var (
	file_goagen_locator_locator_proto_rawDescOnce sync.Once
	file_goagen_locator_locator_proto_rawDescData = file_goagen_locator_locator_proto_rawDesc
)

func file_goagen_locator_locator_proto_rawDescGZIP() []byte {
	file_goagen_locator_locator_proto_rawDescOnce.Do(func() {
		file_goagen_locator_locator_proto_rawDescData = protoimpl.X.CompressGZIP(file_goagen_locator_locator_proto_rawDescData)
	})
	return file_goagen_locator_locator_proto_rawDescData
}

var file_goagen_locator_locator_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_goagen_locator_locator_proto_goTypes = []interface{}{
	(*GetLocationRequest)(nil),  // 0: locator.GetLocationRequest
	(*GetLocationResponse)(nil), // 1: locator.GetLocationResponse
}
var file_goagen_locator_locator_proto_depIdxs = []int32{
	0, // 0: locator.Locator.GetLocation:input_type -> locator.GetLocationRequest
	1, // 1: locator.Locator.GetLocation:output_type -> locator.GetLocationResponse
	1, // [1:2] is the sub-list for method output_type
	0, // [0:1] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_goagen_locator_locator_proto_init() }
func file_goagen_locator_locator_proto_init() {
	if File_goagen_locator_locator_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_goagen_locator_locator_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetLocationRequest); i {
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
		file_goagen_locator_locator_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetLocationResponse); i {
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
			RawDescriptor: file_goagen_locator_locator_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_goagen_locator_locator_proto_goTypes,
		DependencyIndexes: file_goagen_locator_locator_proto_depIdxs,
		MessageInfos:      file_goagen_locator_locator_proto_msgTypes,
	}.Build()
	File_goagen_locator_locator_proto = out.File
	file_goagen_locator_locator_proto_rawDesc = nil
	file_goagen_locator_locator_proto_goTypes = nil
	file_goagen_locator_locator_proto_depIdxs = nil
}
