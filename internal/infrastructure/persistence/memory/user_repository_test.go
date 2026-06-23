package memory

import (
	"testing"

	"github.com/crislerwin/go-clean-template/internal/domain/user"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUserRepository_Save_FindByID_FindAll(t *testing.T) {
	tests := []struct {
		name  string
		users []*user.User
	}{
		{
			name: "saves and retrieves multiple users",
			users: []*user.User{
				mustCreateUser(t, "Margaret Hamilton", "margaret@example.com"),
				mustCreateUser(t, "Donald Knuth", "donald@example.com"),
			},
		},
		{
			name:  "handles empty repository",
			users: []*user.User{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := NewUserRepository()

			for _, u := range tt.users {
				require.NoError(t, repo.Save(u))
			}

			all, err := repo.FindAll()
			require.NoError(t, err)
			assert.Len(t, all, len(tt.users))

			for _, u := range tt.users {
				found, err := repo.FindByID(u.ID)
				require.NoError(t, err)
				assert.Equal(t, u.ID, found.ID)
			}
		})
	}
}

func TestUserRepository_FindByID_NotFound(t *testing.T) {
	repo := NewUserRepository()
	_, err := repo.FindByID("missing-id")
	assert.ErrorIs(t, err, user.ErrUserNotFound)
}

func mustCreateUser(t *testing.T, name, email string) *user.User {
	t.Helper()
	u, err := user.NewUser(name, email)
	require.NoError(t, err)
	return u
}
