package tester

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"regexp"
	"runtime/debug"
	"strings"
	"sync"
	"time"

	"goa.design/clue/log"

	gentester "goa.design/clue/example/weather/services/tester/gen/tester"
)

// Ends a test by calculating duration and appending the results to the test collection
func endTest(tr *gentester.TestResult, start time.Time, tc *TestCollection, results []*gentester.TestResult) {
	elapsed := time.Since(start).Milliseconds()
	tr.Duration = elapsed
	results = append(results, tr)
	tc.AppendTestResult(results...)
}

func getStackTrace(wg *sync.WaitGroup, m *sync.Mutex) (string, error) {
	m.Lock()
	defer wg.Done()
	defer m.Unlock()
	// keep backup of the real stderr
	old := os.Stderr
	f, w, _ := os.Pipe()
	os.Stderr = w

	debug.PrintStack()
	w.Close()

	outC := make(chan string, 1)
	outErr := make(chan error, 1)
	go func() {
		var buf bytes.Buffer
		_, err := io.Copy(&buf, f)
		outErr <- err
		outC <- buf.String()
	}()

	// restoring the real stderr
	os.Stderr = old

	if err := <-outErr; err != nil {
		return "", err
	}
	out := <-outC
	return out, nil
}

// recovers from a panicked test. This is used to ensure that the test
// suite does not crash if a test panics.
func recoverFromTestPanic(ctx context.Context, testNameFunc func() string, testCollection *TestCollection) {
	if r := recover(); r != nil {
		msg := fmt.Sprintf("[Panic Test]: %v", testNameFunc())
		err := errors.New(msg)
		log.Errorf(ctx, err, fmt.Sprintf("%v", r))
		var m sync.Mutex
		var wg sync.WaitGroup
		wg.Add(1)
		trace, err := getStackTrace(&wg, &m)
		wg.Wait()
		var resultMsg string
		if err != nil {
			err = fmt.Errorf("error getting stack trace for panicked test: %v", err)
			resultMsg = err.Error()
		} else {
			err = fmt.Errorf("%v : %v", r, trace)
			// log the error and add the test result to the test collection
			_ = logError(ctx, err)
			resultMsg = fmt.Sprintf("%v | %v", msg, r)
		}
		testCollection.AppendTestResult(&gentester.TestResult{
			Name:     testNameFunc(),
			Passed:   false,
			Error:    &resultMsg,
			Duration: -1,
		})
	}
}

// Converts a wildcard string using * to a regular expression string
func wildCardToRegexp(pattern string) string {
	components := strings.Split(pattern, "*")
	if len(components) == 1 {
		// if len is 1, there are no *'s, return exact match pattern
		return "^" + pattern + "$"
	}
	var result strings.Builder
	for i, literal := range components {

		// Replace * with .*
		if i > 0 {
			result.WriteString(".*")
		}

		// Quote any regular expression meta characters in the
		// literal text.
		result.WriteString(regexp.QuoteMeta(literal))
	}
	return "^" + result.String() + "$"
}

// wraps wildCardToRegexp and returns a bool indicating whether the value
// matches the pattern, the string matched, and an error if one occurred
func match(pattern string, value string) (bool, string, error) {
	r, err := regexp.Compile(wildCardToRegexp(pattern))
	if err != nil {
		return false, "", err
	}
	matches := r.FindStringSubmatch(value)
	if len(matches) > 0 {
		return true, matches[0], nil
	} else {
		return false, "", nil
	}
}

// Filters a testMap based on a test name that is a glob string using only
// * wildcards
func matchTestFilterRegex(ctx context.Context, test string, testMap map[string]func(context.Context, *TestCollection)) (map[string]func(context.Context, *TestCollection), error) {
	retval := make(map[string]func(context.Context, *TestCollection))
	i := 0
	for testName := range testMap {
		_, matchString, err := match(test, testName)
		if err != nil {
			return nil, err
		}
		if matchString != "" {
			retval[matchString] = testMap[testName]
		}
		i++
	}
	return retval, nil
}

