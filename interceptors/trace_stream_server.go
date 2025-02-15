package interceptors

import (
	"context"

	goa "goa.design/goa/v3/pkg"
)

type (
	// TraceBidirectionalStreamServerInterceptor is a server-side interceptor that traces a bidirectional stream by
	// injecting the trace metadata into the streaming payload and extracting it from the streaming result.
	// The injected trace metadata comes from the context passed to the send method of the server stream.
	// The receive method of the server stream returns the extracted trace metadata in its context.
	TraceBidirectionalStreamServerInterceptor[Info TraceBidirectionalStreamServerInfo[Payload, Result], Payload TraceStreamStreamingRecvMessage, Result TraceStreamStreamingSendMessage] struct{}

	// TraceServerToClientStreamServerInterceptor is a server-side interceptor that traces a server to client stream by
	// injecting the trace metadata into the streaming result. The injected trace metadata is returned
	// in the context of the send method of the server stream.
	TraceServerToClientStreamServerInterceptor[Info ServerTraceStreamServerToClientInfo[Result], Result TraceStreamStreamingSendMessage] struct{}

	// TraceClientToServerStreamServerInterceptor is a server-side interceptor that traces a client to server stream by
	// extracting the trace metadata from the streaming payload. The extracted trace metadata is returned
	// in the context of the receive method of the server stream.
	TraceClientToServerStreamServerInterceptor[Info ServerTraceStreamClientToServerInfo[Payload], Payload TraceStreamStreamingRecvMessage] struct{}

	// TraceBidirectionalStreamServerInfo is an interface that matches the interceptor info for a
	// bidirectional stream that can be traced using TraceBidirectionalStreamServer.
	TraceBidirectionalStreamServerInfo[Payload TraceStreamStreamingRecvMessage, Result TraceStreamStreamingSendMessage] interface {
		goa.InterceptorInfo

		ServerStreamingPayload(pay any) Payload
		ServerStreamingResult() Result
	}

	// ServerTraceStreamServerToClientInfo is an interface that matches the interceptor info for a
	// server to client stream that can be traced using TraceServerToClientStreamServer.
	ServerTraceStreamServerToClientInfo[Result TraceStreamStreamingSendMessage] interface {
		goa.InterceptorInfo

		ServerStreamingResult() Result
	}

	// ServerTraceStreamClientToServerInfo is an interface that matches the interceptor info for a
	// client to server stream that can be traced using TraceClientToServerStreamServer.
	ServerTraceStreamClientToServerInfo[Payload TraceStreamStreamingRecvMessage] interface {
		goa.InterceptorInfo

		ServerStreamingPayload(pay any) Payload
	}

	// ServerBidirectionalStream is an interface that matches the server stream for a bidirectional
	// stream that can be wrapped with WrapTraceBidirectionalStreamServerStream.
	ServerBidirectionalStream[Payload any, Result any] interface {
		SendWithContext(ctx context.Context, result Result) error
		RecvWithContext(ctx context.Context) (Payload, error)
		Close() error
	}

	// ServerClientToServerStream is an interface that matches the server stream for a client to server
	// stream that can be wrapped with WrapTraceClientToServerStreamServerStream.
	ServerClientToServerStream[Payload any] interface {
		RecvWithContext(ctx context.Context) (Payload, error)
		Close() error
	}

	// ServerClientToServerStreamWithResult is an interface that matches the server stream for a client to server
	// stream that can be wrapped with WrapTraceClientToServerStreamWithResultServerStream.
	ServerClientToServerStreamWithResult[Payload any, Result any] interface {
		CloseAndSendWithContext(ctx context.Context, result Result) error
		RecvWithContext(ctx context.Context) (Payload, error)
	}

	// WrappedTraceStreamServerBidirectionalStream is a wrapper around a server stream for a bidirectional
	// stream that returns the context with the extracted trace metadata, payload or result, and error after
	// calling the receive method of the stream.
	WrappedTraceStreamServerBidirectionalStream[Payload any, Result any] interface {
		// Send sends a payload on the wrapped server stream.
		Send(ctx context.Context, result Result) error
		// RecvAndReturnContext returns the context with the extracted trace metadata, payload or result, and error after
		// calling the receive method of the wrapped server stream.
		RecvAndReturnContext(ctx context.Context) (context.Context, Payload, error)
		// Close closes the wrapped server stream.
		Close() error
	}

	// WrappedTraceStreamServerClientToServerStream is a wrapper around a server stream for a client to server
	// stream that returns the context with the extracted trace metadata, and payload after
	// calling the receive method of the stream.
	WrappedTraceStreamServerClientToServerStream[Payload any] interface {
		// RecvAndReturnContext returns the context with the extracted trace metadata, and payload after
		// calling the receive method of the wrapped server stream.
		RecvAndReturnContext(ctx context.Context) (context.Context, Payload, error)
		// Close closes the wrapped server stream.
		Close() error
	}

	// WrappedTraceStreamServerClientToServerStreamWithResult is a wrapper around a server stream for a client to server
	// stream that returns the context with the extracted trace metadata, and payload after
	// calling the close and send methods of the stream.
	WrappedTraceStreamServerClientToServerStreamWithResult[Payload any, Result any] interface {
		// CloseAndSend closes the wrapped server stream and sends a result.
		CloseAndSend(ctx context.Context, result Result) error
		// RecvAndReturnContext returns the context with the extracted trace metadata, and payload after
		// calling the receive method of the wrapped server stream.
		RecvAndReturnContext(ctx context.Context) (context.Context, Payload, error)
	}

	// wrappedTraceStreamServerBidirectionalStream is a wrapper around a server stream for a bidirectional
	// stream that returns the context with the extracted trace metadata, payload or result, and error after
	// calling the receive method of the stream.
	wrappedTraceStreamServerBidirectionalStream[Payload any, Result any] struct {
		stream ServerBidirectionalStream[Payload, Result]
	}

	// wrappedTraceStreamServerClientToServerStream is a wrapper around a server stream for a client to server
	// stream that returns the context with the extracted trace metadata, and payload after
	// calling the receive method of the stream.
	wrappedTraceStreamServerClientToServerStream[Payload any] struct {
		stream ServerClientToServerStream[Payload]
	}

	// wrappedTraceStreamServerClientToServerStreamWithResult is a wrapper around a server stream for a client to server
	// stream that returns the context with the extracted trace metadata, and payload after
	// calling the close and send methods of the stream.
	wrappedTraceStreamServerClientToServerStreamWithResult[Payload any, Result any] struct {
		stream ServerClientToServerStreamWithResult[Payload, Result]
	}
)

