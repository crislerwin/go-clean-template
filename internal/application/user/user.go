package user

import (
	"github.com/crislerwin/go-clean-template/internal/domain/user"
	inputports "github.com/crislerwin/go-clean-template/internal/ports/input"
	"github.com/crislerwin/go-clean-template/internal/ports/output"
)

// createUserUseCase orquestra a criação de um usuário.
// A lógica de negócio vive no domínio (NewUser); o caso de uso coordena
// persistência e mapeamento entre DTOs de entrada e entidades.
type createUserUseCase struct {
	repo output.UserRepository
}

// NewCreateUserUseCase expõe a porta de entrada e esconde a implementação.
// O tipo concreto permanece não exportado para garantir acoplamento à interface.
func NewCreateUserUseCase(repo output.UserRepository) inputports.CreateUserUseCase {
	return &createUserUseCase{repo: repo}
}

func (uc *createUserUseCase) Execute(input inputports.CreateUserInput) (inputports.CreateUserOutput, error) {
	u, err := user.NewUser(input.Name, input.Email)
	if err != nil {
		return inputports.CreateUserOutput{}, err
	}

	if err := uc.repo.Save(u); err != nil {
		return inputports.CreateUserOutput{}, err
	}

	return inputports.CreateUserOutput{User: u}, nil
}

// findUserByIDUseCase orquestra a busca por ID.
type findUserByIDUseCase struct {
	repo output.UserRepository
}

func NewFindUserByIDUseCase(repo output.UserRepository) inputports.FindUserByIDUseCase {
	return &findUserByIDUseCase{repo: repo}
}

func (uc *findUserByIDUseCase) Execute(input inputports.FindUserByIDInput) (inputports.FindUserByIDOutput, error) {
	u, err := uc.repo.FindByID(input.ID)
	if err != nil {
		return inputports.FindUserByIDOutput{}, err
	}

	return inputports.FindUserByIDOutput{User: u}, nil
}

// listUsersUseCase orquestra a listagem.
type listUsersUseCase struct {
	repo output.UserRepository
}

func NewListUsersUseCase(repo output.UserRepository) inputports.ListUsersUseCase {
	return &listUsersUseCase{repo: repo}
}

func (uc *listUsersUseCase) Execute() (inputports.ListUsersOutput, error) {
	users, err := uc.repo.FindAll()
	if err != nil {
		return inputports.ListUsersOutput{}, err
	}

	return inputports.ListUsersOutput{Users: users}, nil
}
