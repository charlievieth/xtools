package lspctx

import (
	"context"
	"log"
	"strconv"
	"sync/atomic"
)

var logValue = func() *atomic.Value {
	v := &atomic.Value{}
	v.Store(defaultLoggerFunc)
	return v
}()

func defaultLogger() *log.Logger {
	v := logValue.Load()
	ll, _ := v.(*log.Logger)
	return ll
}

func defaultLoggerFunc(_ context.Context) *log.Logger {
	return log.New(log.Writer(), log.Prefix(), log.Flags())
}

func SetDefaultLoggerFunc(fn func(ctx context.Context) *log.Logger) (previous *log.Logger) {
	return nil
}

type (
	loggerKey    struct{}
	wdKey        struct{}
	requestIDKey struct{}
)

func WithLogger(ctx context.Context, ll *log.Logger) context.Context {
	return context.WithValue(ctx, loggerKey{}, ll)
}

func WithLoggerDefault(ctx context.Context, ll *log.Logger) context.Context {
	if ll == nil {
		ll = defaultLoggerFunc(ctx)
	}
	return context.WithValue(ctx, loggerKey{}, ll)
}

func Logger(ctx context.Context) *log.Logger {
	ll, _ := ctx.Value(loggerKey{}).(*log.Logger)
	return ll
}

func WithWorkingDirectory(ctx context.Context, wd string) context.Context {
	return context.WithValue(ctx, wdKey{}, wd)
}

func WorkingDirectory(ctx context.Context) string {
	s, _ := ctx.Value(wdKey{}).(string)
	return s
}

var requestID uint64

func WithRequestID(ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, requestIDKey{}, id)
}

func WithAutoRequestID(ctx context.Context, id uint64) context.Context {
	return WithRequestID(ctx, strconv.FormatUint(atomic.AddUint64(&requestID, 1), 10))
}

func RequestID(ctx context.Context) string {
	id, _ := ctx.Value(requestIDKey{}).(string)
	return id
}
