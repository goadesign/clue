package design

import (
	. "goa.design/goa/v3/dsl"
)

var _ = API("Tester Service API", func() {
	Title("The Tester Service API")
	Description("A fully instrumented tester service API")
})

// Not Found errors don't make sense for this because if a test is provided in an include or exclude
// it is simply ignored when it doesn't find it. This is not an error.
var _ = Service("tester", func() {
	GRPC(func() {
		// Override Package name to avoid conflicts with the generated code
		Package("weather_tester")
	})
	Description("The Weather System Tester Service is used to manage the integration testing of the weater system")
	Error("include_exclude_both", ErrorResult, "Cannot specify both include and exclude")
	Error("wildcard_compile_error", ErrorResult, "Wildcard glob did not compile")

	Method("test_all", func() {
		Description("Runs all tests in the iam system")
		Payload(SystemTestPayload)
		Result(TestResults)
		GRPC(func() {
			Response(CodeOK)
			Response("include_exclude_both", CodeInvalidArgument)
			Response("wildcard_compile_error", CodeInvalidArgument)
		})
	})

	Method("test_smoke", func() {
		Description("Runs smoke tests in the iam system")
		Result(TestResults)
		GRPC(func() {
			Response(CodeOK)
		})
	})

	Method("test_forecaster", func() {
		Description("Runs tests for the forecaster service")
		Result(TestResults)
		GRPC(func() {
			Response(CodeOK)
			Response("include_exclude_both", CodeInvalidArgument)
		})
	})

	Method("test_locator", func() {
		Description("Runs tests for the locator service")
		Result(TestResults)
		GRPC(func() {
			Response(CodeOK)
			Response("include_exclude_both", CodeInvalidArgument)
		})
	})
})

var TestResult = Type("TestResult", func() {
	Description("Test result for a single test")
	Field(1, "Name", String, "Name of the test", func() {
		Example("TestName")
	})
	Field(2, "Passed", Boolean, "Status of the test", func() {
		Example(true)
	})
	Field(3, "Error", String, "Error message if the test failed", func() {
		Example("Error message")
	})
	Field(4, "Duration", Int64, "Duration of the test in ms", func() {
		Example(100)
	})
	Required("Name", "Passed", "Duration")
})

var TestCollection = Type("TestCollection", func() {
	Description("Collection of test results for grouping by service")
	Field(1, "Name", String, "Name of the test collection", func() {
		Example("TestCollectionName")
	})
	Field(2, "Results", ArrayOf(TestResult), "Test results")
	Field(3, "Duration", Int64, "Duration of the tests in ms", func() {
		Example(100)
	})
	Field(4, "PassCount", Int, "Number of tests that passed", func() {
		Example(10)
	})
	Field(5, "FailCount", Int, "Number of tests that failed", func() {
		Example(1)
	})
	Required("Name", "Duration", "PassCount", "FailCount")
})

var TestResults = Type("TestResults", func() {
	Description("Test results for the iam system integration tests")
	Field(1, "Collections", ArrayOf(TestCollection), "Test collections")
	Field(2, "Duration", Int64, "Duration of the tests in ms", func() {
		Example(100)
	})
	Field(3, "PassCount", Int, "Number of tests that passed", func() {
		Example(10)
	})
	Field(4, "FailCount", Int, "Number of tests that failed", func() {
		Example(1)
	})
	Required("Collections", "Duration", "PassCount", "FailCount")
})

var SystemTestPayload = Type("TesterPayload", func() {
	Description("Payload for the tester service")
	Field(1, "Include", ArrayOf(String), "Tests to run. Allows wildcards.", func() {
		Example([]string{"TestNameToInclude"})
	})
	Field(2, "Exclude", ArrayOf(String), "Tests to exclude. Allows wildcards.", func() {
		Example([]string{"TestNameToExclude"})
	})
})
