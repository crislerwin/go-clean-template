package nethttp

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/crislerwin/go-clean-template/internal/application/user"
	logger "github.com/crislerwin/go-clean-template/internal/infrastructure/logger/slog"
	"github.com/crislerwin/go-clean-template/internal/infrastructure/persistence/memory"
	"github.com/crislerwin/go-clean-template/internal/infrastructure/telemetry/otlp"
	inputports "github.com/crislerwin/go-clean-template/internal/ports/input"
	"github.com/stretchr/testify/assert"
)

func setupHandler() (*UserHandler, *http.ServeMux) {
	repo := memory.NewUserRepository()
	tracer := otlp.NewNoOpTracer()
	logger := logger.NewSlogLogger(nil, "INFO")

	createUC := user.NewCreateUserUseCase(repo, logger)
	findUC := user.NewFindUserByIDUseCase(repo, logger)
	listUC := user.NewListUsersUseCase(repo, logger)

	handler := NewUserHandler(createUC, findUC, listUC, tracer, logger)
	mux := http.NewServeMux()
	handler.RegisterRoutes(mux)
	return handler, mux
}

func TestUserHandler_CreateUser(t *testing.T) {
	_, mux := setupHandler()

	tests := []struct {
		name       string
		body       inputports.CreateUserInput
		wantStatus int
		wantErr    bool
	}{
		{
			name:       "creates user",
			body:       inputports.CreateUserInput{Name: "Ada Lovelace", Email: "ada@example.com"},
			wantStatus: http.StatusCreated,
		},
		{
			name:       "returns bad request for empty name",
			body:       inputports.CreateUserInput{Name: "", Email: "ada@example.com"},
			wantStatus: http.StatusBadRequest,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			payload, _ := json.Marshal(tt.body)
			req := httptest.NewRequest(http.MethodPost, "/api/v1/users", bytes.NewReader(payload))
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()

			mux.ServeHTTP(rec, req)

			assert.Equal(t, tt.wantStatus, rec.Code)
		})
	}
}

func TestUserHandler_ListUsers(t *testing.T) {
	_, mux := setupHandler()

	t.Run("returns empty list", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/v1/users", nil)
		rec := httptest.NewRecorder()

		mux.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Contains(t, rec.Body.String(), "users")
	})

	t.Run("returns health status", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/health", nil)
		rec := httptest.NewRecorder()

		mux.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Contains(t, rec.Body.String(), "healthy")
	})
}

func TestUserHandler_FindByID(t *testing.T) {
	_, mux := setupHandler()

	t.Run("returns not found for missing user", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/v1/users/non-existent-id", nil)
		rec := httptest.NewRecorder()

		mux.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusNotFound, rec.Code)
	})
}
