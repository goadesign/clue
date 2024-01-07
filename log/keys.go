package log

var (
	TraceIDKey      = "trace_id"
	SpanIDKey       = "span_id"
	MessageKey      = "msg"
	ErrorMessageKey = "err"
	TimestampKey    = "time"
	SeverityKey     = "level"
	HTTPMethodKey   = "http.method"
	HTTPURLKey      = "http.url"
	HTTPFromKey     = "http.remote_addr"
	HTTPStatusKey   = "http.status"
	HTTPDurationKey = "http.time_ms"
	HTTPBytesKey    = "http.bytes"
	HTTPBodyKey     = "http.body"
	GRPCServiceKey  = "grpc.service"
	GRPCMethodKey   = "grpc.method"
	GRPCCodeKey     = "grpc.code"
	GRPCStatusKey   = "grpc.status"
	GRPCDurationKey = "grpc.time_ms"
	GoaServiceKey   = "goa.service"
	GoaMethodKey    = "goa.method"
)
