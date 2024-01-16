package tester

import (
	"context"
	"fmt"
	"strings"
	"time"

	gentester "goa.design/clue/example/weather/services/tester/gen/tester"
)

func (svc *Service) TestForecasterValidLatLong(ctx context.Context, tc *TestCollection) {
	results := []*gentester.TestResult{}
	start := time.Now()
	name := "TestForecasterValidLatLong"
	passed := true
	testRes := gentester.TestResult{
		Name:   name,
		Passed: passed,
	}
	// Start Test Logic
	forecastReturn, err := svc.fc.GetForecast(ctx, 34.424019199545704, -119.70487678810295)
	if err != nil {
		testRes.Passed = false
		errorDescription := fmt.Sprintf("error getting forecast for valid lat/long: %v", err)
		testRes.Error = &errorDescription
		endTest(&testRes, start, tc, results)
		return
	}
	if forecastReturn == nil {
		testRes.Passed = false
		errorDescription := fmt.Sprintf("forecast for valid lat/long is nil: %v", err)
		testRes.Error = &errorDescription
		endTest(&testRes, start, tc, results)
		return
	}
	if forecastReturn.Location == nil {
		testRes.Passed = false
		errorDescription := fmt.Sprintf("forecast for valid lat/long has nil location: %v", err)
		testRes.Error = &errorDescription
		endTest(&testRes, start, tc, results)
		return
	}
	if forecastReturn.Location.City == "" {
		testRes.Passed = false
		errorDescription := fmt.Sprintf("forecast for valid lat/long has empty city: %v", err)
		testRes.Error = &errorDescription
		endTest(&testRes, start, tc, results)
		return
	}
	if forecastReturn.Periods == nil {
		testRes.Passed = false
		errorDescription := fmt.Sprintf("forecast for valid lat/long has nil periods: %v", err)
		testRes.Error = &errorDescription
		endTest(&testRes, start, tc, results)
		return
	}
	// End Test Logic
	endTest(&testRes, start, tc, results)
}

func (svc *Service) TestForecasterInvalidLat(ctx context.Context, tc *TestCollection) {
	results := []*gentester.TestResult{}
	start := time.Now()
	name := "TestForecasterInvalidLat"
	passed := true
	testRes := gentester.TestResult{
		Name:   name,
		Passed: passed,
	}
	// Start Test Logic
	forecastReturn, err := svc.fc.GetForecast(ctx, 999, -119.70487678810295)
	if err == nil {
		testRes.Passed = false
		errorDescription := fmt.Sprintf("no error getting forecast for invalid lat: %v", err)
		testRes.Error = &errorDescription
		endTest(&testRes, start, tc, results)
		return
	}
	if forecastReturn != nil {
		testRes.Passed = false
		errorDescription := fmt.Sprintf("forecast for invalid lat is not nil: %v", err)
		testRes.Error = &errorDescription
		endTest(&testRes, start, tc, results)
		return
	}
	if !strings.Contains(err.Error(), "Invalid Parameter") {
		testRes.Passed = false
		errorDescription := fmt.Sprintf("error message for invalid lat is not correct: %v", err)
		testRes.Error = &errorDescription
		endTest(&testRes, start, tc, results)
		return
	}
	// End Test Logic
	endTest(&testRes, start, tc, results)
}

func (svc *Service) TestForecasterInvalidLong(ctx context.Context, tc *TestCollection) {
	results := []*gentester.TestResult{}
	start := time.Now()
	name := "TestForecasterInvalidLong"
	passed := true
	testRes := gentester.TestResult{
		Name:   name,
		Passed: passed,
	}
	// Start Test Logic
	forecastReturn, err := svc.fc.GetForecast(ctx, 34.424019199545704, 999)
	if err == nil {
		testRes.Passed = false
		errorDescription := fmt.Sprintf("no error getting forecast for invalid long: %v", err)
		testRes.Error = &errorDescription
		endTest(&testRes, start, tc, results)
		return
	}
	if forecastReturn != nil {
		testRes.Passed = false
		errorDescription := fmt.Sprintf("forecast for invalid long is not nil: %v", err)
		testRes.Error = &errorDescription
		endTest(&testRes, start, tc, results)
		return
	}
	if !strings.Contains(err.Error(), "Invalid Parameter") {
		testRes.Passed = false
		errorDescription := fmt.Sprintf("error message for invalid long is not correct: %v", err)
		testRes.Error = &errorDescription
		endTest(&testRes, start, tc, results)
		return
	}
	// End Test Logic
	endTest(&testRes, start, tc, results)
}
