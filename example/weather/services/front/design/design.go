package design

import (
	. "goa.design/goa/v3/dsl"
	_ "goa.design/plugins/v3/otel"

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

var TestResult = Type("TestResult", func() {
	Description("Test result for a single test")
	Field(1, "name", String, "Name of the test", func() {
		Example("TestValidIP")
	})
	Field(2, "passed", Boolean, "Status of the test", func() {
		Example(true)
	})
	Field(3, "error", String, "Error message if the test failed", func() {
		Example("error getting location for valid ip: %v")
	})
	Field(4, "duration", Int64, "Duration of the test in ms", func() {
		Example(1234)
	})
	Required("name", "passed", "duration")
})

var TestCollection = Type("TestCollection", func() {
	Description("Collection of test results for grouping by service")
	Field(1, "name", String, "Name of the test collection", func() {
		Example("Locator Tests")
	})
	Field(2, "results", ArrayOf(TestResult), "Test results")
	Field(3, "duration", Int64, "Duration of the tests in ms", func() {
		Example(1234)
	})
	Field(4, "pass_count", Int, "Number of tests that passed", func() {
		Example(12)
	})
	Field(5, "fail_count", Int, "Number of tests that failed", func() {
		Example(1)
	})
	Required("name", "duration", "pass_count", "fail_count")
})

var TestResults = Type("TestResults", func() {
	Description("Test results for the integration tests")
	Field(1, "collections", ArrayOf(TestCollection), "Test collections")
	Field(2, "duration", Int64, "Duration of the tests in ms", func() {
		Example(1234)
	})
	Field(3, "pass_count", Int, "Number of tests that passed", func() {
		Example(12)
	})
	Field(4, "fail_count", Int, "Number of tests that failed", func() {
		Example(1)
	})
	Required("collections", "duration", "pass_count", "fail_count")
})
