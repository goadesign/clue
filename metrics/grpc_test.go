package metrics

import (
	"context"
	"strings"
	"sync"
	"testing"
	"time"

	"goa.design/clue/internal/testsvc"
	"google.golang.org/grpc"
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
			ctx := Context(context.Background(), "testsvc", WithRegisterer(reg), WithDurationBuckets(buckets))
			uinter := UnaryServerInterceptor(ctx)
			cli, stop := testsvc.SetupGRPC(t,
				testsvc.WithServerOptions(grpc.UnaryInterceptor(uinter)),
				testsvc.WithUnaryFunc(noopMethod()))

			_, err := cli.GRPCMethod(context.Background(), &testsvc.Fields{})
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}

			stop()
			reg.AssertHistogram(metricRPCDuration, rpcLabels, 1, c.expectedBucketCounts)
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
			ctx := Context(context.Background(), "testsvc", WithRegisterer(reg), WithRequestSizeBuckets(buckets))
			uinter := UnaryServerInterceptor(ctx)
			cli, stop := testsvc.SetupGRPC(t,
				testsvc.WithServerOptions(grpc.UnaryInterceptor(uinter)),
				testsvc.WithUnaryFunc(noopMethod()))

			_, err := cli.GRPCMethod(context.Background(), &testsvc.Fields{S: &c.str})
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}

			stop()
			reg.AssertHistogram(metricRPCRequestSize, rpcLabels, 1, c.expectedBucketCounts)
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
			ctx := Context(context.Background(), "testsvc", WithRegisterer(reg), WithResponseSizeBuckets(buckets))
			uinter := UnaryServerInterceptor(ctx)
			cli, stop := testsvc.SetupGRPC(t,
				testsvc.WithServerOptions(grpc.UnaryInterceptor(uinter)),
				testsvc.WithUnaryFunc(stringMethod(c.str)))

			_, err := cli.GRPCMethod(context.Background(), &testsvc.Fields{})
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}

			stop()
			reg.AssertHistogram(metricRPCResponseSize, rpcLabels, 1, c.expectedBucketCounts)
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
			ctx := Context(context.Background(), "testsvc", WithRegisterer(reg))
			uinter := UnaryServerInterceptor(ctx)
			chstop := make(chan struct{})
			var running, done sync.WaitGroup
			running.Add(c.numReqs)
			done.Add(c.numReqs)
			cli, stop := testsvc.SetupGRPC(t,
				testsvc.WithServerOptions(grpc.UnaryInterceptor(uinter)),
				testsvc.WithUnaryFunc(waitMethod(&running, &done, chstop)))

			for i := 0; i < c.numReqs; i++ {
				go cli.GRPCMethod(context.Background(), &testsvc.Fields{})
			}

			running.Wait()
			reg.AssertGauge(metricRPCActiveRequests, rpcNoCodeLabels, c.numReqs)
			close(chstop)
			done.Wait()
			stop()
		})
	}
}

func TestStreamServerInterceptorServerDuration(t *testing.T) {
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
			ctx := Context(context.Background(), "testsvc", WithRegisterer(reg), WithDurationBuckets(buckets))
			sinter := StreamServerInterceptor(ctx)
			cli, stop := testsvc.SetupGRPC(t,
				testsvc.WithServerOptions(grpc.StreamInterceptor(sinter)),
				testsvc.WithStreamFunc(echoMethod()))

			stream, err := cli.GRPCStream(context.Background())
			if err != nil {
				t.Errorf("unexpected stream error: %v", err)
			}
			if err := stream.Send(&testsvc.Fields{}); err != nil {
				t.Errorf("unexpected send error: %v", err)
			}
			if _, err := stream.Recv(); err != nil {
				t.Errorf("unexpected recv error: %v", err)
			}

			stop()
			reg.AssertHistogram(metricRPCDuration, rpcLabels, 1, c.expectedBucketCounts)
		})
	}
}

