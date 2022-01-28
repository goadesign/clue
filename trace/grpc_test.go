package trace

import (
	"context"
	"testing"

	"github.com/goadesign/clue/internal/testsvc"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/sdk/trace/tracetest"
	"goa.design/goa/v3/grpc/middleware"
	"google.golang.org/grpc"
)

// NOTE: We are not testing otel here, just make sure a span exists and that the
// request ID is in the attributes on the server.

func TestUnaryServerTrace(t *testing.T) {
	exporter := tracetest.NewInMemoryExporter()
	provider := sdktrace.NewTracerProvider(sdktrace.WithSyncer(exporter))
	traceInterceptor := UnaryServerInterceptor(withProvider(context.Background(), provider, "test"))
	requestIDInterceptor := middleware.UnaryRequestID()
	cli, stop := testsvc.SetupGRPC(t,
		testsvc.WithServerOptions(grpc.ChainUnaryInterceptor(requestIDInterceptor, traceInterceptor)),
		testsvc.WithUnaryFunc(addEventUnaryMethod))
	_, err := cli.GRPCMethod(context.Background(), &testsvc.Fields{})
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	stop()
	spans := exporter.GetSpans()
	if len(spans) != 1 {
		t.Fatalf("got %d spans, want 1", len(spans))
	}
	found := false
	for _, att := range spans[0].Attributes {
		if att.Key == AttributeRequestID {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("request ID not in span attributes")
	}
	events := spans[0].Events
	if len(events) != 3 {
		t.Fatalf("got %d events, want 3", len(events))
	}
	if events[0].Name != "message" {
		t.Errorf("got event name %s, want message", events[0].Name)
	}
	if events[1].Name != "unary method" {
		t.Errorf("unexpected event name: %s", events[0].Name)
	}
	if events[2].Name != "message" {
		t.Errorf("got event name %s, want message", events[0].Name)
	}
}

func TestUnaryServerTraceNoRequestID(t *testing.T) {
	exporter := tracetest.NewInMemoryExporter()
	provider := sdktrace.NewTracerProvider(sdktrace.WithSyncer(exporter))
	traceInterceptor := UnaryServerInterceptor(withProvider(context.Background(), provider, "test"))
	cli, stop := testsvc.SetupGRPC(t,
		testsvc.WithServerOptions(grpc.UnaryInterceptor(traceInterceptor)),
		testsvc.WithUnaryFunc(addEventUnaryMethod))
	_, err := cli.GRPCMethod(context.Background(), &testsvc.Fields{})
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	stop()
	spans := exporter.GetSpans()
	if len(spans) != 1 {
		t.Fatalf("got %d spans, want 1", len(spans))
	}
}

func TestStreamServerTrace(t *testing.T) {
	exporter := tracetest.NewInMemoryExporter()
	provider := sdktrace.NewTracerProvider(sdktrace.WithSyncer(exporter))
	traceInterceptor := StreamServerInterceptor(withProvider(context.Background(), provider, "test"))
	requestIDInterceptor := middleware.StreamRequestID()
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
	stream.Recv()
	stop()
	spans := exporter.GetSpans()
	if len(spans) != 1 {
		t.Fatalf("got %d spans, want 1", len(spans))
	}
	found := false
	for _, att := range spans[0].Attributes {
		if att.Key == AttributeRequestID {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("request ID not in span attributes")
	}
}

func TestStreamServerTraceNoRequestID(t *testing.T) {
	exporter := tracetest.NewInMemoryExporter()
	provider := sdktrace.NewTracerProvider(sdktrace.WithSyncer(exporter))
	traceInterceptor := StreamServerInterceptor(withProvider(context.Background(), provider, "test"))
	cli, stop := testsvc.SetupGRPC(t,
		testsvc.WithServerOptions(grpc.StreamInterceptor(traceInterceptor)),
		testsvc.WithStreamFunc(echoMethod))
	stream, err := cli.GRPCStream(context.Background())
	if err != nil {
		t.Errorf("unexpected stream error: %v", err)
	}
	if err := stream.Send(&testsvc.Fields{}); err != nil {
		t.Errorf("unexpected send error: %v", err)
	}
	stream.Recv()
	stop()
	spans := exporter.GetSpans()
	if len(spans) != 1 {
		t.Fatalf("got %d spans, want 1", len(spans))
	}
}

func TestUnaryClientTrace(t *testing.T) {
	exporter := tracetest.NewInMemoryExporter()
	provider := sdktrace.NewTracerProvider(sdktrace.WithSyncer(exporter))
	cli, stop := testsvc.SetupGRPC(t,
		testsvc.WithDialOptions(grpc.WithUnaryInterceptor(UnaryClientInterceptor(withProvider(context.Background(), provider, "test")))),
		testsvc.WithUnaryFunc(addEventUnaryMethod))
	_, err := cli.GRPCMethod(context.Background(), &testsvc.Fields{})
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	stop()
	spans := exporter.GetSpans()
	if len(spans) != 1 {
		t.Fatalf("got %d spans, want 1", len(spans))
	}
}

func TestStreamClientTrace(t *testing.T) {
	exporter := tracetest.NewInMemoryExporter()
	provider := sdktrace.NewTracerProvider(sdktrace.WithSyncer(exporter))
	cli, stop := testsvc.SetupGRPC(t,
		testsvc.WithDialOptions(grpc.WithStreamInterceptor(StreamClientInterceptor(withProvider(context.Background(), provider, "test")))),
		testsvc.WithStreamFunc(echoMethod))
	ctx, cancel := context.WithCancel(context.Background())
	stream, err := cli.GRPCStream(ctx)
	if err != nil {
		t.Errorf("unexpected stream error: %v", err)
	}
	if err := stream.Send(&testsvc.Fields{}); err != nil {
		t.Errorf("unexpected send error: %v", err)
	}
	cancel()
	stop()
	spans := exporter.GetSpans()
	if len(spans) != 1 {
		t.Fatalf("got %d spans, want 1", len(spans))
	}
}

func addEventUnaryMethod(ctx context.Context, _ *testsvc.Fields) (*testsvc.Fields, error) {
	AddEvent(ctx, "unary method")
	return &testsvc.Fields{}, nil
}

func echoMethod(_ context.Context, stream testsvc.Stream) (err error) {
	f, err := stream.Recv()
	if err != nil {
		return err
	}
	if err := stream.Send(f); err != nil {
		return err
	}
	return stream.Close()
}
