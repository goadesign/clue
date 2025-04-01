package debug

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"

	"goa.design/clue/internal/testsvc"
	"goa.design/clue/log"
)

func TestUnaryServerInterceptor(t *testing.T) {
	var buf bytes.Buffer
	ctx := log.Context(context.Background(), log.WithOutput(&buf), log.WithFormat(logKeyValsOnly))
	cli, stop := testsvc.SetupGRPC(t,
		testsvc.WithServerOptions(grpc.ChainUnaryInterceptor(
			log.UnaryServerInterceptor(ctx, log.WithDisableCallLogging(), log.WithDisableCallID()),
			UnaryServerInterceptor())),
		testsvc.WithUnaryFunc(logUnaryMethod))
	defer stop()
	mux := http.NewServeMux()
	MountDebugLogEnabler(mux)
	ts := httptest.NewServer(mux)
	defer ts.Close()

	steps := []struct {
		name    string
		on      bool
		off     bool
		wantLog string
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
		assert.NoError(t, err)
		assert.Equal(t, step.wantLog, buf.String())
		buf.Reset()
	}
}

func TestStreamServerInterceptor(t *testing.T) {
	var buf bytes.Buffer
	ctx := log.Context(context.Background(), log.WithOutput(&buf), log.WithFormat(logKeyValsOnly))
	cli, stop := testsvc.SetupGRPC(t,
		testsvc.WithServerOptions(grpc.ChainStreamInterceptor(
			log.StreamServerInterceptor(ctx, log.WithDisableCallLogging(), log.WithDisableCallID()),
			StreamServerInterceptor())),
		testsvc.WithStreamFunc(echoMethod))
	defer stop()
	steps := []struct {
		name            string
		enableDebugLogs bool
		wantLog         string
	}{
		{"no debug logs", false, ""},
		{"debug logs", true, "debug=message "},
		{"revert to no debug logs", false, ""},
	}
	for _, step := range steps {
		debugLogs = step.enableDebugLogs
		stream, err := cli.GRPCStream(context.Background())
		assert.NoError(t, err)
		defer func() {
			err := stream.Close()
			if err != nil {
				t.Logf("failed to close stream: %v", err)
			}
		}()
		err = stream.Send(&testsvc.Fields{})
		assert.NoError(t, err)
		_, err = stream.Recv()
		assert.NoError(t, err)
		assert.Equal(t, step.wantLog, buf.String())
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
