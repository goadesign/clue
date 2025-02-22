// Code generated with goa v3.20.0, DO NOT EDIT.
//
// Forecaster protocol buffer definition
//
// Command:
// $ goa gen goa.design/clue/example/weather/services/forecaster/design -o
// services/forecaster

// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.5
// 	protoc        v5.29.3
// source: goagen_forecaster_forecaster.proto

package forecasterpb

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

type ForecastRequest struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	// Latitude
	Lat float64 `protobuf:"fixed64,1,opt,name=lat,proto3" json:"lat,omitempty"`
	// Longitude
	Long          float64 `protobuf:"fixed64,2,opt,name=long,proto3" json:"long,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *ForecastRequest) Reset() {
	*x = ForecastRequest{}
	mi := &file_goagen_forecaster_forecaster_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ForecastRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ForecastRequest) ProtoMessage() {}

func (x *ForecastRequest) ProtoReflect() protoreflect.Message {
	mi := &file_goagen_forecaster_forecaster_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ForecastRequest.ProtoReflect.Descriptor instead.
func (*ForecastRequest) Descriptor() ([]byte, []int) {
	return file_goagen_forecaster_forecaster_proto_rawDescGZIP(), []int{0}
}

func (x *ForecastRequest) GetLat() float64 {
	if x != nil {
		return x.Lat
	}
	return 0
}

func (x *ForecastRequest) GetLong() float64 {
	if x != nil {
		return x.Long
	}
	return 0
}

type ForecastResponse struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	// Forecast location
	Location *Location `protobuf:"bytes,1,opt,name=location,proto3" json:"location,omitempty"`
	// Weather forecast periods
	Periods       []*Period `protobuf:"bytes,2,rep,name=periods,proto3" json:"periods,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *ForecastResponse) Reset() {
	*x = ForecastResponse{}
	mi := &file_goagen_forecaster_forecaster_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ForecastResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ForecastResponse) ProtoMessage() {}

