package logger

import (
	"context"
	"log/slog"
	"os"
	"runtime"
	"time"

	"github.com/google/uuid"
)

type ctxKey string

const (
	RequestIDKey ctxKey = "request_id"
)

type Logger struct {
	*slog.Logger
}

func New(env string, level slog.Level) *Logger {
	var handler slog.Handler

	opts := &slog.HandlerOptions{
		Level: level,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			if env == "development" && a.Key == "time" {
				return slog.Attr{}
			}

			return a
		},
	}

	if env == "production" {
		handler = slog.NewJSONHandler(os.Stdout, opts)
	} else {
		handler = slog.NewTextHandler(os.Stdout, opts)
	}

	logger := slog.New(handler)

	return &Logger{logger}
}

func (l *Logger) WithRequestID(ctx context.Context) (context.Context, *Logger) {
	requestID := uuid.New().String()

	ctx = context.WithValue(ctx, RequestIDKey, requestID)

	return ctx, &Logger{l.With(slog.String("request_id", requestID))}
}

func GetRequestID(ctx context.Context) string {
	if id, ok := ctx.Value(RequestIDKey).(string); ok {
		return id
	}

	return ""
}

func (l *Logger) Error(msg string, err error, attrs ...slog.Attr) {
	if err != nil {
		attrs = append(attrs, slog.String("error", err.Error()))

		if l.Enabled(context.Background(), slog.LevelDebug) {
			stack := make([]byte, 4096)
			stack = stack[:runtime.Stack(stack, false)]
			attrs = append(attrs, slog.String("stack", string(stack)))
		}
	}

	l.LogAttrs(context.Background(), slog.LevelError, msg, attrs...)
}

func (l *Logger) Fatal(msg string, err error, attrs ...slog.Attr) {
	l.Error(msg, err, attrs...)
	os.Exit(1)
}

func (l *Logger) WithComponent(component string) *Logger {
	return &Logger{l.With(slog.String("component", component))}
}

func (l *Logger) WithOperation(operation string) *Logger {
	return &Logger{l.With(slog.String("operation", operation))}
}

func (l *Logger) TimeTrack(start time.Time, name string, attrs ...slog.Attr) {
	elapsed := time.Since(start)
	attrs = append(attrs, slog.Duration("duration", elapsed))
	l.LogAttrs(context.Background(), slog.LevelDebug, name+" completed", attrs...)
}
