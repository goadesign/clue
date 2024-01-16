package tester

import (
	"context"
	"fmt"
	"time"

	gentester "goa.design/clue/example/weather/services/tester/gen/tester"
)

func (svc *Service) TestLocatorValidIP(ctx context.Context, tc *TestCollection) {
	results := []*gentester.TestResult{}
	start := time.Now()
	name := "TestLocatorValidIP"
	passed := true
	testRes := gentester.TestResult{
		Name:   name,
		Passed: passed,
	}
	// Start Test Logic
	locatorReturn, err := svc.lc.GetLocation(ctx, "8.8.8.8")
	if err != nil {
		testRes.Passed = false
		errorDescription := fmt.Sprintf("error getting location for valid ip: %v", err)
		testRes.Error = &errorDescription
		endTest(&testRes, start, tc, results)
		return
	}
	if locatorReturn == nil {
		testRes.Passed = false
		errorDescription := fmt.Sprintf("location for valid ip is nil: %v", err)
		testRes.Error = &errorDescription
		endTest(&testRes, start, tc, results)
		return
	}
	// End Test Logic
	endTest(&testRes, start, tc, results)
}

func (svc *Service) TestLocatorInvalidIP(ctx context.Context, tc *TestCollection) {
	results := []*gentester.TestResult{}
	start := time.Now()
	name := "TestLocatorInvalidIP"
	passed := true
	testRes := gentester.TestResult{
		Name:   name,
		Passed: passed,
	}
	// Start Test Logic
	locatorReturn, err := svc.lc.GetLocation(ctx, "999")
	if err != nil {
		testRes.Passed = false
		errorDescription := fmt.Sprintf("no error getting location for invalid ip: %v", err)
		testRes.Error = &errorDescription
		endTest(&testRes, start, tc, results)
		return
	}
	if locatorReturn == nil {
		testRes.Passed = false
		errorDescription := fmt.Sprintf("location for valid ip is not nil: %v", err)
		testRes.Error = &errorDescription
		endTest(&testRes, start, tc, results)
		return
	}
	// End Test Logic
	endTest(&testRes, start, tc, results)
}
