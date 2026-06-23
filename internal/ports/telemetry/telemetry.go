package telemetry

import "context"

// Tracer é a abstração mínima que a aplicação precisa.
// Qualquer adapter de telemetria (OTLP, Jaeger, stdout) pode implementá-lo.
type Tracer interface {
	Start(ctx context.Context, name string) (context.Context, Span)
	Shutdown(ctx context.Context) error
}

// Span representa uma operação traceada.
type Span interface {
	End()
	RecordError(err error)
}
