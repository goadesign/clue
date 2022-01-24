package log

import (
	"bytes"
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/crossnokaye/micro/internal/testsvc"
	grpcmiddleware "goa.design/goa/v3/grpc/middleware"
	"goa.design/goa/v3/middleware"
	"google.golang.org/grpc"
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

	expected := fmt.Sprintf("{%s,%s,%s,%s%q,%s,%s}\n",
		`"level":"INFO"`,
		`"time":"2022-01-09T20:29:45Z"`,
		`"msg":"hello world"`,
		`"request_id":`,
		*f.S,
		`"key1":"value1"`,
		`"key2":"value2"`)

	if buf.String() != expected {
		t.Errorf("got %s, want %s", buf.String(), expected)
	}
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

	expected := fmt.Sprintf("{%s,%s,%s,%s%q,%s,%s}\n",
		`"level":"INFO"`,
		`"time":"2022-01-09T20:29:45Z"`,
		`"msg":"hello world"`,
		`"request_id":`,
		reqID,
		`"key1":"value1"`,
		`"key2":"value2"`)

	if buf.String() != expected {
		t.Errorf("got %s, want %s", buf.String(), expected)
	}
}

func logUnaryMethod(ctx context.Context, _ *testsvc.Fields) (*testsvc.Fields, error) {
	Print(ctx, "hello world", "key1", "value1", "key2", "value2")
	reqID := ctx.Value(middleware.RequestIDKey).(string)
	return &testsvc.Fields{S: &reqID}, nil
}

func echoMethod(ctx context.Context, stream testsvc.Stream) (err error) {
	Print(ctx, "hello world", "key1", "value1", "key2", "value2")
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
