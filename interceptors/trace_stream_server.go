package interceptors

import (
	"context"

	goa "goa.design/goa/v3/pkg"
)

type (
	// ServerBidirectionalStreamInterceptor is a server-side interceptor that traces a bidirectional stream by
	// injecting the trace metadata into the streaming payload and extracting it from the streaming result.
	// The injected trace metadata comes from the context passed to the send method of the server stream.
	// The receive method of the server stream returns the extracted trace metadata in its context.
	ServerBidirectionalStreamInterceptor[Info ServerTraceBidirectionalStreamInfo[Payload, Result], Payload TraceStreamStreamingRecvMessage, Result TraceStreamStreamingSendMessage] struct{}

	// ServerDownStreamInterceptor is a server-side interceptor that traces a server to client stream by
	// injecting the trace metadata into the streaming result. The injected trace metadata is returned
	// in the context of the send method of the server stream.
	ServerDownStreamInterceptor[Info ServerTraceStreamDownInfo[Result], Result TraceStreamStreamingSendMessage] struct{}

	// ServerUpStreamInterceptor is a server-side interceptor that traces a client to server stream by
	// extracting the trace metadata from the streaming payload. The extracted trace metadata is returned
	// in the context of the receive method of the server stream.
	ServerUpStreamInterceptor[Info ServerTraceStreamUpInfo[Payload], Payload TraceStreamStreamingRecvMessage] struct{}

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

	// ServerBidirectionalStream is an interface that matches the server stream for a bidirectional
	// stream that can be wrapped with WrapTraceStreamServerBidirectionalStream.
	ServerBidirectionalStream[Payload any, Result any] interface {
		SendWithContext(ctx context.Context, result Result) error
		RecvWithContext(ctx context.Context) (Payload, error)
		Close() error
	}

	// ServerUpStream is an interface that matches the server stream for a client to server
	// stream that can be wrapped with WrapTraceStreamServerUpStream.
	ServerUpStream[Payload any] interface {
		RecvWithContext(ctx context.Context) (Payload, error)
		Close() error
	}

	// ServerUpStreamWithResult is an interface that matches the server stream for a client to server
	// stream that can be wrapped with WrapTraceStreamServerUpStreamWithResult.
	ServerUpStreamWithResult[Payload any, Result any] interface {
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

	// WrappedTraceStreamServerUpStream is a wrapper around a server stream for a client to server
	// stream that returns the context with the extracted trace metadata, and payload after
	// calling the receive method of the stream.
	WrappedTraceStreamServerUpStream[Payload any] interface {
		// RecvAndReturnContext returns the context with the extracted trace metadata, and payload after
		// calling the receive method of the wrapped server stream.
		RecvAndReturnContext(ctx context.Context) (context.Context, Payload, error)
		// Close closes the wrapped server stream.
		Close() error
	}

	// WrappedTraceStreamServerUpStreamWithResult is a wrapper around a server stream for a client to server
	// stream that returns the context with the extracted trace metadata, and payload after
	// calling the close and send methods of the stream.
	WrappedTraceStreamServerUpStreamWithResult[Payload any, Result any] interface {
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

	// wrappedTraceStreamServerUpStream is a wrapper around a server stream for a client to server
	// stream that returns the context with the extracted trace metadata, and payload after
	// calling the receive method of the stream.
	wrappedTraceStreamServerUpStream[Payload any] struct {
		stream ServerUpStream[Payload]
	}

	// wrappedTraceStreamServerUpStreamWithResult is a wrapper around a server stream for a client to server
	// stream that returns the context with the extracted trace metadata, and payload after
	// calling the close and send methods of the stream.
	wrappedTraceStreamServerUpStreamWithResult[Payload any, Result any] struct {
		stream ServerUpStreamWithResult[Payload, Result]
	}
)

// TraceBidirectionalStream intercepts a bidirectional stream by injecting the trace metadata into the
// streaming payload and extracting it from the streaming result. The injected trace metadata comes from
// the context passed to the send method of the server stream. The receive method of the server stream
// returns the extracted trace metadata in its context.
func (i *ServerBidirectionalStreamInterceptor[Info, Payload, Result]) TraceBidirectionalStream(ctx context.Context, info Info, next goa.Endpoint) (any, error) {
	return ServerTraceBidirectionalStream(ctx, info, next)
}

// TraceDownStream intercepts a server to client stream by injecting the trace metadata into the
// streaming result. The injected trace metadata is returned in the context of the send method of
// the server stream.
func (i *ServerDownStreamInterceptor[Info, Result]) TraceDownStream(ctx context.Context, info Info, next goa.Endpoint) (any, error) {
	return ServerTraceDownStream(ctx, info, next)
}

// TraceUpStream intercepts a client to server stream by extracting the trace metadata from the
// streaming payload. The extracted trace metadata is returned in the context of the receive method of
// the server stream.
func (i *ServerUpStreamInterceptor[Info, Payload]) TraceUpStream(ctx context.Context, info Info, next goa.Endpoint) (any, error) {
	return ServerTraceUpStream(ctx, info, next)
}

// ServerTraceBidirectionalStream is a server-side interceptor that traces a bidirectional stream by
// injecting the trace metadata into the streaming result and extracting it from the streaming payload.
// The injected trace metadata comes from the context passed to the send method of the server stream.
// The receive method of the server stream returns the extracted trace metadata in its context.
func ServerTraceBidirectionalStream[Payload TraceStreamStreamingRecvMessage, Result TraceStreamStreamingSendMessage](
	ctx context.Context,
	info ServerTraceBidirectionalStreamInfo[Payload, Result],
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

// ServerTraceDownStream is a server-side interceptor that traces a server to client stream by
// injecting the trace metadata into the streaming result. The injected trace metadata is returned
// in the context of the send method of the server stream.
func ServerTraceDownStream[Result TraceStreamStreamingSendMessage](
	ctx context.Context,
	info ServerTraceStreamDownInfo[Result],
	next goa.Endpoint,
) (any, error) {
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
	next goa.Endpoint,
) (any, error) {
	if info.CallType() == goa.InterceptorStreamingRecv {
		return traceStreamRecv(ctx, info, next, info.ServerStreamingPayload)
	}
	return next(ctx, info.RawPayload())
}

// WrapTraceStreamServerBidirectionalStream wraps a server stream for a bidirectional stream with an
// interface that returns the context with the extracted trace metadata, payload or result, and error after
// calling the receive method of the stream.
func WrapTraceStreamServerBidirectionalStream[Payload any, Result any](
	stream ServerBidirectionalStream[Payload, Result],
) WrappedTraceStreamServerBidirectionalStream[Payload, Result] {
	return &wrappedTraceStreamServerBidirectionalStream[Payload, Result]{stream: stream}
}

// WrapTraceStreamServerUpStream wraps a server stream for a client to server stream with an
// interface that returns the context with the extracted trace metadata, and payload after
// calling the receive method of the stream.
func WrapTraceStreamServerUpStream[Payload any](
	stream ServerUpStream[Payload],
) WrappedTraceStreamServerUpStream[Payload] {
	return &wrappedTraceStreamServerUpStream[Payload]{stream: stream}
}

// WrapTraceStreamServerUpStreamWithResult wraps a server stream for a client to server stream with an
// interface that returns the context with the extracted trace metadata, and payload after
// calling the close and send methods of the stream.
func WrapTraceStreamServerUpStreamWithResult[Payload any, Result any](
	stream ServerUpStreamWithResult[Payload, Result],
) WrappedTraceStreamServerUpStreamWithResult[Payload, Result] {
	return &wrappedTraceStreamServerUpStreamWithResult[Payload, Result]{stream: stream}
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
func (w *wrappedTraceStreamServerUpStream[Payload]) RecvAndReturnContext(ctx context.Context) (context.Context, Payload, error) {
	return traceStreamWrapRecvAndReturnContext(ctx, w.stream.RecvWithContext)
}

// Close closes the wrapped server stream.
func (w *wrappedTraceStreamServerUpStream[Payload]) Close() error {
	return w.stream.Close()
}

// CloseAndSend closes the wrapped server stream and sends a result.
func (w *wrappedTraceStreamServerUpStreamWithResult[Payload, Result]) CloseAndSend(ctx context.Context, result Result) error {
	return w.stream.CloseAndSendWithContext(ctx, result)
}

// RecvAndReturnContext returns the context with the extracted trace metadata, and payload after
// calling the receive method of the wrapped server stream.
func (w *wrappedTraceStreamServerUpStreamWithResult[Payload, Result]) RecvAndReturnContext(ctx context.Context) (context.Context, Payload, error) {
	return traceStreamWrapRecvAndReturnContext(ctx, w.stream.RecvWithContext)
}