// Runs the tests from the testmap and handles filtering/exclusion of tests
// Pass in `true` for runSynchronously to run the tests synchronously instead
// of in parallel.
func (svc *Service) runTests(ctx context.Context, p *gentester.TesterPayload, testCollection *TestCollection, testMap map[string]func(context.Context, *TestCollection), runSynchronously bool) (*gentester.TestResults, error) {
	retval := gentester.TestResults{}

	var filtered bool
	testsToRun := testMap
	// we need to filter the tests if there is an include or exclude list
	if p != nil && (len(p.Include) > 0 || len(p.Exclude) > 0) {
		testsToRun = make(map[string]func(context.Context, *TestCollection))
		// If there is an include list, we only run the tests in the include list. This will supersede any exclude list.
		filtered = true
		if len(p.Include) > 0 {
			if len(p.Exclude) > 0 {
				return nil, gentester.MakeIncludeExcludeBoth(errors.New("cannot have both include and exclude lists"))
			}
			for _, test := range p.Include {
				if testFunc, ok := testMap[test]; ok {
					testsToRun[test] = testFunc
				} else { // Test didn't match exactly, so we're gonna try for a wildcard match
					testFuncs, err := matchTestFilterRegex(ctx, test, testMap)
					if err != nil {
						return nil, gentester.MakeWildcardCompileError(err)
					}
					if len(testFuncs) > 0 {
						for testName, testFunc := range testFuncs {
							testsToRun[testName] = testFunc
						}
					} else { // No wildcard match either
						err := fmt.Errorf("test [%v] not found in test map", test)
						_ = logError(ctx, err)
					}
				}
			}
		} else if len(p.Exclude) > 0 { // If there is only an exclude list, we add tests not found in that exclude list to the tests to run
			for testName, test := range testMap {
				wildcardMatch := false
				for _, excludeTest := range p.Exclude {
					wildcardMatchThisTest, _, err := match(excludeTest, testName)
					if err != nil {
						return nil, gentester.MakeWildcardCompileError(err)
					}
					wildcardMatch = wildcardMatch || wildcardMatchThisTest
				}
				if !wildcardMatch {
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
	if runSynchronously {
		// RunSynchronously is used for test collections that need to be run one after in order
		// to avoid single resource contention between tests if they are run in parallel. An
		// example of this is tests that rely on the same cloud resource, such as a spreadsheet,
		// as part of their test functionality.
		//
		// testName is passed to recoverFromTestPanic as a function so that, via a closure, its
		// name can be set before the test is run but after the defer of recoverFromTestPanic is
		// declared. This is done because the test name is not accessible from within the test
		// function itself where it is set.
		log.Infof(ctx, "RUNNING TESTS SYNCHRONOUSLY")
		wg := sync.WaitGroup{}
		wg.Add(1)
		go func() {
			defer wg.Done()
			testName := ""
			testNameFunc := func() string {
				return testName
			}
			defer recoverFromTestPanic(ctx, testNameFunc, testCollection)
			for testNameRunning, test := range testsToRun {
				testName = testNameRunning
				log.Infof(ctx, "RUNNING TEST [%v]", testName)
				test(ctx, testCollection)
			}
		}()
		wg.Wait()
	} else {
		// if not run synchronously, run the tests in parallel and assumed not to have resource
		// contention
		log.Infof(ctx, "RUNNING TESTS IN PARALLEL")
		wg := sync.WaitGroup{}
		for name, test := range testsToRun {
			wg.Add(1)
			go func(f func(context.Context, *TestCollection), testNameRunning string) {
				defer wg.Done()
				testNameFunc := func() string {
					return testNameRunning
				}
				defer recoverFromTestPanic(ctx, testNameFunc, testCollection)
				log.Infof(ctx, "RUNNING TEST [%v]", testNameRunning)
				f(ctx, testCollection)
			}(test, name)
		}
		wg.Wait()
	}

	for _, res := range testCollection.Results {
		if !res.Passed {
			log.Infof(ctx, "[Failed Test] Collection: [%v], Test [%v] failed with message [%s] and a duration of [%v]", testCollection.Name, res.Name, *res.Error, res.Duration)
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
