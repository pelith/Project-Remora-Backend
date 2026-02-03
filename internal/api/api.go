package api

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"remora/internal/config/api"
	"remora/internal/db"
	"remora/internal/liquidity"
	liquidityrepo "remora/internal/liquidity/repository"
	liquiditysvc "remora/internal/liquidity/service"
	"remora/internal/user"
	"remora/internal/user/repository"
	"remora/internal/user/service"
)

type Server struct {
	config        *api.Config
	httpServer    *http.Server
	pool          *pgxpool.Pool
	liquidityRepo *liquidityrepo.Repository
}

type Service struct {
	UserSvc      user.Service
	LiquiditySvc liquidity.Service
}

func NewServer(ctx context.Context, cfg *api.Config) (*Server, error) {
	pool, err := newPgxPool(ctx, cfg.PostgreSQL)
	if err != nil {
		return nil, fmt.Errorf("connect database: %w", err)
	}

	queries := db.New(pool)

	userSvc := service.New(repository.New(queries))

	var liquidityRepo *liquidityrepo.Repository

	if cfg.Ethereum.UseMock {
		slog.InfoContext(ctx, "using mock liquidity repository") //nolint:sloglint // startup config logging, no logger instance available at this scope

		liquidityRepo = liquidityrepo.NewMock()
	} else {
		var err error

		liquidityRepo, err = liquidityrepo.New(liquidityrepo.Config{
			RPCURL:          cfg.Ethereum.RPCURL,
			ContractAddress: cfg.Ethereum.StateViewContractAddr,
		})
		if err != nil {
			pool.Close()

			return nil, fmt.Errorf("create liquidity repository: %w", err)
		}
	}

	liquiditySvc := liquiditysvc.New(liquidityRepo)

	r := chi.NewRouter()
	AddRoutes(r, cfg, userSvc, liquiditySvc)

	return &Server{
		config: cfg,
		httpServer: &http.Server{
			Addr:         cfg.HTTP.Addr,
			ReadTimeout:  cfg.HTTP.ReadTimeout,
			WriteTimeout: cfg.HTTP.WriteTimeout,
			Handler:      r,
		},
		pool:          pool,
		liquidityRepo: liquidityRepo,
	}, nil
}

func (s *Server) Start() func(context.Context) error {
	go func() {
		slog.Info("starting http server", slog.String("addr", s.httpServer.Addr))

		if err := s.httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			slog.Error("start http server failed", slog.Any("error", err))
			os.Exit(1)
		}
	}()

	return func(ctx context.Context) error {
		if s.liquidityRepo != nil {
			s.liquidityRepo.Close()
		}

		if s.pool != nil {
			s.pool.Close()
		}

		return s.httpServer.Shutdown(ctx)
	}
}

func newPgxPool(ctx context.Context, pg api.PostgreSQL) (*pgxpool.Pool, error) {
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
