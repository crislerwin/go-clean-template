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
		{name: "starts and ends span without panic"},
	}

	tracer := NewNoOpTracer()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, span := tracer.Start(context.Background(), "test-span")
			assert.NotNil(t, ctx)
			assert.NotNil(t, span)
			span.End()
			assert.NoError(t, tracer.Shutdown(context.Background()))
		})
	}
}

// Ensure context import is used.
var _ = context.Background
