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
	currentSqrtPriceX96 *big.Int,
) error {
	// 0. Gas Price Check
	gasPrice, err := s.ethClient.SuggestGasPrice(ctx)
	if err != nil {
		return fmt.Errorf("suggest gas price: %w", err)
	}

	// Convert maxGasPriceGwei to wei
	// 1 Gwei = 10^9 wei
	maxGasPriceWei := new(big.Float).Mul(big.NewFloat(s.maxGasPriceGwei), big.NewFloat(1e9))
	maxGasPriceWeiInt, _ := maxGasPriceWei.Int(nil)

	if gasPrice.Cmp(maxGasPriceWeiInt) > 0 {
		s.logger.Warn("gas price too high, skipping rebalance", 
			slog.String("current", gasPrice.String()), 
			slog.String("limit", maxGasPriceWeiInt.String()))
		return fmt.Errorf("gas price too high: %s > %s", gasPrice.String(), maxGasPriceWeiInt.String())
	}

	deadline := big.NewInt(time.Now().Add(20 * time.Minute).Unix()) // 20 min deadline

	// 1. Burn All Old Positions
	// This collects all liquidity and fees back into the vault.
	for _, pos := range oldPositions {
		s.logger.Info("burning position", slog.String("tokenId", pos.TokenID.String()))
		
		// For burn, we use 0 min amounts for now (collecting all available)
		tx, err := vaultClient.BurnPosition(ctx, pos.TokenID, big.NewInt(0), big.NewInt(0), deadline)
		if err != nil {
			return fmt.Errorf("burn position %s: %w", pos.TokenID.String(), err)
		}

		s.logger.Info("burn transaction sent", slog.String("tx", tx.Hash().Hex()))
	}

	// 2. Execute Swap if needed
	if result.SwapAmount != nil && result.SwapAmount.Sign() > 0 {
		// Calculate minAmountOut with slippage protection
		// Expected Out = amountIn * price (if zeroForOne) or amountIn / price (if oneForZero)
		// price = sqrtPriceX96^2 / Q192
		
		sqrtPriceSquared := new(big.Int).Mul(currentSqrtPriceX96, currentSqrtPriceX96)
		var expectedOut *big.Int
		
		if result.SwapToken0To1 {
			// Token0 -> Token1
			// out = in * sqrtP^2 / Q192
			expectedOut = new(big.Int).Mul(result.SwapAmount, sqrtPriceSquared)
			expectedOut.Div(expectedOut, allocation.Q192)
		} else {
			// Token1 -> Token0
			// out = in * Q192 / sqrtP^2
			expectedOut = new(big.Int).Mul(result.SwapAmount, allocation.Q192)
			expectedOut.Div(expectedOut, sqrtPriceSquared)
		}

		// minAmountOut = expectedOut * (10000 - slippageBps) / 10000
		multiplier := big.NewInt(10000 - s.swapSlippageBps)
		minAmountOut := new(big.Int).Mul(expectedOut, multiplier)
		minAmountOut.Div(minAmountOut, big.NewInt(10000))

		s.logger.Info("executing swap", 
			slog.String("amountIn", result.SwapAmount.String()),
			slog.String("minAmountOut", minAmountOut.String()),
			slog.Bool("zeroForOne", result.SwapToken0To1))
		
		tx, err := vaultClient.Swap(ctx, result.SwapToken0To1, result.SwapAmount, minAmountOut, deadline)
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
		
		// amount0Max and amount1Max are used as hard caps (no slippage buffer added as requested)
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
