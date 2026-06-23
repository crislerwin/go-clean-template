package input

import "github.com/crislerwin/go-clean-template/internal/domain/user"

// CreateUserInput é um port de entrada: a infraestrutura HTTP (ou CLI, ou gRPC)
// chama esse contrato sem saber como o caso de uso é implementado.
type CreateUserInput struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

type CreateUserOutput struct {
	User *user.User `json:"user"`
}

type CreateUserUseCase interface {
	Execute(input CreateUserInput) (CreateUserOutput, error)
}

type FindUserByIDInput struct {
	ID string `json:"id"`
}

type FindUserByIDOutput struct {
	User *user.User `json:"user"`
}

type FindUserByIDUseCase interface {
	Execute(input FindUserByIDInput) (FindUserByIDOutput, error)
}

type ListUsersOutput struct {
	Users []*user.User `json:"users"`
}

type ListUsersUseCase interface {
	Execute() (ListUsersOutput, error)
}
