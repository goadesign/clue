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

func TestServerBidirectionalStreamInterceptor(t *testing.T) {
	t.Run("send", func(t *testing.T) {
		var (
			assert      = assert.New(t)
			require     = require.New(t)
			ctx         = log.Context(context.Background(), log.WithFormat(log.FormatText))
			info        = newMockTraceStreamInfo(assert)
			result      = newMockTraceStreamingSendMessage(assert)
			interceptor = &ServerBidirectionalStreamInterceptor[*mockTraceStreamInfo, *mockTraceStreamingRecvMessage, *mockTraceStreamingSendMessage]{}
			nextCalled  = false
			next        = func(ctx context.Context, result any) (any, error) {
				nextCalled = true
				assert.Equal(&traceStreamMessage{String: "abc", Int: 1}, result)
				return nil, nil
			}
		)
		info.addCallType(func() goa.InterceptorCallType {
			return goa.InterceptorStreamingSend
		})
		info.addServerStreamingResult(func() *mockTraceStreamingSendMessage {
			return result
		})
		result.addSetTraceMetadata(func(metadata map[string]string) {
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
		assert.False(result.hasMore(), "missing expected result calls")
	})

	t.Run("receive without prior setup", func(t *testing.T) {
		var (
			assert      = assert.New(t)
			ctx         = log.Context(context.Background(), log.WithFormat(log.FormatText))
			info        = newMockTraceStreamInfo(assert)
			interceptor = &ServerBidirectionalStreamInterceptor[*mockTraceStreamInfo, *mockTraceStreamingRecvMessage, *mockTraceStreamingSendMessage]{}
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
		info.addServerStreamingPayload(func(pay any) *mockTraceStreamingRecvMessage {
			return newMockTraceStreamingRecvMessage(assert)
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
			interceptor = &ServerBidirectionalStreamInterceptor[*mockTraceStreamInfo, *mockTraceStreamingRecvMessage, *mockTraceStreamingSendMessage]{}
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
		info.addServerStreamingPayload(func(pay any) *mockTraceStreamingRecvMessage {
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

	t.Run("unary", func(t *testing.T) {
		var (
			assert      = assert.New(t)
			ctx         = log.Context(context.Background(), log.WithFormat(log.FormatText))
			info        = newMockTraceStreamInfo(assert)
			interceptor = &ServerBidirectionalStreamInterceptor[*mockTraceStreamInfo, *mockTraceStreamingRecvMessage, *mockTraceStreamingSendMessage]{}
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

func TestServerDownStreamInterceptor(t *testing.T) {
	t.Run("send", func(t *testing.T) {
		var (
			assert      = assert.New(t)
			require     = require.New(t)
			ctx         = log.Context(context.Background(), log.WithFormat(log.FormatText))
			info        = newMockTraceStreamInfo(assert)
			result      = newMockTraceStreamingSendMessage(assert)
			interceptor = &ServerDownStreamInterceptor[*mockTraceStreamInfo, *mockTraceStreamingSendMessage]{}
			nextCalled  = false
			next        = func(ctx context.Context, result any) (any, error) {
				nextCalled = true
				assert.Equal(&traceStreamMessage{String: "abc", Int: 1}, result)
				return nil, nil
			}
		)
		info.addCallType(func() goa.InterceptorCallType {
			return goa.InterceptorStreamingSend
		})
		info.addServerStreamingResult(func() *mockTraceStreamingSendMessage {
			return result
		})
		result.addSetTraceMetadata(func(metadata map[string]string) {
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

		_, err = interceptor.TraceDownStream(ctx, info, next)
		assert.NoError(err)

		assert.True(nextCalled, "missing expected next call")
		assert.False(info.hasMore(), "missing expected interceptor info calls")
		assert.False(result.hasMore(), "missing expected result calls")
	})

	t.Run("unary", func(t *testing.T) {
		var (
			assert      = assert.New(t)
			ctx         = log.Context(context.Background(), log.WithFormat(log.FormatText))
			info        = newMockTraceStreamInfo(assert)
			interceptor = &ServerDownStreamInterceptor[*mockTraceStreamInfo, *mockTraceStreamingSendMessage]{}
			nextCalled  = false
			next        = func(ctx context.Context, result any) (any, error) {
				nextCalled = true
				assert.Equal(&traceStreamMessage{String: "abc", Int: 1}, result)
				return &traceStreamMessage{String: "def", Int: 2}, nil
			}
		)
		info.addCallType(func() goa.InterceptorCallType {
			return goa.InterceptorUnary
		})
		info.addRawPayload(func() any {
			return &traceStreamMessage{String: "abc", Int: 1}
		})

		res, err := interceptor.TraceDownStream(ctx, info, next)
		assert.NoError(err)
		assert.Equal(&traceStreamMessage{String: "def", Int: 2}, res)

		assert.True(nextCalled, "missing expected next call")
		assert.False(info.hasMore(), "missing expected interceptor info calls")
	})
}

func TestServerUpStreamInterceptor(t *testing.T) {
	t.Run("receive without prior setup", func(t *testing.T) {
		var (
			assert      = assert.New(t)
			ctx         = log.Context(context.Background(), log.WithFormat(log.FormatText))
			info        = newMockTraceStreamInfo(assert)
			interceptor = &ServerUpStreamInterceptor[*mockTraceStreamInfo, *mockTraceStreamingRecvMessage]{}
			nextCalled  = false
			next        = func(ctx context.Context, payload any) (any, error) {
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
		info.addServerStreamingPayload(func(pay any) *mockTraceStreamingRecvMessage {
			return newMockTraceStreamingRecvMessage(assert)
		})
		info.addService(func() string {
			return "TestService"
		})
		info.addMethod(func() string {
			return "TestMethod"
		})

		assert.PanicsWithError("clue interceptors trace stream receive method called without prior setup (service: TestService, method: TestMethod)", func() {
			_, _ = interceptor.TraceUpStream(ctx, info, next)
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
			interceptor = &ServerUpStreamInterceptor[*mockTraceStreamInfo, *mockTraceStreamingRecvMessage]{}
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
		info.addServerStreamingPayload(func(pay any) *mockTraceStreamingRecvMessage {
			return payload
		})
		payload.addTraceMetadata(func() map[string]string {
			return map[string]string{
				"traceparent": "00-f5cb07fd7c9a0470ebca84a0107f9908-4fd197644c317fed-01",
				"baggage":     "member=123;property=456",
			}
		})

		ctx = SetupTraceStreamRecvContext(ctx)
		res, err := interceptor.TraceUpStream(ctx, info, next)
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

	t.Run("unary", func(t *testing.T) {
		var (
			assert      = assert.New(t)
			ctx         = log.Context(context.Background(), log.WithFormat(log.FormatText))
			info        = newMockTraceStreamInfo(assert)
			interceptor = &ServerUpStreamInterceptor[*mockTraceStreamInfo, *mockTraceStreamingRecvMessage]{}
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

		res, err := interceptor.TraceUpStream(ctx, info, next)
		assert.NoError(err)
		assert.Equal(&traceStreamMessage{String: "def", Int: 2}, res)

		assert.True(nextCalled, "missing expected next call")
		assert.False(info.hasMore(), "missing expected interceptor info calls")
	})
}

func TestWrapTraceStreamServerBidirectionalStream(t *testing.T) {
	var (
		assert        = assert.New(t)
		require       = require.New(t)
		ctx           = log.Context(context.Background(), log.WithFormat(log.FormatText))
		info          = newMockTraceStreamInfo(assert)
		payload       = newMockTraceStreamingRecvMessage(assert)
		stream        = newMockTraceStream(assert)
		wrappedStream = WrapTraceStreamServerBidirectionalStream(stream)
	)
	stream.addSendWithContext(func(ctx context.Context, result *traceStreamMessage) error {
		assert.Equal(&traceStreamMessage{String: "abc", Int: 1}, result)
		return nil
	})
	stream.addRecvWithContext(func(ctx context.Context) (*traceStreamMessage, error) {
		res, err := ServerTraceBidirectionalStream(ctx, info, func(context.Context, any) (any, error) {
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
	info.addServerStreamingPayload(func(pay any) *mockTraceStreamingRecvMessage {
		return payload
	})
	payload.addTraceMetadata(func() map[string]string {
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

	var pay *traceStreamMessage
	assert.NotPanics(func() {
		ctx, pay, err = wrappedStream.RecvAndReturnContext(ctx)
	})
	assert.NoError(err)
	assert.Equal(&traceStreamMessage{String: "def", Int: 2}, pay)

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

	assert.False(info.hasMore(), "missing expected info calls")
	assert.False(payload.hasMore(), "missing expected payload calls")
	assert.False(stream.hasMore(), "missing expected stream calls")
}

func TestWrapTraceStreamServerUpStream(t *testing.T) {
	var (
		assert        = assert.New(t)
		require       = require.New(t)
		ctx           = log.Context(context.Background(), log.WithFormat(log.FormatText))
		info          = newMockTraceStreamInfo(assert)
		payload       = newMockTraceStreamingRecvMessage(assert)
		stream        = newMockTraceStream(assert)
		wrappedStream = WrapTraceStreamServerUpStream(stream)
	)
	stream.addRecvWithContext(func(ctx context.Context) (*traceStreamMessage, error) {
		res, err := ServerTraceUpStream(ctx, info, func(context.Context, any) (any, error) {
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
	info.addServerStreamingPayload(func(pay any) *mockTraceStreamingRecvMessage {
		return payload
	})
	payload.addTraceMetadata(func() map[string]string {
		return map[string]string{
			"traceparent": "00-f5cb07fd7c9a0470ebca84a0107f9908-4fd197644c317fed-01",
			"baggage":     "member=123;property=456",
		}
	})
	stream.addClose(func() error {
		return nil
	})

	var (
		err error
		pay *traceStreamMessage
	)
	assert.NotPanics(func() {
		ctx, pay, err = wrappedStream.RecvAndReturnContext(ctx)
	})
	assert.NoError(err)
	assert.Equal(&traceStreamMessage{String: "def", Int: 2}, pay)

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

	assert.False(info.hasMore(), "missing expected info calls")
	assert.False(payload.hasMore(), "missing expected payload calls")
	assert.False(stream.hasMore(), "missing expected stream calls")
}

func TestWrapTraceStreamServerUpStreamWithResult(t *testing.T) {
	var (
		assert        = assert.New(t)
		require       = require.New(t)
		ctx           = log.Context(context.Background(), log.WithFormat(log.FormatText))
		info          = newMockTraceStreamInfo(assert)
		payload       = newMockTraceStreamingRecvMessage(assert)
		stream        = newMockTraceStream(assert)
		wrappedStream = WrapTraceStreamServerUpStreamWithResult(stream)
	)
	stream.addCloseAndSendWithContext(func(ctx context.Context, result *traceStreamMessage) error {
		assert.Equal(&traceStreamMessage{String: "abc", Int: 1}, result)
		return nil
	})
	stream.addRecvWithContext(func(ctx context.Context) (*traceStreamMessage, error) {
		res, err := ServerTraceUpStream(ctx, info, func(context.Context, any) (any, error) {
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
	info.addServerStreamingPayload(func(pay any) *mockTraceStreamingRecvMessage {
		return payload
	})
	payload.addTraceMetadata(func() map[string]string {
		return map[string]string{
			"traceparent": "00-f5cb07fd7c9a0470ebca84a0107f9908-4fd197644c317fed-01",
			"baggage":     "member=123;property=456",
		}
	})

	err := wrappedStream.CloseAndSend(ctx, &traceStreamMessage{String: "abc", Int: 1})
	assert.NoError(err)

	var pay *traceStreamMessage
	assert.NotPanics(func() {
		ctx, pay, err = wrappedStream.RecvAndReturnContext(ctx)
	})
	assert.NoError(err)
	assert.Equal(&traceStreamMessage{String: "def", Int: 2}, pay)

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
	assert.False(payload.hasMore(), "missing expected payload calls")
	assert.False(stream.hasMore(), "missing expected stream calls")
}
