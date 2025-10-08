package app

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	api "github.com/flexer2006/t-t-ogen-go/generated"
	clientadapter "github.com/flexer2006/t-t-ogen-go/internal/adapters/client"
	"github.com/flexer2006/t-t-ogen-go/internal/adapters/data"
	serveradapter "github.com/flexer2006/t-t-ogen-go/internal/adapters/server"
)

const (
	defaultAddress          = ":42873"
	shutdownTimeout         = 5 * time.Second
	serverReadTimeout       = 15 * time.Second
	serverWriteTimeout      = 15 * time.Second
	serverIdleTimeout       = 60 * time.Second
	serverReadHeaderTimeout = 5 * time.Second
	serverMaxHeaderBytes    = http.DefaultMaxHeaderBytes
)

type Application struct {
	server *http.Server
}

func NewApplication(addr string) (*Application, error) {
	if addr == "" {
		addr = defaultAddress
	}

	storage := data.NewInMemoryUserStorage()

	service, err := newUserService(storage)
	if err != nil {
		return nil, fmt.Errorf("NewApplication: %w", err)
	}

	handler, err := serveradapter.NewUserHandler(service)
	if err != nil {
		return nil, fmt.Errorf("NewApplication: %w", err)
	}

	httpHandler, err := api.NewServer(handler)
	if err != nil {
		return nil, fmt.Errorf("NewApplication: %w", err)
	}

	server := new(http.Server)
	server.Addr = addr
	server.Handler = httpHandler
	server.ReadTimeout = serverReadTimeout
	server.WriteTimeout = serverWriteTimeout
	server.IdleTimeout = serverIdleTimeout
	server.ReadHeaderTimeout = serverReadHeaderTimeout
	server.MaxHeaderBytes = serverMaxHeaderBytes

	return &Application{server: server}, nil
}

func NewClient(baseURL string) (*clientadapter.Client, error) {
	invoker, err := api.NewClient(baseURL)
	if err != nil {
		return nil, fmt.Errorf("NewClient: %w", err)
	}

	client, err := clientadapter.New(invoker)
	if err != nil {
		return nil, fmt.Errorf("NewClient: %w", err)
	}

	return client, nil
}

func (a *Application) Run(ctx context.Context) error {
	serverErrors := make(chan error, 1)

	go func() {
		err := a.server.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			serverErrors <- fmt.Errorf("listen: %w", err)

			return
		}

		serverErrors <- err
	}()

	select {
	case err := <-serverErrors:
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			return err
		}

		return nil
	case <-ctx.Done():
		shutdownCtx, cancel := context.WithTimeout(context.WithoutCancel(ctx), shutdownTimeout)
		defer cancel()

		shutdownErr := a.server.Shutdown(shutdownCtx)
		if shutdownErr != nil {
			return fmt.Errorf("shutdown: %w", shutdownErr)
		}

		serverErr := <-serverErrors
		if serverErr != nil && !errors.Is(serverErr, http.ErrServerClosed) {
			return serverErr
		}

		return nil
	}
}
