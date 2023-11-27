package design

import (
	. "goa.design/goa/v3/dsl"

	. "goa.design/clue/example/weather/design"
)

var _ = API("Weather", func() {
	Title("Weather Forecast Service API")
	Description("The weather forecast service API produces weather forecasts from US-based IPs. It uses IP location to find the appropriate weather station.")
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
	Method("test_all", func() {
		Description("Endpoint for running ALL API Integration Tests for the Weather System, allowing for filtering on included or excluded tests")
		Payload(func() {
			Field(1, "include", ArrayOf(String), "Tests to run")
			Field(2, "exclude", ArrayOf(String), "Tests to exclude")
		})
		Result(TestResults)
		HTTP(func() {
			POST("/tester/all")
		})
	})
	Method("test_smoke", func() {
		Description("Endpoint for running API Integration Tests' Smoke Tests ONLY for the Weather System")
		Result(TestResults)
		HTTP(func() {
			POST("/tester/smoke")
		})
	})
})
