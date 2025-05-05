package interceptors

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/stdout/stdoutmetric"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/trace"

	"goa.design/clue/clue"
	"goa.design/clue/log"
	"goa.design/clue/mock"
	goa "goa.design/goa/v3/pkg"
)

type (
	traceStreamMessage struct {
		String string
		Int    int
	}

	mockTraceStreamInfo struct {
		m      *mock.Mock
		assert *assert.Assertions
	}

	mockTraceStreamInfoService                func() string
	mockTraceStreamInfoMethod                 func() string
	mockTraceStreamInfoCallType               func() goa.InterceptorCallType
	mockTraceStreamInfoRawPayload             func() any
	mockTraceStreamInfoClientStreamingPayload func() *mockTraceStreamingSendMessage
	mockTraceStreamInfoClientStreamingResult  func(res any) *mockTraceStreamingRecvMessage
	mockTraceStreamInfoServerStreamingPayload func(pay any) *mockTraceStreamingRecvMessage
	mockTraceStreamInfoServerStreamingResult  func() *mockTraceStreamingSendMessage

	mockTraceStreamingSendMessage struct {
		m      *mock.Mock
		assert *assert.Assertions
	}

	mockTraceStreamingSendMessageSetTraceMetadata func(map[string]string)

	mockTraceStreamingRecvMessage struct {
		m      *mock.Mock
		assert *assert.Assertions
	}

	mockTraceStreamingRecvMessageTraceMetadata func() map[string]string

	mockTraceStream struct {
		m      *mock.Mock
		assert *assert.Assertions
	}

	mockTraceStreamSendWithContext         func(ctx context.Context, payload *traceStreamMessage) error
	mockTraceStreamRecvWithContext         func(ctx context.Context) (*traceStreamMessage, error)
	mockTraceStreamClose                   func() error
	mockTraceStreamCloseAndSendWithContext func(ctx context.Context, payload *traceStreamMessage) error
	mockTraceStreamCloseAndRecvWithContext func(ctx context.Context) (*traceStreamMessage, error)
)

var tracer trace.Tracer

func init() {
	ctx := log.Context(context.Background(), log.WithFormat(log.FormatTerminal))
	metricExporter, err := stdoutmetric.New()
	if err != nil {
		panic(err)
	}
	traceExporter, err := stdouttrace.New()
	if err != nil {
		panic(err)
	}
	cfg, err := clue.NewConfig(ctx, "test", "0.0.1", metricExporter, traceExporter)
	if err != nil {
		panic(err)
	}
	clue.ConfigureOpenTelemetry(ctx, cfg)
	tracer = otel.Tracer("test")
}

func TestSetupTraceStreamRecvContext(t *testing.T) {
	ctx := context.Background()
	ctx = SetupTraceStreamRecvContext(ctx)
	assert.Equal(t, &traceStreamRecvContext{}, ctx.Value(traceStreamRecvContextKey))
}

func TestGetTraceStreamRecvContext(t *testing.T) {
	type ctxKey string
	ctx := context.Background()

	assert.PanicsWithError(t, "clue interceptors get trace stream receive context method called without prior setup", func() {
		GetTraceStreamRecvContext(ctx)
	})

	ctx = context.WithValue(ctx, traceStreamRecvContextKey, &traceStreamRecvContext{})

	assert.PanicsWithError(t, "clue interceptors get trace stream receive context method called without prior interceptor receive method call", func() {
		GetTraceStreamRecvContext(ctx)
	})

	expectedCtx := context.WithValue(ctx, ctxKey("trace_metadata"), "test")
	ctx = context.WithValue(ctx, traceStreamRecvContextKey, &traceStreamRecvContext{ctx: expectedCtx})
	assert.Equal(t, expectedCtx, GetTraceStreamRecvContext(ctx))
}

func TestInsertTraceStreamRecvContext(t *testing.T) {
	type ctxKey string
	var (
		ctx         = context.Background()
		expectedCtx = context.WithValue(ctx, ctxKey("trace_metadata"), "test")
	)

	assert.PanicsWithError(t, "clue interceptors insert trace stream receive context method called without prior setup", func() {
		InsertTraceStreamRecvContext(ctx, expectedCtx)
	})

	ctx = context.WithValue(ctx, traceStreamRecvContextKey, &traceStreamRecvContext{})
	InsertTraceStreamRecvContext(ctx, expectedCtx)
	assert.Equal(t, expectedCtx, ctx.Value(traceStreamRecvContextKey).(*traceStreamRecvContext).ctx)
}

func newMockTraceStreamInfo(assert *assert.Assertions) *mockTraceStreamInfo {
	var (
		m                                                                                                    = &mockTraceStreamInfo{mock.New(), assert}
		_ TraceBidirectionalStreamClientInfo[*mockTraceStreamingSendMessage, *mockTraceStreamingRecvMessage] = m
		_ ClientTraceStreamServerToClientInfo[*mockTraceStreamingRecvMessage]                                = m
		_ ClientTraceStreamClientToServerInfo[*mockTraceStreamingSendMessage]                                = m
		_ TraceBidirectionalStreamServerInfo[*mockTraceStreamingRecvMessage, *mockTraceStreamingSendMessage] = m
		_ ServerTraceStreamServerToClientInfo[*mockTraceStreamingSendMessage]                                = m
		_ ServerTraceStreamClientToServerInfo[*mockTraceStreamingRecvMessage]                                = m
	)
	return m
}

