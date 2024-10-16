package log

import (
	"bytes"
	"context"
	"fmt"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"goa.design/clue/internal/testsvc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestUnaryServerInterceptor(t *testing.T) {
	now := timeNow
	timeNow = func() time.Time { return time.Date(2022, time.January, 9, 20, 29, 45, 0, time.UTC) }
	defer func() { timeNow = now }()

	shortID = func() string { return "test-request-id" }
	defer func() { shortID = randShortID }()

	prefix := `{"time":"2022-01-09T20:29:45Z","level":"info","request_id":"test-request-id","msg":"start","grpc.service":"test.Test","grpc.method":"GrpcMethod"}`
	logged := `{"time":"2022-01-09T20:29:45Z","level":"info","request_id":"test-request-id","key1":"value1","key2":"value2"}`
	suffix := `{"time":"2022-01-09T20:29:45Z","level":"info","request_id":"test-request-id","msg":"end","grpc.service":"test.Test","grpc.method":"GrpcMethod","grpc.code":"OK","grpc.time_ms":\d+}`
	errors := `{"time":"2022-01-09T20:29:45Z","level":"error","request_id":"test-request-id","err":"rpc error: code = Unknown desc = test-error","grpc.service":"test.Test","grpc.method":"GrpcMethod","grpc.status":"test-error","grpc.code":"Unknown","grpc.time_ms":\d+}`

	cases := []struct {
		name        string
		options     []GRPCLogOption
		method      func(context.Context, *testsvc.Fields) (*testsvc.Fields, error)
		expected    string
		expectedErr string
	}{
		{
			name:     "default",
			options:  nil,
			method:   logUnaryMethod,
			expected: prefix + "\n" + logged + "\n" + suffix + "\n",
		},
		{
			name:        "with error",
			options:     []GRPCLogOption{WithErrorFunc(func(codes.Code) bool { return true })},
			method:      errorMethod,
			expected:    prefix + "\n" + errors + "\n",
			expectedErr: `rpc error: code = Unknown desc = test-error`,
		},
		{
			name:     "with disable call logging",
			options:  []GRPCLogOption{WithDisableCallLogging()},
			method:   logUnaryMethod,
			expected: logged + "\n",
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			var buf bytes.Buffer
			ctx := Context(context.Background(), WithOutput(&buf), WithFormat(FormatJSON))
			logInterceptor := UnaryServerInterceptor(ctx, c.options...)
			cli, stop := testsvc.SetupGRPC(t,
				testsvc.WithServerOptions(grpc.UnaryInterceptor(logInterceptor)),
				testsvc.WithUnaryFunc(c.method))

			_, err := cli.GRPCMethod(context.Background(), &testsvc.Fields{})

			if c.expectedErr != "" {
				require.Error(t, err)
				assert.Equal(t, c.expectedErr, err.Error())
			} else {
				require.NoError(t, err)
			}
			assert.Regexp(t, regexp.MustCompile(strings.ReplaceAll(c.expected, "\n", "\\n")), buf.String())
			stop()
		})
	}
}

