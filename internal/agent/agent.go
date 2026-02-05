package agent

import (
	"context"
	"fmt"
	"log/slog"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"

	"remora/internal/allocation"
	"remora/internal/coverage"
	"remora/internal/liquidity/poolid"
	"remora/internal/signer"
	"remora/internal/strategy"
	"remora/internal/vault"
)

const defaultDeviationThreshold = 0.1

// VaultSource provides access to vault addresses.
type VaultSource interface {
	GetVaultAddresses(ctx context.Context) ([]common.Address, error)
}

// RebalanceResult represents the result of a rebalance operation.
type RebalanceResult struct {
	VaultAddress common.Address
	Rebalanced   bool
	Reason       string // "deviation_exceeded", "skipped", "error"
}

// Service is the main agent orchestrator.
type Service struct {
	vaultSource   VaultSource
	strategySvc   strategy.Service
	signer        *signer.Signer
	ethClient     *ethclient.Client
	logger        *slog.Logger

	deviationThreshold float64
	swapSlippageBps    int64
	mintSlippageBps    int64
	maxGasPriceGwei    float64
	tickRangeOverride  int32
}

// New creates a new agent service.
func New(
	vaultSource VaultSource,
	strategySvc strategy.Service,
	signer *signer.Signer,
	ethClient *ethclient.Client,
	logger *slog.Logger,
) *Service {
	return &Service{
		vaultSource:        vaultSource,
		strategySvc:        strategySvc,
		signer:             signer,
		ethClient:          ethClient,
		logger:             logger,
		deviationThreshold: 0.1,
		swapSlippageBps:    50,  // default: 0.5%
		mintSlippageBps:    50,  // default: 0.5%
		maxGasPriceGwei:    1.0, // default: 1.0 Gwei (suitable for many L2s)
	}
}

// SetProtectionSettings updates the protection settings for the service.
func (s *Service) SetProtectionSettings(swapSlippageBps int64, mintSlippageBps int64, maxGasPriceGwei float64) {
	s.swapSlippageBps = swapSlippageBps
	s.mintSlippageBps = mintSlippageBps
	s.maxGasPriceGwei = maxGasPriceGwei
}

// SetDeviationThreshold updates the threshold for rebalance decision.
func (s *Service) SetDeviationThreshold(threshold float64) {
	s.deviationThreshold = threshold
}

// SetTickRangeAroundCurrent overrides the tick range used for market scan.
// Value is interpreted as +/- ticks around current tick.
func (s *Service) SetTickRangeAroundCurrent(tickRange int32) {
	s.tickRangeOverride = tickRange
}

// Run executes one round of rebalance check for all vaults.
func (s *Service) Run(ctx context.Context) ([]RebalanceResult, error) {
	addresses, err := s.vaultSource.GetVaultAddresses(ctx)
	if err != nil {
		return nil, fmt.Errorf("get vault addresses: %w", err)
	}

	s.logger.InfoContext(ctx, "starting rebalance run", slog.Int("vault_count", len(addresses)))

	var results []RebalanceResult

	for _, addr := range addresses {
		result := s.processVault(ctx, addr)
		results = append(results, result)
	}

	return results, nil
}

