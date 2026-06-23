package input

import (
	"context"

	"github.com/crislerwin/go-clean-template/internal/domain/user"
)

type CreateUserInput struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

type CreateUserOutput struct {
	User *user.User `json:"user"`
}

type CreateUserUseCase interface {
	Execute(ctx context.Context, input CreateUserInput) (*CreateUserOutput, error)
}

type FindUserByIDInput struct {
	ID string `json:"id"`
}

type FindUserByIDOutput struct {
	User *user.User `json:"user"`
}

type FindUserByIDUseCase interface {
	Execute(ctx context.Context, input FindUserByIDInput) (*FindUserByIDOutput, error)
}

type ListUsersOutput struct {
	Users []*user.User `json:"users"`
}

type ListUsersUseCase interface {
	Execute(ctx context.Context) (*ListUsersOutput, error)
}
