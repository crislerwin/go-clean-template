package otlp

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
	"go.opentelemetry.io/otel/trace"

	"github.com/crislerwin/go-clean-template/internal/ports/telemetry"
)

// otelSpan adapta trace.Span para nossa interface interna.
type otelSpan struct {
	span trace.Span
}

func (s *otelSpan) End() {
	s.span.End()
}

func (s *otelSpan) RecordError(err error) {
	s.span.RecordError(err)
}

// otelTracer é o adapter OTLP/HTTP padrão.
// Ele se conecta a um OpenTelemetry Collector via OTEL_EXPORTER_OTLP_ENDPOINT.
// Se a variável não estiver configurada, o tracer funciona como no-op,
// mantendo a aplicação desacoplada da infra de observabilidade.
type otelTracer struct {
	tracer   trace.Tracer
	provider *sdktrace.TracerProvider
}

// NewOTelTracer cria o tracer.
//
// Comportamento:
//   - Se OTEL_EXPORTER_OTLP_ENDPOINT estiver vazio, retorna NoOpTracer.
//   - Se a inicialização do exporter falhar, retorna NoOpTracer para manter
//     a aplicação resiliente (escolha de resiliência sobre rigidez).
//   - TLS é controlado por OTEL_INSECURE: true (padrão) usa WithInsecure();
//     false exige TLS via WithSecure().
func NewOTelTracer(serviceName string) (telemetry.Tracer, error) {
	endpoint := os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT")
	if endpoint == "" {
		return NewNoOpTracer(), nil
	}

	insecure, err := strconv.ParseBool(os.Getenv("OTEL_INSECURE"))
	if err != nil {
		// Padrão do template é insecure=true para facilitar desenvolvimento local.
		insecure = true
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	opts := []otlptracehttp.Option{
		otlptracehttp.WithEndpoint(endpoint),
	}
	if insecure {
		opts = append(opts, otlptracehttp.WithInsecure())
	}

	exporter, err := otlptracehttp.New(ctx, opts...)
	if err != nil {
		return NewNoOpTracer(), fmt.Errorf("failed to create otel exporter, falling back to no-op: %w", err)
	}

	res, err := resource.New(ctx,
		resource.WithAttributes(
			semconv.ServiceNameKey.String(serviceName),
			semconv.ServiceVersionKey.String("1.0.0"),
		),
	)
	if err != nil {
		return NewNoOpTracer(), fmt.Errorf("failed to create otel resource, falling back to no-op: %w", err)
	}

	provider := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(res),
	)

	otel.SetTracerProvider(provider)

	return &otelTracer{
		tracer:   provider.Tracer(serviceName),
		provider: provider,
	}, nil
}

func (t *otelTracer) Start(ctx context.Context, name string) (context.Context, telemetry.Span) {
	ctx, span := t.tracer.Start(ctx, name)
	return ctx, &otelSpan{span: span}
}

func (t *otelTracer) Shutdown(ctx context.Context) error {
	return t.provider.Shutdown(ctx)
}

// noOpTracer é o fallback quando OTLP não está configurado.
// Ele preserva as chamadas de tracing no código sem gerar spans reais.
type noOpTracer struct{}

func NewNoOpTracer() telemetry.Tracer {
	return &noOpTracer{}
}

func (n *noOpTracer) Start(ctx context.Context, name string) (context.Context, telemetry.Span) {
	return ctx, &noOpSpan{}
}

func (n *noOpTracer) Shutdown(ctx context.Context) error {
	return nil
}

type noOpSpan struct{}

func (n *noOpSpan) End()              {}
func (n *noOpSpan) RecordError(error) {}
