package decorator

import (
	"context"

	"github.com/crislerwin/go-clean-template/internal/domain/user"
	"github.com/crislerwin/go-clean-template/internal/ports/output"
	"github.com/crislerwin/go-clean-template/internal/ports/telemetry"
)

// tracedRepository decora um UserRepository para criar spans de tracing.
// É um adapter que adiciona observabilidade sem alterar a interface que a aplicação vê.
type tracedRepository struct {
	repo   output.UserRepository
	tracer telemetry.Tracer
	logger telemetry.Logger
}

// NewTracedRepository envolve um repositório com tracing.
func NewTracedRepository(repo output.UserRepository, tracer telemetry.Tracer, logger telemetry.Logger) output.UserRepository {
	return &tracedRepository{
		repo:   repo,
		tracer: tracer,
		logger: logger.With("component", "tracedRepository"),
	}
}

func (t *tracedRepository) Save(ctx context.Context, u *user.User) error {
	ctx, span := t.tracer.Start(ctx, "UserRepository.Save")
	defer span.End()

	t.logger.Info(ctx, "saving user", "user_id", u.ID)
	err := t.repo.Save(ctx, u)
	if err != nil {
		span.RecordError(err)
		t.logger.Error(ctx, "failed to save user", "error", err)
	}
	return err
}

func (t *tracedRepository) FindByID(ctx context.Context, id string) (*user.User, error) {
	ctx, span := t.tracer.Start(ctx, "UserRepository.FindByID")
	defer span.End()

	t.logger.Info(ctx, "finding user by id", "user_id", id)
	u, err := t.repo.FindByID(ctx, id)
	if err != nil {
		span.RecordError(err)
		t.logger.Error(ctx, "failed to find user by id", "error", err)
	}
	return u, err
}

func (t *tracedRepository) FindAll(ctx context.Context) ([]*user.User, error) {
	ctx, span := t.tracer.Start(ctx, "UserRepository.FindAll")
	defer span.End()

	t.logger.Info(ctx, "listing users")
	users, err := t.repo.FindAll(ctx)
	if err != nil {
		span.RecordError(err)
		t.logger.Error(ctx, "failed to list users", "error", err)
	}
	return users, err
}

var _ output.UserRepository = (*tracedRepository)(nil)
