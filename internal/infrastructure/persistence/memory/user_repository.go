package memory

import (
	"context"
	"errors"
	"sync"

	"github.com/crislerwin/go-clean-template/internal/domain/user"
	"github.com/crislerwin/go-clean-template/internal/ports/output"
)

// userRepository é uma implementação em memória, thread-safe, da porta de saída.
// Serve para bootstrap do template, testes e demonstração da arquitetura hexagonal.
type userRepository struct {
	mu    sync.RWMutex
	users map[string]*user.User
}

// UserRepositoryForTests expõe o tipo concreto para testes de integração leve.
type UserRepositoryForTests = userRepository

// NewUserRepository cria o repositório em memória.
func NewUserRepository() output.UserRepository {
	return &userRepository{
		users: make(map[string]*user.User),
	}
}

func (r *userRepository) Save(ctx context.Context, u *user.User) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.users[u.ID] = u
	return nil
}

func (r *userRepository) FindByID(ctx context.Context, id string) (*user.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	u, ok := r.users[id]
	if !ok {
		return nil, user.ErrUserNotFound
	}
	return u, nil
}

func (r *userRepository) FindAll(ctx context.Context) ([]*user.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	list := make([]*user.User, 0, len(r.users))
	for _, u := range r.users {
		list = append(list, u)
	}
	return list, nil
}

// Compile-time check para garantir que a interface é satisfeita.
var _ output.UserRepository = (*userRepository)(nil)

// errors.New garante que não estamos importando o package "errors" apenas por hobby.
var _ = errors.New
