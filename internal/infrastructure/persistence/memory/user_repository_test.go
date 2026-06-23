package memory

import (
	"context"
	"errors"
	"testing"

	"github.com/crislerwin/go-clean-template/internal/domain/user"
	"github.com/crislerwin/go-clean-template/internal/ports/output"
	"github.com/stretchr/testify/assert"
)

func TestUserRepository_SaveAndFind(t *testing.T) {
	tests := []struct {
		name    string
		setup   func(repo output.UserRepository)
		id      string
		wantErr bool
	}{
		{
			name: "finds saved user",
			setup: func(repo output.UserRepository) {
				u, _ := user.NewUser("Alan Turing", "alan@example.com")
				_ = repo.Save(context.Background(), u)
			},
			wantErr: false,
		},
		{
			name:    "returns error for missing user",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := NewUserRepository()
			if tt.setup != nil {
				tt.setup(repo)
				users, _ := repo.FindAll(context.Background())
				tt.id = users[0].ID
			} else {
				tt.id = "missing-id"
			}

			found, err := repo.FindByID(context.Background(), tt.id)

			if tt.wantErr {
				assert.Error(t, err)
				assert.True(t, errors.Is(err, user.ErrUserNotFound))
				return
			}

			assert.NoError(t, err)
			assert.NotNil(t, found)
			assert.Equal(t, tt.id, found.ID)
		})
	}
}

func TestUserRepository_FindAll(t *testing.T) {
	tests := []struct {
		name      string
		setup     func(repo output.UserRepository)
		wantCount int
	}{
		{
			name:      "empty repository returns empty list",
			wantCount: 0,
		},
		{
			name: "returns all saved users",
			setup: func(repo output.UserRepository) {
				u1, _ := user.NewUser("Alan Turing", "alan@example.com")
				u2, _ := user.NewUser("Grace Hopper", "grace@example.com")
				_ = repo.Save(context.Background(), u1)
				_ = repo.Save(context.Background(), u2)
			},
			wantCount: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := NewUserRepository()
			if tt.setup != nil {
				tt.setup(repo)
			}

			users, err := repo.FindAll(context.Background())

			assert.NoError(t, err)
			assert.Len(t, users, tt.wantCount)
		})
	}
}
