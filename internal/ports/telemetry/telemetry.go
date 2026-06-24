package telemetry

import "context"

// Tracer é a abstração mínima que a aplicação precisa.
// Qualquer adapter de telemetria (OTLP, Jaeger, stdout) pode implementá-lo.
type Tracer interface {
	Start(ctx context.Context, name string) (context.Context, Span)
	Shutdown(ctx context.Context) error
}

// Span representa uma operação dentro de um trace.
type Span interface {
	End()
	RecordError(err error)
}

// Logger é a abstração mínima de logging estruturado.
// A aplicação o usa para emitir logs correlacionados com traces.
// O método With retorna Logger (da mesma interface) para manter o tipo interno.
type Logger interface {
	Debug(ctx context.Context, msg string, args ...any)
	Info(ctx context.Context, msg string, args ...any)
	Warn(ctx context.Context, msg string, args ...any)
	Error(ctx context.Context, msg string, args ...any)
	With(args ...any) Logger
}
