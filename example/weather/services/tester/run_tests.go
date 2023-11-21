package tester

import (
	"context"
	"errors"
	"fmt"
	"runtime/debug"
	"strings"
	"sync"
	"time"

	gentester "goa.design/clue/example/weather/services/tester/gen/tester"
	"goa.design/clue/log"
	"golang.org/x/exp/slices"
)

// Ends a test by calculating duration and appending teh results to the test collection
func endTest(tr *gentester.TestResult, start time.Time, tc *TestCollection, results []*gentester.TestResult) {
	elapsed := time.Since(start).Milliseconds()
	tr.Duration = elapsed
	results = append(results, tr)
	tc.AppendTestResult(results...)
}

// recovers from a panicked test. This is used to ensure that the test
// suite does not crash if a test panics.
func recoverFromTestPanic(ctx context.Context, testName string, testCollection *TestCollection) {
	if r := recover(); r != nil {
		msg := fmt.Sprintf("[Panic Test]: %v", testName)
		err := errors.New(msg)
		log.Errorf(ctx, err, fmt.Sprintf("%v", r))
		// This doesn't work as I'd like because it only prints to stderr, not to a string
		// 	that I can caputre and put in the goa log.
		debug.PrintStack()
		resultMsg := fmt.Sprintf("%v | %v", msg, r)
		testCollection.AppendTestResult(&gentester.TestResult{
			Name:     testName,
			Passed:   false,
			Error:    &resultMsg,
			Duration: -1,
		})
	}
}

// Runs the tests from the testmap and handles filtering/exclusion of tests
func (svc *Service) runTests(ctx context.Context, p *gentester.TesterPayload, testCollection *TestCollection, testMap map[string]func(context.Context, *TestCollection)) (res *gentester.TestResults, err error) {
	retval := gentester.TestResults{}

	testsToRun := make(map[string]func(context.Context, *TestCollection))
	filtered := false
	// If there are no filters, the tests to run is all tests in the map.
	if (p == nil) || (p.Include == nil || len(p.Include) == 0) && (p.Exclude == nil || len(p.Exclude) == 0) {
		testsToRun = testMap
	} else { // Otherwise, we need to filter the tests
		// If there is an include list, we only run the tests in the include list. This will supersede any exclude list.
		filtered = true
		if len(p.Include) > 0 {
			if len(p.Exclude) > 0 {
				return nil, gentester.MakeIncludeExcludeBoth(errors.New("cannot have both include and exclude lists"))
			}
			for _, test := range p.Include {
				if testFunc, exists := testMap[test]; exists {
					testsToRun[test] = testFunc
				} else {
					// QUESTION: Do we want to error the test execution if a test is not found in the test map?
					// 		I'm thinking no, because it's not really an error, it's just a test that doesn't exist.
					log.Infof(ctx, "Test [%v] not found in test map", test)
				}
			}
		} else if len(p.Exclude) > 0 { // If there is only an exclude list, we add tests not found in that exclude list to the tests to run
			for testName, test := range testMap {
				// This is from golang's experimental slices package
				// (https://godoc.org/golang.org/x/exp/slices)
				if !slices.Contains(p.Exclude, testName) {
					testsToRun[testName] = test
				} else {
					log.Debugf(ctx, "Test [%v] excluded", testName)
				}
			}
		}
		// No else because it should never be reached. The top level if catches no filters.
		// len(p.Include)> 0 handles the include case (which supersedes any exclude list)
		// and len(p.Exclude) >0 handles the exclude only case.
	}

	// Run the tests that need to be run and add the results to the testCollection.Results array
	wg := sync.WaitGroup{}
	for n, test := range testsToRun {
		wg.Add(1)
		go func(f func(context.Context, *TestCollection), testName string) {
			defer recoverFromTestPanic(ctx, testName, testCollection)
			defer wg.Done()
			log.Infof(ctx, "RUNNING TEST [%v]", testName)
			f(ctx, testCollection)
		}(test, n)
	}
	wg.Wait()

	// If there was a Panic in the test, sometimes it falls through to and checks results from the test collection
	// before it finishes the deferred recoverFromTestPanic, making it look like the panicked test was not handled
	// or run. This sleep is to give the recoverFromTestPanic function time to finish before we check the results.
	time.Sleep(250 * time.Millisecond)

	for _, res := range testCollection.Results {
		if !res.Passed {
			errorMessage := *res.Error
			log.Infof(ctx, "[Failed Test] Collection: [%v], Test [%v] failed with message [%s] and a duration of [%v]", testCollection.Name, res.Name, errorMessage, res.Duration)
		}
	}

	//Calculate Collection Duration & pass/fail counts
	collectionDuration := int64(0)
	passCount := 0
	failCount := 0
	for _, test := range testCollection.Results {
		collectionDuration += test.Duration
		if test.Passed {
			passCount++
		} else {
			failCount++
		}
	}
	testCollection.Duration = collectionDuration
	testCollection.PassCount = passCount
	testCollection.FailCount = failCount
	returnTc := gentester.TestCollection{
		Name:      testCollection.Name,
		Duration:  testCollection.Duration,
		PassCount: testCollection.PassCount,
		FailCount: testCollection.FailCount,
		Results:   testCollection.Results,
	}
	retval.Collections = append(retval.Collections, &returnTc)

	// Calculate Total Duration & total pass/fail counts
	totalDuration := int64(0)
	totalPassed := 0
	totalFailed := 0
	for _, coll := range retval.Collections {
		totalDuration += coll.Duration
		totalPassed += coll.PassCount
		totalFailed += coll.FailCount
		snake_case_coll_name := strings.Replace(strings.ToLower(coll.Name), " ", "_", -1)
		if filtered {
			snake_case_coll_name = snake_case_coll_name + "_filtered"
		}
		log.Infof(ctx, "Collection: [%v] Duration: [%v]", snake_case_coll_name, coll.Duration)
		log.Infof(ctx, "Collection: [%v] Pass Count: [%v]", snake_case_coll_name, coll.PassCount)
		log.Infof(ctx, "Collection: [%v] Fail Count: [%v]", snake_case_coll_name, coll.FailCount)
	}
	retval.Duration = totalDuration
	retval.PassCount = totalPassed
	retval.FailCount = totalFailed
	return &retval, nil
}
