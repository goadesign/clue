package tester

import (
	"context"

	"goa.design/clue/debug"
	"goa.design/clue/log"
	"google.golang.org/grpc"

	genfront "goa.design/clue/example/weather/services/front/gen/front"
	genclient "goa.design/clue/example/weather/services/tester/gen/grpc/tester/client"
	gentester "goa.design/clue/example/weather/services/tester/gen/tester"
)

type (
	Client interface {
		// Runs ALL API Integration Tests from the Tester service, allowing for filtering on included or excluded tests
		TestAll(ctx context.Context, included, excluded []string) (*genfront.TestResults, error)
		// Runs API Integration Tests' Smoke Tests ONLY from the Tester service
		TestSmoke(ctx context.Context) (*genfront.TestResults, error)
	}

	TestAllPayload struct {
		Include []string
		Exclude []string
	}

	client struct {
		genc *gentester.Client
	}
)

// Creates a new client for the Tester service.
func New(cc *grpc.ClientConn) Client {
	c := genclient.NewClient(cc, grpc.WaitForReady(true))
	testSmoke := debug.LogPayloads(debug.WithClient())(c.TestSmoke())
	testAll := debug.LogPayloads(debug.WithClient())(c.TestAll())
	return &client{genc: gentester.NewClient(testSmoke, testAll, nil, nil)}
}

// TestSmoke runs the Smoke collection as defined in func_map.go of the tester service
func (c *client) TestSmoke(ctx context.Context) (*genfront.TestResults, error) {
	res, err := c.genc.TestSmoke(ctx)
	if err != nil {
		log.Errorf(ctx, err, "failed to run smoke tests: %s", err)
		return nil, err
	}
	return testerTestResultsToFrontTestResults(res), nil
}

// TestAll runs all tests in all collections. Obeys include and exclude filters.
// include and exclude are mutually exclusive and cannot be used together (400 error, bad request)
func (c *client) TestAll(ctx context.Context, included, excluded []string) (*genfront.TestResults, error) {
	gtPayload := &gentester.TesterPayload{
		Include: included,
		Exclude: excluded,
	}
	res, err := c.genc.TestAll(ctx, gtPayload)
	if err != nil {
		log.Errorf(ctx, err, "failed to run all tests: %s", err)
		return nil, err
	}
	return testerTestResultsToFrontTestResults(res), nil
}

func testerTestResultsToFrontTestResults(testResults *gentester.TestResults) *genfront.TestResults {
	var res = &genfront.TestResults{}
	if testResults != nil {
		res.Collections = testerTestCollectionsArrToFrontTestCollectionsArr(testResults.Collections)
		res.Duration = testResults.Duration
		res.PassCount = testResults.PassCount
		res.FailCount = testResults.FailCount
	}
	return res
}

func testerTestCollectionsArrToFrontTestCollectionsArr(testCollection []*gentester.TestCollection) []*genfront.TestCollection {
	var res []*genfront.TestCollection
	for _, v := range testCollection {
		res = append(res, testerTestCollectionToFrontTestCollection(v))
	}
	return res
}

func testerTestCollectionToFrontTestCollection(testCollection *gentester.TestCollection) *genfront.TestCollection {
	var res = &genfront.TestCollection{}
	if testCollection != nil {
		res.Name = testCollection.Name
		res.Results = testerTestResultsArrToFrontTestResultsArr(testCollection.Results)
		res.Duration = testCollection.Duration
		res.PassCount = testCollection.PassCount
		res.FailCount = testCollection.FailCount
	}
	return res
}

func testerTestResultsArrToFrontTestResultsArr(testResults []*gentester.TestResult) []*genfront.TestResult {
	var res []*genfront.TestResult
	for _, v := range testResults {
		res = append(res, (*genfront.TestResult)(v))
	}
	return res
}
