package output

import (
	"context"

	"github.com/crislerwin/go-clean-template/internal/domain/user"
)

// UserRepository define o contrato de persistência da aplicação.
// Qualquer adapter de banco de dados (memory, postgres, mongo) deve implementá-lo.
// A aplicação depende dessa interface, não de tecnologias concretas.
type UserRepository interface {
	Save(ctx context.Context, user *user.User) error
	FindByID(ctx context.Context, id string) (*user.User, error)
	FindAll(ctx context.Context) ([]*user.User, error)
}
