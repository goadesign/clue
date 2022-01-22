// Code generated by goa v3.5.4, DO NOT EDIT.
//
// locator gRPC client
//
// Command:
// $ goa gen
// github.com/crossnokaye/micro/example/weather/services/locator/design -o
// example/weather/services/locator

package client

import (
	"context"

	locatorpb "github.com/crossnokaye/micro/example/weather/services/locator/gen/grpc/locator/pb"
	goagrpc "goa.design/goa/v3/grpc"
	goa "goa.design/goa/v3/pkg"
	"google.golang.org/grpc"
)

// Client lists the service endpoint gRPC clients.
type Client struct {
	grpccli locatorpb.LocatorClient
	opts    []grpc.CallOption
}

// NewClient instantiates gRPC client for all the locator service servers.
func NewClient(cc *grpc.ClientConn, opts ...grpc.CallOption) *Client {
	return &Client{
		grpccli: locatorpb.NewLocatorClient(cc),
		opts:    opts,
	}
}

// GetLocation calls the "GetLocation" function in locatorpb.LocatorClient
// interface.
func (c *Client) GetLocation() goa.Endpoint {
	return func(ctx context.Context, v interface{}) (interface{}, error) {
		inv := goagrpc.NewInvoker(
			BuildGetLocationFunc(c.grpccli, c.opts...),
			EncodeGetLocationRequest,
			DecodeGetLocationResponse)
		res, err := inv.Invoke(ctx, v)
		if err != nil {
			return nil, goa.Fault(err.Error())
		}
		return res, nil
	}
}
