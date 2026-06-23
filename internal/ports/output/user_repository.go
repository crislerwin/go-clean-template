package output

import "github.com/crislerwin/go-clean-template/internal/domain/user"

// UserRepository é a porta de saída do domínio.
// A aplicação diz "preciso de alguém que consiga salvar e recuperar usuários",
// mas não se importa se é PostgreSQL, MongoDB ou memória.
// Isso mantém o domínio livre de detalhes de infraestrutura.
type UserRepository interface {
	Save(u *user.User) error
	FindByID(id string) (*user.User, error)
	FindAll() ([]*user.User, error)
}