// TraceBidirectionalStream intercepts a bidirectional stream by injecting the trace metadata into the
// streaming payload and extracting it from the streaming result. The injected trace metadata comes from
// the context passed to the send method of the server stream. The receive method of the server stream
// returns the extracted trace metadata in its context.
func (i *TraceBidirectionalStreamServerInterceptor[Info, Payload, Result]) TraceBidirectionalStream(ctx context.Context, info Info, next goa.Endpoint) (any, error) {
	return TraceBidirectionalStreamServer(ctx, info, next)
}

// TraceServerToClientStream intercepts a server to client stream by injecting the trace metadata into the
// streaming result. The injected trace metadata is returned in the context of the send method of
// the server stream.
func (i *TraceServerToClientStreamServerInterceptor[Info, Result]) TraceServerToClientStream(ctx context.Context, info Info, next goa.Endpoint) (any, error) {
	return TraceServerToClientStreamServer(ctx, info, next)
}

// TraceClientToServerStream intercepts a client to server stream by extracting the trace metadata from the
// streaming payload. The extracted trace metadata is returned in the context of the receive method of
// the server stream.
func (i *TraceClientToServerStreamServerInterceptor[Info, Payload]) TraceClientToServerStream(ctx context.Context, info Info, next goa.Endpoint) (any, error) {
	return TraceClientToServerStreamServer(ctx, info, next)
}

// TraceBidirectionalStreamServer is a server-side interceptor that traces a bidirectional stream by
// injecting the trace metadata into the streaming result and extracting it from the streaming payload.
// The injected trace metadata comes from the context passed to the send method of the server stream.
// The receive method of the server stream returns the extracted trace metadata in its context.
func TraceBidirectionalStreamServer[Payload TraceStreamStreamingRecvMessage, Result TraceStreamStreamingSendMessage](
	ctx context.Context,
	info TraceBidirectionalStreamServerInfo[Payload, Result],
	next goa.Endpoint,
) (any, error) {
	switch info.CallType() {
	case goa.InterceptorStreamingRecv:
		return traceStreamRecv(ctx, info, next, info.ServerStreamingPayload)
	case goa.InterceptorStreamingSend:
		return traceStreamSend(ctx, info, next, info.ServerStreamingResult)
	}
	return next(ctx, info.RawPayload())
}