func TestStreamServerInterceptorMessageSize(t *testing.T) {
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
			ctx := Context(context.Background(), "testsvc", WithRegisterer(reg), WithRequestSizeBuckets(buckets))
			sinter := StreamServerInterceptor(ctx)
			cli, stop := testsvc.SetupGRPC(t,
				testsvc.WithServerOptions(grpc.StreamInterceptor(sinter)),
				testsvc.WithStreamFunc(echoMethod()))

			stream, err := cli.GRPCStream(context.Background())
			if err != nil {
				t.Errorf("unexpected stream error: %v", err)
			}
			if err := stream.Send(&testsvc.Fields{S: &c.str}); err != nil {
				t.Errorf("unexpected send error: %v", err)
			}
			if _, err := stream.Recv(); err != nil {
				t.Errorf("unexpected recv error: %v", err)
			}

			stop()
			reg.AssertHistogram(metricRPCStreamMessageSize, rpcNoCodeLabels, 1, c.expectedBucketCounts)
		})
	}
}

func TestStreamServerInterceptorResponseSize(t *testing.T) {
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
			ctx := Context(context.Background(), "testsvc", WithRegisterer(reg), WithResponseSizeBuckets(buckets))
			sinter := StreamServerInterceptor(ctx)
			cli, stop := testsvc.SetupGRPC(t,
				testsvc.WithServerOptions(grpc.StreamInterceptor(sinter)),
				testsvc.WithStreamFunc(echoMethod()))

			stream, err := cli.GRPCStream(context.Background())
			if err != nil {
				t.Errorf("unexpected stream error: %v", err)
			}
			if err := stream.Send(&testsvc.Fields{S: &c.str}); err != nil {
				t.Errorf("unexpected send error: %v", err)
			}
			if _, err := stream.Recv(); err != nil {
				t.Errorf("unexpected recv error: %v", err)
			}

			stop()
			reg.AssertHistogram(metricRPCStreamResponseSize, rpcNoCodeLabels, 1, c.expectedBucketCounts)
		})
	}
}

func TestStreamServerInterceptorActiveRequests(t *testing.T) {
	cases := []struct {
		name    string
		numReqs int
	}{
		{"one", 1},
		{"ten", 10},
		{"one hundred", 100},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			reg := NewTestRegistry(t)
			ctx := Context(context.Background(), "testsvc", WithRegisterer(reg))
			sinter := StreamServerInterceptor(ctx)
			chstop := make(chan struct{})
			var running, done sync.WaitGroup
			running.Add(c.numReqs)
			done.Add(c.numReqs)
			cli, stop := testsvc.SetupGRPC(t,
				testsvc.WithServerOptions(grpc.StreamInterceptor(sinter)),
				testsvc.WithStreamFunc(recvWaitMethod(&running, &done, chstop)))

			for i := 0; i < c.numReqs; i++ {
				stream, err := cli.GRPCStream(context.Background())
				if err != nil {
					t.Errorf("unexpected stream error: %v", err)
				}
				go func() {
					if err := stream.Send(&testsvc.Fields{}); err != nil {
						t.Errorf("unexpected send error: %v", err)
					}
				}()
			}
			running.Wait()
			reg.AssertGauge(metricRPCActiveRequests, rpcNoCodeLabels, c.numReqs)
			close(chstop)
			done.Wait()
			stop()
		})
	}
}

func noopMethod() testsvc.UnaryFunc {
	return func(_ context.Context, _ *testsvc.Fields) (*testsvc.Fields, error) {
		return &testsvc.Fields{}, nil
	}
}

func stringMethod(str string) testsvc.UnaryFunc {
	return func(_ context.Context, _ *testsvc.Fields) (*testsvc.Fields, error) {
		return &testsvc.Fields{S: &str}, nil
	}
}

func waitMethod(running, done *sync.WaitGroup, stop chan struct{}) testsvc.UnaryFunc {
	return func(_ context.Context, _ *testsvc.Fields) (*testsvc.Fields, error) {
		running.Done()
		defer done.Done()
		<-stop
		return &testsvc.Fields{}, nil
	}
}

func echoMethod() testsvc.StreamFunc {
	return func(_ context.Context, stream testsvc.Stream) (err error) {
		f, err := stream.Recv()
		if err != nil {
			return err
		}
		if err := stream.Send(f); err != nil {
			return err
		}
		return stream.Close()
	}
}

func recvWaitMethod(running, done *sync.WaitGroup, stop chan struct{}) testsvc.StreamFunc {
	return func(_ context.Context, stream testsvc.Stream) (err error) {
		running.Done()
		defer done.Done()
		if _, err := stream.Recv(); err != nil {
			return err
		}
		<-stop
		return stream.Close()
	}
}
