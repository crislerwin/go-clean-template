package nethttp

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/crislerwin/go-clean-template/internal/domain/user"
	"github.com/crislerwin/go-clean-template/internal/ports/input"
)

// UserHandler é o adapter HTTP. Ele converte requisições/respostas
// da biblioteca padrão net/http em chamadas aos ports de entrada da aplicação.
// Usar net/http mantém o template agnóstico: você pode trocar por Gin, Echo,
// Fiber ou qualquer outro framework sem tocar no domínio ou nos casos de uso.
type UserHandler struct {
	createUserUC input.CreateUserUseCase
	findUserUC   input.FindUserByIDUseCase
	listUsersUC  input.ListUsersUseCase
}

func NewUserHandler(
	create input.CreateUserUseCase,
	find input.FindUserByIDUseCase,
	list input.ListUsersUseCase,
) *UserHandler {
	return &UserHandler{
		createUserUC: create,
		findUserUC:   find,
		listUsersUC:  list,
	}
}

func (h *UserHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/v1/users", h.createUser)
	mux.HandleFunc("GET /api/v1/users/{id}", h.findUserByID)
	mux.HandleFunc("GET /api/v1/users", h.listUsers)
}

func (h *UserHandler) createUser(w http.ResponseWriter, r *http.Request) {
	var req input.CreateUserInput
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	if err := validateCreateInput(req); err != nil {
		respondJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	output, err := h.createUserUC.Execute(req)
	if err != nil {
		respondJSON(w, http.StatusUnprocessableEntity, map[string]string{"error": err.Error()})
		return
	}

	respondJSON(w, http.StatusCreated, output)
}

func (h *UserHandler) findUserByID(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if strings.TrimSpace(id) == "" {
		respondJSON(w, http.StatusBadRequest, map[string]string{"error": "id is required"})
		return
	}

	output, err := h.findUserUC.Execute(input.FindUserByIDInput{ID: id})
	if errors.Is(err, user.ErrUserNotFound) {
		respondJSON(w, http.StatusNotFound, map[string]string{"error": err.Error()})
		return
	}
	if err != nil {
		respondJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	respondJSON(w, http.StatusOK, output)
}

func (h *UserHandler) listUsers(w http.ResponseWriter, r *http.Request) {
	output, err := h.listUsersUC.Execute()
	if err != nil {
		respondJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	respondJSON(w, http.StatusOK, output)
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
