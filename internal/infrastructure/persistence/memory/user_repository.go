package memory

import (
	"sync"

	"github.com/crislerwin/go-clean-template/internal/domain/user"
	"github.com/crislerwin/go-clean-template/internal/ports/output"
)

// userRepository é um adapter de repositório em memória.
// Ele implementa a porta output.UserRepository sem expor detalhes de
// implementação. A escolha de memória permite testes rápidos e um
// bootstrap agnóstico de banco de dados.
type userRepository struct {
	mu     sync.RWMutex
	users  map[string]*user.User
}

// NewUserRepository cria um repositório em memória zerado.
// Pode ser usado em testes, demos e ambientes que não exigem persistência real.
func NewUserRepository() output.UserRepository {
	return &userRepository{
		users: make(map[string]*user.User),
	}
}

func (r *userRepository) Save(u *user.User) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	// Salvamos uma cópia rasa para evitar mutações externas no ponteiro armazenado.
	stored := *u
	r.users[u.ID] = &stored
	return nil
}

func (r *userRepository) FindByID(id string) (*user.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	u, ok := r.users[id]
	if !ok {
		return nil, user.ErrUserNotFound
	}
	return u, nil
}

func (r *userRepository) FindAll() ([]*user.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	result := make([]*user.User, 0, len(r.users))
	for _, u := range r.users {
		result = append(result, u)
	}
	return result, nil
}
