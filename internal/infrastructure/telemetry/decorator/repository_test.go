package decorator

import (
	"context"
	"errors"
	"testing"

	"github.com/crislerwin/go-clean-template/internal/domain/user"
	logger "github.com/crislerwin/go-clean-template/internal/infrastructure/logger/slog"
	"github.com/crislerwin/go-clean-template/internal/infrastructure/persistence/memory"
	"github.com/crislerwin/go-clean-template/internal/infrastructure/telemetry/otlp"
	"github.com/crislerwin/go-clean-template/internal/ports/output"
	"github.com/stretchr/testify/assert"
)

func TestTracedRepository_PropagatesContext(t *testing.T) {
	tests := []struct {
		name string
		fn   func(repo output.UserRepository) error
	}{
		{
			name: "save",
			fn: func(repo output.UserRepository) error {
				u, _ := user.NewUser("Test", "test@example.com")
				return repo.Save(context.Background(), u)
			},
		},
		{
			name: "find by id",
			fn: func(repo output.UserRepository) error {
				_, err := repo.FindByID(context.Background(), "any-id")
				return err
			},
		},
		{
			name: "find all",
			fn: func(repo output.UserRepository) error {
				_, err := repo.FindAll(context.Background())
				return err
			},
		},
	}

	tracer := otlp.NewNoOpTracer()
	logger := logger.NewSlogLogger(nil, "INFO")

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			base := memory.NewUserRepository()
			repo := NewTracedRepository(base, tracer, logger)

			err := tt.fn(repo)

			if tt.name == "find by id" {
				assert.True(t, errors.Is(err, user.ErrUserNotFound))
				return
			}
			assert.NoError(t, err)
		})
	}
}
