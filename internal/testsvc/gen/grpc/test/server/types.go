// Code generated by goa v3.8.3, DO NOT EDIT.
//
// test gRPC server types
//
// Command:
// $ goa gen goa.design/clue/internal/testsvc/design

package server

import (
	testpb "goa.design/clue/internal/testsvc/gen/grpc/test/pb"
	test "goa.design/clue/internal/testsvc/gen/test"
)

// NewGrpcMethodPayload builds the payload of the "grpc_method" endpoint of the
// "test" service from the gRPC request type.
func NewGrpcMethodPayload(message *testpb.GrpcMethodRequest) *test.Fields {
	v := &test.Fields{}
	if message.S != "" {
		v.S = &message.S
	}
	if message.I != 0 {
		iptr := int(message.I)
		v.I = &iptr
	}
	return v
}

// NewProtoGrpcMethodResponse builds the gRPC response type from the result of
// the "grpc_method" endpoint of the "test" service.
func NewProtoGrpcMethodResponse(result *test.Fields) *testpb.GrpcMethodResponse {
	message := &testpb.GrpcMethodResponse{}
	if result.S != nil {
		message.S = *result.S
	}
	if result.I != nil {
		message.I = int32(*result.I)
	}
	return message
}

// NewProtoGrpcStreamResponse builds the gRPC response type from the result of
// the "grpc_stream" endpoint of the "test" service.
func NewProtoGrpcStreamResponse(result *test.Fields) *testpb.GrpcStreamResponse {
	message := &testpb.GrpcStreamResponse{}
	if result.S != nil {
		message.S = *result.S
	}
	if result.I != nil {
		message.I = int32(*result.I)
	}
	return message
}

func NewFields(v *testpb.GrpcStreamStreamingRequest) *test.Fields {
	spayload := &test.Fields{}
	if v.S != "" {
		spayload.S = &v.S
	}
	if v.I != 0 {
		iptr := int(v.I)
		spayload.I = &iptr
	}
	return spayload
}
