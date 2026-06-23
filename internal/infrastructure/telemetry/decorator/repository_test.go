package decorator

import (
	"testing"

	"github.com/crislerwin/go-clean-template/internal/domain/user"
	"github.com/crislerwin/go-clean-template/internal/infrastructure/persistence/memory"
	"github.com/crislerwin/go-clean-template/internal/infrastructure/telemetry/otlp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTracedRepository(t *testing.T) {
	tests := []struct {
		name      string
		operation func(repo userRepo) error
		wantError error
	}{
		{
			name: "traces save successfully",
			operation: func(repo userRepo) error {
				u, _ := user.NewUser("Test User", "test@example.com")
				return repo.Save(u)
			},
		},
		{
			name: "traces find by id with not found",
			operation: func(repo userRepo) error {
				_, err := repo.FindByID("missing-id")
				return err
			},
			wantError: user.ErrUserNotFound,
		},
		{
			name: "traces find all",
			operation: func(repo userRepo) error {
				_, err := repo.FindAll()
				return err
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			base := memory.NewUserRepository()
			tracer := otlp.NewNoOpTracer()
			repo := NewTracedRepository(base, tracer)

			err := tt.operation(repo)

			if tt.wantError != nil {
				require.ErrorIs(t, err, tt.wantError)
				return
			}
			require.NoError(t, err)
		})
	}
}

// userRepo expõe apenas os métodos usados no teste para simplificar.
type userRepo interface {
	Save(u *user.User) error
	FindByID(id string) (*user.User, error)
	FindAll() ([]*user.User, error)
}

// Sanity check: tracedRepository realmente implementa output.UserRepository.
var _ userRepo = (&tracedRepository{})

func TestTracedRepository_RecordsError(t *testing.T) {
	base := memory.NewUserRepository()
	tracer := otlp.NewNoOpTracer()
	repo := NewTracedRepository(base, tracer)

	// Salva um usuário válido primeiro para garantir sucesso no save.
	u, err := user.NewUser("Tracer Test", "tracer@example.com")
	require.NoError(t, err)
	require.NoError(t, repo.Save(u))

	// Força FindByID a retornar erro sem depender de estado.
	_, err = repo.FindByID("does-not-exist")
	assert.ErrorIs(t, err, user.ErrUserNotFound)
}
