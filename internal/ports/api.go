package ports

import (
	"context"
	"errors"

	"github.com/google/uuid"

	"github.com/flexer2006/t-t-ogen-go/internal/domain"
)

var ErrUserNotFound = errors.New("user not found")

type UserService interface {
	ListUsers(ctx context.Context) ([]domain.User, error)
	CreateUser(ctx context.Context, name, username string) (domain.User, error)
	GetUser(ctx context.Context, id uuid.UUID) (domain.User, error)
	UpdateUser(ctx context.Context, id uuid.UUID, name *string, username *string) (domain.User, error)
	DeleteUser(ctx context.Context, id uuid.UUID) error
}
