package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/joho/godotenv"
	"github.com/robfig/cron/v3"

	"remora/internal/agent"
	liquidityrepo "remora/internal/liquidity/repository"
	liquidityservice "remora/internal/liquidity/service"
	"remora/internal/signer"
	strategyservice "remora/internal/strategy/service"
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

	// Initialize vault source from factory
	factoryAddr := os.Getenv("FACTORY_ADDRESS")
	if factoryAddr == "" {
		logger.Error("FACTORY_ADDRESS not set")
		os.Exit(1)
	}

	vaultSource := agent.NewFactoryVaultSource(
		ethClient,
		common.HexToAddress(factoryAddr),
	)

	// Initialize liquidity repository + strategy service
	stateViewAddr := os.Getenv("STATEVIEW_CONTRACT_ADDR")
	if stateViewAddr == "" {
		logger.Error("STATEVIEW_CONTRACT_ADDR not set")
		os.Exit(1)
	}

	liqRepo, err := liquidityrepo.New(liquidityrepo.Config{
		RPCURL:          rpcURL,
		ContractAddress: stateViewAddr,
	})
	if err != nil {
		logger.Error("failed to init liquidity repository", slog.Any("error", err))
		os.Exit(1)
	}
	defer liqRepo.Close()

	liqSvc := liquidityservice.New(liqRepo)
	strategySvc := strategyservice.New(liqSvc)

	// Initialize agent service
	agentSvc := agent.New(
		vaultSource,
		strategySvc,
		sgn,
		ethClient,
		logger,
		liqRepo,
	)

	// Load protection settings from env
	swapSlippage := os.Getenv("SWAP_SLIPPAGE_BPS")
	mintSlippage := os.Getenv("MINT_SLIPPAGE_BPS")
	maxGasPrice := os.Getenv("MAX_GAS_PRICE_GWEI")
	devThreshold := os.Getenv("DEVIATION_THRESHOLD")
	tickRangeEnv := os.Getenv("TICK_RANGE_AROUND_CURRENT")

	sSlippage := int64(50) // default 0.5%
	if swapSlippage != "" {
		if val, err := strconv.ParseInt(swapSlippage, 10, 64); err == nil {
			sSlippage = val
		}
	}

	mSlippage := int64(50) // default 0.5%
	if mintSlippage != "" {
		if val, err := strconv.ParseInt(mintSlippage, 10, 64); err == nil {
			mSlippage = val
		}
	}

	mGasPrice := 1.0 // default 1.0 Gwei
	if maxGasPrice != "" {
		if val, err := strconv.ParseFloat(maxGasPrice, 64); err == nil {
			mGasPrice = val
		}
	}

	dThreshold := 0.1 // default 10%
	if devThreshold != "" {
		if val, err := strconv.ParseFloat(devThreshold, 64); err == nil {
			dThreshold = val
		}
	}

	if tickRangeEnv != "" {
		if val, err := strconv.ParseInt(tickRangeEnv, 10, 32); err == nil && val > 0 {
			agentSvc.SetTickRangeAroundCurrent(int32(val))
			logger.Info("tick range override set", slog.String("raw", tickRangeEnv), slog.Int64("value", val))
		} else {
			logger.Warn("invalid TICK_RANGE_AROUND_CURRENT", slog.String("raw", tickRangeEnv))
		}
	} else {
		logger.Info("no TICK_RANGE_AROUND_CURRENT provided")
	}

	agentSvc.SetProtectionSettings(sSlippage, mSlippage, mGasPrice)
	agentSvc.SetDeviationThreshold(dThreshold)

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
