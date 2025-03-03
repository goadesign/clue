// Code generated by goa v3.20.0, DO NOT EDIT.
//
// locator gRPC client encoders and decoders
//
// Command:
// $ goa gen goa.design/clue/example/weather/services/locator/design -o
// services/locator

package client

import (
	"context"

	locatorpb "goa.design/clue/example/weather/services/locator/gen/grpc/locator/pb"
	goagrpc "goa.design/goa/v3/grpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

// BuildGetLocationFunc builds the remote method to invoke for "locator"
// service "get_location" endpoint.
func BuildGetLocationFunc(grpccli locatorpb.LocatorClient, cliopts ...grpc.CallOption) goagrpc.RemoteFunc {
	return func(ctx context.Context, reqpb any, opts ...grpc.CallOption) (any, error) {
		for _, opt := range cliopts {
			opts = append(opts, opt)
		}
		if reqpb != nil {
			return grpccli.GetLocation(ctx, reqpb.(*locatorpb.GetLocationRequest), opts...)
		}
		return grpccli.GetLocation(ctx, &locatorpb.GetLocationRequest{}, opts...)
	}
}

// EncodeGetLocationRequest encodes requests sent to locator get_location
// endpoint.
func EncodeGetLocationRequest(ctx context.Context, v any, md *metadata.MD) (any, error) {
	payload, ok := v.(string)
	if !ok {
		return nil, goagrpc.ErrInvalidType("locator", "get_location", "string", v)
	}
	return NewProtoGetLocationRequest(payload), nil
}

// DecodeGetLocationResponse decodes responses from the locator get_location
// endpoint.
func DecodeGetLocationResponse(ctx context.Context, v any, hdr, trlr metadata.MD) (any, error) {
	message, ok := v.(*locatorpb.GetLocationResponse)
	if !ok {
		return nil, goagrpc.ErrInvalidType("locator", "get_location", "*locatorpb.GetLocationResponse", v)
	}
	res := NewGetLocationResult(message)
	return res, nil
}
