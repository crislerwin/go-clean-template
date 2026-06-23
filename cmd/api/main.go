package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/crislerwin/go-clean-template/internal/application/user"
	"github.com/crislerwin/go-clean-template/internal/config"
	netHttpAdapter "github.com/crislerwin/go-clean-template/internal/infrastructure/http/nethttp"
	"github.com/crislerwin/go-clean-template/internal/infrastructure/persistence/memory"
	"github.com/crislerwin/go-clean-template/internal/infrastructure/telemetry/decorator"
	"github.com/crislerwin/go-clean-template/internal/infrastructure/telemetry/otlp"
)

// main monta a aplicação via injeção manual de dependências.
// Repositório em memória é usado por padrão, mantendo o template
// agnóstico de banco de dados e pronto para rodar sem infra externa.
// OpenTelemetry é carregado de forma opcional: se OTEL_EXPORTER_OTLP_ENDPOINT
// não estiver definido, o tracer no-op é usado e a aplicação continua funcionando.
func main() {
	tracer, err := otlp.NewOTelTracer("go-clean-template")
	if err != nil {
		log.Fatalf("failed to initialize tracer: %v", err)
	}

	userRepo := memory.NewUserRepository()
	// Decoramos o repositório com tracing sem alterar a interface que a aplicação vê.
	tracedUserRepo := decorator.NewTracedRepository(userRepo, tracer)

	cfg := &config.AppConfig{
		UserRepository: tracedUserRepo,
		Tracer:         tracer,
	}

	createUC := user.NewCreateUserUseCase(cfg.UserRepository, cfg.Tracer)
	findUC := user.NewFindUserByIDUseCase(cfg.UserRepository, cfg.Tracer)
	listUC := user.NewListUsersUseCase(cfg.UserRepository, cfg.Tracer)

	handler := netHttpAdapter.NewUserHandler(createUC, findUC, listUC, cfg.Tracer)

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

	go func() {
		log.Printf("Server running on :%s", port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("failed to start server: %v", err)
		}
	}()

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)
	<-shutdown

	log.Println("shutting down gracefully...")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Printf("failed to shutdown server: %v", err)
	}
	if err := tracer.Shutdown(ctx); err != nil {
		log.Printf("failed to shutdown tracer: %v", err)
	}
}
