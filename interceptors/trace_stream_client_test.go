package interceptors

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel/baggage"
	"go.opentelemetry.io/otel/trace"

	"goa.design/clue/log"
	goa "goa.design/goa/v3/pkg"
)

func TestTraceBidirectionalStreamClientInterceptor(t *testing.T) {
	t.Run("send", func(t *testing.T) {
		var (
			assert      = assert.New(t)
			require     = require.New(t)
			ctx         = log.Context(context.Background(), log.WithFormat(log.FormatText))
			info        = newMockTraceStreamInfo(assert)
			payload     = newMockTraceStreamingSendMessage(assert)
			interceptor = &TraceBidirectionalStreamClientInterceptor[*mockTraceStreamInfo, *mockTraceStreamingSendMessage, *mockTraceStreamingRecvMessage]{}
			nextCalled  = false
			next        = func(ctx context.Context, payload any) (any, error) {
				nextCalled = true
				assert.Equal(&traceStreamMessage{String: "abc", Int: 1}, payload)
				return nil, nil
			}
		)
		info.addCallType(func() goa.InterceptorCallType {
			return goa.InterceptorStreamingSend
		})
		info.addClientStreamingPayload(func() *mockTraceStreamingSendMessage {
			return payload
		})
		payload.addSetTraceMetadata(func(metadata map[string]string) {
			assert.Contains(metadata, "traceparent")
			if assert.Contains(metadata, "baggage") {
				assert.Equal("member=123;property=456", metadata["baggage"])
			}
		})
		info.addRawPayload(func() any {
			return &traceStreamMessage{String: "abc", Int: 1}
		})

		ctx, span := tracer.Start(ctx, "TestService.TestMethod")
		defer span.End()

		property, err := baggage.NewKeyValuePropertyRaw("property", "456")
		require.NoError(err)
		member, err := baggage.NewMemberRaw("member", "123", property)
		require.NoError(err)
		bag, err := baggage.New(member)
		require.NoError(err)
		ctx = baggage.ContextWithBaggage(ctx, bag)

		_, err = interceptor.TraceBidirectionalStream(ctx, info, next)
		assert.NoError(err)

		assert.True(nextCalled, "missing expected next call")
		assert.False(info.hasMore(), "missing expected interceptor info calls")
		assert.False(payload.hasMore(), "missing expected payload calls")
	})

	t.Run("receive without prior setup", func(t *testing.T) {
		var (
			assert      = assert.New(t)
			ctx         = log.Context(context.Background(), log.WithFormat(log.FormatText))
			info        = newMockTraceStreamInfo(assert)
			interceptor = &TraceBidirectionalStreamClientInterceptor[*mockTraceStreamInfo, *mockTraceStreamingSendMessage, *mockTraceStreamingRecvMessage]{}
			nextCalled  = false
			next        = func(ctx context.Context, _ any) (any, error) {
				nextCalled = true
				return &traceStreamMessage{String: "abc", Int: 1}, nil
			}
		)
		info.addCallType(func() goa.InterceptorCallType {
			return goa.InterceptorStreamingRecv
		})
		info.addRawPayload(func() any {
			return nil
		})
		info.addService(func() string {
			return "TestService"
		})
		info.addMethod(func() string {
			return "TestMethod"
		})

		assert.PanicsWithError("clue interceptors trace stream receive method called without prior setup (service: TestService, method: TestMethod)", func() {
			_, _ = interceptor.TraceBidirectionalStream(ctx, info, next)
		})

		assert.True(nextCalled, "missing expected next call")
		assert.False(info.hasMore(), "missing expected interceptor info calls")
	})

	t.Run("receive", func(t *testing.T) {
		var (
			assert      = assert.New(t)
			ctx         = log.Context(context.Background(), log.WithFormat(log.FormatText))
			info        = newMockTraceStreamInfo(assert)
			payload     = newMockTraceStreamingRecvMessage(assert)
			interceptor = &TraceBidirectionalStreamClientInterceptor[*mockTraceStreamInfo, *mockTraceStreamingSendMessage, *mockTraceStreamingRecvMessage]{}
			nextCalled  = false
			next        = func(ctx context.Context, _ any) (any, error) {
				nextCalled = true
				return &traceStreamMessage{String: "abc", Int: 1}, nil
			}
		)
		info.addCallType(func() goa.InterceptorCallType {
			return goa.InterceptorStreamingRecv
		})
		info.addRawPayload(func() any {
			return nil
		})
		info.addClientStreamingResult(func(res any) *mockTraceStreamingRecvMessage {
			return payload
		})
		payload.addTraceMetadata(func() map[string]string {
			return map[string]string{
				"traceparent": "00-f5cb07fd7c9a0470ebca84a0107f9908-4fd197644c317fed-01",
				"baggage":     "member=123;property=456",
			}
		})

		ctx = SetupTraceStreamRecvContext(ctx)
		res, err := interceptor.TraceBidirectionalStream(ctx, info, next)
		assert.NoError(err)
		assert.Equal(&traceStreamMessage{String: "abc", Int: 1}, res)

		assert.NotPanics(func() {
			ctx = GetTraceStreamRecvContext(ctx)
		})
		span := trace.SpanFromContext(ctx)
		assert.Equal("f5cb07fd7c9a0470ebca84a0107f9908", span.SpanContext().TraceID().String())
		assert.Equal("4fd197644c317fed", span.SpanContext().SpanID().String())
		bag := baggage.FromContext(ctx)
		member := bag.Member("member")
		assert.Equal("123", member.Value())
		properties := member.Properties()
		if assert.Len(properties, 1) {
			property := properties[0]
			assert.Equal("property", property.Key())
			value, ok := property.Value()
			assert.True(ok)
			assert.Equal("456", value)
		}

		assert.True(nextCalled, "missing expected next call")
		assert.False(info.hasMore(), "missing expected interceptor info calls")
		assert.False(payload.hasMore(), "missing expected payload calls")
	})

	t.Run("receive with error", func(t *testing.T) {
		var (
			ctx         = log.Context(context.Background(), log.WithFormat(log.FormatText))
			info        = newMockTraceStreamInfo(assert.New(t))
			interceptor = &TraceBidirectionalStreamClientInterceptor[*mockTraceStreamInfo, *mockTraceStreamingSendMessage, *mockTraceStreamingRecvMessage]{}
			nextCalled  = false
			next        = func(ctx context.Context, _ any) (any, error) {
				nextCalled = true
				return nil, assert.AnError
			}
		)
		info.addCallType(func() goa.InterceptorCallType {
			return goa.InterceptorStreamingRecv
		})
		info.addRawPayload(func() any {
			return nil
		})

		ctx = SetupTraceStreamRecvContext(ctx)
		res, err := interceptor.TraceBidirectionalStream(ctx, info, next)
		assert.ErrorIs(t, err, assert.AnError)
		assert.Nil(t, res)

		assert.NotPanics(t, func() {
			ctx = GetTraceStreamRecvContext(ctx)
		})

		assert.True(t, nextCalled, "missing expected next call")
		assert.False(t, info.hasMore(), "missing expected interceptor info calls")
	})

	t.Run("unary", func(t *testing.T) {
		var (
			assert      = assert.New(t)
			ctx         = log.Context(context.Background(), log.WithFormat(log.FormatText))
			info        = newMockTraceStreamInfo(assert)
			interceptor = &TraceBidirectionalStreamClientInterceptor[*mockTraceStreamInfo, *mockTraceStreamingSendMessage, *mockTraceStreamingRecvMessage]{}
			nextCalled  = false
			next        = func(ctx context.Context, payload any) (any, error) {
				nextCalled = true
				assert.Equal(&traceStreamMessage{String: "abc", Int: 1}, payload)
				return &traceStreamMessage{String: "def", Int: 2}, nil
			}
		)
		info.addCallType(func() goa.InterceptorCallType {
			return goa.InterceptorUnary
		})
		info.addRawPayload(func() any {
			return &traceStreamMessage{String: "abc", Int: 1}
		})

		res, err := interceptor.TraceBidirectionalStream(ctx, info, next)
		assert.NoError(err)
		assert.Equal(&traceStreamMessage{String: "def", Int: 2}, res)

		assert.True(nextCalled, "missing expected next call")
		assert.False(info.hasMore(), "missing expected interceptor info calls")
	})
}

