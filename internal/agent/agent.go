package agent

import (
	"context"
	"log/slog"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"

	"remora/internal/allocation"
	"remora/internal/coverage"
	"remora/internal/liquidity"
	"remora/internal/signer"
	"remora/internal/strategy"
	"remora/internal/vault"
)

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
	stateViewAddr common.Address

	deviationThreshold float64
}

// New creates a new agent service.
func New(
	vaultSource VaultSource,
	strategySvc strategy.Service,
	signer *signer.Signer,
	ethClient *ethclient.Client,
	logger *slog.Logger,
	stateViewAddr common.Address,
) *Service {
	return &Service{
		vaultSource:        vaultSource,
		strategySvc:        strategySvc,
		signer:             signer,
		ethClient:          ethClient,
		logger:             logger,
		stateViewAddr:      stateViewAddr,
		deviationThreshold: 0.1,
	}
}

// Run executes one round of rebalance check for all vaults.
func (s *Service) Run(ctx context.Context) ([]RebalanceResult, error) {
	addresses, err := s.vaultSource.GetVaultAddresses(ctx)
	if err != nil {
		return nil, err
	}

	s.logger.Info("starting rebalance run", slog.Int("vault_count", len(addresses)))

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
	// Convert vault.PoolKey to liquidity.PoolKey
	liqPoolKey := liquidity.PoolKey{
		Currency0:   state.PoolKey.Currency0.Hex(),
		Currency1:   state.PoolKey.Currency1.Hex(),
		Fee:         uint32(state.PoolKey.Fee.Uint64()),         //nolint:gosec // fee fits in uint24
		TickSpacing: int32(state.PoolKey.TickSpacing.Int64()),   //nolint:gosec // tickSpacing fits in int24
		Hooks:       state.PoolKey.Hooks.Hex(),
	}

	computeParams := &strategy.ComputeParams{
		PoolKey:      liqPoolKey,
		BinSizeTicks: 200,  // TODO: Configurable
		TickRange:    1000, // TODO: Configurable
		AlgoConfig:   coverage.DefaultConfig(),
	}

	targetResult, err := s.strategySvc.ComputeTargetPositions(ctx, computeParams)
	if err != nil {
		s.logger.Error("failed to compute target", slog.Any("error", err))
		return RebalanceResult{VaultAddress: vaultAddr, Reason: "strategy_error"}
	}

	// Step 4: Calculate Total Assets (Idle + Invested)
	// We need decimals and balances
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

	invested0 := big.NewInt(0)
	invested1 := big.NewInt(0)

	// We use the Strategy's SqrtPriceX96 to estimate current position value
	// Note: accurate value requires getting the real positions info including uncollected fees,
	// but here we just estimate principal from liquidity.
	for _, pos := range positions {
		// Fetch real liquidity from POSM/StateView
		liquidity, err := s.getPositionLiquidity(ctx, state.Posm, state.PoolID, pos.TokenID)
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

		// Calculate amounts for this position
		// allocation.GetAmount0ForLiquidity needs sqrtPriceX96, sqrtPriceA, sqrtPriceB, liquidity
		sqrtPriceAX96 := allocation.TickToSqrtPriceX96(int(pos.TickLower))
		sqrtPriceBX96 := allocation.TickToSqrtPriceX96(int(pos.TickUpper))

		amt0 := allocation.GetAmount0ForLiquidity(targetResult.SqrtPriceX96, sqrtPriceAX96, sqrtPriceBX96, pos.Liquidity)
		amt1 := allocation.GetAmount1ForLiquidity(targetResult.SqrtPriceX96, sqrtPriceAX96, sqrtPriceBX96, pos.Liquidity)

		invested0.Add(invested0, amt0)
		invested1.Add(invested1, amt1)
	}

	// Sum total
	total0 := new(big.Int).Add(idle0, invested0)
	total1 := new(big.Int).Add(idle1, invested1)

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

	s.logger.Info("allocation computed",
		slog.String("swap_amount", allocationResult.SwapAmount.String()),
		slog.Bool("zero_for_one", allocationResult.SwapToken0To1),
		slog.Int("new_positions", len(allocationResult.Positions)),
	)

	// Step 6: Execute rebalance
	err = s.executeRebalance(ctx, vaultClient, positions, allocationResult)
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
