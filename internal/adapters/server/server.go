// Package server provides HTTP handlers for the user service.
package server

import (
	"context"
	"errors"
	"fmt"

	api "github.com/flexer2006/t-t-ogen-go/generated"
	"github.com/flexer2006/t-t-ogen-go/internal/domain"
	"github.com/flexer2006/t-t-ogen-go/internal/ports"
)

var ErrNilUserService = errors.New("nil user service")

var errNilRequest = errors.New("nil request")

type UserHandler struct {
	service ports.UserService
}

var _ api.Handler = (*UserHandler)(nil)

func NewUserHandler(service ports.UserService) (*UserHandler, error) {
	if service == nil {
		return nil, ErrNilUserService
	}

	return &UserHandler{service: service}, nil
}

func (h *UserHandler) ListUsers(ctx context.Context) ([]api.User, error) {
	users, err := h.service.ListUsers(ctx)
	if err != nil {
		return nil, fmt.Errorf("UserHandler.ListUsers: %w", err)
	}

	result := make([]api.User, len(users))

	for i, user := range users {
		result[i] = toAPIUser(user)
	}

	return result, nil
}

func (h *UserHandler) CreateUser(ctx context.Context, req *api.NewUser) (*api.User, error) {
	if req == nil {
		return nil, errNilRequest
	}

	user, err := h.service.CreateUser(ctx, req.GetName(), req.GetUsername())
	if err != nil {
		return nil, fmt.Errorf("UserHandler.CreateUser: %w", err)
	}

	apiUser := toAPIUser(user)

	return &apiUser, nil
}

func (h *UserHandler) GetUser(ctx context.Context, params api.GetUserParams) (api.GetUserRes, error) {
	user, err := h.service.GetUser(ctx, params.ID)
	if err != nil {
		if errors.Is(err, ports.ErrUserNotFound) {
			return &api.GetUserNotFound{}, nil
		}

		return nil, fmt.Errorf("UserHandler.GetUser: %w", err)
	}

	apiUser := toAPIUser(user)

	return &apiUser, nil
}

func (h *UserHandler) UpdateUser(ctx context.Context,
	req *api.UpdateUser,
	params api.UpdateUserParams) (api.UpdateUserRes, error) {
	if req == nil {
		return nil, errNilRequest
	}

	name := optStringToPtr(req.GetName())
	username := optStringToPtr(req.GetUsername())

	updated, err := h.service.UpdateUser(ctx, params.ID, name, username)
	if err != nil {
		if errors.Is(err, ports.ErrUserNotFound) {
			return &api.UpdateUserNotFound{}, nil
		}

		return nil, fmt.Errorf("UserHandler.UpdateUser: %w", err)
	}

	apiUser := toAPIUser(updated)

	return &apiUser, nil
}

func (h *UserHandler) DeleteUser(ctx context.Context, params api.DeleteUserParams) (api.DeleteUserRes, error) {
	err := h.service.DeleteUser(ctx, params.ID)
	if err != nil {
		if errors.Is(err, ports.ErrUserNotFound) {
			return &api.DeleteUserNotFound{}, nil
		}

		return nil, fmt.Errorf("UserHandler.DeleteUser: %w", err)
	}

	return &api.DeleteUserNoContent{}, nil
}

func toAPIUser(user domain.User) api.User {
	return api.User{
		ID:       user.ID,
		Name:     user.Name,
		Username: user.Username,
	}
}

func optStringToPtr(opt api.OptString) *string {
	if value, ok := opt.Get(); ok {
		return &value
	}

	return nil
}