func TestTraceServerToClientStreamClientInterceptor(t *testing.T) {
	t.Run("receive without prior setup", func(t *testing.T) {
		var (
			assert      = assert.New(t)
			ctx         = log.Context(context.Background(), log.WithFormat(log.FormatText))
			info        = newMockTraceStreamInfo(assert)
			interceptor = &TraceServerToClientStreamClientInterceptor[*mockTraceStreamInfo, *mockTraceStreamingRecvMessage]{}
			nextCalled  = false
			next        = func(ctx context.Context, _ any) (any, error) {
				nextCalled = true
				return &traceStreamMessage{String: "abc", Int: 1}, nil
			}
		)
		info.addCallType(func() goa.InterceptorCallType {
			return goa.InterceptorStreamingRecv
		})
		info.addRawPayload(func() any {
			return nil
		})
		info.addService(func() string {
			return "TestService"
		})
		info.addMethod(func() string {
			return "TestMethod"
		})

		assert.PanicsWithError("clue interceptors trace stream receive method called without prior setup (service: TestService, method: TestMethod)", func() {
			_, _ = interceptor.TraceServerToClientStream(ctx, info, next)
		})

		assert.True(nextCalled, "missing expected next call")
		assert.False(info.hasMore(), "missing expected interceptor info calls")
	})

	t.Run("receive", func(t *testing.T) {
		var (
			assert      = assert.New(t)
			ctx         = log.Context(context.Background(), log.WithFormat(log.FormatText))
			info        = newMockTraceStreamInfo(assert)
			payload     = newMockTraceStreamingRecvMessage(assert)
			interceptor = &TraceServerToClientStreamClientInterceptor[*mockTraceStreamInfo, *mockTraceStreamingRecvMessage]{}
			nextCalled  = false
			next        = func(ctx context.Context, _ any) (any, error) {
				nextCalled = true
				return &traceStreamMessage{String: "abc", Int: 1}, nil
			}
		)
		info.addCallType(func() goa.InterceptorCallType {
			return goa.InterceptorStreamingRecv
		})
		info.addRawPayload(func() any {
			return nil
		})
		info.addClientStreamingResult(func(res any) *mockTraceStreamingRecvMessage {
			return payload
		})
		payload.addTraceMetadata(func() map[string]string {
			return map[string]string{
				"traceparent": "00-f5cb07fd7c9a0470ebca84a0107f9908-4fd197644c317fed-01",
				"baggage":     "member=123;property=456",
			}
		})

		ctx = SetupTraceStreamRecvContext(ctx)
		res, err := interceptor.TraceServerToClientStream(ctx, info, next)
		assert.NoError(err)
		assert.Equal(&traceStreamMessage{String: "abc", Int: 1}, res)

		assert.NotPanics(func() {
			ctx = GetTraceStreamRecvContext(ctx)
		})
		span := trace.SpanFromContext(ctx)
		assert.Equal("f5cb07fd7c9a0470ebca84a0107f9908", span.SpanContext().TraceID().String())
		assert.Equal("4fd197644c317fed", span.SpanContext().SpanID().String())
		bag := baggage.FromContext(ctx)
		member := bag.Member("member")
		assert.Equal("123", member.Value())
		properties := member.Properties()
		if assert.Len(properties, 1) {
			property := properties[0]
			assert.Equal("property", property.Key())
			value, ok := property.Value()
			assert.True(ok)
			assert.Equal("456", value)
		}

		assert.True(nextCalled, "missing expected next call")
		assert.False(info.hasMore(), "missing expected interceptor info calls")
		assert.False(payload.hasMore(), "missing expected payload calls")
	})

	t.Run("receive with error", func(t *testing.T) {
		var (
			ctx         = log.Context(context.Background(), log.WithFormat(log.FormatText))
			info        = newMockTraceStreamInfo(assert.New(t))
			interceptor = &TraceServerToClientStreamClientInterceptor[*mockTraceStreamInfo, *mockTraceStreamingRecvMessage]{}
			nextCalled  = false
			next        = func(ctx context.Context, _ any) (any, error) {
				nextCalled = true
				return nil, assert.AnError
			}
		)
		info.addCallType(func() goa.InterceptorCallType {
			return goa.InterceptorStreamingRecv
		})
		info.addRawPayload(func() any {
			return nil
		})

		ctx = SetupTraceStreamRecvContext(ctx)
		res, err := interceptor.TraceServerToClientStream(ctx, info, next)
		assert.ErrorIs(t, err, assert.AnError)
		assert.Nil(t, res)

		assert.NotPanics(t, func() {
			ctx = GetTraceStreamRecvContext(ctx)
		})

		assert.True(t, nextCalled, "missing expected next call")
		assert.False(t, info.hasMore(), "missing expected interceptor info calls")
	})

	t.Run("unary", func(t *testing.T) {
		var (
			assert      = assert.New(t)
			ctx         = log.Context(context.Background(), log.WithFormat(log.FormatText))
			info        = newMockTraceStreamInfo(assert)
			interceptor = &TraceServerToClientStreamClientInterceptor[*mockTraceStreamInfo, *mockTraceStreamingRecvMessage]{}
			nextCalled  = false
			next        = func(ctx context.Context, payload any) (any, error) {
				nextCalled = true
				assert.Equal(&traceStreamMessage{String: "abc", Int: 1}, payload)
				return &traceStreamMessage{String: "def", Int: 2}, nil
			}
		)
		info.addCallType(func() goa.InterceptorCallType {
			return goa.InterceptorUnary
		})
		info.addRawPayload(func() any {
			return &traceStreamMessage{String: "abc", Int: 1}
		})

		res, err := interceptor.TraceServerToClientStream(ctx, info, next)
		assert.NoError(err)
		assert.Equal(&traceStreamMessage{String: "def", Int: 2}, res)

		assert.True(nextCalled, "missing expected next call")
		assert.False(info.hasMore(), "missing expected interceptor info calls")
	})
}

