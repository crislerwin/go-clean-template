package user

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Testes table-driven revelam intenção e cobrem casos de borda de forma compacta.
func TestNewUser(t *testing.T) {
	tests := []struct {
		name      string
		userName  string
		email     string
		wantError error
	}{
		{
			name:      "creates user with valid data",
			userName:  "Ada Lovelace",
			email:     "ada@example.com",
			wantError: nil,
		},
		{
			name:      "rejects empty name",
			userName:  "",
			email:     "ada@example.com",
			wantError: ErrInvalidName,
		},
		{
			name:      "rejects empty email",
			userName:  "Ada Lovelace",
			email:     "",
			wantError: ErrInvalidEmail,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user, err := NewUser(tt.userName, tt.email)

			if tt.wantError != nil {
				require.ErrorIs(t, err, tt.wantError)
				assert.Nil(t, user)
				return
			}

			require.NoError(t, err)
			require.NotNil(t, user)
			assert.NotEmpty(t, user.ID)
			assert.Equal(t, tt.userName, user.Name)
			assert.Equal(t, tt.email, user.Email)
			assert.False(t, user.CreatedAt.IsZero())
			assert.Equal(t, user.CreatedAt, user.UpdatedAt)
		})
	}
}

func TestUser_UpdateName(t *testing.T) {
	user, err := NewUser("Grace Hopper", "grace@example.com")
	require.NoError(t, err)

	tests := []struct {
		name      string
		newName   string
		wantError error
	}{
		{
			name:      "updates name successfully",
			newName:   "Grace Murray Hopper",
			wantError: nil,
		},
		{
			name:      "rejects empty name",
			newName:   "",
			wantError: ErrInvalidName,
		},
	}

	previousUpdatedAt := user.UpdatedAt
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := user.UpdateName(tt.newName)

			if tt.wantError != nil {
				require.ErrorIs(t, err, tt.wantError)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tt.newName, user.Name)
			assert.True(t, user.UpdatedAt.After(previousUpdatedAt))
		})
	}
}
