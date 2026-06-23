package otlp

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNoOpTracer(t *testing.T) {
	tests := []struct {
		name string
	}{
		{name: "starts and ends span without side effects"},
		{name: "records error without panic"},
	}

	tracer := NewNoOpTracer()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			ctx, span := tracer.Start(ctx, "operation")
			assert.NotNil(t, ctx)
			span.RecordError(assert.AnError)
			span.End()

			err := tracer.Shutdown(ctx)
			assert.NoError(t, err)
		})
	}
}
