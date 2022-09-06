package log

var (
	TraceIDKey      = "trace-id"
	SpanIDKey       = "span-id"
	RequestIDKey    = "request-id"
	MessageKey      = "msg"
	ErrorMessageKey = "err"
	TimestampKey    = "time"
	SeverityKey     = "level"
	HTTPMethodKey   = "http.method"
	HTTPURLKey      = "http.url"
	HTTPStatusKey   = "http.status"
	HTTPDurationKey = "http.time_ms"
	GRPCServiceKey  = "grpc.service"
	GRPCMethodKey   = "grpc.method"
	GRPCCodeKey     = "grpc.code"
	GRPCStatusKey   = "grpc.status"
	GRPCDurationKey = "grpc.time_ms"
)
