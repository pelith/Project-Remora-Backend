package agent

import (
	"context"
	"fmt"
	"log/slog"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"

	"remora/internal/allocation"
	"remora/internal/vault"
)

// sendAndWait is a helper to send a transaction and wait for it to be mined.
func (s *Service) sendAndWait(ctx context.Context, txName string, tx *types.Transaction) error {
	s.logger.Info(fmt.Sprintf("%s transaction sent", txName), slog.String("tx", tx.Hash().Hex()))

	receipt, err := bind.WaitMined(ctx, s.ethClient, tx)
	if err != nil {
		return fmt.Errorf("wait for %s tx mined: %w", txName, err)
	}

	if receipt.Status != types.ReceiptStatusSuccessful {
		return fmt.Errorf("%s transaction failed: receipt status %v", txName, receipt.Status)
	}

	s.logger.Info(fmt.Sprintf("%s transaction confirmed", txName), 
		slog.String("tx", tx.Hash().Hex()), 
		slog.Uint64("block", receipt.BlockNumber.Uint64()))
	return nil
}

// executeRebalance orchestrates the execution of a rebalance plan.
// Flow: Burn old positions -> Swap tokens -> Mint new positions.
func (s *Service) executeRebalance(
	ctx context.Context,
	vaultClient vault.Vault,
	oldPositions []vault.Position,
	result *allocation.AllocationResult,
	currentSqrtPriceX96 *big.Int,
	token0 common.Address,
	token1 common.Address,
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
		tx, err := vaultClient.BurnPosition(ctx, pos.TokenID, big.NewInt(0), big.NewInt(0), deadline)
		if err != nil {
			return fmt.Errorf("burn position %s: %w", pos.TokenID.String(), err)
		}

		if err := s.sendAndWait(ctx, fmt.Sprintf("burn-%s", pos.TokenID.String()), tx); err != nil {
			return err
		}
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
		
		if err := s.sendAndWait(ctx, "swap", tx); err != nil {
			return err
		}
	}

	// Post-swap: compare actual vault balance vs allocation expectation
	postSwap0, _ := s.getTokenBalance(ctx, token0, vaultClient.Address())
	postSwap1, _ := s.getTokenBalance(ctx, token1, vaultClient.Address())
	surplus0 := new(big.Int).Sub(postSwap0, result.TotalAmount0)
	surplus1 := new(big.Int).Sub(postSwap1, result.TotalAmount1)
	s.logger.Info("post-swap balance vs allocation",
		slog.String("vault_token0", postSwap0.String()),
		slog.String("vault_token1", postSwap1.String()),
		slog.String("alloc_needs_token0", result.TotalAmount0.String()),
		slog.String("alloc_needs_token1", result.TotalAmount1.String()),
		slog.String("surplus_token0", surplus0.String()),
		slog.String("surplus_token1", surplus1.String()))

	// 3. Mint New Positions
	for i, posPlan := range result.Positions {
		// Pre-mint balance — use actual vault balance as max (not inflated planned amount)
		// This allows POSM to use any "reserve" from the allocation buffer if price moved.
		preMint0, _ := s.getTokenBalance(ctx, token0, vaultClient.Address())
		preMint1, _ := s.getTokenBalance(ctx, token1, vaultClient.Address())

		liquidityToMint := posPlan.Liquidity
		amount0Required := posPlan.Amount0
		amount1Required := posPlan.Amount1

		// Add tiny rounding buffer (1 bps = 0.01%) to handle POSM rounding up.
		// The vault contract consumes exactly amount0Max, so keep it minimal.
		roundingBuffer := big.NewInt(10001) // 1.0001x
		amount0Max := new(big.Int).Mul(amount0Required, roundingBuffer)
		amount0Max.Div(amount0Max, big.NewInt(10000))
		amount1Max := new(big.Int).Mul(amount1Required, roundingBuffer)
		amount1Max.Div(amount1Max, big.NewInt(10000))

		s.logger.Info("minting new position",
			slog.Int("index", i),
			slog.Int("tickLower", posPlan.TickLower),
			slog.Int("tickUpper", posPlan.TickUpper),
			slog.String("liquidity", liquidityToMint.String()),
			slog.String("amount0_required", amount0Required.String()),
			slog.String("amount1_required", amount1Required.String()),
			slog.String("amount0_max", amount0Max.String()),
			slog.String("amount1_max", amount1Max.String()))

		tx, err := vaultClient.MintPosition(
			ctx,
			int32(posPlan.TickLower),
			int32(posPlan.TickUpper),
			liquidityToMint,
			amount0Max,
			amount1Max,
			deadline,
		)
		if err != nil {
			return fmt.Errorf("mint position %d: %w", i, err)
		}

		if err := s.sendAndWait(ctx, fmt.Sprintf("mint-%d", i), tx); err != nil {
			return err
		}

		// Post-mint: compare actual consumed vs expected
		postMint0, _ := s.getTokenBalance(ctx, token0, vaultClient.Address())
		postMint1, _ := s.getTokenBalance(ctx, token1, vaultClient.Address())
		consumed0 := new(big.Int).Sub(preMint0, postMint0)
		consumed1 := new(big.Int).Sub(preMint1, postMint1)
		s.logger.Info("mint consumed vs expected",
			slog.Int("index", i),
			slog.String("expected0", amount0Required.String()),
			slog.String("consumed0", consumed0.String()),
			slog.String("expected1", amount1Required.String()),
			slog.String("consumed1", consumed1.String()))
	}

	return nil
}
