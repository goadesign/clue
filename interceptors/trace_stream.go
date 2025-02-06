package interceptors

import (
	"context"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	goa "goa.design/goa/v3/pkg"
)

type (
	// ClientTraceBidirectionalStreamInfo is an interface that matches the interceptor info for a
	// bidirectional stream that can be traced using ClientTraceBidirectionalStream.
	ClientTraceBidirectionalStreamInfo[Payload TraceStreamStreamingSendMessage, Result TraceStreamStreamingRecvMessage] interface {
		goa.InterceptorInfo

		ClientStreamingPayload() Payload
		ClientStreamingResult(res any) Result
	}

	// ClientTraceStreamDownInfo is an interface that matches the interceptor info for a
	// server to client stream that can be traced using ClientTraceDownStream.
	ClientTraceStreamDownInfo[Result TraceStreamStreamingRecvMessage] interface {
		goa.InterceptorInfo

		ClientStreamingResult(res any) Result
	}

	// ClientTraceStreamUpInfo is an interface that matches the interceptor info for a
	// client to server stream that can be traced using ClientTraceUpStream.
	ClientTraceStreamUpInfo[Payload TraceStreamStreamingSendMessage] interface {
		goa.InterceptorInfo

		ClientStreamingPayload() Payload
	}

	// ServerTraceBidirectionalStreamInfo is an interface that matches the interceptor info for a
	// bidirectional stream that can be traced using ServerTraceBidirectionalStream.
	ServerTraceBidirectionalStreamInfo[Payload TraceStreamStreamingRecvMessage, Result TraceStreamStreamingSendMessage] interface {
		goa.InterceptorInfo

		ServerStreamingPayload(pay any) Payload
		ServerStreamingResult() Result
	}

	// ServerTraceStreamDownInfo is an interface that matches the interceptor info for a
	// server to client stream that can be traced using ServerTraceDownStream.
	ServerTraceStreamDownInfo[Result TraceStreamStreamingSendMessage] interface {
		goa.InterceptorInfo

		ServerStreamingResult() Result
	}

	// ServerTraceStreamUpInfo is an interface that matches the interceptor info for a
	// client to server stream that can be traced using ServerTraceUpStream.
	ServerTraceStreamUpInfo[Payload TraceStreamStreamingRecvMessage] interface {
		goa.InterceptorInfo

		ServerStreamingPayload(pay any) Payload
	}

	// TraceStreamStreamingSendMessage is an interface that matches the streaming send payload or result
	// for a stream that can be traced using ClientTraceBidirectionalStream, ClientTraceUpStream,
	// ServerTraceBidirectionalStream, or ServerTraceDownStream.
	TraceStreamStreamingSendMessage interface {
		SetTraceMetadata(map[string]string)
	}

	// TraceStreamStreamingRecvMessage is an interface that matches the streaming receive payload or result
	// for a stream that can be traced using ClientTraceBidirectionalStream, ClientTraceDownStream,
	// ServerTraceBidirectionalStream, or ServerTraceUpStream.
	TraceStreamStreamingRecvMessage interface {
		TraceMetadata() map[string]string
	}
)

// ClientTraceBidirectionalStream is a client-side interceptor that traces a bidirectional stream by
// injecting the trace metadata into the streaming payload and extracting it from the streaming result.
// The injected trace metadata comes from the context passed to the send method of the client stream.
// The receive method of the client stream returns the extracted trace metadata in its context.
func ClientTraceBidirectionalStream[Payload TraceStreamStreamingSendMessage, Result TraceStreamStreamingRecvMessage](
	ctx context.Context,
	info ClientTraceBidirectionalStreamInfo[Payload, Result],
	next goa.InterceptorEndpoint,
) (any, context.Context, error) {
	switch info.CallType() {
	case goa.InterceptorStreamingRecv:
		return traceStreamRecv(ctx, info, next, info.ClientStreamingResult)
	case goa.InterceptorStreamingSend:
		return traceStreamSend(ctx, info, next, info.ClientStreamingPayload)
	}
	return next(ctx, info.RawPayload())
}

