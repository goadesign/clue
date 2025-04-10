// Code generated by goa v3.20.0, DO NOT EDIT.
//
// locator gRPC server types
//
// Command:
// $ goa gen goa.design/clue/example/weather/services/locator/design -o
// services/locator

package server

import (
	locatorpb "goa.design/clue/example/weather/services/locator/gen/grpc/locator/pb"
	locator "goa.design/clue/example/weather/services/locator/gen/locator"
)

// NewGetLocationPayload builds the payload of the "get_location" endpoint of
// the "locator" service from the gRPC request type.
func NewGetLocationPayload(message *locatorpb.GetLocationRequest) string {
	v := message.Field
	return v
}

// NewProtoGetLocationResponse builds the gRPC response type from the result of
// the "get_location" endpoint of the "locator" service.
func NewProtoGetLocationResponse(result *locator.WorldLocation) *locatorpb.GetLocationResponse {
	message := &locatorpb.GetLocationResponse{
		Lat:     result.Lat,
		Long:    result.Long,
		City:    result.City,
		Region:  result.Region,
		Country: result.Country,
	}
	return message
}
