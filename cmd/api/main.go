package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"remora/internal/api"
	"remora/internal/config"
	apiCfg "remora/internal/config/api"
	"remora/internal/db"
	"remora/internal/user"
	"remora/internal/user/repository"
	"remora/internal/user/service"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: true,
	}))

	ctx := context.Background()
	if err := run(ctx, logger); err != nil {
		logger.ErrorContext(ctx, "api failed", slog.Any("error", err))
		os.Exit(1)
	}
}

func run(ctx context.Context, logger *slog.Logger) error {
	env := os.Getenv("ENV")
	if env == "" {
		env = "local"
	}

	cfg, err := config.LoadFromDir[*apiCfg.Config](env, "./config/api")
	if err != nil {
		return fmt.Errorf("load config: %w", err)
	}

	slog.SetLogLoggerLevel(cfg.Log.Level)
	logger = logger.With(slog.String("env", cfg.Env))

	runCtx, cancel := context.WithCancel(ctx)
	defer cancel()

	pool, err := newPgxPool(runCtx, cfg.AppConfig.PostgreSQL)
	if err != nil {
		return fmt.Errorf("connect database: %w", err)
	}
	defer pool.Close()

	queries := db.New(pool)

	var userSvc user.Service = service.New(repository.New(queries))

	handler := api.NewMux(userSvc)

	srv := &http.Server{
		Addr:         cfg.AppConfig.HTTP.Addr,
		Handler:      handler,
		ReadTimeout:  cfg.AppConfig.HTTP.ReadTimeout,
		WriteTimeout: cfg.AppConfig.HTTP.WriteTimeout,
	}

	errCh := make(chan error, 1)

	go func() {
		logger.InfoContext(runCtx, "starting http server", slog.String("addr", srv.Addr))

		errCh <- srv.ListenAndServe()
	}()

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM)

	select {
	case <-interrupt:
		logger.InfoContext(runCtx, "stopping...")
	case err = <-errCh:
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.ErrorContext(runCtx, "http server failed", slog.Any("error", err))
		}
	}

	const shutdownTimeout = 10 * time.Second

	shutdownCtx, shutdownCancel := context.WithTimeout(runCtx, shutdownTimeout)
	defer shutdownCancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		logger.ErrorContext(runCtx, "shutdown failed", slog.Any("error", err))
	}

	return nil
}

func newPgxPool(ctx context.Context, pg apiCfg.PostgreSQL) (*pgxpool.Pool, error) {
	hostAndPort := net.JoinHostPort(pg.Host, pg.Port)
	connectURI := fmt.Sprintf(
		"postgres://%s:%s@%s/%s?sslmode=disable",
		pg.User,
		pg.Password,
		hostAndPort,
		pg.Database,
	)

	pool, err := pgxpool.New(ctx, connectURI)
	if err != nil {
		return nil, fmt.Errorf("pgxpool new: %w", err)
	}

	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("pgx ping: %w", err)
	}

	return pool, nil
}