func TestTraceClientToServerStreamClientInterceptor(t *testing.T) {
	t.Run("send", func(t *testing.T) {
		var (
			assert      = assert.New(t)
			require     = require.New(t)
			ctx         = log.Context(context.Background(), log.WithFormat(log.FormatText))
			info        = newMockTraceStreamInfo(assert)
			payload     = newMockTraceStreamingSendMessage(assert)
			interceptor = &TraceClientToServerStreamClientInterceptor[*mockTraceStreamInfo, *mockTraceStreamingSendMessage]{}
			nextCalled  = false
			next        = func(ctx context.Context, payload any) (any, error) {
				nextCalled = true
				assert.Equal(&traceStreamMessage{String: "abc", Int: 1}, payload)
				return nil, nil
			}
		)
		info.addCallType(func() goa.InterceptorCallType {
			return goa.InterceptorStreamingSend
		})
		info.addClientStreamingPayload(func() *mockTraceStreamingSendMessage {
			return payload
		})
		payload.addSetTraceMetadata(func(metadata map[string]string) {
			assert.Contains(metadata, "traceparent")
			if assert.Contains(metadata, "baggage") {
				assert.Equal("member=123;property=456", metadata["baggage"])
			}
		})
		info.addRawPayload(func() any {
			return &traceStreamMessage{String: "abc", Int: 1}
		})

		ctx, span := tracer.Start(ctx, "TestService.TestMethod")
		defer span.End()

		property, err := baggage.NewKeyValuePropertyRaw("property", "456")
		require.NoError(err)
		member, err := baggage.NewMemberRaw("member", "123", property)
		require.NoError(err)
		bag, err := baggage.New(member)
		require.NoError(err)
		ctx = baggage.ContextWithBaggage(ctx, bag)

		_, err = interceptor.TraceClientToServerStream(ctx, info, next)
		assert.NoError(err)

		assert.True(nextCalled, "missing expected next call")
		assert.False(info.hasMore(), "missing expected interceptor info calls")
		assert.False(payload.hasMore(), "missing expected payload calls")
	})

	t.Run("unary", func(t *testing.T) {
		var (
			assert      = assert.New(t)
			ctx         = log.Context(context.Background(), log.WithFormat(log.FormatText))
			info        = newMockTraceStreamInfo(assert)
			interceptor = &TraceClientToServerStreamClientInterceptor[*mockTraceStreamInfo, *mockTraceStreamingSendMessage]{}
			nextCalled  = false
			next        = func(ctx context.Context, payload any) (any, error) {
				nextCalled = true
				assert.Equal(&traceStreamMessage{String: "abc", Int: 1}, payload)
				return &traceStreamMessage{String: "def", Int: 2}, nil
			}
		)
		info.addCallType(func() goa.InterceptorCallType {
			return goa.InterceptorUnary
		})
		info.addRawPayload(func() any {
			return &traceStreamMessage{String: "abc", Int: 1}
		})

		res, err := interceptor.TraceClientToServerStream(ctx, info, next)
		assert.NoError(err)
		assert.Equal(&traceStreamMessage{String: "def", Int: 2}, res)

		assert.True(nextCalled, "missing expected next call")
		assert.False(info.hasMore(), "missing expected interceptor info calls")
	})
}

