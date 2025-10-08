// Package data provides storage adapters.
package data

import (
	"context"
	"fmt"
	"sync"

	"github.com/google/uuid"

	"github.com/flexer2006/t-t-ogen-go/internal/domain"
	"github.com/flexer2006/t-t-ogen-go/internal/ports"
)

type InMemoryUserStorage struct {
	mu    sync.RWMutex
	users map[uuid.UUID]domain.User
}

func NewInMemoryUserStorage() *InMemoryUserStorage {
	return &InMemoryUserStorage{
		mu:    sync.RWMutex{},
		users: make(map[uuid.UUID]domain.User),
	}
}

func (s *InMemoryUserStorage) ListUsers(ctx context.Context) ([]domain.User, error) {
	err := ctx.Err()
	if err != nil {
		return nil, fmt.Errorf("InMemoryUserStorage.ListUsers: %w", err)
	}

	s.mu.RLock()
	defer s.mu.RUnlock()

	result := make([]domain.User, 0, len(s.users))

	for _, user := range s.users {
		result = append(result, user)
	}

	return result, nil
}

func (s *InMemoryUserStorage) CreateUser(ctx context.Context, name, username string) (domain.User, error) {
	err := ctx.Err()
	if err != nil {
		return domain.User{}, fmt.Errorf("InMemoryUserStorage.CreateUser: %w", err)
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	user := domain.User{
		ID:       uuid.New(),
		Name:     name,
		Username: username,
	}

	s.users[user.ID] = user

	return user, nil
}

func (s *InMemoryUserStorage) GetUser(ctx context.Context, userID uuid.UUID) (domain.User, error) {
	err := ctx.Err()
	if err != nil {
		return domain.User{}, fmt.Errorf("InMemoryUserStorage.GetUser: %w", err)
	}

	s.mu.RLock()
	defer s.mu.RUnlock()

	user, ok := s.users[userID]
	if !ok {
		return domain.User{}, ports.ErrUserNotFound
	}

	return user, nil
}

func (s *InMemoryUserStorage) UpdateUser(
	ctx context.Context,
	userID uuid.UUID,
	name *string,
	username *string,
) (domain.User, error) {
	err := ctx.Err()
	if err != nil {
		return domain.User{}, fmt.Errorf("InMemoryUserStorage.UpdateUser: %w", err)
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	user, ok := s.users[userID]
	if !ok {
		return domain.User{}, ports.ErrUserNotFound
	}

	if name != nil {
		user.Name = *name
	}

	if username != nil {
		user.Username = *username
	}

	s.users[userID] = user

	return user, nil
}

func (s *InMemoryUserStorage) DeleteUser(ctx context.Context, userID uuid.UUID) error {
	err := ctx.Err()
	if err != nil {
		return fmt.Errorf("InMemoryUserStorage.DeleteUser: %w", err)
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.users[userID]; !ok {
		return ports.ErrUserNotFound
	}

	delete(s.users, userID)

	return nil
}
