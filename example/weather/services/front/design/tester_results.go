package design

import (
	. "goa.design/goa/v3/dsl"
)

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
	Description("Test results for the iam system integration tests")
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
