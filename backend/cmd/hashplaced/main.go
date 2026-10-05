// Command hashplaced runs the hashplace HTTP server.
package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/nem0z/hashplace/backend/internal/api"
	"github.com/nem0z/hashplace/backend/internal/config"
)

const shutdownTimeout = 10 * time.Second

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	if err := run(logger); err != nil {
		logger.Error("server failed", "error", err)
		os.Exit(1)
	}
}

func run(logger *slog.Logger) error {
	cfg, err := config.Load(os.Getenv)
	if err != nil {
		return err
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	srv := newServer(cfg)
	errChan := start(srv)

	logger.Info("server started", "addr", cfg.Addr)

	select {
	case err := <-errChan:
		return err
	case <-ctx.Done():
	}

	logger.Info("server shutting down")

	return shutdown(srv, errChan)
}

func newServer(cfg config.Config) *http.Server {
	return &http.Server{
		Addr:         cfg.Addr,
		Handler:      api.NewHandler(),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}
}

// start runs srv in the background. The returned channel receives the error that stopped it.
func start(srv *http.Server) <-chan error {
	errChan := make(chan error, 1)
	go func() { errChan <- srv.ListenAndServe() }()

	return errChan
}

// shutdown stops srv gracefully and waits for it to return.
func shutdown(srv *http.Server, errChan <-chan error) error {
	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		return err
	}

	if err := <-errChan; !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	return nil
}
