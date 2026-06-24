package user

import (
	"context"
	"testing"

	"github.com/crislerwin/go-clean-template/internal/domain/user"
	logger "github.com/crislerwin/go-clean-template/internal/infrastructure/logger/slog"
	"github.com/crislerwin/go-clean-template/internal/infrastructure/persistence/memory"
	inputports "github.com/crislerwin/go-clean-template/internal/ports/input"
	"github.com/crislerwin/go-clean-template/internal/ports/output"
	"github.com/stretchr/testify/assert"
)

func TestCreateUserUseCase_Execute(t *testing.T) {
	tests := []struct {
		name      string
		input     inputports.CreateUserInput
		wantErr   bool
		errIs     error
		wantSaved bool
	}{
		{
			name:      "creates valid user",
			input:     inputports.CreateUserInput{Name: "Ada Lovelace", Email: "ada@example.com"},
			wantSaved: true,
		},
		{
			name:    "returns error for empty name",
			input:   inputports.CreateUserInput{Name: "", Email: "ada@example.com"},
			wantErr: true,
			errIs:   user.ErrInvalidName,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := memory.NewUserRepository()
			logger := logger.NewSlogLogger(nil, "INFO")
			uc := NewCreateUserUseCase(repo, logger)

			output, err := uc.Execute(context.Background(), tt.input)

			if tt.wantErr {
				assert.Error(t, err)
				if tt.errIs != nil {
					assert.ErrorIs(t, err, tt.errIs)
				}
				return
			}

			assert.NoError(t, err)
			assert.NotNil(t, output)
			assert.NotEmpty(t, output.User.ID)
			assert.Equal(t, tt.input.Name, output.User.Name)
			assert.Equal(t, tt.input.Email, output.User.Email)

			saved, err := repo.FindByID(context.Background(), output.User.ID)
			assert.NoError(t, err)
			assert.Equal(t, output.User.ID, saved.ID)
		})
	}
}

func TestFindUserByIDUseCase_Execute(t *testing.T) {
	tests := []struct {
		name    string
		setup   func(repo output.UserRepository)
		input   inputports.FindUserByIDInput
		wantErr bool
		errIs   error
	}{
		{
			name: "finds existing user",
			setup: func(repo output.UserRepository) {
				u, _ := user.NewUser("Grace Hopper", "grace@example.com")
				_ = repo.Save(context.Background(), u)
			},
			input: inputports.FindUserByIDInput{ID: ""},
		},
		{
			name:    "returns not found for missing user",
			input:   inputports.FindUserByIDInput{ID: "non-existent-id"},
			wantErr: true,
			errIs:   user.ErrUserNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := memory.NewUserRepository()
			logger := logger.NewSlogLogger(nil, "INFO")
			uc := NewFindUserByIDUseCase(repo, logger)

			var userID string
			if tt.setup != nil {
				tt.setup(repo)
				users, _ := repo.FindAll(context.Background())
				userID = users[0].ID
				tt.input.ID = userID
			}

			output, err := uc.Execute(context.Background(), tt.input)

			if tt.wantErr {
				assert.Error(t, err)
				if tt.errIs != nil {
					assert.ErrorIs(t, err, tt.errIs)
				}
				return
			}

			assert.NoError(t, err)
			assert.NotNil(t, output)
			assert.Equal(t, userID, output.User.ID)
		})
	}
}

func TestListUsersUseCase_Execute(t *testing.T) {
	tests := []struct {
		name      string
		setup     func(repo output.UserRepository)
		wantCount int
	}{
		{
			name:      "returns empty list",
			wantCount: 0,
		},
		{
			name: "returns all users",
			setup: func(repo output.UserRepository) {
				u1, _ := user.NewUser("Grace Hopper", "grace@example.com")
				u2, _ := user.NewUser("Ada Lovelace", "ada@example.com")
				_ = repo.Save(context.Background(), u1)
				_ = repo.Save(context.Background(), u2)
			},
			wantCount: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := memory.NewUserRepository()
			logger := logger.NewSlogLogger(nil, "INFO")
			uc := NewListUsersUseCase(repo, logger)

			if tt.setup != nil {
				tt.setup(repo)
			}

			output, err := uc.Execute(context.Background())

			assert.NoError(t, err)
			assert.Len(t, output.Users, tt.wantCount)
		})
	}
}
