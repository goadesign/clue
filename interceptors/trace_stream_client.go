package interceptors

import (
	"context"

	goa "goa.design/goa/v3/pkg"
)

type (
	// ClientBidirectionalStreamInterceptor is a client-side interceptor that traces a bidirectional stream by
	// injecting the trace metadata into the streaming payload and extracting it from the streaming result.
	// The injected trace metadata comes from the context passed to the send method of the client stream.
	// The receive method of the client stream returns the extracted trace metadata in its context.
	ClientBidirectionalStreamInterceptor[Info ClientTraceBidirectionalStreamInfo[Payload, Result], Payload TraceStreamStreamingSendMessage, Result TraceStreamStreamingRecvMessage] struct{}

	// ClientDownStreamInterceptor is a client-side interceptor that traces a server to client stream by
	// extracting the trace metadata from the streaming result. The extracted trace metadata is returned
	// in the context of the receive method of the client stream.
	ClientDownStreamInterceptor[Info ClientTraceStreamDownInfo[Result], Result TraceStreamStreamingRecvMessage] struct{}

	// ClientUpStreamInterceptor is a client-side interceptor that traces a client to server stream by
	// injecting the trace metadata into the streaming payload. The injected trace metadata is returned
	// in the context of the send method of the client stream.
	ClientUpStreamInterceptor[Info ClientTraceStreamUpInfo[Payload], Payload TraceStreamStreamingSendMessage] struct{}

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

	// ClientBidirectionalStream is an interface that matches the client stream for a bidirectional
	// stream that can be wrapped with WrapTraceStreamClientBidirectionalStream.
	ClientBidirectionalStream[Payload any, Result any] interface {
		SendWithContext(ctx context.Context, payload Payload) error
		RecvWithContext(ctx context.Context) (Result, error)
		Close() error
	}

	// ClientDownStream is an interface that matches the client stream for a server to client
	// stream that can be wrapped with WrapTraceStreamClientDownStream.
	ClientDownStream[Result any] interface {
		RecvWithContext(ctx context.Context) (Result, error)
	}

	// ClientUpStreamWithResult is an interface that matches the client stream for a client to server
	// stream that can be wrapped with WrapTraceStreamClientUpStreamWithResult.
	ClientUpStreamWithResult[Payload any, Result any] interface {
		SendWithContext(ctx context.Context, payload Payload) error
		CloseAndRecvWithContext(ctx context.Context) (Result, error)
	}

	// WrappedTraceStreamClientBidirectionalStream is a wrapper around a client stream for a bidirectional
	// stream that returns the context with the extracted trace metadata, payload or result, and error after
	// calling the receive method of the stream.
	WrappedTraceStreamClientBidirectionalStream[Payload any, Result any] interface {
		// Send sends a payload on the wrapped client stream.
		Send(ctx context.Context, payload Payload) error
		// RecvAndReturnContext returns the context with the extracted trace metadata, payload or result, and error after
		// calling the receive method of the wrapped client stream.
		RecvAndReturnContext(ctx context.Context) (context.Context, Result, error)
		// Close closes the wrapped client stream.
		Close() error
	}

	// WrappedTraceStreamClientDownStream is a wrapper around a client stream for a server to client
	// stream that returns the context with the extracted trace metadata, and result after
	// calling the receive method of the stream.
	WrappedTraceStreamClientDownStream[Result any] interface {
		// RecvAndReturnContext returns the context with the extracted trace metadata, and result after
		// calling the receive method of the wrapped client stream.
		RecvAndReturnContext(ctx context.Context) (context.Context, Result, error)
	}

	// WrappedTraceStreamClientUpStreamWithResult is a wrapper around a client stream for a client to server
	// stream that returns the context with the extracted trace metadata, and result after
	// calling the close and receive methods of the stream.
	WrappedTraceStreamClientUpStreamWithResult[Payload any, Result any] interface {
		// Send sends a payload on the wrapped client stream.
		Send(ctx context.Context, payload Payload) error
		// CloseAndRecvAndReturnContext returns the context with the extracted trace metadata, and result after
		// calling the close and receive methods of the wrapped client stream.
		CloseAndRecvAndReturnContext(ctx context.Context) (context.Context, Result, error)
	}

	// wrappedTraceStreamClientBidirectionalStream is a wrapper around a client stream for a bidirectional
	// stream that returns the context with the extracted trace metadata, payload or result, and error after
	// calling the receive method of the stream.
	wrappedTraceStreamClientBidirectionalStream[Payload any, Result any] struct {
		stream ClientBidirectionalStream[Payload, Result]
	}

	// wrappedTraceStreamClientDownStream is a wrapper around a client stream for a server to client
	// stream that returns the context with the extracted trace metadata, and result after
	// calling the receive method of the stream.
	wrappedTraceStreamClientDownStream[Result any] struct {
		stream ClientDownStream[Result]
	}

	// wrappedTraceStreamClientUpStreamWithResult is a wrapper around a client stream for a client to server
	// stream that returns the context with the extracted trace metadata, and result after
	// calling the close and receive methods of the stream.
	wrappedTraceStreamClientUpStreamWithResult[Payload any, Result any] struct {
		stream ClientUpStreamWithResult[Payload, Result]
	}
)

