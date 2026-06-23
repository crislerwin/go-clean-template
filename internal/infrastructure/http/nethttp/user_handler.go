package nethttp

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/crislerwin/go-clean-template/internal/domain/user"
	"github.com/crislerwin/go-clean-template/internal/ports/input"
	"github.com/crislerwin/go-clean-template/internal/ports/telemetry"
)

// UserHandler é o adapter HTTP. Ele converte requisições/respostas
// da biblioteca padrão net/http em chamadas aos ports de entrada da aplicação.
// Usar net/http mantém o template agnóstico: você pode trocar por Gin, Echo,
// Fiber ou qualquer outro framework sem tocar no domínio ou nos casos de uso.
type UserHandler struct {
	createUserUC input.CreateUserUseCase
	findUserUC   input.FindUserByIDUseCase
	listUsersUC  input.ListUsersUseCase
	tracer       telemetry.Tracer
}

func NewUserHandler(
	create input.CreateUserUseCase,
	find input.FindUserByIDUseCase,
	list input.ListUsersUseCase,
	tracer telemetry.Tracer,
) *UserHandler {
	return &UserHandler{
		createUserUC: create,
		findUserUC:   find,
		listUsersUC:  list,
		tracer:       tracer,
	}
}

func (h *UserHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/v1/users", h.createUser)
	mux.HandleFunc("GET /api/v1/users/{id}", h.findUserByID)
	mux.HandleFunc("GET /api/v1/users", h.listUsers)
}

func (h *UserHandler) createUser(w http.ResponseWriter, r *http.Request) {
	ctx, span := h.tracer.Start(r.Context(), "HTTP.CreateUser")
	defer span.End()

	var req input.CreateUserInput
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		span.RecordError(err)
		respondJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	if err := validateCreateInput(req); err != nil {
		span.RecordError(err)
		respondJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	output, err := h.createUserUC.Execute(req)
	if err != nil {
		span.RecordError(err)
		respondJSON(w, http.StatusUnprocessableEntity, map[string]string{"error": err.Error()})
		return
	}

	respondJSON(w, http.StatusCreated, output)
	_ = ctx
}

func (h *UserHandler) findUserByID(w http.ResponseWriter, r *http.Request) {
	ctx, span := h.tracer.Start(r.Context(), "HTTP.FindUserByID")
	defer span.End()

	id := r.PathValue("id")
	if strings.TrimSpace(id) == "" {
		err := errors.New("id is required")
		span.RecordError(err)
		respondJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	output, err := h.findUserUC.Execute(input.FindUserByIDInput{ID: id})
	if errors.Is(err, user.ErrUserNotFound) {
		span.RecordError(err)
		respondJSON(w, http.StatusNotFound, map[string]string{"error": err.Error()})
		return
	}
	if err != nil {
		span.RecordError(err)
		respondJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	respondJSON(w, http.StatusOK, output)
	_ = ctx
}

func (h *UserHandler) listUsers(w http.ResponseWriter, r *http.Request) {
	ctx, span := h.tracer.Start(r.Context(), "HTTP.ListUsers")
	defer span.End()

	output, err := h.listUsersUC.Execute()
	if err != nil {
		span.RecordError(err)
		respondJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	respondJSON(w, http.StatusOK, output)
	_ = ctx
}

func validateCreateInput(req input.CreateUserInput) error {
	if strings.TrimSpace(req.Name) == "" {
		return errors.New("name is required")
	}
	if strings.TrimSpace(req.Email) == "" {
		return errors.New("email is required")
	}
	return nil
}

func respondJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}