func (m *mockTraceStreamInfo) addService(service mockTraceStreamInfoService) {
	m.m.Add("Service", service)
}

func (m *mockTraceStreamInfo) Service() string {
	if f := m.m.Next("Service"); f != nil {
		return f.(mockTraceStreamInfoService)()
	}
	m.assert.Fail("unexpected Service call")
	return ""
}

func (m *mockTraceStreamInfo) addMethod(method mockTraceStreamInfoMethod) {
	m.m.Add("Method", method)
}

func (m *mockTraceStreamInfo) Method() string {
	if f := m.m.Next("Method"); f != nil {
		return f.(mockTraceStreamInfoMethod)()
	}
	m.assert.Fail("unexpected Method call")
	return ""
}

func (m *mockTraceStreamInfo) addCallType(callType mockTraceStreamInfoCallType) {
	m.m.Add("CallType", callType)
}

func (m *mockTraceStreamInfo) CallType() goa.InterceptorCallType {
	if f := m.m.Next("CallType"); f != nil {
		return f.(mockTraceStreamInfoCallType)()
	}
	m.assert.Fail("unexpected CallType call")
	return 0
}

func (m *mockTraceStreamInfo) addRawPayload(rawPayload mockTraceStreamInfoRawPayload) {
	m.m.Add("RawPayload", rawPayload)
}

func (m *mockTraceStreamInfo) RawPayload() any {
	if f := m.m.Next("RawPayload"); f != nil {
		return f.(mockTraceStreamInfoRawPayload)()
	}
	m.assert.Fail("unexpected RawPayload call")
	return nil
}

func (m *mockTraceStreamInfo) addClientStreamingPayload(clientStreamingPayload mockTraceStreamInfoClientStreamingPayload) {
	m.m.Add("ClientStreamingPayload", clientStreamingPayload)
}

func (m *mockTraceStreamInfo) ClientStreamingPayload() *mockTraceStreamingSendMessage {
	if f := m.m.Next("ClientStreamingPayload"); f != nil {
		return f.(mockTraceStreamInfoClientStreamingPayload)()
	}
	m.assert.Fail("unexpected ClientStreamingPayload call")
	return &mockTraceStreamingSendMessage{mock.New(), m.assert}
}

func (m *mockTraceStreamInfo) addClientStreamingResult(clientStreamingResult mockTraceStreamInfoClientStreamingResult) {
	m.m.Add("ClientStreamingResult", clientStreamingResult)
}

func (m *mockTraceStreamInfo) ClientStreamingResult(res any) *mockTraceStreamingRecvMessage {
	if f := m.m.Next("ClientStreamingResult"); f != nil {
		return f.(mockTraceStreamInfoClientStreamingResult)(res)
	}
	m.assert.Fail("unexpected ClientStreamingResult call")
	return &mockTraceStreamingRecvMessage{mock.New(), m.assert}
}

func (m *mockTraceStreamInfo) addServerStreamingPayload(serverStreamingPayload mockTraceStreamInfoServerStreamingPayload) {
	m.m.Add("ServerStreamingPayload", serverStreamingPayload)
}

func (m *mockTraceStreamInfo) ServerStreamingPayload(pay any) *mockTraceStreamingRecvMessage {
	if f := m.m.Next("ServerStreamingPayload"); f != nil {
		return f.(mockTraceStreamInfoServerStreamingPayload)(pay)
	}
	m.assert.Fail("unexpected ServerStreamingPayload call")
	return &mockTraceStreamingRecvMessage{mock.New(), m.assert}
}

func (m *mockTraceStreamInfo) addServerStreamingResult(serverStreamingResult mockTraceStreamInfoServerStreamingResult) {
	m.m.Add("ServerStreamingResult", serverStreamingResult)
}

func (m *mockTraceStreamInfo) ServerStreamingResult() *mockTraceStreamingSendMessage {
	if f := m.m.Next("ServerStreamingResult"); f != nil {
		return f.(mockTraceStreamInfoServerStreamingResult)()
	}
	m.assert.Fail("unexpected ServerStreamingResult call")
	return &mockTraceStreamingSendMessage{mock.New(), m.assert}
}

func (m *mockTraceStreamInfo) hasMore() bool {
	return m.m.HasMore()
}

func newMockTraceStreamingSendMessage(assert *assert.Assertions) *mockTraceStreamingSendMessage {
	var (
		m                                 = &mockTraceStreamingSendMessage{mock.New(), assert}
		_ TraceStreamStreamingSendMessage = m
	)
	return m
}

