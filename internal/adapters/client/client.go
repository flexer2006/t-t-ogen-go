// Package client provides the API client adapter.
package client

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"

	api "github.com/flexer2006/t-t-ogen-go/generated"
	"github.com/flexer2006/t-t-ogen-go/internal/domain"
	"github.com/flexer2006/t-t-ogen-go/internal/ports"
)

var ErrNilInvoker = errors.New("nil invoker")

var errUnexpectedResponse = errors.New("unexpected response type")

type Client struct {
	invoker api.Invoker
}

var _ ports.UserService = (*Client)(nil)

func New(invoker api.Invoker) (*Client, error) {
	if invoker == nil {
		return nil, ErrNilInvoker
	}

	return &Client{invoker: invoker}, nil
}

func (c *Client) ListUsers(ctx context.Context) ([]domain.User, error) {
	users, err := c.invoker.ListUsers(ctx)
	if err != nil {
		return nil, fmt.Errorf("Client.ListUsers: %w", err)
	}

	result := make([]domain.User, len(users))

	for i, user := range users {
		result[i] = toDomainUser(user)
	}

	return result, nil
}

func (c *Client) CreateUser(ctx context.Context, name, username string) (domain.User, error) {
	resp, err := c.invoker.CreateUser(ctx, &api.NewUser{
		Name:     name,
		Username: username,
	})
	if err != nil {
		return domain.User{}, fmt.Errorf("Client.CreateUser: %w", err)
	}

	return toDomainUser(*resp), nil
}

func (c *Client) GetUser(ctx context.Context, id uuid.UUID) (domain.User, error) {
	resp, err := c.invoker.GetUser(ctx, api.GetUserParams{ID: id})
	if err != nil {
		return domain.User{}, fmt.Errorf("Client.GetUser: %w", err)
	}

	switch result := resp.(type) {
	case *api.User:
		return toDomainUser(*result), nil
	case *api.GetUserNotFound:
		return domain.User{}, ports.ErrUserNotFound
	default:
		return domain.User{}, fmt.Errorf("Client.GetUser: %w (%T)", errUnexpectedResponse, result)
	}
}

func (c *Client) UpdateUser(
	ctx context.Context,
	userID uuid.UUID,
	name *string,
	username *string,
) (domain.User, error) {
	var payload api.UpdateUser

	if name != nil {
		payload.Name.SetTo(*name)
	}

	if username != nil {
		payload.Username.SetTo(*username)
	}

	resp, err := c.invoker.UpdateUser(ctx, &payload, api.UpdateUserParams{ID: userID})
	if err != nil {
		return domain.User{}, fmt.Errorf("Client.UpdateUser: %w", err)
	}

	switch result := resp.(type) {
	case *api.User:
		return toDomainUser(*result), nil
	case *api.UpdateUserNotFound:
		return domain.User{}, ports.ErrUserNotFound
	default:
		return domain.User{}, fmt.Errorf("Client.UpdateUser: %w (%T)", errUnexpectedResponse, result)
	}
}

func (c *Client) DeleteUser(ctx context.Context, id uuid.UUID) error {
	resp, err := c.invoker.DeleteUser(ctx, api.DeleteUserParams{ID: id})
	if err != nil {
		return fmt.Errorf("Client.DeleteUser: %w", err)
	}

	switch result := resp.(type) {
	case *api.DeleteUserNoContent:
		return nil
	case *api.DeleteUserNotFound:
		return ports.ErrUserNotFound
	default:
		return fmt.Errorf("Client.DeleteUser: %w (%T)", errUnexpectedResponse, result)
	}
}

func toDomainUser(user api.User) domain.User {
	return domain.User{
		ID:       user.GetID(),
		Name:     user.GetName(),
		Username: user.GetUsername(),
	}
}
