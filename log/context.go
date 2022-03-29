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
	if l.options.disableBuffering != nil && l.options.disableBuffering(ctx) {
		l.flush()
	}
	return context.WithValue(ctx, ctxLogger, l)
}

func WithContext(parentCtx, logCtx context.Context) context.Context {
	logger, ok := logCtx.Value(ctxLogger).(*logger)
	if !ok {
		return parentCtx
	}
	return context.WithValue(parentCtx, ctxLogger, logger)
}

func MustContainLogger(logCtx context.Context) {
	_, ok := logCtx.Value(ctxLogger).(*logger)
	if !ok {
		panic("log.HTTP called without log.Context")
	}
}