// TraceBidirectionalStream intercepts a bidirectional stream by injecting the trace metadata into the
// streaming payload and extracting it from the streaming result. The injected trace metadata comes from
// the context passed to the send method of the client stream. The receive method of the client stream
// returns the extracted trace metadata in its context.
func (i *ClientBidirectionalStreamInterceptor[Info, Payload, Result]) TraceBidirectionalStream(ctx context.Context, info Info, next goa.Endpoint) (any, error) {

	return ClientTraceBidirectionalStream(ctx, info, next)
}

// TraceDownStream intercepts a server to client stream by extracting the trace metadata from the
// streaming result. The extracted trace metadata is returned in the context of the receive method of
// the client stream.
func (i *ClientDownStreamInterceptor[Info, Result]) TraceDownStream(ctx context.Context, info Info, next goa.Endpoint) (any, error) {
	return ClientTraceDownStream(ctx, info, next)
}

// TraceUpStream intercepts a client to server stream by injecting the trace metadata into the
// streaming payload. The injected trace metadata is returned in the context of the send method of
// the client stream.
func (i *ClientUpStreamInterceptor[Info, Payload]) TraceUpStream(ctx context.Context, info Info, next goa.Endpoint) (any, error) {
	return ClientTraceUpStream(ctx, info, next)
}

// ClientTraceBidirectionalStream is a client-side interceptor that traces a bidirectional stream by
// injecting the trace metadata into the streaming payload and extracting it from the streaming result.
// The injected trace metadata comes from the context passed to the send method of the client stream.
// The receive method of the client stream returns the extracted trace metadata in its context.
func ClientTraceBidirectionalStream[Payload TraceStreamStreamingSendMessage, Result TraceStreamStreamingRecvMessage](
	ctx context.Context,
	info ClientTraceBidirectionalStreamInfo[Payload, Result],
	next goa.Endpoint,
) (any, error) {
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
	next goa.Endpoint,
) (any, error) {
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
	next goa.Endpoint,
) (any, error) {
	if info.CallType() == goa.InterceptorStreamingSend {
		return traceStreamSend(ctx, info, next, info.ClientStreamingPayload)
	}
	return next(ctx, info.RawPayload())
}

// WrapTraceStreamClientBidirectionalStream wraps a client stream for a bidirectional stream with an
// interface that returns the context with the extracted trace metadata, payload or result, and error after
// calling the receive method of the stream.
func WrapTraceStreamClientBidirectionalStream[Payload any, Result any](
	stream ClientBidirectionalStream[Payload, Result],
) WrappedTraceStreamClientBidirectionalStream[Payload, Result] {
	return &wrappedTraceStreamClientBidirectionalStream[Payload, Result]{stream: stream}
}

// WrapTraceStreamClientDownStream wraps a client stream for a server to client stream with an
// interface that returns the context with the extracted trace metadata, and result after
// calling the receive method of the stream.
func WrapTraceStreamClientDownStream[Result any](
	stream ClientDownStream[Result],
) WrappedTraceStreamClientDownStream[Result] {
	return &wrappedTraceStreamClientDownStream[Result]{stream: stream}
}

// WrapTraceStreamClientUpStreamWithResult wraps a client stream for a client to server stream with an
// interface that returns the context with the extracted trace metadata, and result after
// calling the close and receive methods of the stream.
func WrapTraceStreamClientUpStreamWithResult[Payload any, Result any](
	stream ClientUpStreamWithResult[Payload, Result],
) WrappedTraceStreamClientUpStreamWithResult[Payload, Result] {
	return &wrappedTraceStreamClientUpStreamWithResult[Payload, Result]{stream: stream}
}

// Send sends a payload on the wrapped client stream.
func (w *wrappedTraceStreamClientBidirectionalStream[Payload, Result]) Send(ctx context.Context, payload Payload) error {
	return w.stream.SendWithContext(ctx, payload)
}

// RecvAndReturnContext returns the context with the extracted trace metadata, payload or result, and error after
// calling the receive method of the wrapped client stream.
func (w *wrappedTraceStreamClientBidirectionalStream[Payload, Result]) RecvAndReturnContext(ctx context.Context) (context.Context, Result, error) {
	return traceStreamWrapRecvAndReturnContext(ctx, w.stream.RecvWithContext)
}

// Close closes the wrapped client stream.
func (w *wrappedTraceStreamClientBidirectionalStream[Payload, Result]) Close() error {
	return w.stream.Close()
}

// RecvAndReturnContext returns the context with the extracted trace metadata, and result after
// calling the receive method of the wrapped client stream.
func (w *wrappedTraceStreamClientDownStream[Result]) RecvAndReturnContext(ctx context.Context) (context.Context, Result, error) {
	return traceStreamWrapRecvAndReturnContext(ctx, w.stream.RecvWithContext)
}

// Send sends a payload on the wrapped client stream.
func (w *wrappedTraceStreamClientUpStreamWithResult[Payload, Result]) Send(ctx context.Context, payload Payload) error {
	return w.stream.SendWithContext(ctx, payload)
}

// CloseAndRecvAndReturnContext returns the context with the extracted trace metadata, and result after
// calling the close and receive methods of the wrapped client stream.
func (w *wrappedTraceStreamClientUpStreamWithResult[Payload, Result]) CloseAndRecvAndReturnContext(ctx context.Context) (context.Context, Result, error) {
	return traceStreamWrapRecvAndReturnContext(ctx, w.stream.CloseAndRecvWithContext)
}