func TestStreamServerTrace(t *testing.T) {
	now := timeNow
	timeNow = func() time.Time { return time.Date(2022, time.January, 9, 20, 29, 45, 0, time.UTC) }
	defer func() { timeNow = now }()

	shortID = func() string { return "test-request-id" }
	defer func() { shortID = randShortID }()

	prefix := `{"time":"2022-01-09T20:29:45Z","level":"info","request_id":"test-request-id","msg":"start","grpc.service":"test.Test","grpc.method":"GrpcStream"}`
	logged := `{"time":"2022-01-09T20:29:45Z","level":"info","request_id":"test-request-id","key1":"value1","key2":"value2"}`
	suffix := `{"time":"2022-01-09T20:29:45Z","level":"info","request_id":"test-request-id","msg":"end","grpc.service":"test.Test","grpc.method":"GrpcStream","grpc.code":"OK","grpc.time_ms":XXX}`
	errors := `{"time":"2022-01-09T20:29:45Z","level":"error","request_id":"test-request-id","err":"rpc error: code = Unknown desc = test-error","grpc.service":"test.Test","grpc.method":"GrpcStream","grpc.status":"test-error","grpc.code":"Unknown","grpc.time_ms":XXX}`

	cases := []struct {
		name        string
		options     []GRPCLogOption
		method      func(context.Context, testsvc.Stream) error
		expected    string
		expectedErr string
	}{
		{
			name:     "default",
			options:  nil,
			method:   echoMethod,
			expected: prefix + "\n" + logged + "\n" + suffix + "\n",
		},
		{
			name:        "with error",
			options:     []GRPCLogOption{WithErrorFunc(func(codes.Code) bool { return true })},
			method:      echoErrorMethod,
			expected:    prefix + "\n" + errors + "\n",
			expectedErr: `rpc error: code = Unknown desc = test-error`,
		},
		{
			name:     "with disable call logging",
			options:  []GRPCLogOption{WithDisableCallLogging()},
			method:   echoMethod,
			expected: logged + "\n",
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			var buf bytes.Buffer
			ctx := Context(context.Background(), WithOutput(&buf), WithFormat(FormatJSON))
			logInterceptor := StreamServerInterceptor(ctx, c.options...)
			cli, stop := testsvc.SetupGRPC(t,
				testsvc.WithServerOptions(grpc.StreamInterceptor(logInterceptor)),
				testsvc.WithStreamFunc(c.method))

			stream, err := cli.GRPCStream(context.Background())

			require.NoError(t, err)
			err = stream.Send(&testsvc.Fields{})
			assert.NoError(t, err)
			_, err = stream.Recv()
			if c.expectedErr != "" {
				require.Error(t, err)
				assert.Equal(t, c.expectedErr, err.Error())
			} else {
				assert.NoError(t, err)
			}
			stop()
			want := regexp.QuoteMeta(c.expected)
			want = strings.ReplaceAll(want, "XXX", `\d+`)
			assert.Regexp(t, want, buf.String())
		})
	}
}

func TestUnaryClientInterceptor(t *testing.T) {
	startLog := `time=2022-01-09T20:29:45Z level=info msg=start grpc.service=test.Test grpc.method=GrpcMethod`
	endLog := `time=2022-01-09T20:29:45Z level=info msg=end grpc.service=test.Test grpc.method=GrpcMethod grpc.code=OK grpc.time_ms=42`
	errorLog := `time=2022-01-09T20:29:45Z level=error err="rpc error: code = Unknown desc = error" grpc.service=test.Test grpc.method=GrpcMethod grpc.status=error grpc.code=Unknown grpc.time_ms=42`
	statusLog := `time=2022-01-09T20:29:45Z level=error err="rpc error: code = Unknown desc = error" grpc.service=test.Test grpc.method=GrpcMethod grpc.status=error grpc.code=Unknown grpc.time_ms=42`

	cases := []struct {
		name      string
		noLog     bool
		clientErr error
		opt       GRPCLogOption
		expected  string
	}{
		{
			name:     "no logger",
			noLog:    true,
			expected: "",
		},
		{
			name:     "success",
			expected: startLog + "\n" + endLog,
		},
		{
			name:      "error",
			clientErr: fmt.Errorf("error"),
			expected:  startLog + "\n" + errorLog,
		},
		{
			name:      "with status",
			clientErr: fmt.Errorf("error"),
			opt:       WithErrorFunc(func(codes.Code) bool { return true }),
			expected:  startLog + "\n" + statusLog,
		},
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
			cli.GRPCMethod(ctx, &testsvc.Fields{}) // nolint:errcheck
			stop()

			assert.Equal(t, c.expected, strings.TrimSpace(buf.String()))
		})
	}
}

func TestStreamClientInterceptor(t *testing.T) {
	startLog := `time=2022-01-09T20:29:45Z level=info msg=start grpc.service=test.Test grpc.method=GrpcStream`
	endLog := `time=2022-01-09T20:29:45Z level=info msg=end grpc.service=test.Test grpc.method=GrpcStream grpc.code=OK grpc.time_ms=42`

	cases := []struct {
		name     string
		noLog    bool
		expected string
	}{
		{"no logger", true, ""},
		{"success", false, startLog + "\n" + endLog},
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
	return &testsvc.Fields{}, nil
}

func errorMethod(ctx context.Context, _ *testsvc.Fields) (*testsvc.Fields, error) {
	return nil, status.Error(codes.Unknown, "test-error")
}

func echoMethod(ctx context.Context, stream testsvc.Stream) (err error) {
	Print(ctx, KV{"key1", "value1"}, KV{"key2", "value2"})
	f, err := stream.Recv()
	if err != nil {
		return err
	}
	if err := stream.Send(f); err != nil {
		return err
	}
	return stream.Close()
}

func echoErrorMethod(ctx context.Context, stream testsvc.Stream) error {
	if _, err := stream.Recv(); err != nil {
		return err
	}
	return status.Error(codes.Unknown, "test-error")
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
