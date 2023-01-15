package debug

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"google.golang.org/grpc"

	"goa.design/clue/internal/testsvc"
	"goa.design/clue/log"
)

func TestUnaryServerInterceptor(t *testing.T) {
	var buf bytes.Buffer
	ctx := log.Context(context.Background(), log.WithOutput(&buf), log.WithFormat(logKeyValsOnly))
	cli, stop := testsvc.SetupGRPC(t,
		testsvc.WithServerOptions(grpc.ChainUnaryInterceptor(log.UnaryServerInterceptor(ctx), UnaryServerInterceptor())),
		testsvc.WithUnaryFunc(logUnaryMethod))
	defer stop()
	mux := http.NewServeMux()
	MountDebugLogEnabler(mux)
	ts := httptest.NewServer(mux)
	defer ts.Close()

	steps := []struct {
		name         string
		on           bool
		off          bool
		expectedLogs string
	}{
		{"start", false, false, ""},
		{"turn debug logs on", true, false, "debug=message "},
		{"with debug logs on", false, false, "debug=message "},
		{"turn debug logs off", false, true, ""},
		{"with debug logs off", false, false, ""},
	}
	for _, step := range steps {
		if step.on {
			makeRequest(t, ts.URL+"/debug?debug-logs=on")
		}
		if step.off {
			makeRequest(t, ts.URL+"/debug?debug-logs=off")
		}
		_, err := cli.GRPCMethod(context.Background(), nil)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if buf.String() != step.expectedLogs {
			t.Errorf("expected log %q, got %q", step.expectedLogs, buf.String())
		}
		buf.Reset()
	}
}

func TestStreamServerInterceptor(t *testing.T) {
	var buf bytes.Buffer
	ctx := log.Context(context.Background(), log.WithOutput(&buf), log.WithFormat(logKeyValsOnly))
	cli, stop := testsvc.SetupGRPC(t,
		testsvc.WithServerOptions(grpc.ChainStreamInterceptor(log.StreamServerInterceptor(ctx), StreamServerInterceptor())),
		testsvc.WithStreamFunc(echoMethod))
	defer stop()
	steps := []struct {
		name            string
		enableDebugLogs bool
		expectedLogs    string
	}{
		{"no debug logs", false, ""},
		{"debug logs", true, "debug=message "},
		{"revert to no debug logs", false, ""},
	}
	for _, step := range steps {
		debugLogs = step.enableDebugLogs
		stream, err := cli.GRPCStream(context.Background())
		if err != nil {
			t.Errorf("%s: unexpected error: %v", step.name, err)
		}
		defer stream.Close()
		if err = stream.Send(&testsvc.Fields{}); err != nil {
			t.Errorf("%s: unexpected send error: %v", step.name, err)
		}
		if _, err = stream.Recv(); err != nil {
			t.Errorf("%s: unexpected recv error: %v", step.name, err)
		}
		if buf.String() != step.expectedLogs {
			t.Errorf("%s: unexpected log %q", step.name, buf.String())
		}
		buf.Reset()
	}
}

func logUnaryMethod(ctx context.Context, _ *testsvc.Fields) (*testsvc.Fields, error) {
	log.Debug(ctx, log.KV{K: "debug", V: "message"})
	return &testsvc.Fields{}, nil
}

func echoMethod(ctx context.Context, stream testsvc.Stream) (err error) {
	for {
		_, err := stream.Recv()
		if err != nil {
			return err
		}
		log.Debug(ctx, log.KV{K: "debug", V: "message"})
		if err = stream.Send(&testsvc.Fields{}); err != nil {
			return err
		}
	}
}
