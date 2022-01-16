package instrument

import (
	"context"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/crossnokaye/micro/instrument/testsvc"
)

func TestUnaryServerInterceptorServerDuration(t *testing.T) {
	buckets := []float64{10, 110}
	cases := []struct {
		name                 string
		d                    time.Duration
		expectedBucketCounts []int
	}{
		{"fast", 1 * time.Millisecond, []int{1, 1}},
		{"slow", 100 * time.Millisecond, []int{0, 1}},
		{"very slow", 1000 * time.Millisecond, []int{0, 0}},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			restore := timeSince
			defer func() { timeSince = restore }()
			timeSince = func(time.Time) time.Duration { return c.d }

			reg := NewTestRegistry(t)
			uinter := UnaryServerInterceptor(context.Background(), "testsvc", WithRegisterer(reg), WithDurationBuckets(buckets))
			cli, stop := testsvc.SetupGRPC(t,
				testsvc.WithUnaryInterceptor(uinter),
				testsvc.WithUnaryFunc(noopUnaryMethod()))

			_, err := cli.GRPCMethod(context.Background(), &testsvc.Fields{})
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}

			stop()
			reg.AssertHistogram(MetricRPCDuration, RPCLabels, 1, c.expectedBucketCounts)
		})
	}
}

func TestUnaryServerInterceptorRequestSize(t *testing.T) {
	buckets := []float64{10, 110}
	cases := []struct {
		name                 string
		str                  string
		expectedBucketCounts []int
	}{
		{"small", "1", []int{1, 1}},
		{"large", strings.Repeat("1", 100), []int{0, 1}},
		{"very large", strings.Repeat("1", 1000), []int{0, 0}},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			reg := NewTestRegistry(t)
			uinter := UnaryServerInterceptor(context.Background(), "testsvc", WithRegisterer(reg), WithRequestSizeBuckets(buckets))
			cli, stop := testsvc.SetupGRPC(t,
				testsvc.WithUnaryInterceptor(uinter),
				testsvc.WithUnaryFunc(noopUnaryMethod()))

			_, err := cli.GRPCMethod(context.Background(), &testsvc.Fields{S: &c.str})
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}

			stop()
			reg.AssertHistogram(MetricRPCRequestSize, RPCLabels, 1, c.expectedBucketCounts)
		})
	}
}

func TestUnaryServerInterceptorResponseSize(t *testing.T) {
	buckets := []float64{10, 110}
	cases := []struct {
		name                 string
		str                  string
		expectedBucketCounts []int
	}{
		{"small", "1", []int{1, 1}},
		{"large", strings.Repeat("1", 100), []int{0, 1}},
		{"very large", strings.Repeat("1", 1000), []int{0, 0}},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			reg := NewTestRegistry(t)
			uinter := UnaryServerInterceptor(context.Background(), "testsvc", WithRegisterer(reg), WithResponseSizeBuckets(buckets))
			cli, stop := testsvc.SetupGRPC(t,
				testsvc.WithUnaryInterceptor(uinter),
				testsvc.WithUnaryFunc(stringUnaryMethod(c.str)))

			_, err := cli.GRPCMethod(context.Background(), &testsvc.Fields{})
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}

			stop()
			reg.AssertHistogram(MetricRPCResponseSize, RPCLabels, 1, c.expectedBucketCounts)
		})
	}
}

func TestUnaryServerInterceptorActiveRequests(t *testing.T) {
	cases := []struct {
		name    string
		numReqs int
	}{
		{"one", 1},
		{"ten", 10},
		{"one thousand", 1000},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			reg := NewTestRegistry(t)
			uinter := UnaryServerInterceptor(context.Background(), "testsvc", WithRegisterer(reg))
			chstop := make(chan struct{})
			var running, done sync.WaitGroup
			running.Add(c.numReqs)
			done.Add(c.numReqs)
			cli, stop := testsvc.SetupGRPC(t,
				testsvc.WithUnaryInterceptor(uinter),
				testsvc.WithUnaryFunc(waitUnaryMethod(&running, &done, chstop)))

			for i := 0; i < c.numReqs; i++ {
				go cli.GRPCMethod(context.Background(), &testsvc.Fields{})
			}

			running.Wait()
			reg.AssertGauge(MetricRPCActiveRequests, RPCActiveRequestsLabels, c.numReqs)
			close(chstop)
			done.Wait()
			reg.AssertGauge(MetricRPCActiveRequests, RPCActiveRequestsLabels, 0)
			stop()
		})
	}
}

func noopUnaryMethod() testsvc.UnaryFunc {
	return func(ctx context.Context, _ *testsvc.Fields) (*testsvc.Fields, error) {
		return &testsvc.Fields{}, nil
	}
}

func stringUnaryMethod(str string) testsvc.UnaryFunc {
	return func(ctx context.Context, _ *testsvc.Fields) (*testsvc.Fields, error) {
		return &testsvc.Fields{S: &str}, nil
	}
}

func waitUnaryMethod(running, done *sync.WaitGroup, stop chan struct{}) testsvc.UnaryFunc {
	return func(ctx context.Context, _ *testsvc.Fields) (*testsvc.Fields, error) {
		running.Done()
		defer done.Done()
		<-stop
		return &testsvc.Fields{}, nil
	}
}
