// Code generated by goa v3.5.3, DO NOT EDIT.
//
// test gRPC server encoders and decoders
//
// Command:
// $ goa gen github.com/crossnokaye/micro/instrument/testsvc/design

package server

import (
	"context"

	testpb "github.com/crossnokaye/micro/instrument/testsvc/gen/grpc/test/pb"
	test "github.com/crossnokaye/micro/instrument/testsvc/gen/test"
	goagrpc "goa.design/goa/v3/grpc"
	"google.golang.org/grpc/metadata"
)

// EncodeGrpcMethodResponse encodes responses from the "test" service
// "grpc_method" endpoint.
func EncodeGrpcMethodResponse(ctx context.Context, v interface{}, hdr, trlr *metadata.MD) (interface{}, error) {
	result, ok := v.(*test.Fields)
	if !ok {
		return nil, goagrpc.ErrInvalidType("test", "grpc_method", "*test.Fields", v)
	}
	resp := NewGrpcMethodResponse(result)
	return resp, nil
}

// DecodeGrpcMethodRequest decodes requests sent to "test" service
// "grpc_method" endpoint.
func DecodeGrpcMethodRequest(ctx context.Context, v interface{}, md metadata.MD) (interface{}, error) {
	var (
		message *testpb.GrpcMethodRequest
		ok      bool
	)
	{
		if message, ok = v.(*testpb.GrpcMethodRequest); !ok {
			return nil, goagrpc.ErrInvalidType("test", "grpc_method", "*testpb.GrpcMethodRequest", v)
		}
	}
	var payload *test.Fields
	{
		payload = NewGrpcMethodPayload(message)
	}
	return payload, nil
}

// EncodeGrpcStreamResponse encodes responses from the "test" service
// "grpc_stream" endpoint.
func EncodeGrpcStreamResponse(ctx context.Context, v interface{}, hdr, trlr *metadata.MD) (interface{}, error) {
	result, ok := v.(*test.Fields)
	if !ok {
		return nil, goagrpc.ErrInvalidType("test", "grpc_stream", "*test.Fields", v)
	}
	resp := NewGrpcStreamResponse(result)
	return resp, nil
}
