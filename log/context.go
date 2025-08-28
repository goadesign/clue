package log

import "context"

const (
	ctxLogger ctxKey = iota + 1
)

// Context initializes a context for logging.
func Context(ctx context.Context, opts ...LogOption) context.Context {
	l, ok := ctx.Value(ctxLogger).(*logger)
	if !ok {
		l = &logger{options: defaultOptions()}
	}
	l.lock.Lock()
	defer l.lock.Unlock()
	for _, opt := range opts {
		opt(l.options)
	}
	// Copy OTEL logger from options to logger struct
	l.otellog = l.options.otellog
	// Reset the OTEL logging guard for new loggers
	l.otelLogging = false
	if l.options.disableBuffering != nil && l.options.disableBuffering(ctx) {
		l.flush()
	}
	return context.WithValue(ctx, ctxLogger, l)
}

// WithContext will inject the second context in the given one.
//
// It is useful when building middleware handlers such as log.HTTP
func WithContext(parentCtx, logCtx context.Context) context.Context {
	logger, ok := logCtx.Value(ctxLogger).(*logger)
	if !ok {
		return parentCtx
	}
	return context.WithValue(parentCtx, ctxLogger, logger)
}

// MustContainLogger will panic if the given context is missing the logger.
//
// It can be used during server initialisation when you have a function or
// middleware that you want to ensure receives a context with a logger.
func MustContainLogger(logCtx context.Context) {
	_, ok := logCtx.Value(ctxLogger).(*logger)
	if !ok {
		panic("provided a context without a logger. Use log.Context")
	}
}

// DebugEnabled returns true if the given context has debug logging enabled.
func DebugEnabled(ctx context.Context) bool {
	v := ctx.Value(ctxLogger)
	if v == nil {
		return false
	}
	l := v.(*logger)
	l.lock.Lock()
	defer l.lock.Unlock()
	return l.options.debug
}
