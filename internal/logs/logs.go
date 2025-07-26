package logs

import (
	"context"

	"go.uber.org/zap"
)

type logKeyType string

const logKey logKeyType = "log"

func GetFromContext(ctx context.Context) *zap.Logger {
	return ctx.Value(logKey).(*zap.Logger)
}

func WithLog(ctx context.Context, log *zap.Logger, fields ...zap.Field) context.Context {
	log = log.With(fields...)
	return context.WithValue(ctx, logKey, log)
}