// TraceServerToClientStreamServer is a server-side interceptor that traces a server to client stream by
// injecting the trace metadata into the streaming result. The injected trace metadata is returned
// in the context of the send method of the server stream.
func TraceServerToClientStreamServer[Result TraceStreamStreamingSendMessage](
	ctx context.Context,
	info ServerTraceStreamServerToClientInfo[Result],
	next goa.Endpoint,
) (any, error) {
	if info.CallType() == goa.InterceptorStreamingSend {
		return traceStreamSend(ctx, info, next, info.ServerStreamingResult)
	}
	return next(ctx, info.RawPayload())
}

// TraceClientToServerStreamServer is a server-side interceptor that traces a client to server stream by
// extracting the trace metadata from the streaming payload. The extracted trace metadata is returned
// in the context of the receive method of the server stream.
func TraceClientToServerStreamServer[Result TraceStreamStreamingRecvMessage](
	ctx context.Context,
	info ServerTraceStreamClientToServerInfo[Result],
	next goa.Endpoint,
) (any, error) {
	if info.CallType() == goa.InterceptorStreamingRecv {
		return traceStreamRecv(ctx, info, next, info.ServerStreamingPayload)
	}
	return next(ctx, info.RawPayload())
}

// WrapTraceBidirectionalStreamServerStream wraps a server stream for a bidirectional stream with an
// interface that returns the context with the extracted trace metadata, payload or result, and error after
// calling the receive method of the stream.
func WrapTraceBidirectionalStreamServerStream[Payload any, Result any](
	stream ServerBidirectionalStream[Payload, Result],
) WrappedTraceStreamServerBidirectionalStream[Payload, Result] {
	return &wrappedTraceStreamServerBidirectionalStream[Payload, Result]{stream: stream}
}

// WrapTraceClientToServerStreamServerStream wraps a server stream for a client to server stream with an
// interface that returns the context with the extracted trace metadata, and payload after
// calling the receive method of the stream.
func WrapTraceClientToServerStreamServerStream[Payload any](
	stream ServerClientToServerStream[Payload],
) WrappedTraceStreamServerClientToServerStream[Payload] {
	return &wrappedTraceStreamServerClientToServerStream[Payload]{stream: stream}
}

// WrapTraceClientToServerStreamWithResultServerStream wraps a server stream for a client to server stream with an
// interface that returns the context with the extracted trace metadata, and payload after
// calling the close and send methods of the stream.
func WrapTraceClientToServerStreamWithResultServerStream[Payload any, Result any](
	stream ServerClientToServerStreamWithResult[Payload, Result],
) WrappedTraceStreamServerClientToServerStreamWithResult[Payload, Result] {
	return &wrappedTraceStreamServerClientToServerStreamWithResult[Payload, Result]{stream: stream}
}

// Send sends a result on the wrapped server stream.
func (w *wrappedTraceStreamServerBidirectionalStream[Payload, Result]) Send(ctx context.Context, result Result) error {
	return w.stream.SendWithContext(ctx, result)
}

// RecvAndReturnContext returns the context with the extracted trace metadata, payload or result, and error after
// calling the receive method of the wrapped server stream.
func (w *wrappedTraceStreamServerBidirectionalStream[Payload, Result]) RecvAndReturnContext(ctx context.Context) (context.Context, Payload, error) {
	return traceStreamWrapRecvAndReturnContext(ctx, w.stream.RecvWithContext)
}

// Close closes the wrapped server stream.
func (w *wrappedTraceStreamServerBidirectionalStream[Payload, Result]) Close() error {
	return w.stream.Close()
}

// RecvAndReturnContext returns the context with the extracted trace metadata, and payload after
// calling the receive method of the wrapped server stream.
func (w *wrappedTraceStreamServerClientToServerStream[Payload]) RecvAndReturnContext(ctx context.Context) (context.Context, Payload, error) {
	return traceStreamWrapRecvAndReturnContext(ctx, w.stream.RecvWithContext)
}

// Close closes the wrapped server stream.
func (w *wrappedTraceStreamServerClientToServerStream[Payload]) Close() error {
	return w.stream.Close()
}

// CloseAndSend closes the wrapped server stream and sends a result.
func (w *wrappedTraceStreamServerClientToServerStreamWithResult[Payload, Result]) CloseAndSend(ctx context.Context, result Result) error {
	return w.stream.CloseAndSendWithContext(ctx, result)
}

// RecvAndReturnContext returns the context with the extracted trace metadata, and payload after
// calling the receive method of the wrapped server stream.
func (w *wrappedTraceStreamServerClientToServerStreamWithResult[Payload, Result]) RecvAndReturnContext(ctx context.Context) (context.Context, Payload, error) {
	return traceStreamWrapRecvAndReturnContext(ctx, w.stream.RecvWithContext)
}
