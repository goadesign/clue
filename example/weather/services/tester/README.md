# Tester
The Tester service is a part of the weather application that is responsible for running various tests to ensure the functionality and reliability of the application. It is located in the example/weather/services/tester directory.

The service includes a set of predefined tests that are designed to validate different aspects of the application, such as the ability to locate IP addresses and forecast weather conditions. These tests are organized into three categories: smoke tests, forecaster tests, and locator tests.

Port: `8090`

## Tester gRPC API Design
The Tester service exposes a gRPC API that allows other services (such as the externally facing `front` service) to run tests, or it be called directly.

### TestResults return object
`TestResults`: This is the main return object for all test methods in the `tester` service and it represents the results of the system integration tests. It has four fields: `Collections` (an array of `TestCollection`), `Duration` (the total duration of all the tests in milliseconds), `PassCount` (the total number of tests that passed), and `FailCount` (the total number of tests that failed).

`TestCollection`: This type represents a collection of test results, typically grouped by service. It has five fields: `Name` (the name of the test collection), `Results` (an array of `TestResult`), `Duration` (the total duration of the tests in the `Collection` in milliseconds), `PassCount` (the number of tests that passed), and `FailCount` (the number of tests that failed).

`TestResult`: This type represents the result of a single test. It has four fields: `Name` (the name of the test), `Passed` (a boolean indicating whether the test passed), `Error` (an error message if the test failed), and `Duration` (the duration of the test in milliseconds).

### TestSmoke method
`TestSmoke` takes no inputs and returns `TestResults``. It will run just the tests defined as smoke tests in `services/tester/func_map.go`.

`Smoke Tests` are a subset of all tests that are designed to provide a basic level of confidence that the application is functioning properly.

### TestAll method
`TestAll` takes a `TesterPayload` for optional filtering of tests and returns TestResults. It will run all tests defined in `services/tester/func_map.go`.

`TestPayload`: This type represents a payload that can be passed to the `TestAll` method to filter the tests that are run. It has two fields: `Include` (an array of strings that represent the names of tests to include) and `Exclude` (an array of strings that represent the names of tests to exclude).

### TestForecaster & TestLocator methods
These methods run all tests defined for those services as found in `services/tester/func_map.go`.

## Tester exposed via Front
The `front` service exposes the `TestAll` and `TestSmoke` methods via gRPC. This allows the `front` service to tests in the application and return the results to the caller. 

This is useful when your gRPC services are not exposed publicly but your application still needs to be tested. In this case, the `front` service can be exposed publicly (with appropriate authentication required) and used to run tests on the application. This allows you to run tests on your application without exposing your gRPC services publicly.

An example of usage like this is to use the `front` service to run tests on the application from a CI/CD pipeline. A bash script can be written that calls the `front` service to run tests on the application and then exits with an error code if any tests fail. This script can then be called from a GitHub Actions workflow to run tests on the application after each commit, and even parse the results using a cli like `jq` to provide more detailed information about the test results in a comment to a PR.