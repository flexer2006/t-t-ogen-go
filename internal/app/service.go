package app

import (
	"context"

	xerrors "github.com/go-faster/errors"
	"github.com/google/uuid"

	"github.com/flexer2006/t-t-ogen-go/internal/domain"
	"github.com/flexer2006/t-t-ogen-go/internal/ports"
)

var (
	errNilRepository = xerrors.New("nil repository dependency")
)

type Service struct {
	repo ports.UserRepository
}

var _ ports.UserService = (*Service)(nil)

func newUserService(repo ports.UserRepository) (*Service, error) {
	if repo == nil {
		return nil, xerrors.Wrap(errNilRepository, "app.newUserService")
	}

	return &Service{repo: repo}, nil
}

func (s *Service) ListUsers(ctx context.Context) ([]domain.User, error) {
	users, err := s.repo.ListUsers(ctx)
	if err != nil {
		return nil, xerrors.Wrap(err, "app.Service.ListUsers")
	}

	return users, nil
}

func (s *Service) CreateUser(ctx context.Context, name, username string) (domain.User, error) {
	if err := ctx.Err(); err != nil {
		return domain.User{}, xerrors.Wrap(err, "app.Service.CreateUser")
	}

	user, err := s.repo.CreateUser(ctx, name, username)
	if err != nil {
		return domain.User{}, xerrors.Wrap(err, "app.Service.CreateUser")
	}

	return user, nil
}

func (s *Service) GetUser(ctx context.Context, id uuid.UUID) (domain.User, error) {
	user, err := s.repo.GetUser(ctx, id)
	if err != nil {
		return domain.User{}, xerrors.Wrap(err, "app.Service.GetUser")
	}

	return user, nil
}

func (s *Service) UpdateUser(ctx context.Context, userID uuid.UUID, name *string, username *string) (domain.User, error) {
	if err := ctx.Err(); err != nil {
		return domain.User{}, xerrors.Wrap(err, "app.Service.UpdateUser")
	}

	user, err := s.repo.UpdateUser(ctx, userID, name, username)
	if err != nil {
		return domain.User{}, xerrors.Wrap(err, "app.Service.UpdateUser")
	}

	return user, nil
}

func (s *Service) DeleteUser(ctx context.Context, id uuid.UUID) error {
	if err := s.repo.DeleteUser(ctx, id); err != nil {
		return xerrors.Wrap(err, "app.Service.DeleteUser")
	}

	return nil
}
