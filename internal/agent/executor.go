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
	"remora/internal/liquidity/poolid"
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
	poolKey *poolid.PoolKey,
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

	// Read actual vault balances after swap for refit calculation.
	postSwap0, _ := s.getTokenBalance(ctx, token0, vaultClient.Address())
	postSwap1, _ := s.getTokenBalance(ctx, token1, vaultClient.Address())

	// Refresh sqrtPriceX96 after swap (price may move).
	effectiveSqrtPriceX96 := currentSqrtPriceX96
	if s.slot0Fetcher != nil && poolKey != nil {
		if slot0, err := s.slot0Fetcher.GetSlot0(ctx, poolKey); err != nil {
			s.logger.Warn("failed to refresh slot0 after swap", slog.Any("error", err))
		} else if slot0 != nil && slot0.SqrtPriceX96 != nil {
			effectiveSqrtPriceX96 = slot0.SqrtPriceX96
		}
	}

	// Refit positions to actual post-swap balances.
	// Out-of-range positions (single-token) are unaffected by price movement — keep as-is.
	// Only the in-range position (contains current price) needs recalculation.
	remaining0 := new(big.Int).Set(postSwap0)
	remaining1 := new(big.Int).Set(postSwap1)
	inRangeIdx := -1

	for i := range result.Positions {
		p := &result.Positions[i]
		sqrtA := allocation.TickToSqrtPriceX96(p.TickLower)
		sqrtB := allocation.TickToSqrtPriceX96(p.TickUpper)

		if effectiveSqrtPriceX96.Cmp(sqrtA) > 0 && effectiveSqrtPriceX96.Cmp(sqrtB) < 0 {
			inRangeIdx = i
			continue
		}

		// Out-of-range: liquidity unchanged, recompute amounts at new price for consistency.
		p.Amount0 = allocation.GetAmount0ForLiquidity(effectiveSqrtPriceX96, sqrtA, sqrtB, p.Liquidity)
		p.Amount1 = allocation.GetAmount1ForLiquidity(effectiveSqrtPriceX96, sqrtA, sqrtB, p.Liquidity)
		remaining0.Sub(remaining0, p.Amount0)
		remaining1.Sub(remaining1, p.Amount1)
	}

	if inRangeIdx >= 0 {
		p := &result.Positions[inRangeIdx]
		sqrtA := allocation.TickToSqrtPriceX96(p.TickLower)
		sqrtB := allocation.TickToSqrtPriceX96(p.TickUpper)

		if remaining0.Sign() < 0 {
			remaining0.SetInt64(0)
		}
		if remaining1.Sign() < 0 {
			remaining1.SetInt64(0)
		}

		newLiq := allocation.GetLiquidityForAmounts(effectiveSqrtPriceX96, sqrtA, sqrtB, remaining0, remaining1)
		p.Liquidity = newLiq
		p.Amount0 = allocation.GetAmount0ForLiquidity(effectiveSqrtPriceX96, sqrtA, sqrtB, newLiq)
		p.Amount1 = allocation.GetAmount1ForLiquidity(effectiveSqrtPriceX96, sqrtA, sqrtB, newLiq)

		s.logger.Info("refitted in-range position to post-swap balances",
			slog.Int("index", inRangeIdx),
			slog.String("remaining0", remaining0.String()),
			slog.String("remaining1", remaining1.String()),
			slog.String("new_liquidity", newLiq.String()),
			slog.String("amount0", p.Amount0.String()),
			slog.String("amount1", p.Amount1.String()))
	}

	// Update totals.
	totalAmount0 := big.NewInt(0)
	totalAmount1 := big.NewInt(0)
	for _, p := range result.Positions {
		totalAmount0.Add(totalAmount0, p.Amount0)
		totalAmount1.Add(totalAmount1, p.Amount1)
	}
	result.TotalAmount0 = totalAmount0
	result.TotalAmount1 = totalAmount1

	// 3. Mint New Positions
	for i, posPlan := range result.Positions {
		// Pre-mint balance — use actual vault balance as max (not inflated planned amount)
		// This allows POSM to use any "reserve" from the allocation buffer if price moved.
		preMint0, _ := s.getTokenBalance(ctx, token0, vaultClient.Address())
		preMint1, _ := s.getTokenBalance(ctx, token1, vaultClient.Address())

		liquidityToMint := posPlan.Liquidity
		amount0Required := posPlan.Amount0
		amount1Required := posPlan.Amount1

		// Compute amountMax with +1 wei rounding buffer for POSM rounding up.
		amount0Max := new(big.Int).Set(amount0Required)
		amount1Max := new(big.Int).Set(amount1Required)
		if amount0Max.Sign() > 0 {
			amount0Max.Add(amount0Max, big.NewInt(1))
		}
		if amount1Max.Sign() > 0 {
			amount1Max.Add(amount1Max, big.NewInt(1))
		}

		// If vault balance can't cover amountMax (including rounding buffer),
		// recalculate L from (balance - 1 wei) so that amount + 1 <= balance.
		if amount0Max.Cmp(preMint0) > 0 || amount1Max.Cmp(preMint1) > 0 {
			available0 := new(big.Int).Set(preMint0)
			available1 := new(big.Int).Set(preMint1)
			if available0.Sign() > 0 {
				available0.Sub(available0, big.NewInt(1))
			}
			if available1.Sign() > 0 {
				available1.Sub(available1, big.NewInt(1))
			}
			sqrtA := allocation.TickToSqrtPriceX96(posPlan.TickLower)
			sqrtB := allocation.TickToSqrtPriceX96(posPlan.TickUpper)
			liquidityToMint = allocation.GetLiquidityForAmounts(effectiveSqrtPriceX96, sqrtA, sqrtB, available0, available1)
			amount0Required = allocation.GetAmount0ForLiquidity(effectiveSqrtPriceX96, sqrtA, sqrtB, liquidityToMint)
			amount1Required = allocation.GetAmount1ForLiquidity(effectiveSqrtPriceX96, sqrtA, sqrtB, liquidityToMint)

			amount0Max.Set(amount0Required)
			amount1Max.Set(amount1Required)
			if amount0Max.Sign() > 0 {
				amount0Max.Add(amount0Max, big.NewInt(1))
			}
			if amount1Max.Sign() > 0 {
				amount1Max.Add(amount1Max, big.NewInt(1))
			}

			s.logger.Info("adjusted liquidity to fit vault balance",
				slog.Int("index", i),
				slog.String("original_liquidity", posPlan.Liquidity.String()),
				slog.String("adjusted_liquidity", liquidityToMint.String()))
		}

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

	}

	// Final vault balance after all mints.
	final0, _ := s.getTokenBalance(ctx, token0, vaultClient.Address())
	final1, _ := s.getTokenBalance(ctx, token1, vaultClient.Address())
	s.logger.Info("rebalance complete — vault balances",
		slog.String("token0_remaining", final0.String()),
		slog.String("token1_remaining", final1.String()))

	return nil
}