// processVault handles rebalance logic for a single vault.
func (s *Service) processVault(ctx context.Context, vaultAddr common.Address) RebalanceResult {

	s.logger.Info("processing vault", slog.String("address", vaultAddr.Hex()))

	// Step 1: Create vault client
	auth, err := s.signer.TransactOpts()
	if err != nil {
		return RebalanceResult{VaultAddress: vaultAddr, Reason: "signer_error"}
	}

	vaultClient, err := vault.NewClient(vaultAddr, s.ethClient, auth)
	if err != nil {
		return RebalanceResult{VaultAddress: vaultAddr, Reason: "vault_client_error"}
	}

	// Step 2: Get vault state and current positions
	state, err := vaultClient.GetState(ctx)
	if err != nil {
		s.logger.Error("failed to get vault state", slog.Any("error", err))
		return RebalanceResult{VaultAddress: vaultAddr, Reason: "get_state_error"}
	}

	// Step 3: Compute target positions using strategy service
	// Convert vault.PoolKey to poolid.PoolKey
	liqPoolKey := poolid.PoolKey{
		Currency0:   state.PoolKey.Currency0.Hex(),
		Currency1:   state.PoolKey.Currency1.Hex(),
		Fee:         uint32(state.PoolKey.Fee.Uint64()),         //nolint:gosec // fee fits in uint24
		TickSpacing: int32(state.PoolKey.TickSpacing.Int64()),   //nolint:gosec // tickSpacing fits in int24
		Hooks:       state.PoolKey.Hooks.Hex(),
	}

	tickSpacing := int32(state.PoolKey.TickSpacing.Int64()) //nolint:gosec // tickSpacing fits in int24
	
	// Determine effective scan radius. 
	// Start with the full width of the vault's allowed range.
	vaultRange := state.AllowedTickUpper - state.AllowedTickLower
	tickRange := vaultRange
	
	// If an override is provided (TICK_RANGE_AROUND_CURRENT) and the vault range is wider,
	// cap the scan radius to the override value to optimize performance.
	if s.tickRangeOverride > 0 && tickRange > s.tickRangeOverride {
		s.logger.Info("capping tick range with override", 
			slog.Int("vault_range", int(vaultRange)), 
			slog.Int("override_limit", int(s.tickRangeOverride)))
		tickRange = s.tickRangeOverride
	}

	s.logger.Info("tick range selection",
		slog.Int("vault_allowed_width", int(vaultRange)),
		slog.Int("override_setting", int(s.tickRangeOverride)),
		slog.Int("final_scan_radius", int(tickRange)),
	)

	// Build coverage config: use vault's MaxPositionsK if set, otherwise default
	algoConfig := coverage.DefaultConfig()
	if state.MaxPositionsK != nil && state.MaxPositionsK.Sign() > 0 {
		algoConfig.N = int(state.MaxPositionsK.Int64())
	}

	computeParams := &strategy.ComputeParams{
		PoolKey:          liqPoolKey,
		BinSizeTicks:     tickSpacing,
		TickRange:        tickRange,
		AlgoConfig:       algoConfig,
		AllowedTickLower: state.AllowedTickLower,
		AllowedTickUpper: state.AllowedTickUpper,
	}

	s.logger.Info("computing target positions",
		slog.String("pool_currency0", liqPoolKey.Currency0),
		slog.String("pool_currency1", liqPoolKey.Currency1),
		slog.Int("tick_spacing", int(tickSpacing)),
		slog.Int("tick_range", int(tickRange)),
		slog.Int("max_positions", algoConfig.N),
		slog.Int("allowed_lower", int(state.AllowedTickLower)),
		slog.Int("allowed_upper", int(state.AllowedTickUpper)),
	)

	targetResult, err := s.strategySvc.ComputeTargetPositions(ctx, computeParams)
	if err != nil {
		s.logger.Error("failed to compute target", slog.Any("error", err))
		return RebalanceResult{VaultAddress: vaultAddr, Reason: "strategy_error"}
	}

	s.logger.Info("target positions computed",
		slog.Int("segments", len(targetResult.Segments)),
		slog.Int("bins", len(targetResult.Bins)),
		slog.Int("current_tick", int(targetResult.CurrentTick)),
	)

	// Step 4: Calculate Total Assets (Idle + Invested)
	token0 := state.PoolKey.Currency0
	token1 := state.PoolKey.Currency1

	// TODO: Cache decimals
	decimals0, err := s.getTokenDecimals(ctx, token0)
	if err != nil {
		s.logger.Error("failed to get token0 decimals", slog.Any("error", err))
		return RebalanceResult{VaultAddress: vaultAddr, Reason: "token_error"}
	}
	decimals1, err := s.getTokenDecimals(ctx, token1)
	if err != nil {
		s.logger.Error("failed to get token1 decimals", slog.Any("error", err))
		return RebalanceResult{VaultAddress: vaultAddr, Reason: "token_error"}
	}

	// Get Idle Balances
	idle0, err := s.getTokenBalance(ctx, token0, vaultAddr)
	if err != nil {
		s.logger.Error("failed to get token0 balance", slog.Any("error", err))
		return RebalanceResult{VaultAddress: vaultAddr, Reason: "balance_error"}
	}
	idle1, err := s.getTokenBalance(ctx, token1, vaultAddr)
	if err != nil {
		s.logger.Error("failed to get token1 balance", slog.Any("error", err))
		return RebalanceResult{VaultAddress: vaultAddr, Reason: "balance_error"}
	}

	// Get Invested Balances (from current positions)
	positions, err := vaultClient.GetPositions(ctx)
	if err != nil {
		s.logger.Error("failed to get positions", slog.Any("error", err))
		return RebalanceResult{VaultAddress: vaultAddr, Reason: "get_positions_error"}
	}

	s.logger.Info("fetched positions from vault", slog.Any("token_ids", func() []string {
		ids := make([]string, len(positions))
		for i, p := range positions {
			ids[i] = p.TokenID.String()
		}
		return ids
	}()))

	invested0 := big.NewInt(0)
	invested1 := big.NewInt(0)

	for i := range positions {
		pos := &positions[i]
		// Fetch real liquidity from POSM/StateView
		liquidity, err := s.getPositionLiquidity(ctx, state.Posm, pos.TokenID)
		if err != nil {
			s.logger.Warn("failed to get position liquidity",
				slog.String("tokenID", pos.TokenID.String()),
				slog.Any("error", err))
			continue
		}
		pos.Liquidity = liquidity

		if pos.Liquidity == nil || pos.Liquidity.Sign() == 0 {
			continue
		}

		sqrtPriceAX96 := allocation.TickToSqrtPriceX96(int(pos.TickLower))
		sqrtPriceBX96 := allocation.TickToSqrtPriceX96(int(pos.TickUpper))

		amt0 := allocation.GetAmount0ForLiquidity(targetResult.SqrtPriceX96, sqrtPriceAX96, sqrtPriceBX96, pos.Liquidity)
		amt1 := allocation.GetAmount1ForLiquidity(targetResult.SqrtPriceX96, sqrtPriceAX96, sqrtPriceBX96, pos.Liquidity)

		invested0.Add(invested0, amt0)
		invested1.Add(invested1, amt1)
	}

	// Sum total assets
	total0 := new(big.Int).Add(idle0, invested0)
	total1 := new(big.Int).Add(idle1, invested1)

	// Apply Safety Buffer
	// Reduce available funds by mint slippage tolerance to ensure successful minting even if price moves.
	bufferBps := big.NewInt(s.mintSlippageBps)
	multiplier := big.NewInt(10000)
	multiplier.Sub(multiplier, bufferBps)

	total0.Mul(total0, multiplier).Div(total0, big.NewInt(10000))
	total1.Mul(total1, multiplier).Div(total1, big.NewInt(10000))

	s.logger.Info("preparing allocation",
		slog.Int("decimals0", int(decimals0)),
		slog.Int("decimals1", int(decimals1)),
		slog.String("total0_buffered", total0.String()),
		slog.String("total1_buffered", total1.String()),
	)

	// Step 5: Allocate
	poolState := allocation.PoolState{
		SqrtPriceX96:   targetResult.SqrtPriceX96,
		CurrentTick:    int(targetResult.CurrentTick),
		Token0Decimals: int(decimals0),
		Token1Decimals: int(decimals1),
	}

	userFunds := allocation.UserFunds{
		Amount0: total0,
		Amount1: total1,
	}

	allocationResult, err := allocation.Allocate(targetResult.Segments, userFunds, poolState, state.SwapAllowed)
	if err != nil {
		s.logger.Error("failed to allocate", slog.Any("error", err))
		return RebalanceResult{VaultAddress: vaultAddr, Reason: "allocation_error"}
	}

	// Step 4.5: Deviation Check
	// Decide if we really need to rebalance based on post-allocation plan
	deviation := s.calculateDeviation(positions, allocationResult.Positions)
	s.logger.Info("deviation calculated", slog.Float64("deviation", deviation), slog.Float64("threshold", s.deviationThreshold))

	if deviation < s.deviationThreshold {
		return RebalanceResult{
			VaultAddress: vaultAddr,
			Rebalanced:   false,
			Reason:       "deviation_below_threshold",
		}
	}

	s.logger.Info("allocation computed",
		slog.String("swap_amount", allocationResult.SwapAmount.String()),
		slog.Bool("zero_for_one", allocationResult.SwapToken0To1),
		slog.Int("new_positions", len(allocationResult.Positions)),
		slog.String("total_amount0", allocationResult.TotalAmount0.String()),
		slog.String("total_amount1", allocationResult.TotalAmount1.String()),
		slog.String("available_token0", total0.String()),
		slog.String("available_token1", total1.String()),
	)

	for i, pos := range allocationResult.Positions {
		s.logger.Info("allocation position detail",
			slog.Int("index", i),
			slog.Int("tickLower", pos.TickLower),
			slog.Int("tickUpper", pos.TickUpper),
			slog.String("liquidity", pos.Liquidity.String()),
			slog.String("amount0", pos.Amount0.String()),
			slog.String("amount1", pos.Amount1.String()),
			slog.Float64("weight", pos.Weight),
		)
	}

	// Step 6: Execute rebalance
	err = s.executeRebalance(ctx, vaultClient, positions, allocationResult, targetResult.SqrtPriceX96, token0, token1)
	if err != nil {
		s.logger.Error("failed to execute rebalance", slog.Any("error", err))
		return RebalanceResult{VaultAddress: vaultAddr, Reason: "execution_error"}
	}

	return RebalanceResult{
		VaultAddress: vaultAddr,
		Rebalanced:   true,
		Reason:       "success",
	}
}
