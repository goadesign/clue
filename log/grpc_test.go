package log

import (
	"bytes"
	"context"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"goa.design/clue/internal/testsvc"
	grpcmiddleware "goa.design/goa/v3/grpc/middleware"
	"goa.design/goa/v3/middleware"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

func TestUnaryServerInterceptor(t *testing.T) {
	now := timeNow
	timeNow = func() time.Time { return time.Date(2022, time.January, 9, 20, 29, 45, 0, time.UTC) }
	defer func() { timeNow = now }()

	var buf bytes.Buffer
	ctx := Context(context.Background(), WithOutput(&buf), WithFormat(FormatJSON))
	logInterceptor := UnaryServerInterceptor(ctx)
	requestIDInterceptor := grpcmiddleware.UnaryRequestID()
	cli, stop := testsvc.SetupGRPC(t,
		testsvc.WithServerOptions(grpc.ChainUnaryInterceptor(requestIDInterceptor, logInterceptor)),
		testsvc.WithUnaryFunc(logUnaryMethod))
	f, err := cli.GRPCMethod(context.Background(), &testsvc.Fields{})
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	stop()

	expected := fmt.Sprintf("{%s,%s,%s%q,%s,%s}\n",
		`"time":"2022-01-09T20:29:45Z"`,
		`"level":"info"`,
		`"request-id":`,
		*f.S,
		`"key1":"value1"`,
		`"key2":"value2"`)

	assert.Equal(t, expected, buf.String())
}

func TestStreamServerTrace(t *testing.T) {
	now := timeNow
	timeNow = func() time.Time { return time.Date(2022, time.January, 9, 20, 29, 45, 0, time.UTC) }
	defer func() { timeNow = now }()

	var buf bytes.Buffer
	ctx := Context(context.Background(), WithOutput(&buf), WithFormat(FormatJSON))
	traceInterceptor := StreamServerInterceptor(ctx)
	requestIDInterceptor := grpcmiddleware.StreamRequestID()
	cli, stop := testsvc.SetupGRPC(t,
		testsvc.WithServerOptions(grpc.ChainStreamInterceptor(requestIDInterceptor, traceInterceptor)),
		testsvc.WithStreamFunc(echoMethod))
	stream, err := cli.GRPCStream(context.Background())
	if err != nil {
		t.Errorf("unexpected stream error: %v", err)
	}
	if err := stream.Send(&testsvc.Fields{}); err != nil {
		t.Errorf("unexpected send error: %v", err)
	}
	f, err := stream.Recv()
	if err != nil {
		t.Errorf("unexpected recv error: %v", err)
	}
	reqID := *f.S
	stop()

	expected := fmt.Sprintf("{%s,%s,%s%q,%s,%s}\n",
		`"time":"2022-01-09T20:29:45Z"`,
		`"level":"info"`,
		`"request-id":`,
		reqID,
		`"key1":"value1"`,
		`"key2":"value2"`)

	assert.Equal(t, expected, buf.String())
}

func TestUnaryClientInterceptor(t *testing.T) {
	successLogs := `time=2022-01-09T20:29:45Z level=info msg="finished client unary call" grpc.service=test.Test grpc.method=GrpcMethod grpc.code=OK grpc.time_ms=42`
	errorLogs := `time=2022-01-09T20:29:45Z level=error err="rpc error: code = Unknown desc = error" msg="finished client unary call" grpc.service=test.Test grpc.method=GrpcMethod grpc.status=error grpc.code=Unknown grpc.time_ms=42`
	statusLogs := `time=2022-01-09T20:29:45Z level=error err="rpc error: code = Unknown desc = error" msg="finished client unary call" grpc.service=test.Test grpc.method=GrpcMethod grpc.status=error grpc.code=Unknown grpc.time_ms=42`
	cases := []struct {
		name      string
		noLog     bool
		clientErr error
		opt       GRPCClientLogOption
		expected  string
	}{
		{"no logger", true, nil, nil, ""},
		{"success", false, nil, nil, successLogs},
		{"error", false, fmt.Errorf("error"), nil, errorLogs},
		{"with status", false, fmt.Errorf("error"), WithErrorFunc(func(codes.Code) bool { return true }), statusLogs},
	}
	now := timeNow
	timeNow = func() time.Time { return time.Date(2022, time.January, 9, 20, 29, 45, 0, time.UTC) }
	defer func() { timeNow = now }()
	duration := 42 * time.Millisecond
	since := timeSince
	timeSince = func(_ time.Time) time.Duration { return duration }
	defer func() { timeSince = since }()

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			var buf bytes.Buffer
			ctx := Context(context.Background(), WithOutput(&buf))
			if c.noLog {
				ctx = context.Background()
			}
			opts := []testsvc.GRPCOption{testsvc.WithUnaryFunc(dummyMethod(c.clientErr))}
			if c.opt != nil {
				opts = append(opts, testsvc.WithDialOptions(
					grpc.WithUnaryInterceptor(UnaryClientInterceptor(c.opt))))
			} else {
				opts = append(opts, testsvc.WithDialOptions(
					grpc.WithUnaryInterceptor(UnaryClientInterceptor())))
			}
			cli, stop := testsvc.SetupGRPC(t, opts...)
			cli.GRPCMethod(ctx, &testsvc.Fields{})
			stop()

			assert.Equal(t, strings.TrimSpace(buf.String()), c.expected)
		})
	}
}

