package ports

import (
	"context"

	"github.com/google/uuid"

	"github.com/flexer2006/t-t-ogen-go/internal/domain"
)

type UserRepository interface {
	ListUsers(ctx context.Context) ([]domain.User, error)
	CreateUser(ctx context.Context, name, username string) (domain.User, error)
	GetUser(ctx context.Context, id uuid.UUID) (domain.User, error)
	UpdateUser(ctx context.Context, id uuid.UUID, name *string, username *string) (domain.User, error)
	DeleteUser(ctx context.Context, id uuid.UUID) error
}