func (x *ForecastResponse) ProtoReflect() protoreflect.Message {
	mi := &file_goagen_forecaster_forecaster_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ForecastResponse.ProtoReflect.Descriptor instead.
func (*ForecastResponse) Descriptor() ([]byte, []int) {
	return file_goagen_forecaster_forecaster_proto_rawDescGZIP(), []int{1}
}

func (x *ForecastResponse) GetLocation() *Location {
	if x != nil {
		return x.Location
	}
	return nil
}

func (x *ForecastResponse) GetPeriods() []*Period {
	if x != nil {
		return x.Periods
	}
	return nil
}

// Geographical location
type Location struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	// Latitude
	Lat float64 `protobuf:"fixed64,1,opt,name=lat,proto3" json:"lat,omitempty"`
	// Longitude
	Long float64 `protobuf:"fixed64,2,opt,name=long,proto3" json:"long,omitempty"`
	// City
	City string `protobuf:"bytes,3,opt,name=city,proto3" json:"city,omitempty"`
	// State
	State         string `protobuf:"bytes,4,opt,name=state,proto3" json:"state,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *Location) Reset() {
	*x = Location{}
	mi := &file_goagen_forecaster_forecaster_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Location) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Location) ProtoMessage() {}

func (x *Location) ProtoReflect() protoreflect.Message {
	mi := &file_goagen_forecaster_forecaster_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Location.ProtoReflect.Descriptor instead.
func (*Location) Descriptor() ([]byte, []int) {
	return file_goagen_forecaster_forecaster_proto_rawDescGZIP(), []int{2}
}

func (x *Location) GetLat() float64 {
	if x != nil {
		return x.Lat
	}
	return 0
}

func (x *Location) GetLong() float64 {
	if x != nil {
		return x.Long
	}
	return 0
}

func (x *Location) GetCity() string {
	if x != nil {
		return x.City
	}
	return ""
}

func (x *Location) GetState() string {
	if x != nil {
		return x.State
	}
	return ""
}

// Weather forecast period
type Period struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	// Period name
	Name string `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	// Start time
	StartTime string `protobuf:"bytes,2,opt,name=start_time,json=startTime,proto3" json:"start_time,omitempty"`
	// End time
	EndTime string `protobuf:"bytes,3,opt,name=end_time,json=endTime,proto3" json:"end_time,omitempty"`
	// Temperature
	Temperature int32 `protobuf:"zigzag32,4,opt,name=temperature,proto3" json:"temperature,omitempty"`
	// Temperature unit
	TemperatureUnit string `protobuf:"bytes,5,opt,name=temperature_unit,json=temperatureUnit,proto3" json:"temperature_unit,omitempty"`
	// Summary
	Summary       string `protobuf:"bytes,6,opt,name=summary,proto3" json:"summary,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *Period) Reset() {
	*x = Period{}
	mi := &file_goagen_forecaster_forecaster_proto_msgTypes[3]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Period) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Period) ProtoMessage() {}

func (x *Period) ProtoReflect() protoreflect.Message {
	mi := &file_goagen_forecaster_forecaster_proto_msgTypes[3]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Period.ProtoReflect.Descriptor instead.
func (*Period) Descriptor() ([]byte, []int) {
	return file_goagen_forecaster_forecaster_proto_rawDescGZIP(), []int{3}
}

func (x *Period) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *Period) GetStartTime() string {
	if x != nil {
		return x.StartTime
	}
	return ""
}

func (x *Period) GetEndTime() string {
	if x != nil {
		return x.EndTime
	}
	return ""
}

func (x *Period) GetTemperature() int32 {
	if x != nil {
		return x.Temperature
	}
	return 0
}

func (x *Period) GetTemperatureUnit() string {
	if x != nil {
		return x.TemperatureUnit
	}
	return ""
}

func (x *Period) GetSummary() string {
	if x != nil {
		return x.Summary
	}
	return ""
}

var File_goagen_forecaster_forecaster_proto protoreflect.FileDescriptor

var file_goagen_forecaster_forecaster_proto_rawDesc = string([]byte{
	0x0a, 0x22, 0x67, 0x6f, 0x61, 0x67, 0x65, 0x6e, 0x5f, 0x66, 0x6f, 0x72, 0x65, 0x63, 0x61, 0x73,
	0x74, 0x65, 0x72, 0x5f, 0x66, 0x6f, 0x72, 0x65, 0x63, 0x61, 0x73, 0x74, 0x65, 0x72, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x12, 0x0a, 0x66, 0x6f, 0x72, 0x65, 0x63, 0x61, 0x73, 0x74, 0x65, 0x72,
	0x22, 0x37, 0x0a, 0x0f, 0x46, 0x6f, 0x72, 0x65, 0x63, 0x61, 0x73, 0x74, 0x52, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x12, 0x10, 0x0a, 0x03, 0x6c, 0x61, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x01,
	0x52, 0x03, 0x6c, 0x61, 0x74, 0x12, 0x12, 0x0a, 0x04, 0x6c, 0x6f, 0x6e, 0x67, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x01, 0x52, 0x04, 0x6c, 0x6f, 0x6e, 0x67, 0x22, 0x72, 0x0a, 0x10, 0x46, 0x6f, 0x72,
	0x65, 0x63, 0x61, 0x73, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x30, 0x0a,
	0x08, 0x6c, 0x6f, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32,
	0x14, 0x2e, 0x66, 0x6f, 0x72, 0x65, 0x63, 0x61, 0x73, 0x74, 0x65, 0x72, 0x2e, 0x4c, 0x6f, 0x63,
	0x61, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x08, 0x6c, 0x6f, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x12,
	0x2c, 0x0a, 0x07, 0x70, 0x65, 0x72, 0x69, 0x6f, 0x64, 0x73, 0x18, 0x02, 0x20, 0x03, 0x28, 0x0b,
	0x32, 0x12, 0x2e, 0x66, 0x6f, 0x72, 0x65, 0x63, 0x61, 0x73, 0x74, 0x65, 0x72, 0x2e, 0x50, 0x65,
	0x72, 0x69, 0x6f, 0x64, 0x52, 0x07, 0x70, 0x65, 0x72, 0x69, 0x6f, 0x64, 0x73, 0x22, 0x5a, 0x0a,
	0x08, 0x4c, 0x6f, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x10, 0x0a, 0x03, 0x6c, 0x61, 0x74,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x01, 0x52, 0x03, 0x6c, 0x61, 0x74, 0x12, 0x12, 0x0a, 0x04, 0x6c,
	0x6f, 0x6e, 0x67, 0x18, 0x02, 0x20, 0x01, 0x28, 0x01, 0x52, 0x04, 0x6c, 0x6f, 0x6e, 0x67, 0x12,
	0x12, 0x0a, 0x04, 0x63, 0x69, 0x74, 0x79, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x63,
	0x69, 0x74, 0x79, 0x12, 0x14, 0x0a, 0x05, 0x73, 0x74, 0x61, 0x74, 0x65, 0x18, 0x04, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x05, 0x73, 0x74, 0x61, 0x74, 0x65, 0x22, 0xbd, 0x01, 0x0a, 0x06, 0x50, 0x65,
	0x72, 0x69, 0x6f, 0x64, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x1d, 0x0a, 0x0a, 0x73, 0x74, 0x61, 0x72,
	0x74, 0x5f, 0x74, 0x69, 0x6d, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x73, 0x74,
	0x61, 0x72, 0x74, 0x54, 0x69, 0x6d, 0x65, 0x12, 0x19, 0x0a, 0x08, 0x65, 0x6e, 0x64, 0x5f, 0x74,
	0x69, 0x6d, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x65, 0x6e, 0x64, 0x54, 0x69,
	0x6d, 0x65, 0x12, 0x20, 0x0a, 0x0b, 0x74, 0x65, 0x6d, 0x70, 0x65, 0x72, 0x61, 0x74, 0x75, 0x72,
	0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x11, 0x52, 0x0b, 0x74, 0x65, 0x6d, 0x70, 0x65, 0x72, 0x61,
	0x74, 0x75, 0x72, 0x65, 0x12, 0x29, 0x0a, 0x10, 0x74, 0x65, 0x6d, 0x70, 0x65, 0x72, 0x61, 0x74,
	0x75, 0x72, 0x65, 0x5f, 0x75, 0x6e, 0x69, 0x74, 0x18, 0x05, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0f,
	0x74, 0x65, 0x6d, 0x70, 0x65, 0x72, 0x61, 0x74, 0x75, 0x72, 0x65, 0x55, 0x6e, 0x69, 0x74, 0x12,
	0x18, 0x0a, 0x07, 0x73, 0x75, 0x6d, 0x6d, 0x61, 0x72, 0x79, 0x18, 0x06, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x07, 0x73, 0x75, 0x6d, 0x6d, 0x61, 0x72, 0x79, 0x32, 0x53, 0x0a, 0x0a, 0x46, 0x6f, 0x72,
	0x65, 0x63, 0x61, 0x73, 0x74, 0x65, 0x72, 0x12, 0x45, 0x0a, 0x08, 0x46, 0x6f, 0x72, 0x65, 0x63,
	0x61, 0x73, 0x74, 0x12, 0x1b, 0x2e, 0x66, 0x6f, 0x72, 0x65, 0x63, 0x61, 0x73, 0x74, 0x65, 0x72,
	0x2e, 0x46, 0x6f, 0x72, 0x65, 0x63, 0x61, 0x73, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x1a, 0x1c, 0x2e, 0x66, 0x6f, 0x72, 0x65, 0x63, 0x61, 0x73, 0x74, 0x65, 0x72, 0x2e, 0x46, 0x6f,
	0x72, 0x65, 0x63, 0x61, 0x73, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x42, 0x0f,
	0x5a, 0x0d, 0x2f, 0x66, 0x6f, 0x72, 0x65, 0x63, 0x61, 0x73, 0x74, 0x65, 0x72, 0x70, 0x62, 0x62,
	0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
})

var (
	file_goagen_forecaster_forecaster_proto_rawDescOnce sync.Once
	file_goagen_forecaster_forecaster_proto_rawDescData []byte
)

func file_goagen_forecaster_forecaster_proto_rawDescGZIP() []byte {
	file_goagen_forecaster_forecaster_proto_rawDescOnce.Do(func() {
		file_goagen_forecaster_forecaster_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_goagen_forecaster_forecaster_proto_rawDesc), len(file_goagen_forecaster_forecaster_proto_rawDesc)))
	})
	return file_goagen_forecaster_forecaster_proto_rawDescData
}

var file_goagen_forecaster_forecaster_proto_msgTypes = make([]protoimpl.MessageInfo, 4)
var file_goagen_forecaster_forecaster_proto_goTypes = []any{
	(*ForecastRequest)(nil),  // 0: forecaster.ForecastRequest
	(*ForecastResponse)(nil), // 1: forecaster.ForecastResponse
	(*Location)(nil),         // 2: forecaster.Location
	(*Period)(nil),           // 3: forecaster.Period
}
var file_goagen_forecaster_forecaster_proto_depIdxs = []int32{
	2, // 0: forecaster.ForecastResponse.location:type_name -> forecaster.Location
	3, // 1: forecaster.ForecastResponse.periods:type_name -> forecaster.Period
	0, // 2: forecaster.Forecaster.Forecast:input_type -> forecaster.ForecastRequest
	1, // 3: forecaster.Forecaster.Forecast:output_type -> forecaster.ForecastResponse
	3, // [3:4] is the sub-list for method output_type
	2, // [2:3] is the sub-list for method input_type
	2, // [2:2] is the sub-list for extension type_name
	2, // [2:2] is the sub-list for extension extendee
	0, // [0:2] is the sub-list for field type_name
}

func init() { file_goagen_forecaster_forecaster_proto_init() }
func file_goagen_forecaster_forecaster_proto_init() {
	if File_goagen_forecaster_forecaster_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_goagen_forecaster_forecaster_proto_rawDesc), len(file_goagen_forecaster_forecaster_proto_rawDesc)),
			NumEnums:      0,
			NumMessages:   4,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_goagen_forecaster_forecaster_proto_goTypes,
		DependencyIndexes: file_goagen_forecaster_forecaster_proto_depIdxs,
		MessageInfos:      file_goagen_forecaster_forecaster_proto_msgTypes,
	}.Build()
	File_goagen_forecaster_forecaster_proto = out.File
	file_goagen_forecaster_forecaster_proto_goTypes = nil
	file_goagen_forecaster_forecaster_proto_depIdxs = nil
}
