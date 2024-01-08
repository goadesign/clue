# Tester

The Tester service is a part of the weather application that is responsible for running various 
tests to ensure the functionality and reliability of the application. It is located in the 
example/weather/services/tester directory.

The service includes a set of predefined tests that are designed to validate different aspects 
of the application, such as the ability to locate IP addresses and forecast weather conditions. 
These tests are organized into three categories: smoke tests, forecaster tests, and locator 
tests.

Port: `8090`

## Tester gRPC API Design

The Tester service exposes a gRPC API that allows other services (such as the externally facing 
`front` service) to run tests, or for it to be called directly.

### TestResults return object

`TestResults` serves as the key output for the `tester` service's test methods,
encapsulating the system integration test outcomes. It details the collective
test results with fields like `Collections`, an array encompassing various
`TestCollection` instances, the `Duration` for total test time in milliseconds,
and counts of both passed (`PassCount`) and failed tests (`FailCount`).

Each `TestCollection` within `TestResults` groups test outcomes, often by
service. Its structure includes the collection's Name, an array of individual
`TestResult` items, the total `Duration` for all tests in the collection, and
counts of successful (`PassCount`) and unsuccessful (`FailCount`) tests.

In `TestResult`, the focus narrows down to individual test performance. This
includes the specific `Name` of the test, a boolean `Passed` status, a potential
`Error` message for failures, and the `Duration` each test takes.

### TestSmoke method

`TestSmoke` takes no inputs and returns `TestResults`. It will run just the tests defined as smoke 
tests in `services/tester/func_map.go`.

`Smoke Tests` are a subset of all tests that are designed to provide a basic level of confidence 
that the application is functioning properly.

### TestAll method

`TestAll` takes a `TesterPayload` for optional filtering of tests and returns TestResults. It will 
run all tests defined in `services/tester/func_map.go`.

`TestPayload`: This type represents a payload that can be passed to the `TestAll` method to filter 
the tests that are run. It has two fields: `Include` (an array of strings that represent the names 
of tests to include) and `Exclude` (an array of strings that represent the names of tests to 
exclude).

`Include` and `Exclude` are mutually exclusive and cannot be used together. If that is done then
the `TestAll` method will return a `400 Bad Request` error.

### TestForecaster & TestLocator methods

These methods run all tests defined for those services as found in `services/tester/func_map.go`.

## Tester exposed via Front

The `front` service exposes the `TestAll` and `TestSmoke` methods via gRPC. This allows the `front` 
service to tests in the application and return the results to the caller. 

This is useful when your gRPC services are not exposed publicly but your application still needs to 
be tested. In this case, the `front` service can be exposed publicly (with appropriate 
authentication required) and used to run tests on the application. This allows you to run tests on 
your application without exposing your gRPC services publicly.

An example of usage like this is to use the `front` service to run tests on the application from a 
CI/CD pipeline. A bash script can be written that calls the `front` service to run tests on the 
application and then exits with an error code if any tests fail. This script can then be called 
from a GitHub Actions workflow to run tests on the application after each commit, and even parse 
the results using a cli like `jq` to provide more detailed information about the test results in a 
comment to a PR.