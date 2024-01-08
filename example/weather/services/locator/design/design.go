package design

import (
	. "goa.design/goa/v3/dsl"
)

var _ = API("IP Location API", func() {
	Title("IP Location Service API")
	Description("A fully instrumented IP location service API")
	Version("1.0.0")
})

var _ = Service("locator", func() {
	Description("Public HTTP frontend")

	Method("get_location", func() {
		Description("Retrieve location information for a given IP address")
		Payload(String, "IP address", func() {
			Format(FormatIP)
		})
		Result(WorldLocation)
		GRPC(func() {})
	})
})

var WorldLocation = Type("WorldLocation", func() {
	Description("Geographical location")
	Field(1, "lat", Float64, "Latitude", func() {
		Example(37.8267)
	})
	Field(2, "long", Float64, "Longitude", func() {
		Example(-122.4233)
	})
	Field(3, "city", String, "City", func() {
		Example("San Francisco")
	})
	Field(4, "region", String, "State, region etc.", func() {
		Example("CA")
	})
	Field(5, "country", String, "Country", func() {
		Example("USA")
	})
	Required("lat", "long", "city", "region", "country")
})