func TestWrapTraceBidirectionalStreamClientStream(t *testing.T) {
	var (
		assert        = assert.New(t)
		require       = require.New(t)
		ctx           = log.Context(context.Background(), log.WithFormat(log.FormatText))
		info          = newMockTraceStreamInfo(assert)
		result        = newMockTraceStreamingRecvMessage(assert)
		stream        = newMockTraceStream(assert)
		wrappedStream = WrapTraceBidirectionalStreamClientStream(stream)
	)
	stream.addSendWithContext(func(ctx context.Context, payload *traceStreamMessage) error {
		assert.Equal(&traceStreamMessage{String: "abc", Int: 1}, payload)
		return nil
	})
	stream.addRecvWithContext(func(ctx context.Context) (*traceStreamMessage, error) {
		res, err := TraceBidirectionalStreamClient(ctx, info, func(context.Context, any) (any, error) {
			return &traceStreamMessage{String: "def", Int: 2}, nil
		})
		require.IsType(&traceStreamMessage{}, res)
		return res.(*traceStreamMessage), err
	})
	info.addCallType(func() goa.InterceptorCallType {
		return goa.InterceptorStreamingRecv
	})
	info.addRawPayload(func() any {
		return nil
	})
	info.addClientStreamingResult(func(res any) *mockTraceStreamingRecvMessage {
		return result
	})
	result.addTraceMetadata(func() map[string]string {
		return map[string]string{
			"traceparent": "00-f5cb07fd7c9a0470ebca84a0107f9908-4fd197644c317fed-01",
			"baggage":     "member=123;property=456",
		}
	})
	stream.addClose(func() error {
		return nil
	})

	err := wrappedStream.Send(ctx, &traceStreamMessage{String: "abc", Int: 1})
	assert.NoError(err)

	var res *traceStreamMessage
	assert.NotPanics(func() {
		ctx, res, err = wrappedStream.RecvAndReturnContext(ctx)
	})
	assert.NoError(err)
	assert.Equal(&traceStreamMessage{String: "def", Int: 2}, res)

	span := trace.SpanFromContext(ctx)
	assert.Equal("f5cb07fd7c9a0470ebca84a0107f9908", span.SpanContext().TraceID().String())
	assert.Equal("4fd197644c317fed", span.SpanContext().SpanID().String())
	bag := baggage.FromContext(ctx)
	member := bag.Member("member")
	assert.Equal("123", member.Value())
	properties := member.Properties()
	if assert.Len(properties, 1) {
		property := properties[0]
		assert.Equal("property", property.Key())
		value, ok := property.Value()
		assert.True(ok)
		assert.Equal("456", value)
	}

	err = wrappedStream.Close()
	assert.NoError(err)

	assert.False(result.hasMore(), "missing expected payload calls")
	assert.False(info.hasMore(), "missing expected info calls")
	assert.False(stream.hasMore(), "missing expected stream calls")
}

