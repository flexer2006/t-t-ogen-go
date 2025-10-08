package data

import (
	"context"
	"sync"

	xerrors "github.com/go-faster/errors"
	"github.com/google/uuid"

	"github.com/flexer2006/t-t-ogen-go/internal/domain"
	"github.com/flexer2006/t-t-ogen-go/internal/ports"
)

type InMemoryUserStorage struct {
	mu    sync.RWMutex
	users map[uuid.UUID]domain.User
}

var _ ports.UserRepository = (*InMemoryUserStorage)(nil)

func NewInMemoryUserStorage() *InMemoryUserStorage {
	return &InMemoryUserStorage{
		mu:    sync.RWMutex{},
		users: make(map[uuid.UUID]domain.User),
	}
}

func (s *InMemoryUserStorage) ListUsers(ctx context.Context) ([]domain.User, error) {
	if err := ctx.Err(); err != nil {
		return nil, xerrors.Wrap(err, "data.InMemoryUserStorage.ListUsers")
	}

	s.mu.RLock()
	defer s.mu.RUnlock()

	result := make([]domain.User, 0, len(s.users))

	for _, user := range s.users {
		result = append(result, cloneUser(user))
	}

	return result, nil
}

func (s *InMemoryUserStorage) CreateUser(ctx context.Context, name, username string) (domain.User, error) {
	if err := ctx.Err(); err != nil {
		return domain.User{}, xerrors.Wrap(err, "data.InMemoryUserStorage.CreateUser")
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	user := domain.User{
		ID:       uuid.New(),
		Name:     name,
		Username: username,
	}

	s.users[user.ID] = user

	return cloneUser(user), nil
}

func (s *InMemoryUserStorage) GetUser(ctx context.Context, userID uuid.UUID) (domain.User, error) {
	if err := ctx.Err(); err != nil {
		return domain.User{}, xerrors.Wrap(err, "data.InMemoryUserStorage.GetUser")
	}

	s.mu.RLock()
	defer s.mu.RUnlock()

	user, ok := s.users[userID]
	if !ok {
		return domain.User{}, ports.ErrUserNotFound
	}

	return cloneUser(user), nil
}

func (s *InMemoryUserStorage) UpdateUser(ctx context.Context, userID uuid.UUID, name *string, username *string) (domain.User, error) {
	if err := ctx.Err(); err != nil {
		return domain.User{}, xerrors.Wrap(err, "data.InMemoryUserStorage.UpdateUser")
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

	return cloneUser(user), nil
}

func (s *InMemoryUserStorage) DeleteUser(ctx context.Context, userID uuid.UUID) error {
	if err := ctx.Err(); err != nil {
		return xerrors.Wrap(err, "data.InMemoryUserStorage.DeleteUser")
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.users[userID]; !ok {
		return ports.ErrUserNotFound
	}

	delete(s.users, userID)

	return nil
}

func cloneUser(user domain.User) domain.User {
	return domain.User{
		ID:       user.ID,
		Name:     user.Name,
		Username: user.Username,
	}
}
