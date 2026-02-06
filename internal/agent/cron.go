package agent

import (
	"context"
	"log/slog"
	"os"
	"strconv"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/joho/godotenv"
	"github.com/robfig/cron/v3"

	liquidityrepo "remora/internal/liquidity/repository"
	liquidityservice "remora/internal/liquidity/service"
	"remora/internal/signer"
	strategyservice "remora/internal/strategy/service"
)

// StartCron starts the rebalance cron from environment variables.
// If useDefaultSchedule is false and REBALANCE_SCHEDULE is not set, returns (nil, nil) and no cron is run.
// If useDefaultSchedule is true (e.g. when running the rebalance binary), default schedule "*/5 * * * *" is used when unset.
// Call the returned stop function on shutdown to stop the cron and release resources.
func StartCron(ctx context.Context, logger *slog.Logger, useDefaultSchedule bool) (stop func(), err error) {
	_ = godotenv.Load()

	schedule := os.Getenv("REBALANCE_SCHEDULE")
	if schedule == "" {
		if !useDefaultSchedule {
			return func() {}, nil
		}

		schedule = "*/5 * * * *"
	}

	sgn, err := signer.NewFromEnv()
	if err != nil {
		return nil, err
	}

	logger.Info("signer initialized for rebalance", slog.String("address", sgn.Address().Hex()))

	rpcURL := os.Getenv("RPC_URL")
	if rpcURL == "" {
		return nil, errEnv("RPC_URL")
	}

	ethClient, err := ethclient.Dial(rpcURL)
	if err != nil {
		return nil, err
	}

	factoryAddr := os.Getenv("FACTORY_ADDRESS")
	if factoryAddr == "" {
		ethClient.Close()
		return nil, errEnv("FACTORY_ADDRESS")
	}

	vaultSource := NewFactoryVaultSource(ethClient, common.HexToAddress(factoryAddr))

	stateViewAddr := os.Getenv("STATEVIEW_CONTRACT_ADDR")
	if stateViewAddr == "" {
		ethClient.Close()
		return nil, errEnv("STATEVIEW_CONTRACT_ADDR")
	}

	liqRepo, err := liquidityrepo.New(liquidityrepo.Config{
		RPCURL:          rpcURL,
		ContractAddress: stateViewAddr,
	})
	if err != nil {
		ethClient.Close()
		return nil, err
	}

	liqSvc := liquidityservice.New(liqRepo)
	strategySvc := strategyservice.New(liqSvc)

	agentSvc := New(
		vaultSource,
		strategySvc,
		sgn,
		ethClient,
		logger,
		liqRepo,
	)

	applyProtectionFromEnv(agentSvc, logger)

	ctxCron, cancel := context.WithCancel(ctx)
	c := cron.New()

	_, err = c.AddFunc(schedule, func() {
		runOnce(ctxCron, agentSvc, logger)
	})
	if err != nil {
		cancel()
		ethClient.Close()
		liqRepo.Close()

		return nil, err
	}

	c.Start()
	logger.Info("rebalance cron started", slog.String("schedule", schedule))

	runOnce(ctxCron, agentSvc, logger)

	stop = func() {
		c.Stop()
		cancel()
		ethClient.Close()
		liqRepo.Close()
	}

	return stop, nil
}

func runOnce(ctx context.Context, agentSvc *Service, logger *slog.Logger) {
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

func applyProtectionFromEnv(svc *Service, logger *slog.Logger) {
	swapSlippage := parseInt64(os.Getenv("SWAP_SLIPPAGE_BPS"), 50)
	mintSlippage := parseInt64(os.Getenv("MINT_SLIPPAGE_BPS"), 50)
	maxGasPrice := parseFloat64(os.Getenv("MAX_GAS_PRICE_GWEI"), 1.0)
	devThreshold := parseFloat64(os.Getenv("DEVIATION_THRESHOLD"), 0.1)

	svc.SetProtectionSettings(swapSlippage, mintSlippage, maxGasPrice)
	svc.SetDeviationThreshold(devThreshold)

	if tickRange := os.Getenv("TICK_RANGE_AROUND_CURRENT"); tickRange != "" {
		if val, err := strconv.ParseInt(tickRange, 10, 32); err == nil && val > 0 {
			svc.SetTickRangeAroundCurrent(int32(val))
			logger.Info("tick range override set", slog.String("raw", tickRange), slog.Int64("value", val))
		} else {
			logger.Warn("invalid TICK_RANGE_AROUND_CURRENT", slog.String("raw", tickRange))
		}
	}
}

func parseInt64(s string, defaultVal int64) int64 {
	if s == "" {
		return defaultVal
	}

	val, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return defaultVal
	}

	return val
}

func parseFloat64(s string, defaultVal float64) float64 {
	if s == "" {
		return defaultVal
	}

	val, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return defaultVal
	}

	return val
}

func errEnv(name string) error {
	return &envError{name: name}
}

type envError struct {
	name string
}

func (e *envError) Error() string {
	return e.name + " not set"
}