// ClientTraceDownStream is a client-side interceptor that traces a server to client stream by
// extracting the trace metadata from the streaming result. The extracted trace metadata is returned
// in the context of the receive method of the client stream.
func ClientTraceDownStream[Result TraceStreamStreamingRecvMessage](
	ctx context.Context,
	info ClientTraceStreamDownInfo[Result],
	next goa.InterceptorEndpoint,
) (any, context.Context, error) {
	if info.CallType() == goa.InterceptorStreamingRecv {
		return traceStreamRecv(ctx, info, next, info.ClientStreamingResult)
	}
	return next(ctx, info.RawPayload())
}

// ClientTraceUpStream is a client-side interceptor that traces a client to server stream by
// injecting the trace metadata into the streaming payload. The injected trace metadata is returned
// in the context of the send method of the client stream.
func ClientTraceUpStream[Payload TraceStreamStreamingSendMessage](
	ctx context.Context,
	info ClientTraceStreamUpInfo[Payload],
	next goa.InterceptorEndpoint,
) (any, context.Context, error) {
	if info.CallType() == goa.InterceptorStreamingSend {
		return traceStreamSend(ctx, info, next, info.ClientStreamingPayload)
	}
	return next(ctx, info.RawPayload())
}

// ServerTraceBidirectionalStream is a server-side interceptor that traces a bidirectional stream by
// injecting the trace metadata into the streaming result and extracting it from the streaming payload.
// The injected trace metadata comes from the context passed to the send method of the server stream.
// The receive method of the server stream returns the extracted trace metadata in its context.
func ServerTraceBidirectionalStream[Payload TraceStreamStreamingRecvMessage, Result TraceStreamStreamingSendMessage](
	ctx context.Context,
	info ServerTraceBidirectionalStreamInfo[Payload, Result],
	next goa.InterceptorEndpoint,
) (any, context.Context, error) {
	switch info.CallType() {
	case goa.InterceptorStreamingRecv:
		return traceStreamRecv(ctx, info, next, info.ServerStreamingPayload)
	case goa.InterceptorStreamingSend:
		return traceStreamSend(ctx, info, next, info.ServerStreamingResult)
	}
	return next(ctx, info.RawPayload())
}

// ServerTraceDownStream is a server-side interceptor that traces a server to client stream by
// injecting the trace metadata into the streaming result. The injected trace metadata is returned
// in the context of the send method of the server stream.
func ServerTraceDownStream[Result TraceStreamStreamingSendMessage](
	ctx context.Context,
	info ServerTraceStreamDownInfo[Result],
	next goa.InterceptorEndpoint,
) (any, context.Context, error) {
	if info.CallType() == goa.InterceptorStreamingSend {
		return traceStreamSend(ctx, info, next, info.ServerStreamingResult)
	}
	return next(ctx, info.RawPayload())
}

// ServerTraceUpStream is a server-side interceptor that traces a client to server stream by
// extracting the trace metadata from the streaming payload. The extracted trace metadata is returned
// in the context of the receive method of the server stream.
func ServerTraceUpStream[Result TraceStreamStreamingRecvMessage](
	ctx context.Context,
	info ServerTraceStreamUpInfo[Result],
	next goa.InterceptorEndpoint,
) (any, context.Context, error) {
	if info.CallType() == goa.InterceptorStreamingRecv {
		return traceStreamRecv(ctx, info, next, info.ServerStreamingPayload)
	}
	return next(ctx, info.RawPayload())
}

// traceStreamSend is a helper function that traces a stream by injecting the trace metadata
// into the streaming payload or result. The injected trace metadata comes from the context
// passed to the send method of the stream.
func traceStreamSend[Message TraceStreamStreamingSendMessage](
	ctx context.Context,
	info goa.InterceptorInfo,
	next goa.InterceptorEndpoint,
	streamingMessage func() Message,
) (any, context.Context, error) {
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
	next goa.InterceptorEndpoint,
	streamingMessage func(any) Message,
) (any, context.Context, error) {
	msg, ctx, err := next(ctx, info.RawPayload())
	propagator := otel.GetTextMapPropagator()
	sm := streamingMessage(msg)
	ctx = propagator.Extract(ctx, propagation.MapCarrier(sm.TraceMetadata()))
	return msg, ctx, err
}