func TestWrapTraceServerToClientStreamClientStream(t *testing.T) {
	var (
		assert        = assert.New(t)
		require       = require.New(t)
		ctx           = log.Context(context.Background(), log.WithFormat(log.FormatText))
		info          = newMockTraceStreamInfo(assert)
		result        = newMockTraceStreamingRecvMessage(assert)
		stream        = newMockTraceStream(assert)
		wrappedStream = WrapTraceServerToClientStreamClientStream(stream)
	)
	stream.addRecvWithContext(func(ctx context.Context) (*traceStreamMessage, error) {
		res, err := TraceServerToClientStreamClient(ctx, info, func(context.Context, any) (any, error) {
			return &traceStreamMessage{String: "def", Int: 2}, nil
		})
		require.IsType(&traceStreamMessage{}, res)
		return res.(*traceStreamMessage), err
	})
	info.addCallType(func() goa.InterceptorCallType {
		return goa.InterceptorStreamingRecv
	})
	info.addRawPayload(func() any {
		return nil
	})
	info.addClientStreamingResult(func(res any) *mockTraceStreamingRecvMessage {
		return result
	})
	result.addTraceMetadata(func() map[string]string {
		return map[string]string{
			"traceparent": "00-f5cb07fd7c9a0470ebca84a0107f9908-4fd197644c317fed-01",
			"baggage":     "member=123;property=456",
		}
	})

	var (
		err error
		res *traceStreamMessage
	)
	assert.NotPanics(func() {
		ctx, res, err = wrappedStream.RecvAndReturnContext(ctx)
	})
	assert.NoError(err)
	assert.Equal(&traceStreamMessage{String: "def", Int: 2}, res)

	span := trace.SpanFromContext(ctx)
	assert.Equal("f5cb07fd7c9a0470ebca84a0107f9908", span.SpanContext().TraceID().String())
	assert.Equal("4fd197644c317fed", span.SpanContext().SpanID().String())
	bag := baggage.FromContext(ctx)
	member := bag.Member("member")
	assert.Equal("123", member.Value())
	properties := member.Properties()
	if assert.Len(properties, 1) {
		property := properties[0]
		assert.Equal("property", property.Key())
		value, ok := property.Value()
		assert.True(ok)
		assert.Equal("456", value)
	}

	assert.False(info.hasMore(), "missing expected info calls")
	assert.False(result.hasMore(), "missing expected payload calls")
	assert.False(stream.hasMore(), "missing expected stream calls")
}

