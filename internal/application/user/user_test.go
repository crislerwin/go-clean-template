package user

import (
	"testing"

	"github.com/crislerwin/go-clean-template/internal/domain/user"
	"github.com/crislerwin/go-clean-template/internal/infrastructure/persistence/memory"
	"github.com/crislerwin/go-clean-template/internal/infrastructure/telemetry/otlp"
	"github.com/crislerwin/go-clean-template/internal/ports/input"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Testes de aplicação usam o repositório em memória e um tracer no-op
// para provar que o caso de uso funciona independentemente de tecnologia.
func TestCreateUserUseCase_Execute(t *testing.T) {
	tests := []struct {
		name      string
		input     input.CreateUserInput
		wantError error
	}{
		{
			name:  "creates and persists user",
			input: input.CreateUserInput{Name: "Alan Turing", Email: "alan@example.com"},
		},
		{
			name:      "rejects invalid input",
			input:     input.CreateUserInput{Name: "", Email: "alan@example.com"},
			wantError: user.ErrInvalidName,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := memory.NewUserRepository()
			tracer := otlp.NewNoOpTracer()
			uc := NewCreateUserUseCase(repo, tracer)

			output, err := uc.Execute(tt.input)

			if tt.wantError != nil {
				require.ErrorIs(t, err, tt.wantError)
				return
			}

			require.NoError(t, err)
			require.NotNil(t, output.User)
			assert.Equal(t, tt.input.Name, output.User.Name)
			assert.Equal(t, tt.input.Email, output.User.Email)

			found, err := repo.FindByID(output.User.ID)
			require.NoError(t, err)
			assert.Equal(t, output.User.ID, found.ID)
		})
	}
}

func TestFindUserByIDUseCase_Execute(t *testing.T) {
	repo := memory.NewUserRepository()
	tracer := otlp.NewNoOpTracer()
	created, err := user.NewUser("Tim Berners-Lee", "tim@example.com")
	require.NoError(t, err)
	require.NoError(t, repo.Save(created))

	uc := NewFindUserByIDUseCase(repo, tracer)

	output, err := uc.Execute(input.FindUserByIDInput{ID: created.ID})
	require.NoError(t, err)
	assert.Equal(t, created.ID, output.User.ID)

	_, err = uc.Execute(input.FindUserByIDInput{ID: "non-existent-id"})
	assert.ErrorIs(t, err, user.ErrUserNotFound)
}

func TestListUsersUseCase_Execute(t *testing.T) {
	repo := memory.NewUserRepository()
	tracer := otlp.NewNoOpTracer()
	uc := NewListUsersUseCase(repo, tracer)

	output, err := uc.Execute()
	require.NoError(t, err)
	assert.Empty(t, output.Users)

	created, err := user.NewUser("Linus Torvalds", "linus@example.com")
	require.NoError(t, err)
	require.NoError(t, repo.Save(created))

	output, err = uc.Execute()
	require.NoError(t, err)
	assert.Len(t, output.Users, 1)
	assert.Equal(t, created.ID, output.Users[0].ID)
}
