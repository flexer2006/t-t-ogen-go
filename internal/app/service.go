// Package app provides the application logic and dependency injection.
package app

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"

	"github.com/flexer2006/t-t-ogen-go/internal/domain"
	"github.com/flexer2006/t-t-ogen-go/internal/ports"
)

var ErrNilStorage = errors.New("nil storage dependency")

type Service struct {
	storage ports.UserService
}

var _ ports.UserService = (*Service)(nil)

func newUserService(storage ports.UserService) (*Service, error) {
	if storage == nil {
		return nil, fmt.Errorf("NewUserService: %w", ErrNilStorage)
	}

	return &Service{storage: storage}, nil
}

func (s *Service) ListUsers(ctx context.Context) ([]domain.User, error) {
	users, err := s.storage.ListUsers(ctx)
	if err != nil {
		return nil, fmt.Errorf("ListUsers: %w", err)
	}

	return users, nil
}

func (s *Service) CreateUser(ctx context.Context, name, username string) (domain.User, error) {
	if ctx.Err() != nil {
		return domain.User{}, fmt.Errorf("CreateUser: %w", ctx.Err())
	}

	user, err := s.storage.CreateUser(ctx, name, username)
	if err != nil {
		return domain.User{}, fmt.Errorf("CreateUser: %w", err)
	}

	return user, nil
}

func (s *Service) GetUser(ctx context.Context, id uuid.UUID) (domain.User, error) {
	user, err := s.storage.GetUser(ctx, id)
	if err != nil {
		return domain.User{}, fmt.Errorf("GetUser: %w", err)
	}

	return user, nil
}

func (s *Service) UpdateUser(
	ctx context.Context,
	userID uuid.UUID,
	name *string,
	username *string,
) (domain.User, error) {
	if ctx.Err() != nil {
		return domain.User{}, fmt.Errorf("UpdateUser: %w", ctx.Err())
	}

	user, err := s.storage.UpdateUser(ctx, userID, name, username)
	if err != nil {
		return domain.User{}, fmt.Errorf("UpdateUser: %w", err)
	}

	return user, nil
}

func (s *Service) DeleteUser(ctx context.Context, id uuid.UUID) error {
	err := s.storage.DeleteUser(ctx, id)
	if err != nil {
		return fmt.Errorf("DeleteUser: %w", err)
	}

	return nil
}
