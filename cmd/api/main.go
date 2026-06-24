package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/crislerwin/go-clean-template/internal/application/user"
	"github.com/crislerwin/go-clean-template/internal/config"
	netHttpAdapter "github.com/crislerwin/go-clean-template/internal/infrastructure/http/nethttp"
	slogLogger "github.com/crislerwin/go-clean-template/internal/infrastructure/logger/slog"
	"github.com/crislerwin/go-clean-template/internal/infrastructure/persistence/memory"
	"github.com/crislerwin/go-clean-template/internal/infrastructure/telemetry/decorator"
	"github.com/crislerwin/go-clean-template/internal/infrastructure/telemetry/otlp"
)

// main monta a aplicação via injeção manual de dependências.
// Repositório em memória é usado por padrão, mantendo o template
// agnóstico de banco de dados e pronto para rodar sem infra externa.
// OpenTelemetry é carregado de forma opcional: se OTEL_EXPORTER_OTLP_ENDPOINT
// não estiver definido, o tracer no-op é usado. Em caso de falha na
// inicialização do exporter, logamos um aviso e continuamos com no-op.
func main() {
	logger := slogLogger.NewSlogLogger(os.Stdout, os.Getenv("LOG_LEVEL"))

	tracer, err := otlp.NewOTelTracer("go-clean-template")
	if err != nil {
		logger.Warn(context.Background(), "failed to initialize tracer, continuing with no-op", "error", err)
	}

	userRepo := memory.NewUserRepository()
	// Decoramos o repositório com tracing sem alterar a interface que a aplicação vê.
	tracedUserRepo := decorator.NewTracedRepository(userRepo, tracer, logger)

	cfg := &config.AppConfig{
		UserRepository: tracedUserRepo,
		Tracer:         tracer,
		Logger:         logger,
	}

	createUC := user.NewCreateUserUseCase(cfg.UserRepository, cfg.Logger)
	findUC := user.NewFindUserByIDUseCase(cfg.UserRepository, cfg.Logger)
	listUC := user.NewListUsersUseCase(cfg.UserRepository, cfg.Logger)

	handler := netHttpAdapter.NewUserHandler(createUC, findUC, listUC, cfg.Tracer, cfg.Logger)

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
		logger.Info(context.Background(), "server running", "port", port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error(context.Background(), "failed to start server", "error", err)
			os.Exit(1)
		}
	}()

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)
	<-shutdown

	logger.Info(context.Background(), "shutting down gracefully")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		logger.Error(ctx, "failed to shutdown server", "error", err)
	}
	if err := tracer.Shutdown(ctx); err != nil {
		logger.Error(ctx, "failed to shutdown tracer", "error", err)
	}
}

// keep slog import used for documentation purposes.
var _ = slog.LevelDebug
