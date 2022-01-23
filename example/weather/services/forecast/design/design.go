package design

import (
	. "goa.design/goa/v3/dsl"

	. "github.com/crossnokaye/micro/example/weather/design"
)

var _ = API("Weather Service API", func() {
	Title("The Weather Service API")
	Description("A fully instrumented weather service API")
})

var _ = Service("Forecast", func() {
	Description("Service that provides weather forecasts")
	Method("forecast", func() {
		Description("Retrieve weather forecast for a given location")
		Payload(func() {
			Field(1, "lat", Float64, "Latitude", func() {
				Example(37.8267)
			})
			Field(2, "long", Float64, "Longitude", func() {
				Example(-122.4233)
			})
			Required("lat", "long")
		})
		Result(Forecast)
		GRPC(func() {})
	})
})
