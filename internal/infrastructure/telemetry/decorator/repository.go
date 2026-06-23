package decorator

import (
	"context"

	"github.com/crislerwin/go-clean-template/internal/domain/user"
	"github.com/crislerwin/go-clean-template/internal/ports/output"
	"github.com/crislerwin/go-clean-template/internal/ports/telemetry"
)

// tracedRepository decora um UserRepository para criar spans de tracing.
// É um adapter que implementa output.UserRepository e demonstra como
// infraestrutura transversal (observabilidade) pode ser adicionada sem
// alterar domínio ou casos de uso.
type tracedRepository struct {
	repo   output.UserRepository
	tracer telemetry.Tracer
}

func NewTracedRepository(repo output.UserRepository, tracer telemetry.Tracer) output.UserRepository {
	return &tracedRepository{repo: repo, tracer: tracer}
}

func (r *tracedRepository) Save(u *user.User) error {
	_, span := r.tracer.Start(context.Background(), "UserRepository.Save")
	defer span.End()

	if err := r.repo.Save(u); err != nil {
		span.RecordError(err)
		return err
	}
	return nil
}

func (r *tracedRepository) FindByID(id string) (*user.User, error) {
	_, span := r.tracer.Start(context.Background(), "UserRepository.FindByID")
	defer span.End()

	u, err := r.repo.FindByID(id)
	if err != nil {
		span.RecordError(err)
		return nil, err
	}
	return u, nil
}

func (r *tracedRepository) FindAll() ([]*user.User, error) {
	_, span := r.tracer.Start(context.Background(), "UserRepository.FindAll")
	defer span.End()

	users, err := r.repo.FindAll()
	if err != nil {
		span.RecordError(err)
		return nil, err
	}
	return users, nil
}
