package logger

import (
	"bytes"
	"context"
	"encoding/json"
	"log/slog"
	"testing"

	"github.com/stretchr/testify/assert"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"

	"github.com/crislerwin/go-clean-template/internal/ports/telemetry"
)

func TestSlogLogger_JSONOutput(t *testing.T) {
	tests := []struct {
		name     string
		level    string
		logFn    func(ctx context.Context, l telemetry.Logger)
		expected []string
	}{
		{
			name:     "info log contains message",
			level:    "INFO",
			logFn:    func(ctx context.Context, l telemetry.Logger) { l.Info(ctx, "user created", "user_id", "123") },
			expected: []string{"user created", "user_id", "123"},
		},
		{
			name:     "warn log is emitted when level allows",
			level:    "WARN",
			logFn:    func(ctx context.Context, l telemetry.Logger) { l.Warn(ctx, "slow request", "duration_ms", 42) },
			expected: []string{"slow request", "duration_ms", "42"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			l := NewSlogLogger(&buf, tt.level)
			tt.logFn(context.Background(), l)

			for _, exp := range tt.expected {
				assert.Contains(t, buf.String(), exp)
			}
		})
	}
}

func TestSlogLogger_SkipsLowerLevel(t *testing.T) {
	var buf bytes.Buffer
	l := NewSlogLogger(&buf, "ERROR")

	l.Info(context.Background(), "should not appear")

	assert.Empty(t, buf.String())
}

func TestSlogLogger_WithFields(t *testing.T) {
	var buf bytes.Buffer
	l := NewSlogLogger(&buf, "INFO").With("service", "go-clean-template")

	l.Info(context.Background(), "hello")

	assert.Contains(t, buf.String(), "service")
	assert.Contains(t, buf.String(), "go-clean-template")
	assert.Contains(t, buf.String(), "hello")
}

func TestSlogLogger_TraceCorrelation(t *testing.T) {
	var buf bytes.Buffer
	l := NewSlogLogger(&buf, "INFO")

	tp := sdktrace.NewTracerProvider()
	tracer := tp.Tracer("test")
	ctx, span := tracer.Start(context.Background(), "test-span")
	defer span.End()

	l.Info(ctx, "correlated log")

	var payload map[string]any
	assert.NoError(t, json.Unmarshal(buf.Bytes(), &payload))
	assert.NotEmpty(t, payload["trace_id"])
	assert.NotEmpty(t, payload["span_id"])
}

var _ telemetry.Logger = (*slogLogger)(nil)
var _ = slog.LevelDebug
