package interceptors

import (
	"context"

	goa "goa.design/goa/v3/pkg"
)

type (
	// TraceBidirectionalStreamClientInterceptor is a client-side interceptor that traces a bidirectional stream by
	// injecting the trace metadata into the streaming payload and extracting it from the streaming result.
	// The injected trace metadata comes from the context passed to the send method of the client stream.
	// The receive method of the client stream returns the extracted trace metadata in its context.
	TraceBidirectionalStreamClientInterceptor[Info TraceBidirectionalStreamClientInfo[Payload, Result], Payload TraceStreamStreamingSendMessage, Result TraceStreamStreamingRecvMessage] struct{}

	// TraceServerToClientStreamClientInterceptor is a client-side interceptor that traces a server to client stream by
	// extracting the trace metadata from the streaming result. The extracted trace metadata is returned
	// in the context of the receive method of the client stream.
	TraceServerToClientStreamClientInterceptor[Info ClientTraceStreamServerToClientInfo[Result], Result TraceStreamStreamingRecvMessage] struct{}

	// TraceClientToServerStreamClientInterceptor is a client-side interceptor that traces a client to server stream by
	// injecting the trace metadata into the streaming payload. The injected trace metadata is returned
	// in the context of the send method of the client stream.
	TraceClientToServerStreamClientInterceptor[Info ClientTraceStreamClientToServerInfo[Payload], Payload TraceStreamStreamingSendMessage] struct{}

	// TraceBidirectionalStreamClientInfo is an interface that matches the interceptor info for a
	// bidirectional stream that can be traced using TraceBidirectionalStreamClient.
	TraceBidirectionalStreamClientInfo[Payload TraceStreamStreamingSendMessage, Result TraceStreamStreamingRecvMessage] interface {
		goa.InterceptorInfo

		ClientStreamingPayload() Payload
		ClientStreamingResult(res any) Result
	}

	// ClientTraceStreamServerToClientInfo is an interface that matches the interceptor info for a
	// server to client stream that can be traced using TraceServerToClientStreamClient.
	ClientTraceStreamServerToClientInfo[Result TraceStreamStreamingRecvMessage] interface {
		goa.InterceptorInfo

		ClientStreamingResult(res any) Result
	}

	// ClientTraceStreamClientToServerInfo is an interface that matches the interceptor info for a
	// client to server stream that can be traced using TraceClientToServerStreamClient.
	ClientTraceStreamClientToServerInfo[Payload TraceStreamStreamingSendMessage] interface {
		goa.InterceptorInfo

		ClientStreamingPayload() Payload
	}

	// ClientBidirectionalStream is an interface that matches the client stream for a bidirectional
	// stream that can be wrapped with WrapTraceBidirectionalStreamClientStream.
	ClientBidirectionalStream[Payload any, Result any] interface {
		SendWithContext(ctx context.Context, payload Payload) error
		RecvWithContext(ctx context.Context) (Result, error)
		Close() error
	}

	// ClientServerToClientStream is an interface that matches the client stream for a server to client
	// stream that can be wrapped with WrapTraceServerToClientStreamClientStream.
	ClientServerToClientStream[Result any] interface {
		RecvWithContext(ctx context.Context) (Result, error)
	}

	// ClientClientToServerStreamWithResult is an interface that matches the client stream for a client to server
	// stream that can be wrapped with WrapTraceClientToServerStreamWithResultClientStream.
	ClientClientToServerStreamWithResult[Payload any, Result any] interface {
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

	// WrappedTraceStreamClientServerToClientStream is a wrapper around a client stream for a server to client
	// stream that returns the context with the extracted trace metadata, and result after
	// calling the receive method of the stream.
	WrappedTraceStreamClientServerToClientStream[Result any] interface {
		// RecvAndReturnContext returns the context with the extracted trace metadata, and result after
		// calling the receive method of the wrapped client stream.
		RecvAndReturnContext(ctx context.Context) (context.Context, Result, error)
	}

	// WrappedTraceStreamClientClientToServerStreamWithResult is a wrapper around a client stream for a client to server
	// stream that returns the context with the extracted trace metadata, and result after
	// calling the close and receive methods of the stream.
	WrappedTraceStreamClientClientToServerStreamWithResult[Payload any, Result any] interface {
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

	// wrappedTraceStreamClientServerToClientStream is a wrapper around a client stream for a server to client
	// stream that returns the context with the extracted trace metadata, and result after
	// calling the receive method of the stream.
	wrappedTraceStreamClientServerToClientStream[Result any] struct {
		stream ClientServerToClientStream[Result]
	}

	// wrappedTraceStreamClientClientToServerStreamWithResult is a wrapper around a client stream for a client to server
	// stream that returns the context with the extracted trace metadata, and result after
	// calling the close and receive methods of the stream.
	wrappedTraceStreamClientClientToServerStreamWithResult[Payload any, Result any] struct {
		stream ClientClientToServerStreamWithResult[Payload, Result]
	}
)

// TraceBidirectionalStream intercepts a bidirectional stream by injecting the trace metadata into the
// streaming payload and extracting it from the streaming result. The injected trace metadata comes from
// the context passed to the send method of the client stream. The receive method of the client stream
// returns the extracted trace metadata in its context.
func (i *TraceBidirectionalStreamClientInterceptor[Info, Payload, Result]) TraceBidirectionalStream(ctx context.Context, info Info, next goa.Endpoint) (any, error) {
	return TraceBidirectionalStreamClient(ctx, info, next)
}

// TraceServerToClientStream intercepts a server to client stream by extracting the trace metadata from the
// streaming result. The extracted trace metadata is returned in the context of the receive method of
// the client stream.
func (i *TraceServerToClientStreamClientInterceptor[Info, Result]) TraceServerToClientStream(ctx context.Context, info Info, next goa.Endpoint) (any, error) {
	return TraceServerToClientStreamClient(ctx, info, next)
}

// TraceClientToServerStream intercepts a client to server stream by injecting the trace metadata into the
// streaming payload. The injected trace metadata is returned in the context of the send method of
// the client stream.
func (i *TraceClientToServerStreamClientInterceptor[Info, Payload]) TraceClientToServerStream(ctx context.Context, info Info, next goa.Endpoint) (any, error) {
	return TraceClientToServerStreamClient(ctx, info, next)
}

// TraceBidirectionalStreamClient is a client-side interceptor that traces a bidirectional stream by
// injecting the trace metadata into the streaming payload and extracting it from the streaming result.
// The injected trace metadata comes from the context passed to the send method of the client stream.
// The receive method of the client stream returns the extracted trace metadata in its context.
func TraceBidirectionalStreamClient[Payload TraceStreamStreamingSendMessage, Result TraceStreamStreamingRecvMessage](
	ctx context.Context,
	info TraceBidirectionalStreamClientInfo[Payload, Result],
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

// TraceServerToClientStreamClient is a client-side interceptor that traces a server to client stream by
// extracting the trace metadata from the streaming result. The extracted trace metadata is returned
// in the context of the receive method of the client stream.
func TraceServerToClientStreamClient[Result TraceStreamStreamingRecvMessage](
	ctx context.Context,
	info ClientTraceStreamServerToClientInfo[Result],
	next goa.Endpoint,
) (any, error) {
	if info.CallType() == goa.InterceptorStreamingRecv {
		return traceStreamRecv(ctx, info, next, info.ClientStreamingResult)
	}
	return next(ctx, info.RawPayload())
}

// TraceClientToServerStreamClient is a client-side interceptor that traces a client to server stream by
// injecting the trace metadata into the streaming payload. The injected trace metadata is returned
// in the context of the send method of the client stream.
func TraceClientToServerStreamClient[Payload TraceStreamStreamingSendMessage](
	ctx context.Context,
	info ClientTraceStreamClientToServerInfo[Payload],
	next goa.Endpoint,
) (any, error) {
	if info.CallType() == goa.InterceptorStreamingSend {
		return traceStreamSend(ctx, info, next, info.ClientStreamingPayload)
	}
	return next(ctx, info.RawPayload())
}

// WrapTraceBidirectionalStreamClientStream wraps a client stream for a bidirectional stream with an
// interface that returns the context with the extracted trace metadata, payload or result, and error after
// calling the receive method of the stream.
func WrapTraceBidirectionalStreamClientStream[Payload any, Result any](
	stream ClientBidirectionalStream[Payload, Result],
) WrappedTraceStreamClientBidirectionalStream[Payload, Result] {
	return &wrappedTraceStreamClientBidirectionalStream[Payload, Result]{stream: stream}
}

// WrapTraceServerToClientStreamClientStream wraps a client stream for a server to client stream with an
// interface that returns the context with the extracted trace metadata, and result after
// calling the receive method of the stream.
func WrapTraceServerToClientStreamClientStream[Result any](
	stream ClientServerToClientStream[Result],
) WrappedTraceStreamClientServerToClientStream[Result] {
	return &wrappedTraceStreamClientServerToClientStream[Result]{stream: stream}
}

// WrapTraceClientToServerStreamWithResultClientStream wraps a client stream for a client to server stream with an
// interface that returns the context with the extracted trace metadata, and result after
// calling the close and receive methods of the stream.
func WrapTraceClientToServerStreamWithResultClientStream[Payload any, Result any](
	stream ClientClientToServerStreamWithResult[Payload, Result],
) WrappedTraceStreamClientClientToServerStreamWithResult[Payload, Result] {
	return &wrappedTraceStreamClientClientToServerStreamWithResult[Payload, Result]{stream: stream}
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
func (w *wrappedTraceStreamClientServerToClientStream[Result]) RecvAndReturnContext(ctx context.Context) (context.Context, Result, error) {
	return traceStreamWrapRecvAndReturnContext(ctx, w.stream.RecvWithContext)
}

// Send sends a payload on the wrapped client stream.
func (w *wrappedTraceStreamClientClientToServerStreamWithResult[Payload, Result]) Send(ctx context.Context, payload Payload) error {
	return w.stream.SendWithContext(ctx, payload)
}

// CloseAndRecvAndReturnContext returns the context with the extracted trace metadata, and result after
// calling the close and receive methods of the wrapped client stream.
func (w *wrappedTraceStreamClientClientToServerStreamWithResult[Payload, Result]) CloseAndRecvAndReturnContext(ctx context.Context) (context.Context, Result, error) {
	return traceStreamWrapRecvAndReturnContext(ctx, w.stream.CloseAndRecvWithContext)
}
