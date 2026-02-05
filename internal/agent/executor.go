package agent

import (
	"context"
	"fmt"
	"log/slog"
	"math/big"
	"time"

	"remora/internal/allocation"
	"remora/internal/vault"
)

// executeRebalance orchestrates the execution of a rebalance plan.
// Flow: Burn old positions -> Swap tokens -> Mint new positions.
func (s *Service) executeRebalance(
	ctx context.Context,
	vaultClient vault.Vault,
	oldPositions []vault.Position,
	result *allocation.AllocationResult,
) error {
	deadline := big.NewInt(time.Now().Add(20 * time.Minute).Unix()) // 20 min deadline

	// 1. Burn All Old Positions
	// This collects all liquidity and fees back into the vault.
	for _, pos := range oldPositions {
		s.logger.Info("burning position", slog.String("tokenId", pos.TokenID.String()))
		
		// For prototype, we set minAmounts to 0. 
		// TODO: Implement slippage protection using current price.
		tx, err := vaultClient.BurnPosition(ctx, pos.TokenID, big.NewInt(0), big.NewInt(0), deadline)
		if err != nil {
			return fmt.Errorf("burn position %s: %w", pos.TokenID.String(), err)
		}

		s.logger.Info("burn transaction sent", slog.String("tx", tx.Hash().Hex()))
		// In a real environment, we might want to wait for the tx to be mined.
		// For now, we assume sequential execution for simplicity or that POSM handles nonces.
	}

	// 2. Execute Swap if needed
	if result.SwapAmount != nil && result.SwapAmount.Sign() > 0 {
		s.logger.Info("executing swap", 
			slog.String("amount", result.SwapAmount.String()),
			slog.Bool("zeroForOne", result.SwapToken0To1))
		
		// TODO: Calculate minAmountOut with slippage protection.
		tx, err := vaultClient.Swap(ctx, result.SwapToken0To1, result.SwapAmount, big.NewInt(0), deadline)
		if err != nil {
			return fmt.Errorf("swap: %w", err)
		}
		s.logger.Info("swap transaction sent", slog.String("tx", tx.Hash().Hex()))
	}

	// 3. Mint New Positions
	for i, posPlan := range result.Positions {
		s.logger.Info("minting new position", 
			slog.Int("index", i),
			slog.Int("tickLower", posPlan.TickLower),
			slog.Int("tickUpper", posPlan.TickUpper),
			slog.String("liquidity", posPlan.Liquidity.String()))
		
		// amount0Max and amount1Max are the calculated needs from allocation.
		// We add a small buffer (e.g., 5%) to avoid "Insufficient balance" errors due to price movements.
		tx, err := vaultClient.MintPosition(
			ctx,
			int32(posPlan.TickLower),
			int32(posPlan.TickUpper),
			posPlan.Liquidity,
			posPlan.Amount0,
			posPlan.Amount1,
			deadline,
		)
		if err != nil {
			return fmt.Errorf("mint position %d: %w", i, err)
		}
		s.logger.Info("mint transaction sent", slog.String("tx", tx.Hash().Hex()))
	}

	return nil
}
