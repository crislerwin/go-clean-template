package nethttp

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/crislerwin/go-clean-template/internal/application/user"
	"github.com/crislerwin/go-clean-template/internal/infrastructure/persistence/memory"
	"github.com/crislerwin/go-clean-template/internal/ports/input"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Testes de infraestrutura HTTP verificam apenas a tradução de status/corpo.
// Regras de negócio já foram testadas nos níveis internos.
func TestUserHandler(t *testing.T) {
	repo := memory.NewUserRepository()
	createUC := user.NewCreateUserUseCase(repo)
	findUC := user.NewFindUserByIDUseCase(repo)
	listUC := user.NewListUsersUseCase(repo)
	handler := NewUserHandler(createUC, findUC, listUC)

	mux := http.NewServeMux()
	handler.RegisterRoutes(mux)

	t.Run("creates a user", func(t *testing.T) {
		body, _ := json.Marshal(input.CreateUserInput{Name: "Dennis Ritchie", Email: "dennis@example.com"})
		req := httptest.NewRequest(http.MethodPost, "/api/v1/users", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		mux.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusCreated, rec.Code)
		assert.Contains(t, rec.Body.String(), "Dennis Ritchie")
	})

	t.Run("returns bad request on invalid create payload", func(t *testing.T) {
		body, _ := json.Marshal(input.CreateUserInput{Name: "", Email: "bad"})
		req := httptest.NewRequest(http.MethodPost, "/api/v1/users", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		mux.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("lists users", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/v1/users", nil)
		rec := httptest.NewRecorder()

		mux.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Contains(t, rec.Body.String(), "users")
	})
}

func TestUserHandler_FindByID(t *testing.T) {
	repo := memory.NewUserRepository()
	createUC := user.NewCreateUserUseCase(repo)
	findUC := user.NewFindUserByIDUseCase(repo)
	listUC := user.NewListUsersUseCase(repo)
	handler := NewUserHandler(createUC, findUC, listUC)

	mux := http.NewServeMux()
	handler.RegisterRoutes(mux)

	created, err := createUC.Execute(input.CreateUserInput{Name: "Barbara Liskov", Email: "barbara@example.com"})
	require.NoError(t, err)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/users/"+created.User.ID, nil)
	rec := httptest.NewRecorder()

	mux.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Contains(t, rec.Body.String(), "Barbara Liskov")
}