func (m *mockTraceStreamingSendMessage) addSetTraceMetadata(setTraceMetadata mockTraceStreamingSendMessageSetTraceMetadata) {
	m.m.Add("SetTraceMetadata", setTraceMetadata)
}

func (m *mockTraceStreamingSendMessage) SetTraceMetadata(metadata map[string]string) {
	if f := m.m.Next("SetTraceMetadata"); f != nil {
		f.(mockTraceStreamingSendMessageSetTraceMetadata)(metadata)
		return
	}
	m.assert.Fail("unexpected SetTraceMetadata call")
}

func (m *mockTraceStreamingSendMessage) hasMore() bool {
	return m.m.HasMore()
}

func newMockTraceStreamingRecvMessage(assert *assert.Assertions) *mockTraceStreamingRecvMessage {
	var (
		m                                 = &mockTraceStreamingRecvMessage{mock.New(), assert}
		_ TraceStreamStreamingRecvMessage = m
	)
	return m
}

func (m *mockTraceStreamingRecvMessage) TraceMetadata() map[string]string {
	if f := m.m.Next("TraceMetadata"); f != nil {
		return f.(mockTraceStreamingRecvMessageTraceMetadata)()
	}
	m.assert.Fail("unexpected TraceMetadata call")
	return nil
}

func (m *mockTraceStreamingRecvMessage) addTraceMetadata(traceMetadata mockTraceStreamingRecvMessageTraceMetadata) {
	m.m.Add("TraceMetadata", traceMetadata)
}

func (m *mockTraceStreamingRecvMessage) hasMore() bool {
	return m.m.HasMore()
}

func newMockTraceStream(assert *assert.Assertions) *mockTraceStream {
	var (
		m                                                                                = &mockTraceStream{mock.New(), assert}
		_ ClientBidirectionalStream[*traceStreamMessage, *traceStreamMessage]            = m
		_ ClientServerToClientStream[*traceStreamMessage]                                = m
		_ ClientClientToServerStreamWithResult[*traceStreamMessage, *traceStreamMessage] = m
		_ ServerBidirectionalStream[*traceStreamMessage, *traceStreamMessage]            = m
		_ ServerClientToServerStream[*traceStreamMessage]                                = m
		_ ServerClientToServerStreamWithResult[*traceStreamMessage, *traceStreamMessage] = m
	)
	return m
}

func (m *mockTraceStream) addSendWithContext(sendWithContext mockTraceStreamSendWithContext) {
	m.m.Add("SendWithContext", sendWithContext)
}

func (m *mockTraceStream) SendWithContext(ctx context.Context, payload *traceStreamMessage) error {
	if f := m.m.Next("SendWithContext"); f != nil {
		return f.(mockTraceStreamSendWithContext)(ctx, payload)
	}
	m.assert.Fail("unexpected SendWithContext call")
	return nil
}

func (m *mockTraceStream) addRecvWithContext(recvWithContext mockTraceStreamRecvWithContext) {
	m.m.Add("RecvWithContext", recvWithContext)
}

func (m *mockTraceStream) RecvWithContext(ctx context.Context) (*traceStreamMessage, error) {
	if f := m.m.Next("RecvWithContext"); f != nil {
		return f.(mockTraceStreamRecvWithContext)(ctx)
	}
	m.assert.Fail("unexpected RecvWithContext call")
	return nil, nil
}

func (m *mockTraceStream) addClose(close mockTraceStreamClose) {
	m.m.Add("Close", close)
}

func (m *mockTraceStream) Close() error {
	if f := m.m.Next("Close"); f != nil {
		return f.(mockTraceStreamClose)()
	}
	m.assert.Fail("unexpected Close call")
	return nil
}

func (m *mockTraceStream) addCloseAndSendWithContext(closeAndSendWithContext mockTraceStreamCloseAndSendWithContext) {
	m.m.Add("CloseAndSendWithContext", closeAndSendWithContext)
}

func (m *mockTraceStream) CloseAndSendWithContext(ctx context.Context, payload *traceStreamMessage) error {
	if f := m.m.Next("CloseAndSendWithContext"); f != nil {
		return f.(mockTraceStreamCloseAndSendWithContext)(ctx, payload)
	}
	m.assert.Fail("unexpected CloseAndSendWithContext call")
	return nil
}

func (m *mockTraceStream) addCloseAndRecvWithContext(closeAndRecvWithContext mockTraceStreamCloseAndRecvWithContext) {
	m.m.Add("CloseAndRecvWithContext", closeAndRecvWithContext)
}

func (m *mockTraceStream) CloseAndRecvWithContext(ctx context.Context) (*traceStreamMessage, error) {
	if f := m.m.Next("CloseAndRecvWithContext"); f != nil {
		return f.(mockTraceStreamCloseAndRecvWithContext)(ctx)
	}
	m.assert.Fail("unexpected CloseAndRecvWithContext call")
	return nil, nil
}

func (m *mockTraceStream) hasMore() bool {
	return m.m.HasMore()
}
