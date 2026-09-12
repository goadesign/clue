package debug

import (
	"bytes"
	"context"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"

	"goa.design/clue/log"
)

func TestDebugFailureKeepsBufferedLogs(t *testing.T) {
	debugLogs.Store(false)
	for _, transport := range []string{"unary", "stream"} {
		t.Run(transport, func(t *testing.T) {
			var output bytes.Buffer
			base := log.Context(context.Background(), log.WithOutputs(log.Output{Writer: &output, Format: logKeyValsOnly}))
			failure := errors.New("request failed")
			switch transport {
			case "unary":
				info := &grpc.UnaryServerInfo{FullMethod: "/test/call"}
				_, err := log.UnaryServerInterceptor(base)(context.Background(), nil, info, func(ctx context.Context, req any) (any, error) {
					return UnaryServerInterceptor()(ctx, req, info, func(ctx context.Context, _ any) (any, error) {
						log.Info(ctx, log.KV{K: "step", V: "before failure"})
						return nil, failure
					})
				})
				assert.ErrorIs(t, err, failure)
			case "stream":
				info := &grpc.StreamServerInfo{FullMethod: "/test/call"}
				err := log.StreamServerInterceptor(base)(nil, &streamWithContext{ctx: context.Background()}, info, func(srv any, stream grpc.ServerStream) error {
					return StreamServerInterceptor()(srv, stream, info, func(_ any, stream grpc.ServerStream) error {
						log.Info(stream.Context(), log.KV{K: "step", V: "before failure"})
						return failure
					})
				})
				assert.ErrorIs(t, err, failure)
			}
			assert.Contains(t, output.String(), "step=before failure")
			assert.Contains(t, output.String(), "err=request failed")
		})
	}
}

func TestDebugInitialState(t *testing.T) {
	t.Cleanup(func() {
		debugLogs.Store(false)
	})
	for _, enabled := range []bool{true, false} {
		mux := http.NewServeMux()
		MountDebugLogEnabler(mux, WithInitialState(enabled))
		for _, transport := range []string{"http", "unary", "stream"} {
			serveDebugRequest(t, transport, log.Context(context.Background()), func(ctx context.Context) {
				assert.Equal(t, enabled, log.DebugEnabled(ctx), transport)
			})
		}
		// Without a configured logger, debug middleware creates one itself.
		_, err := UnaryServerInterceptor()(context.Background(), nil, nil, func(ctx context.Context, _ any) (any, error) {
			assert.Equal(t, enabled, log.DebugEnabled(ctx))
			return nil, nil
		})
		assert.NoError(t, err)
		err = StreamServerInterceptor()(nil, &streamWithContext{ctx: context.Background()}, nil, func(_ any, stream grpc.ServerStream) error {
			assert.Equal(t, enabled, log.DebugEnabled(stream.Context()))
			return nil
		})
		assert.NoError(t, err)
		HTTP()(http.HandlerFunc(func(_ http.ResponseWriter, r *http.Request) {
			assert.Equal(t, enabled, log.DebugEnabled(r.Context()))
		})).ServeHTTP(httptest.NewRecorder(), httptest.NewRequest(http.MethodGet, "/", nil))
		// Another listener without a startup choice uses the same process setting.
		other := http.NewServeMux()
		MountDebugLogEnabler(other)
		response := httptest.NewRecorder()
		other.ServeHTTP(response, httptest.NewRequest(http.MethodGet, "/debug", nil))
		want := `{"debug-logs":"off"}`
		if enabled {
			want = `{"debug-logs":"on"}`
		}
		assert.Equal(t, want, response.Body.String())
	}
}

func TestDebugRequestIsolation(t *testing.T) {
	for _, transport := range []string{"http", "unary", "stream"} {
		t.Run(transport, func(t *testing.T) {
			mux := http.NewServeMux()
			MountDebugLogEnabler(mux)
			base := log.Context(context.Background())
			mux.ServeHTTP(httptest.NewRecorder(), httptest.NewRequest(http.MethodGet, "/debug?debug-logs=on", nil))
			var active context.Context
			serveDebugRequest(t, transport, base, func(ctx context.Context) {
				active = ctx
				assert.True(t, log.DebugEnabled(ctx))
			})
			mux.ServeHTTP(httptest.NewRecorder(), httptest.NewRequest(http.MethodGet, "/debug?debug-logs=off", nil))
			serveDebugRequest(t, transport, base, func(ctx context.Context) {
				assert.False(t, log.DebugEnabled(ctx))
			})
			assert.True(t, log.DebugEnabled(active), "later requests must not change an earlier request's setting")
			assert.False(t, log.DebugEnabled(base), "request settings must not change the startup logger")
		})
	}
}

func TestConcurrentDebugRequests(t *testing.T) {
	mux := http.NewServeMux()
	MountDebugLogEnabler(mux)
	base := log.Context(context.Background(), log.WithOutputs(log.Output{Writer: io.Discard, Format: log.FormatJSON}))
	var workers sync.WaitGroup
	for _, transport := range []string{"http", "unary", "stream"} {
		workers.Go(func() {
			for range 100 {
				serveDebugRequest(t, transport, log.With(base), func(ctx context.Context) {
					log.DebugEnabled(ctx)
					log.Debug(ctx, log.KV{K: "message", V: "concurrent request"})
				})
			}
		})
	}
	workers.Go(func() {
		for range 100 {
			for _, state := range []string{"on", "off"} {
				mux.ServeHTTP(httptest.NewRecorder(), httptest.NewRequest(http.MethodGet, "/debug?debug-logs="+state, nil))
			}
		}
	})
	workers.Wait()
	debugLogs.Store(false)
}

// serveDebugRequest applies the real middleware to a request that already has
// a logger, matching the order used by HTTP and gRPC services.
func serveDebugRequest(t *testing.T, transport string, ctx context.Context, handler func(context.Context)) {
	t.Helper()
	switch transport {
	case "http":
		handler := HTTP()(http.HandlerFunc(func(_ http.ResponseWriter, r *http.Request) {
			handler(r.Context())
		}))
		log.HTTP(ctx, log.WithDisableRequestID(), log.WithDisableRequestLogging())(handler).ServeHTTP(httptest.NewRecorder(), httptest.NewRequest(http.MethodGet, "/", nil))
	case "unary":
		info := &grpc.UnaryServerInfo{FullMethod: "/test/call"}
		_, err := log.UnaryServerInterceptor(ctx, log.WithDisableCallID(), log.WithDisableCallLogging())(context.Background(), nil, info, func(ctx context.Context, req any) (any, error) {
			return UnaryServerInterceptor()(ctx, req, info, func(ctx context.Context, _ any) (any, error) {
				handler(ctx)
				return nil, nil
			})
		})
		assert.NoError(t, err)
	case "stream":
		info := &grpc.StreamServerInfo{FullMethod: "/test/call"}
		err := log.StreamServerInterceptor(ctx, log.WithDisableCallID(), log.WithDisableCallLogging())(nil, &streamWithContext{ctx: context.Background()}, info, func(srv any, stream grpc.ServerStream) error {
			return StreamServerInterceptor()(srv, stream, info, func(_ any, stream grpc.ServerStream) error {
				handler(stream.Context())
				return nil
			})
		})
		assert.NoError(t, err)
	}
}
