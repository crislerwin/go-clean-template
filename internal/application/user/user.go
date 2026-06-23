package user

import (
	"context"
	"errors"

	"github.com/crislerwin/go-clean-template/internal/domain/user"
	inputports "github.com/crislerwin/go-clean-template/internal/ports/input"
	"github.com/crislerwin/go-clean-template/internal/ports/output"
	"github.com/crislerwin/go-clean-template/internal/ports/telemetry"
)

// createUserUseCase orquestra a criação de um usuário.
type createUserUseCase struct {
	repo   output.UserRepository
	logger telemetry.Logger
}

// NewCreateUserUseCase cria o caso de uso com as dependências necessárias.
func NewCreateUserUseCase(repo output.UserRepository, logger telemetry.Logger) inputports.CreateUserUseCase {
	return &createUserUseCase{
		repo:   repo,
		logger: logger.With("usecase", "CreateUser"),
	}
}

func (uc *createUserUseCase) Execute(ctx context.Context, input inputports.CreateUserInput) (*inputports.CreateUserOutput, error) {
	uc.logger.Info(ctx, "executing create user use case", "email", input.Email)

	u, err := user.NewUser(input.Name, input.Email)
	if err != nil {
		uc.logger.Error(ctx, "failed to build user entity", "error", err)
		return nil, err
	}

	if err := uc.repo.Save(ctx, u); err != nil {
		uc.logger.Error(ctx, "failed to save user", "error", err)
		return nil, err
	}

	uc.logger.Info(ctx, "user created successfully", "user_id", u.ID)
	return &inputports.CreateUserOutput{User: u}, nil
}

// findUserByIDUseCase busca um usuário pelo ID.
type findUserByIDUseCase struct {
	repo   output.UserRepository
	logger telemetry.Logger
}

// NewFindUserByIDUseCase cria o caso de uso de busca.
func NewFindUserByIDUseCase(repo output.UserRepository, logger telemetry.Logger) inputports.FindUserByIDUseCase {
	return &findUserByIDUseCase{
		repo:   repo,
		logger: logger.With("usecase", "FindUserByID"),
	}
}

func (uc *findUserByIDUseCase) Execute(ctx context.Context, input inputports.FindUserByIDInput) (*inputports.FindUserByIDOutput, error) {
	uc.logger.Info(ctx, "executing find user by id use case", "user_id", input.ID)

	u, err := uc.repo.FindByID(ctx, input.ID)
	if err != nil {
		if errors.Is(err, user.ErrUserNotFound) {
			uc.logger.Warn(ctx, "user not found", "user_id", input.ID)
		} else {
			uc.logger.Error(ctx, "failed to find user", "error", err)
		}
		return nil, err
	}

	uc.logger.Info(ctx, "user found", "user_id", u.ID)
	return &inputports.FindUserByIDOutput{User: u}, nil
}

// listUsersUseCase lista todos os usuários.
type listUsersUseCase struct {
	repo   output.UserRepository
	logger telemetry.Logger
}

// NewListUsersUseCase cria o caso de uso de listagem.
func NewListUsersUseCase(repo output.UserRepository, logger telemetry.Logger) inputports.ListUsersUseCase {
	return &listUsersUseCase{
		repo:   repo,
		logger: logger.With("usecase", "ListUsers"),
	}
}

func (uc *listUsersUseCase) Execute(ctx context.Context) (*inputports.ListUsersOutput, error) {
	uc.logger.Info(ctx, "executing list users use case")

	users, err := uc.repo.FindAll(ctx)
	if err != nil {
		uc.logger.Error(ctx, "failed to list users", "error", err)
		return nil, err
	}

	uc.logger.Info(ctx, "users listed", "count", len(users))
	return &inputports.ListUsersOutput{Users: users}, nil
}
