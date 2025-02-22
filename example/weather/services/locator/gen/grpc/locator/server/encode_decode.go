// Code generated by goa v3.20.0, DO NOT EDIT.
//
// locator gRPC server encoders and decoders
//
// Command:
// $ goa gen goa.design/clue/example/weather/services/locator/design -o
// services/locator

package server

import (
	"context"

	locatorpb "goa.design/clue/example/weather/services/locator/gen/grpc/locator/pb"
	locator "goa.design/clue/example/weather/services/locator/gen/locator"
	goagrpc "goa.design/goa/v3/grpc"
	"google.golang.org/grpc/metadata"
)

// EncodeGetLocationResponse encodes responses from the "locator" service
// "get_location" endpoint.
func EncodeGetLocationResponse(ctx context.Context, v any, hdr, trlr *metadata.MD) (any, error) {
	result, ok := v.(*locator.WorldLocation)
	if !ok {
		return nil, goagrpc.ErrInvalidType("locator", "get_location", "*locator.WorldLocation", v)
	}
	resp := NewProtoGetLocationResponse(result)
	return resp, nil
}

// DecodeGetLocationRequest decodes requests sent to "locator" service
// "get_location" endpoint.
func DecodeGetLocationRequest(ctx context.Context, v any, md metadata.MD) (any, error) {
	var (
		message *locatorpb.GetLocationRequest
		ok      bool
	)
	{
		if message, ok = v.(*locatorpb.GetLocationRequest); !ok {
			return nil, goagrpc.ErrInvalidType("locator", "get_location", "*locatorpb.GetLocationRequest", v)
		}
	}
	var payload string
	{
		payload = NewGetLocationPayload(message)
	}
	return payload, nil
}
