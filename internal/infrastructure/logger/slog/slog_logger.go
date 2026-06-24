package logger

import (
	"context"
	"io"
	"log/slog"
	"os"

	"go.opentelemetry.io/otel/trace"

	"github.com/crislerwin/go-clean-template/internal/ports/telemetry"
)

// slogLogger é a implementação padrão usando log/slog do stdlib.
// Emite JSON estruturado, ideal para Loki/Grafana.
type slogLogger struct {
	logger *slog.Logger
}

// NewSlogLogger cria um logger estruturado em JSON.
func NewSlogLogger(output io.Writer, level string) telemetry.Logger {
	if output == nil {
		output = os.Stdout
	}

	var slogLevel slog.Level
	switch level {
	case "DEBUG":
		slogLevel = slog.LevelDebug
	case "INFO":
		slogLevel = slog.LevelInfo
	case "WARN":
		slogLevel = slog.LevelWarn
	case "ERROR":
		slogLevel = slog.LevelError
	default:
		slogLevel = slog.LevelInfo
	}

	handler := slog.NewJSONHandler(output, &slog.HandlerOptions{
		Level: slogLevel,
	})

	return &slogLogger{logger: slog.New(handler)}
}

func (l *slogLogger) Debug(ctx context.Context, msg string, args ...any) {
	l.logger.DebugContext(ctx, msg, append(args, traceContextArgs(ctx)...)...)
}

func (l *slogLogger) Info(ctx context.Context, msg string, args ...any) {
	l.logger.InfoContext(ctx, msg, append(args, traceContextArgs(ctx)...)...)
}

func (l *slogLogger) Warn(ctx context.Context, msg string, args ...any) {
	l.logger.WarnContext(ctx, msg, append(args, traceContextArgs(ctx)...)...)
}

func (l *slogLogger) Error(ctx context.Context, msg string, args ...any) {
	l.logger.ErrorContext(ctx, msg, append(args, traceContextArgs(ctx)...)...)
}

func (l *slogLogger) With(args ...any) telemetry.Logger {
	return &slogLogger{logger: l.logger.With(args...)}
}

// traceContextArgs extrai trace/span IDs do contexto para correlação logs↔traces.
func traceContextArgs(ctx context.Context) []any {
	span := trace.SpanFromContext(ctx)
	if span == nil || !span.SpanContext().IsValid() {
		return nil
	}
	sc := span.SpanContext()
	return []any{
		"trace_id", sc.TraceID().String(),
		"span_id", sc.SpanID().String(),
		"trace_flags", sc.TraceFlags().String(),
	}
}

var _ telemetry.Logger = (*slogLogger)(nil)
