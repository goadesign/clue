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
		// filteringPayload needs reset as TestAll calls the OTHER test methods
		// and we only want a filteringPayload if it came via TestAll.
		// if the other test methods are called directly (e.g. TestForecaster)
		// it doesn't accept a gentester.TesterPayload
		filteringPayload = &gentester.TesterPayload{}
		return nil, err
	}
	locatorResults, err := svc.TestLocator(ctx)
	if err != nil {
		_ = logError(ctx, err)
		// filteringPayload needs reset as TestAll calls the OTHER test methods
		// and we only want a filteringPayload if it came via TestAll.
		// if the other test methods are called directly (e.g. TestForecaster)
		// it doesn't accept a gentester.TesterPayload
		filteringPayload = &gentester.TesterPayload{}
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

	// filteringPayload needs reset as TestAll calls the OTHER test methods
	// and we only want a filteringPayload if it came via TestAll.
	// if the other test methods are called directly (e.g. TestForecaster)
	// it doesn't accept a gentester.TesterPayload
	filteringPayload = &gentester.TesterPayload{}
	return &retval, nil
}

// Runs the Smoke tests as a collection in parallel
func (svc *Service) TestSmoke(ctx context.Context) (res *gentester.TestResults, err error) {
	// Smoke tests
	smokeCollection := TestCollection{
		Name: "Smoke Tests",
	}
	return svc.runTests(ctx, filteringPayload, &smokeCollection, svc.smokeTestMap, false)
}

// Runs the Forecaster Service tests as a collection in parallel
func (svc *Service) TestForecaster(ctx context.Context) (res *gentester.TestResults, err error) {
	// Forecaster tests
	forecasterCollection := TestCollection{
		Name: "Forecaster Tests",
	}
	return svc.runTests(ctx, filteringPayload, &forecasterCollection, svc.forecasterTestMap, false)
}

// Runs the Locator Service tests as a collection synchronously
func (svc *Service) TestLocator(ctx context.Context) (res *gentester.TestResults, err error) {
	// Locator tests
	locatorCollection := TestCollection{
		Name: "Locator Tests",
	}
	return svc.runTests(ctx, filteringPayload, &locatorCollection, svc.locatorTestMap, true)
}