func TestStreamClientInterceptor(t *testing.T) {
	successLogs := `time=2022-01-09T20:29:45Z level=info msg="finished client streaming call" grpc.service=test.Test grpc.method=GrpcStream grpc.code=OK grpc.time_ms=42`
	cases := []struct {
		name     string
		noLog    bool
		expected string
	}{
		{"no logger", true, ""},
		{"success", false, successLogs},
	}
	now := timeNow
	timeNow = func() time.Time { return time.Date(2022, time.January, 9, 20, 29, 45, 0, time.UTC) }
	defer func() { timeNow = now }()
	duration := 42 * time.Millisecond
	since := timeSince
	timeSince = func(_ time.Time) time.Duration { return duration }
	defer func() { timeSince = since }()

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			var buf bytes.Buffer
			ctx := Context(context.Background(), WithOutput(&buf))
			if c.noLog {
				ctx = context.Background()
			}
			cli, stop := testsvc.SetupGRPC(t,
				testsvc.WithDialOptions(grpc.WithStreamInterceptor(StreamClientInterceptor())),
				testsvc.WithStreamFunc(dummyStreamMethod()))

			_, err := cli.GRPCStream(ctx)
			if err != nil {
				t.Errorf("unexpected stream error: %v", err)
			}
			stop()

			assert.Equal(t, strings.TrimSpace(buf.String()), c.expected)
		})
	}
}

func logUnaryMethod(ctx context.Context, _ *testsvc.Fields) (*testsvc.Fields, error) {
	Print(ctx, KV{"key1", "value1"}, KV{"key2", "value2"})
	reqID := ctx.Value(middleware.RequestIDKey).(string)
	return &testsvc.Fields{S: &reqID}, nil
}

func echoMethod(ctx context.Context, stream testsvc.Stream) (err error) {
	Print(ctx, KV{"key1", "value1"}, KV{"key2", "value2"})
	f, err := stream.Recv()
	if err != nil {
		return err
	}
	reqID := ctx.Value(middleware.RequestIDKey).(string)
	f.S = &reqID
	if err := stream.Send(f); err != nil {
		return err
	}
	return stream.Close()
}

func dummyMethod(err error) func(context.Context, *testsvc.Fields) (*testsvc.Fields, error) {
	return func(ctx context.Context, _ *testsvc.Fields) (*testsvc.Fields, error) {
		return &testsvc.Fields{}, err
	}
}

func dummyStreamMethod() func(context.Context, testsvc.Stream) error {
	return func(ctx context.Context, stream testsvc.Stream) error {
		return stream.Close()
	}
}
