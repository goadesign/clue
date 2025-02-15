package interceptors

import (
	"context"
	"fmt"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"

	goa "goa.design/goa/v3/pkg"
)

type (
	// TraceStreamStreamingSendMessage is an interface that matches the streaming send payload or result
	// for a stream that can be traced using TraceBidirectionalStreamClient, TraceClientToServerStreamClient,
	// TraceBidirectionalStreamServer, or TraceServerToClientStreamServer.
	TraceStreamStreamingSendMessage interface {
		SetTraceMetadata(map[string]string)
	}

	// TraceStreamStreamingRecvMessage is an interface that matches the streaming receive payload or result
	// for a stream that can be traced using TraceBidirectionalStreamClient, TraceServerToClientStreamClient,
	// TraceBidirectionalStreamServer, or TraceClientToServerStreamServer.
	TraceStreamStreamingRecvMessage interface {
		TraceMetadata() map[string]string
	}

	// traceStreamRecvContextKeyType is the type of the context key for the trace metadata extracted from the
	// streaming payload.
	traceStreamRecvContextKeyType struct{}

	// traceStreamRecvContext is a struct that contains the context of the receive method of the stream.
	traceStreamRecvContext struct {
		ctx context.Context
	}
)

// traceStreamRecvContextKey is the context key for the trace metadata extracted from the streaming payload.
var traceStreamRecvContextKey = traceStreamRecvContextKeyType{}

// SetupTraceStreamRecvContext returns a copy of the context that is set up for use with the receive
// method of a stream so that the trace metadata can be extracted from the streaming payload or result.
// After the receive method of the stream returns, the GetTraceStreamRecvContext function can be used to
// retrieve the context with the extracted trace metadata.
func SetupTraceStreamRecvContext(ctx context.Context) context.Context {
	return context.WithValue(ctx, traceStreamRecvContextKey, &traceStreamRecvContext{})
}

// GetTraceStreamRecvContext returns the context with the extracted trace metadata after calling the
// receive method of a stream. The context must have been set up using the SetupTraceStreamRecvContext
// function.
func GetTraceStreamRecvContext(ctx context.Context) context.Context {
	rc, ok := ctx.Value(traceStreamRecvContextKey).(*traceStreamRecvContext)
	if !ok {
		panic(fmt.Errorf("clue interceptors get trace stream receive context method called without prior setup"))
	} else if rc.ctx == nil {
		panic(fmt.Errorf("clue interceptors get trace stream receive context method called without prior interceptor receive method call"))
	}
	return rc.ctx
}

// traceStreamSend is a helper function that traces a stream by injecting the trace metadata
// into the streaming payload or result. The injected trace metadata comes from the context
// passed to the send method of the stream.
func traceStreamSend[Message TraceStreamStreamingSendMessage](
	ctx context.Context,
	info goa.InterceptorInfo,
	next goa.Endpoint,
	streamingMessage func() Message,
) (any, error) {
	propagator := otel.GetTextMapPropagator()
	md := make(propagation.MapCarrier)
	propagator.Inject(ctx, md)
	sm := streamingMessage()
	sm.SetTraceMetadata(md)
	return next(ctx, info.RawPayload())
}

// traceStreamRecv is a helper function that traces a stream by extracting the trace metadata
// from the streaming payload or result. The extracted trace metadata is returned in the context
// of the receive method of the stream.
func traceStreamRecv[Message TraceStreamStreamingRecvMessage](
	ctx context.Context,
	info goa.InterceptorInfo,
	next goa.Endpoint,
	streamingMessage func(any) Message,
) (any, error) {
	msg, err := next(ctx, info.RawPayload())
	propagator := otel.GetTextMapPropagator()
	sm := streamingMessage(msg)
	rc, ok := ctx.Value(traceStreamRecvContextKey).(*traceStreamRecvContext)
	if !ok {
		panic(fmt.Errorf("clue interceptors trace stream receive method called without prior setup (service: %v, method: %v)", info.Service(), info.Method()))
	}
	rc.ctx = propagator.Extract(ctx, propagation.MapCarrier(sm.TraceMetadata()))
	return msg, err
}

// traceStreamWrapRecvAndReturnContext is a helper function for wrapped trace stream receive methods
// that returns the context with the extracted trace metadata, payload or result, and error after
// calling the receive method of the stream.
func traceStreamWrapRecvAndReturnContext[Message any](ctx context.Context, recv func(context.Context) (Message, error)) (context.Context, Message, error) {
	ctx = SetupTraceStreamRecvContext(ctx)
	msg, err := recv(ctx)
	return GetTraceStreamRecvContext(ctx), msg, err
}
