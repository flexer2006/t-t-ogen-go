package app

import (
	"context"
	"errors"
	"net/http"
	"strings"
	"time"

	xerrors "github.com/go-faster/errors"

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
	server  *http.Server
	baseURL string
}

func NewApplication(addr string) (*Application, error) {
	if addr == "" {
		addr = defaultAddress
	}

	repo := data.NewInMemoryUserStorage()

	service, err := newUserService(repo)
	if err != nil {
		return nil, xerrors.Wrap(err, "app.NewApplication: user service")
	}

	handler, err := serveradapter.NewUserHandler(service)
	if err != nil {
		return nil, xerrors.Wrap(err, "app.NewApplication: handler")
	}

	httpHandler, err := api.NewServer(handler)
	if err != nil {
		return nil, xerrors.Wrap(err, "app.NewApplication: http server")
	}

	server := &http.Server{
		Addr:              addr,
		Handler:           httpHandler,
		ReadTimeout:       serverReadTimeout,
		WriteTimeout:      serverWriteTimeout,
		IdleTimeout:       serverIdleTimeout,
		ReadHeaderTimeout: serverReadHeaderTimeout,
		MaxHeaderBytes:    serverMaxHeaderBytes,
	}

	return &Application{
		server:  server,
		baseURL: inferBaseURL(server.Addr),
	}, nil
}

func NewClient(baseURL string) (*clientadapter.Client, error) {
	invoker, err := api.NewClient(baseURL)
	if err != nil {
		return nil, xerrors.Wrap(err, "app.NewClient: invoker")
	}

	client, err := clientadapter.New(invoker)
	if err != nil {
		return nil, xerrors.Wrap(err, "app.NewClient: adapter")
	}

	return client, nil
}

func (a *Application) Client() (*clientadapter.Client, error) {
	return NewClient(a.baseURL)
}

func (a *Application) Run(ctx context.Context) error {
	serverErrors := make(chan error, 1)

	go func() {
		if err := a.server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			serverErrors <- xerrors.Wrap(err, "app.Application.Run: listen")
			return
		}

		serverErrors <- nil
	}()

	select {
	case err := <-serverErrors:
		return err
	case <-ctx.Done():
		shutdownCtx, cancel := context.WithTimeout(context.WithoutCancel(ctx), shutdownTimeout)
		defer cancel()

		if err := a.server.Shutdown(shutdownCtx); err != nil {
			return xerrors.Wrap(err, "app.Application.Run: shutdown")
		}

		return <-serverErrors
	}
}

func inferBaseURL(addr string) string {
	if addr == "" {
		return ""
	}

	if strings.HasPrefix(addr, ":") {
		return "http://localhost" + addr
	}

	if strings.HasPrefix(addr, "http://") || strings.HasPrefix(addr, "https://") {
		return addr
	}

	return "http://" + addr
}
