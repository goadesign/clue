// Code generated by goa v3.5.3, DO NOT EDIT.
//
// test gRPC client types
//
// Command:
// $ goa gen github.com/crossnokaye/micro/instrument/testsvc/design

package client

import (
	testpb "github.com/crossnokaye/micro/instrument/testsvc/gen/grpc/test/pb"
	test "github.com/crossnokaye/micro/instrument/testsvc/gen/test"
)

// NewGrpcMethodRequest builds the gRPC request type from the payload of the
// "grpc_method" endpoint of the "test" service.
func NewGrpcMethodRequest(payload *test.Fields) *testpb.GrpcMethodRequest {
	message := &testpb.GrpcMethodRequest{}
	if payload.S != nil {
		message.S = *payload.S
	}
	if payload.I != nil {
		message.I = int32(*payload.I)
	}
	return message
}

// NewGrpcMethodResult builds the result type of the "grpc_method" endpoint of
// the "test" service from the gRPC response type.
func NewGrpcMethodResult(message *testpb.GrpcMethodResponse) *test.Fields {
	result := &test.Fields{}
	if message.S != "" {
		result.S = &message.S
	}
	if message.I != 0 {
		iptr := int(message.I)
		result.I = &iptr
	}
	return result
}

func NewFields(v *testpb.GrpcStreamingResponse) *test.Fields {
	result := &test.Fields{}
	if v.S != "" {
		result.S = &v.S
	}
	if v.I != 0 {
		iptr := int(v.I)
		result.I = &iptr
	}
	return result
}

func NewGrpcStreamingStreamingRequest(spayload *test.Fields) *testpb.GrpcStreamingStreamingRequest {
	v := &testpb.GrpcStreamingStreamingRequest{}
	if spayload.S != nil {
		v.S = *spayload.S
	}
	if spayload.I != nil {
		v.I = int32(*spayload.I)
	}
	return v
}
