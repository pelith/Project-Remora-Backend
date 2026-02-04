package agent

import (
	"context"
	"log/slog"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"

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
	vaultSource VaultSource
	strategySvc strategy.Service
	signer      *signer.Signer
	ethClient   *ethclient.Client
	logger      *slog.Logger

	deviationThreshold float64
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
	// TODO: state, err := vaultClient.GetState(ctx)
	// TODO: currentPositions, err := vaultClient.GetPositions(ctx)
	_ = vaultClient

	// Step 3: Compute target positions using strategy service
	// TODO: targetResult, err := s.computeTargetPositions(ctx, state.PoolKey)

	// Step 4: Calculate deviation between current and target
	// TODO: deviation := s.calculateDeviation(currentPositions, targetResult)

	// Step 5: Check if rebalance is needed
	// TODO: if deviation < s.deviationThreshold { return skipped }

	// Step 6: Execute rebalance
	// TODO: err := s.executeRebalance(ctx, vaultClient, targetResult)

	return RebalanceResult{
		VaultAddress: vaultAddr,
		Rebalanced:   false,
		Reason:       "not_implemented",
	}
}

// =============================================================================
// Private methods to implement
// =============================================================================

// computeTargetPositions computes target positions for a vault.
// Flow: PoolKey -> liquidity.GetDistribution -> strategy.ComputeTargetPositions
// func (s *Service) computeTargetPositions(ctx context.Context, poolKey vault.PoolKey) (*strategy.ComputeResult, error)

// calculateDeviation calculates deviation between current and target positions.
// func (s *Service) calculateDeviation(current []vault.Position, target *strategy.ComputeResult) float64

// executeRebalance executes rebalance transactions.
// Flow: 1. Burn all existing positions  2. Mint new positions
// func (s *Service) executeRebalance(ctx context.Context, client vault.Vault, target *strategy.ComputeResult) error
