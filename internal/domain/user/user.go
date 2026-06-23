package user

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

// Domain errors são declarados no domínio para que qualquer adapter
// possa reconhecê-los sem depender de detalhes de infraestrutura.
var (
	ErrInvalidName  = errors.New("user name must not be empty")
	ErrInvalidEmail = errors.New("user email must not be empty")
	ErrUserNotFound = errors.New("user not found")
)

// User é a entidade de domínio. Ela não conhece HTTP, banco de dados
// ou qualquer framework. A "theory" aqui é: um usuário é identificado
// por um UUID imutável e possui dados que podem evoluir, mas sua
// identidade permanece a mesma.
type User struct {
	ID        string
	Name      string
	Email     string
	CreatedAt time.Time
	UpdatedAt time.Time
}

// NewUser cria uma nova instância válida de User.
// A criação de entidades é uma responsabilidade de domínio, não de um adapter.
func NewUser(name, email string) (*User, error) {
	if err := validateUserFields(name, email); err != nil {
		return nil, err
	}

	now := time.Now()
	return &User{
		ID:        uuid.NewString(),
		Name:      name,
		Email:     email,
		CreatedAt: now,
		UpdatedAt: now,
	}, nil
}

// UpdateName muda o nome do usuário, preservando a identidade.
// Pequenas funções com nomes que revelam intenção mantêm a teoria clara.
func (u *User) UpdateName(name string) error {
	if name == "" {
		return ErrInvalidName
	}
	u.Name = name
	u.UpdatedAt = time.Now()
	return nil
}

func validateUserFields(name, email string) error {
	if name == "" {
		return ErrInvalidName
	}
	if email == "" {
		return ErrInvalidEmail
	}
	return nil
}
