package user

import (
	"context"

	"github.com/crislerwin/go-clean-template/internal/domain/user"
	inputports "github.com/crislerwin/go-clean-template/internal/ports/input"
	"github.com/crislerwin/go-clean-template/internal/ports/output"
	"github.com/crislerwin/go-clean-template/internal/ports/telemetry"
)

// createUserUseCase orquestra a criação de um usuário.
// A lógica de negócio vive no domínio (NewUser); o caso de uso coordena
// persistência e mapeamento entre DTOs de entrada e entidades.
type createUserUseCase struct {
	repo   output.UserRepository
	tracer telemetry.Tracer
}

// NewCreateUserUseCase expõe a porta de entrada e esconde a implementação.
// O tipo concreto permanece não exportado para garantir acoplamento à interface.
func NewCreateUserUseCase(repo output.UserRepository, tracer telemetry.Tracer) inputports.CreateUserUseCase {
	return &createUserUseCase{repo: repo, tracer: tracer}
}

func (uc *createUserUseCase) Execute(input inputports.CreateUserInput) (inputports.CreateUserOutput, error) {
	_, span := uc.tracer.Start(context.Background(), "CreateUserUseCase.Execute")
	defer span.End()

	u, err := user.NewUser(input.Name, input.Email)
	if err != nil {
		span.RecordError(err)
		return inputports.CreateUserOutput{}, err
	}

	if err := uc.repo.Save(u); err != nil {
		span.RecordError(err)
		return inputports.CreateUserOutput{}, err
	}

	return inputports.CreateUserOutput{User: u}, nil
}

// findUserByIDUseCase orquestra a busca por ID.
type findUserByIDUseCase struct {
	repo   output.UserRepository
	tracer telemetry.Tracer
}

func NewFindUserByIDUseCase(repo output.UserRepository, tracer telemetry.Tracer) inputports.FindUserByIDUseCase {
	return &findUserByIDUseCase{repo: repo, tracer: tracer}
}

func (uc *findUserByIDUseCase) Execute(input inputports.FindUserByIDInput) (inputports.FindUserByIDOutput, error) {
	_, span := uc.tracer.Start(context.Background(), "FindUserByIDUseCase.Execute")
	defer span.End()

	u, err := uc.repo.FindByID(input.ID)
	if err != nil {
		span.RecordError(err)
		return inputports.FindUserByIDOutput{}, err
	}

	return inputports.FindUserByIDOutput{User: u}, nil
}

// listUsersUseCase orquestra a listagem.
type listUsersUseCase struct {
	repo   output.UserRepository
	tracer telemetry.Tracer
}

func NewListUsersUseCase(repo output.UserRepository, tracer telemetry.Tracer) inputports.ListUsersUseCase {
	return &listUsersUseCase{repo: repo, tracer: tracer}
}

func (uc *listUsersUseCase) Execute() (inputports.ListUsersOutput, error) {
	_, span := uc.tracer.Start(context.Background(), "ListUsersUseCase.Execute")
	defer span.End()

	users, err := uc.repo.FindAll()
	if err != nil {
		span.RecordError(err)
		return inputports.ListUsersOutput{}, err
	}

	return inputports.ListUsersOutput{Users: users}, nil
}
