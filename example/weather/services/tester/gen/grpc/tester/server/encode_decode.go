// Code generated by goa v3.20.0, DO NOT EDIT.
//
// tester gRPC server encoders and decoders
//
// Command:
// $ goa gen goa.design/clue/example/weather/services/tester/design -o
// services/tester

package server

import (
	"context"

	testerpb "goa.design/clue/example/weather/services/tester/gen/grpc/tester/pb"
	tester "goa.design/clue/example/weather/services/tester/gen/tester"
	goagrpc "goa.design/goa/v3/grpc"
	"google.golang.org/grpc/metadata"
)

// EncodeTestAllResponse encodes responses from the "tester" service "test_all"
// endpoint.
func EncodeTestAllResponse(ctx context.Context, v any, hdr, trlr *metadata.MD) (any, error) {
	result, ok := v.(*tester.TestResults)
	if !ok {
		return nil, goagrpc.ErrInvalidType("tester", "test_all", "*tester.TestResults", v)
	}
	resp := NewProtoTestAllResponse(result)
	return resp, nil
}

// DecodeTestAllRequest decodes requests sent to "tester" service "test_all"
// endpoint.
func DecodeTestAllRequest(ctx context.Context, v any, md metadata.MD) (any, error) {
	var (
		message *testerpb.TestAllRequest
		ok      bool
	)
	{
		if message, ok = v.(*testerpb.TestAllRequest); !ok {
			return nil, goagrpc.ErrInvalidType("tester", "test_all", "*testerpb.TestAllRequest", v)
		}
	}
	var payload *tester.TesterPayload
	{
		payload = NewTestAllPayload(message)
	}
	return payload, nil
}

// EncodeTestSmokeResponse encodes responses from the "tester" service
// "test_smoke" endpoint.
func EncodeTestSmokeResponse(ctx context.Context, v any, hdr, trlr *metadata.MD) (any, error) {
	result, ok := v.(*tester.TestResults)
	if !ok {
		return nil, goagrpc.ErrInvalidType("tester", "test_smoke", "*tester.TestResults", v)
	}
	resp := NewProtoTestSmokeResponse(result)
	return resp, nil
}

// EncodeTestForecasterResponse encodes responses from the "tester" service
// "test_forecaster" endpoint.
func EncodeTestForecasterResponse(ctx context.Context, v any, hdr, trlr *metadata.MD) (any, error) {
	result, ok := v.(*tester.TestResults)
	if !ok {
		return nil, goagrpc.ErrInvalidType("tester", "test_forecaster", "*tester.TestResults", v)
	}
	resp := NewProtoTestForecasterResponse(result)
	return resp, nil
}

// EncodeTestLocatorResponse encodes responses from the "tester" service
// "test_locator" endpoint.
func EncodeTestLocatorResponse(ctx context.Context, v any, hdr, trlr *metadata.MD) (any, error) {
	result, ok := v.(*tester.TestResults)
	if !ok {
		return nil, goagrpc.ErrInvalidType("tester", "test_locator", "*tester.TestResults", v)
	}
	resp := NewProtoTestLocatorResponse(result)
	return resp, nil
}
