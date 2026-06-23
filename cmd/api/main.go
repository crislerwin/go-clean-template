package main

import (
	"log"
	"net/http"
	"os"

	"github.com/crislerwin/go-clean-template/internal/application/user"
	"github.com/crislerwin/go-clean-template/internal/config"
	netHttpAdapter "github.com/crislerwin/go-clean-template/internal/infrastructure/http/nethttp"
	"github.com/crislerwin/go-clean-template/internal/infrastructure/persistence/memory"
)

// main monta a aplicação via injeção manual de dependências.
// Repositório em memória é usado por padrão, mantendo o template
// agnóstico de banco de dados e pronto para rodar sem infra externa.
func main() {
	userRepo := memory.NewUserRepository()
	cfg := &config.AppConfig{UserRepository: userRepo}

	createUC := user.NewCreateUserUseCase(cfg.UserRepository)
	findUC := user.NewFindUserByIDUseCase(cfg.UserRepository)
	listUC := user.NewListUsersUseCase(cfg.UserRepository)

	handler := netHttpAdapter.NewUserHandler(createUC, findUC, listUC)

	mux := http.NewServeMux()
	handler.RegisterRoutes(mux)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	server := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	log.Printf("Server running on :%s", port)
	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
