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

// Runs all test collections EXCEPT smoke tests (those are in their own collections as well)
func (svc *Service) TestAll(ctx context.Context, p *gentester.TesterPayload) (res *gentester.TestResults, err error) {
	retval := gentester.TestResults{}

	// Forecaster tests
	forecasterCollection := TestCollection{
		Name: "Forecaster Tests",
	}
	forecasterResults, err := svc.runTests(ctx, p, &forecasterCollection, svc.forecasterTestMap, false)
	if err != nil {
		_ = logError(ctx, err)
		return nil, err
	}

	// Locator tests
	locatorCollection := TestCollection{
		Name: "Locator Tests",
	}
	locatorResults, err := svc.runTests(ctx, p, &locatorCollection, svc.locatorTestMap, true)
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
	// Smoke tests
	smokeCollection := TestCollection{
		Name: "Smoke Tests",
	}
	return svc.runTests(ctx, nil, &smokeCollection, svc.smokeTestMap, false)
}

// Runs the Forecaster Service tests as a collection in parallel
func (svc *Service) TestForecaster(ctx context.Context) (res *gentester.TestResults, err error) {
	// Forecaster tests
	forecasterCollection := TestCollection{
		Name: "Forecaster Tests",
	}
	return svc.runTests(ctx, nil, &forecasterCollection, svc.forecasterTestMap, false)
}

// Runs the Locator Service tests as a collection synchronously
func (svc *Service) TestLocator(ctx context.Context) (res *gentester.TestResults, err error) {
	// Locator tests
	locatorCollection := TestCollection{
		Name: "Locator Tests",
	}
	return svc.runTests(ctx, nil, &locatorCollection, svc.locatorTestMap, true)
}