func TestWrapTraceClientToServerStreamWithResultClientStream(t *testing.T) {
	var (
		assert        = assert.New(t)
		require       = require.New(t)
		ctx           = log.Context(context.Background(), log.WithFormat(log.FormatText))
		info          = newMockTraceStreamInfo(assert)
		result        = newMockTraceStreamingRecvMessage(assert)
		stream        = newMockTraceStream(assert)
		wrappedStream = WrapTraceClientToServerStreamWithResultClientStream(stream)
	)
	stream.addSendWithContext(func(ctx context.Context, payload *traceStreamMessage) error {
		assert.Equal(&traceStreamMessage{String: "abc", Int: 1}, payload)
		return nil
	})
	stream.addCloseAndRecvWithContext(func(ctx context.Context) (*traceStreamMessage, error) {
		res, err := TraceBidirectionalStreamClient(ctx, info, func(context.Context, any) (any, error) {
			return &traceStreamMessage{String: "def", Int: 2}, nil
		})
		require.IsType(&traceStreamMessage{}, res)
		return res.(*traceStreamMessage), err
	})
	info.addCallType(func() goa.InterceptorCallType {
		return goa.InterceptorStreamingRecv
	})
	info.addRawPayload(func() any {
		return nil
	})
	info.addClientStreamingResult(func(res any) *mockTraceStreamingRecvMessage {
		assert.Equal(&traceStreamMessage{String: "def", Int: 2}, res)
		return result
	})
	result.addTraceMetadata(func() map[string]string {
		return map[string]string{
			"traceparent": "00-f5cb07fd7c9a0470ebca84a0107f9908-4fd197644c317fed-01",
			"baggage":     "member=123;property=456",
		}
	})

	err := wrappedStream.Send(ctx, &traceStreamMessage{String: "abc", Int: 1})
	assert.NoError(err)

	var res *traceStreamMessage
	assert.NotPanics(func() {
		ctx, res, err = wrappedStream.CloseAndRecvAndReturnContext(ctx)
	})
	assert.NoError(err)
	assert.Equal(&traceStreamMessage{String: "def", Int: 2}, res)

	span := trace.SpanFromContext(ctx)
	assert.Equal("f5cb07fd7c9a0470ebca84a0107f9908", span.SpanContext().TraceID().String())
	assert.Equal("4fd197644c317fed", span.SpanContext().SpanID().String())
	bag := baggage.FromContext(ctx)
	member := bag.Member("member")
	assert.Equal("123", member.Value())
	properties := member.Properties()
	if assert.Len(properties, 1) {
		property := properties[0]
		assert.Equal("property", property.Key())
		value, ok := property.Value()
		assert.True(ok)
		assert.Equal("456", value)
	}

	assert.False(info.hasMore(), "missing expected info calls")
	assert.False(result.hasMore(), "missing expected result calls")
	assert.False(stream.hasMore(), "missing expected stream calls")
}
