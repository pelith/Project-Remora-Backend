package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/joho/godotenv"
	"github.com/robfig/cron/v3"

	"remora/internal/agent"
	"remora/internal/signer"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: true,
	}))
	slog.SetDefault(logger)

	// Load .env file
	if err := godotenv.Load(); err != nil {
		logger.Warn("no .env file found, using environment variables")
	}

	// Initialize signer
	sgn, err := signer.NewFromEnv()
	if err != nil {
		logger.Error("failed to create signer", slog.Any("error", err))
		os.Exit(1)
	}

	logger.Info("signer initialized", slog.String("address", sgn.Address().Hex()))

	// Initialize eth client
	rpcURL := os.Getenv("RPC_URL")
	if rpcURL == "" {
		logger.Error("RPC_URL not set")
		os.Exit(1)
	}

	ethClient, err := ethclient.Dial(rpcURL)
	if err != nil {
		logger.Error("failed to connect to RPC", slog.Any("error", err))
		os.Exit(1)
	}

	defer ethClient.Close()

	logger.Info("connected to RPC", slog.String("url", rpcURL))

	// Initialize vault source (mock for now)
	// TODO: Replace with real VaultFactory implementation
	vaultSource := agent.NewMockVaultSource([]common.Address{
		// Add test vault addresses here
	})

	// Initialize agent service
	agentSvc := agent.New(
		vaultSource,
		nil, // TODO: strategySvc
		sgn,
		ethClient,
		logger,
	)

	ctx, cancel := context.WithCancel(context.Background())
	exitCode := 0

	defer func() {
		cancel()

		if exitCode != 0 {
			os.Exit(exitCode)
		}
	}()

	// Setup cron scheduler
	schedule := parseRebalanceSchedule()
	c := cron.New()

	_, err = c.AddFunc(schedule, func() {
		runAgent(ctx, agentSvc, logger)
	})
	if err != nil {
		logger.Error("invalid cron schedule", slog.Any("error", err))

		exitCode = 1

		return
	}

	c.Start()
	logger.Info("agent started", slog.String("schedule", schedule))

	// Run immediately on startup
	runAgent(ctx, agentSvc, logger)

	// Handle shutdown
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM)

	<-interrupt
	logger.Info("shutting down...")
	c.Stop()
	cancel()
}

func runAgent(ctx context.Context, agentSvc *agent.Service, logger *slog.Logger) {
	logger.InfoContext(ctx, "running rebalance check")

	results, err := agentSvc.Run(ctx)
	if err != nil {
		logger.ErrorContext(ctx, "rebalance run failed", slog.Any("error", err))
		return
	}

	for _, r := range results {
		logger.InfoContext(ctx, "vault processed",
			slog.String("address", r.VaultAddress.Hex()),
			slog.Bool("rebalanced", r.Rebalanced),
			slog.String("reason", r.Reason),
		)
	}

	logger.InfoContext(ctx, "rebalance check completed", slog.Int("vaults", len(results)))
}

func parseRebalanceSchedule() string {
	schedule := os.Getenv("REBALANCE_SCHEDULE")
	if schedule == "" {
		return "*/5 * * * *" // default: every 5 minutes
	}

	return schedule
}
