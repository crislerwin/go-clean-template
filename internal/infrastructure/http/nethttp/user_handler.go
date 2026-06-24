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

// UserHandler é o adapter HTTP. Ele converte requisições/respostas HTTP
// em chamadas aos casos de uso (ports de entrada), mantendo o domínio
// e a aplicação completamente livres de detalhes de transporte.
type UserHandler struct {
	createUserUC   input.CreateUserUseCase
	findUserByIDUC input.FindUserByIDUseCase
	listUsersUC    input.ListUsersUseCase
	tracer         telemetry.Tracer
	logger         telemetry.Logger
}

// NewUserHandler injeta as portas de entrada no adapter HTTP.
func NewUserHandler(
	createUC input.CreateUserUseCase,
	findUC input.FindUserByIDUseCase,
	listUC input.ListUsersUseCase,
	tracer telemetry.Tracer,
	logger telemetry.Logger,
) *UserHandler {
	return &UserHandler{
		createUserUC:   createUC,
		findUserByIDUC: findUC,
		listUsersUC:    listUC,
		tracer:         tracer,
		logger:         logger.With("component", "UserHandler"),
	}
}

// RegisterRoutes conecta as rotas ao mux do net/http padrão.
// Usar o mux nativo mantém o template agnóstico de frameworks.
func (h *UserHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/v1/users", h.createUser)
	mux.HandleFunc("GET /api/v1/users/{id}", h.findUserByID)
	mux.HandleFunc("GET /api/v1/users", h.listUsers)
	mux.HandleFunc("GET /health", h.healthCheck)
}

func (h *UserHandler) createUser(w http.ResponseWriter, r *http.Request) {
	ctx, span := h.tracer.Start(r.Context(), "HTTP.CreateUser")
	defer span.End()

	h.logger.Info(ctx, "received create user request")

	var req input.CreateUserInput
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		span.RecordError(err)
		h.logger.Error(ctx, "failed to decode request", "error", err)
		respondJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	if err := validateCreateInput(req); err != nil {
		span.RecordError(err)
		h.logger.Warn(ctx, "invalid create user request", "error", err)
		respondJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	output, err := h.createUserUC.Execute(ctx, req)
	if err != nil {
		span.RecordError(err)
		h.logger.Error(ctx, "failed to execute create user use case", "error", err)
		respondJSON(w, http.StatusUnprocessableEntity, map[string]string{"error": err.Error()})
		return
	}

	h.logger.Info(ctx, "user created via HTTP", "user_id", output.User.ID)
	respondJSON(w, http.StatusCreated, output)
}

func (h *UserHandler) findUserByID(w http.ResponseWriter, r *http.Request) {
	ctx, span := h.tracer.Start(r.Context(), "HTTP.FindUserByID")
	defer span.End()

	id := r.PathValue("id")
	h.logger.Info(ctx, "received find user by id request", "user_id", id)

	if strings.TrimSpace(id) == "" {
		err := errors.New("id is required")
		span.RecordError(err)
		h.logger.Warn(ctx, "missing user id in request")
		respondJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	output, err := h.findUserByIDUC.Execute(ctx, input.FindUserByIDInput{ID: id})
	if errors.Is(err, user.ErrUserNotFound) {
		span.RecordError(err)
		h.logger.Warn(ctx, "user not found", "user_id", id)
		respondJSON(w, http.StatusNotFound, map[string]string{"error": err.Error()})
		return
	}
	if err != nil {
		span.RecordError(err)
		h.logger.Error(ctx, "failed to execute find user by id use case", "error", err)
		respondJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	h.logger.Info(ctx, "user found via HTTP", "user_id", output.User.ID)
	respondJSON(w, http.StatusOK, output)
}

func (h *UserHandler) listUsers(w http.ResponseWriter, r *http.Request) {
	ctx, span := h.tracer.Start(r.Context(), "HTTP.ListUsers")
	defer span.End()

	h.logger.Info(ctx, "received list users request")

	output, err := h.listUsersUC.Execute(ctx)
	if err != nil {
		span.RecordError(err)
		h.logger.Error(ctx, "failed to execute list users use case", "error", err)
		respondJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	h.logger.Info(ctx, "users listed via HTTP", "count", len(output.Users))
	respondJSON(w, http.StatusOK, output)
}

// healthCheck expõe um endpoint simples de readiness/liveness.
// Ele não depende de repositório ou lógica de negócio, apenas confirma
// que o servidor HTTP está aceitando conexões.
func (h *UserHandler) healthCheck(w http.ResponseWriter, r *http.Request) {
	respondJSON(w, http.StatusOK, map[string]string{"status": "healthy"})
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

// respondJSON serializa a resposta em JSON e define o Content-Type.
func respondJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}
