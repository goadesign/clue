package tester

import (
	"context"

	"goa.design/clue/debug"
	"goa.design/clue/log"
	goa "goa.design/goa/v3/pkg"
	"google.golang.org/grpc"

	genfront "goa.design/clue/example/weather/services/front/gen/front"
	genclient "goa.design/clue/example/weather/services/tester/gen/grpc/tester/client"
	gentester "goa.design/clue/example/weather/services/tester/gen/tester"
)

type (
	Client interface {
		// Run smoke tests
		TestSmoke(ctx context.Context) (*genfront.TestResults, error)
		TestAll(ctx context.Context, payload *TestAllPayload) (*genfront.TestResults, error)
	}

	TestAllPayload struct {
		Include []string
		Exclude []string
	}

	client struct {
		testSmoke goa.Endpoint
		testAll   goa.Endpoint
	}
)

func New(cc *grpc.ClientConn) Client {
	c := genclient.NewClient(cc, grpc.WaitForReady(true))
	return &client{
		debug.LogPayloads(debug.WithClient())(c.TestSmoke()),
		debug.LogPayloads(debug.WithClient())(c.TestAll()),
	}
}

func (c *client) TestSmoke(ctx context.Context) (*genfront.TestResults, error) {
	res, err := c.testSmoke(ctx, nil)
	if err != nil {
		log.Errorf(ctx, err, "failed to run smoke tests: %s", err)
		return nil, err
	}
	return testerTestResultsToFrontTestResults(res.(*gentester.TestResults)), nil
}

func (c *client) TestAll(ctx context.Context, payload *TestAllPayload) (*genfront.TestResults, error) {
	gtPayload := &gentester.TesterPayload{
		Include: payload.Include,
		Exclude: payload.Exclude,
	}
	res, err := c.testAll(ctx, gtPayload)
	if err != nil {
		log.Errorf(ctx, err, "failed to run all tests: %s", err)
		return nil, err
	}
	return testerTestResultsToFrontTestResults(res.(*gentester.TestResults)), nil
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
		res = append(res, testerTestResultToFrontTestResult(v))
	}
	return res
}

func testerTestResultToFrontTestResult(testResult *gentester.TestResult) *genfront.TestResult {
	var res = &genfront.TestResult{}
	if testResult != nil {
		res.Name = testResult.Name
		res.Passed = testResult.Passed
		res.Error = testResult.Error
		res.Duration = testResult.Duration
	}
	return res
}
