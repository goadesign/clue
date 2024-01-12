package tester

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"runtime/debug"
	"strings"
	"sync"
	"time"

	"github.com/gobwas/glob"
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

func getStackTrace(wg *sync.WaitGroup, m *sync.Mutex) string {
	m.Lock()
	defer wg.Done()
	defer m.Unlock()
	// keep backup of the real stderr
	old := os.Stderr
	f, w, _ := os.Pipe()
	os.Stderr = w
	defer w.Close()

	debug.PrintStack()

	outC := make(chan string)
	go func() {
		var buf bytes.Buffer
		io.Copy(&buf, f) // nolint: errcheck
		outC <- buf.String()
	}()

	// restoring the real stderr
	os.Stderr = old
	out := <-outC

	return out
}

// recovers from a panicked test. This is used to ensure that the test
// suite does not crash if a test panics.
func recoverFromTestPanic(ctx context.Context, testName string, testCollection *TestCollection) {
	if r := recover(); r != nil {
		msg := fmt.Sprintf("[Panic Test]: %v", testName)
		err := errors.New(msg)
		log.Errorf(ctx, err, fmt.Sprintf("%v", r))
		var m sync.Mutex
		var wg sync.WaitGroup
		wg.Add(1)
		trace := getStackTrace(&wg, &m)
		wg.Wait()
		err = fmt.Errorf("%v : %v", r, trace)
		// log the error and add the test result to the test collection
		_ = logError(ctx, err)
		resultMsg := fmt.Sprintf("%v | %v", msg, r)
		testCollection.AppendTestResult(&gentester.TestResult{
			Name:     testName,
			Passed:   false,
			Error:    &resultMsg,
			Duration: -1,
		})
	}
}

// Filters a testMap based on a test name that is a glob string
// using standard wildcards https://tldp.org/LDP/GNU-Linux-Tools-Summary/html/x11655.htm
func matchTestFilter(ctx context.Context, test string, testMap map[string]func(context.Context, *TestCollection)) (bool, []func(context.Context, *TestCollection)) {
	match := false
	testMatches := []func(context.Context, *TestCollection){}
	var g glob.Glob
	g = glob.MustCompile(test)
	i := 0
	for testName, _ := range testMap {
		match = g.Match(testName)
		if match {
			testMatches = append(testMatches, testMap[testName])
		}
		i++
	}
	return match, testMatches
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
					findByWildcard, testFuncs := matchTestFilter(ctx, test, testMap)
					if findByWildcard {
						for i, testFunc := range testFuncs {
							testsToRun[fmt.Sprintf("%s_%d", test, i)] = testFunc
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
					var g glob.Glob
					g = glob.MustCompile(excludeTest)
					wildcardMatch = wildcardMatch || g.Match(testName)
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
		log.Infof(ctx, "RUNNING TESTS SYNCHRONOUSLY")
		for n, test := range testsToRun {
			log.Infof(ctx, "RUNNING TEST [%v]", n)
			test(ctx, testCollection)
		}
	} else {
		log.Infof(ctx, "RUNNING TESTS IN PARALLEL")
		wg := sync.WaitGroup{}
		for n, test := range testsToRun {
			wg.Add(1)
			go func(f func(context.Context, *TestCollection), testName string) {
				defer wg.Done()
				defer recoverFromTestPanic(ctx, testName, testCollection)
				log.Infof(ctx, "RUNNING TEST [%v]", testName)
				f(ctx, testCollection)
			}(test, n)
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
