package design

import (
	. "goa.design/goa/v3/dsl"

	. "goa.design/clue/example/weather/design"
)

var _ = API("Weather", func() {
	Title("Weather Forecast Service API")
	Description("The weather forecast service API produces weather forecasts from US-based IPs. It uses IP location to find the appropriate weather station.")
	Version("1.0.0")
})

var _ = Service("front", func() {
	Description("Public HTTP frontend")
	Method("forecast", func() {
		Description("Retrieve weather forecast for given IP")
		Payload(String, func() {
			Format(FormatIP)
		})
		Result(Forecast)
		Error("not_usa", ErrorResult, "IP address is not in the US")
		HTTP(func() {
			GET("/forecast/{ip}")
			Response("not_usa", StatusBadRequest)
		})
	})
})
