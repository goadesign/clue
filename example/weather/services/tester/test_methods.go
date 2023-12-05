package tester

import (
	"context"
	"sync"

	gentester "goa.design/clue/example/weather/services/tester/gen/tester"
)

// Alias type for gen/tester.TestCollection to allow for adding an Append func
type TestCollection gentester.TestCollection

// Appends a slice of TestResults to the TestCollection
func (t *TestCollection) AppendTestResult(tr ...*gentester.TestResult) {
	mutex := sync.Mutex{}
	mutex.Lock()
	defer mutex.Unlock()
	t.Results = append(t.Results, tr...)
}

var filteringPayload = &gentester.TesterPayload{}

// Runs all test collections EXCEPT smoke tests (those are in their own collections as well)
func (svc *Service) TestAll(ctx context.Context, p *gentester.TesterPayload) (res *gentester.TestResults, err error) {
	retval := gentester.TestResults{}
	filteringPayload = p
	forecasterResults, err := svc.TestForecaster(ctx)
	if err != nil {
		_ = logError(ctx, err)
		return nil, err
	}
	locatorResults, err := svc.TestLocator(ctx)
	if err != nil {
		_ = logError(ctx, err)
		return nil, err
	}

	//Merge Disparate Collection Results
	allResults := []*gentester.TestResults{}
	allResults = append(allResults, forecasterResults)
	allResults = append(allResults, locatorResults)
	for _, r := range allResults {
		retval.Collections = append(retval.Collections, r.Collections...)
		retval.Duration += r.Duration
		retval.PassCount += r.PassCount
		retval.FailCount += r.FailCount
	}

	return &retval, nil
}

// Runs the Smoke tests as a collection in parallel
func (svc *Service) TestSmoke(ctx context.Context) (res *gentester.TestResults, err error) {
	svc.smokeTestMapInit(ctx)

	// Smoke tests
	smokeCollection := TestCollection{
		Name: "Smoke Tests",
	}
	return svc.runTests(ctx, filteringPayload, &smokeCollection, smokeTestMap, false)
}

// Runs the ACL Service tests as a collection in parallel
func (svc *Service) TestForecaster(ctx context.Context) (res *gentester.TestResults, err error) {
	svc.forecasterTestMapInit(ctx)

	// ACL tests
	aclCollection := TestCollection{
		Name: "Forecaster Tests",
	}
	return svc.runTests(ctx, filteringPayload, &aclCollection, forecasterTestMap, false)
}

// Runs the Login Service tests as a collection synchronously
func (svc *Service) TestLocator(ctx context.Context) (res *gentester.TestResults, err error) {
	svc.locatorTestMapInit(ctx)

	// Login tests
	loginCollection := TestCollection{
		Name: "Locator Tests",
	}
	return svc.runTests(ctx, filteringPayload, &loginCollection, locatorTestMap, true)
}
