package allocation

import (
	"math/big"

	"remora/internal/coverage"
)

// Allocate distributes user funds across segments based on weights
// Main entry point for allocation calculation
func Allocate(segments []coverage.Segment, funds UserFunds, pool PoolState, swapAllowed bool) (*AllocationResult, error) {
	if len(segments) == 0 {
		return &AllocationResult{
			Positions:    []PositionPlan{},
			TotalAmount0: big.NewInt(0),
			TotalAmount1: big.NewInt(0),
			SwapAmount:   big.NewInt(0),
		}, nil
	}

	// 1. Calculate weights from LiquidityAdded
	weights := normalizeWeights(segments)

	if swapAllowed {
		return allocateWithSwap(segments, weights, funds, pool)
	}

	return allocateWithoutSwap(segments, weights, funds, pool)
}

// allocateWithSwap implements the original logic: total value -> weights -> calculate swap
func allocateWithSwap(segments []coverage.Segment, weights []float64, funds UserFunds, pool PoolState) (*AllocationResult, error) {
	totalValue := calculateTotalValue(funds, pool)

	positions := make([]PositionPlan, len(segments))
	totalAmount0 := big.NewInt(0)
	totalAmount1 := big.NewInt(0)

	for i, seg := range segments {
		totalValueFloat := new(big.Float).SetInt(totalValue)
		allocatedFloat := new(big.Float).Mul(totalValueFloat, big.NewFloat(weights[i]))
		allocatedValue, _ := allocatedFloat.Int(nil)

		pos := calculatePosition(allocatedValue, int(seg.TickLower), int(seg.TickUpper), weights[i], pool)
		positions[i] = *pos

		totalAmount0.Add(totalAmount0, pos.Amount0)
		totalAmount1.Add(totalAmount1, pos.Amount1)
	}

	swapAmount, token0To1 := calculateSwapNeeded(totalAmount0, totalAmount1, funds, pool)

	return &AllocationResult{
		Positions:     positions,
		TotalAmount0:  totalAmount0,
		TotalAmount1:  totalAmount1,
		SwapAmount:    swapAmount,
		SwapToken0To1: token0To1,
	}, nil
}

// allocateWithoutSwap implements "Fit-to-Balance" logic: distribute existing tokens by weight
func allocateWithoutSwap(segments []coverage.Segment, weights []float64, funds UserFunds, pool PoolState) (*AllocationResult, error) {
	positions := make([]PositionPlan, len(segments))
	totalAmount0 := big.NewInt(0)
	totalAmount1 := big.NewInt(0)

	for i, seg := range segments {
		// Distribute current balances by weight
		budget0Float := new(big.Float).Mul(new(big.Float).SetInt(funds.Amount0), big.NewFloat(weights[i]))
		budget1Float := new(big.Float).Mul(new(big.Float).SetInt(funds.Amount1), big.NewFloat(weights[i]))
		
		budget0, _ := budget0Float.Int(nil)
		budget1, _ := budget1Float.Int(nil)

		sqrtPriceAX96 := TickToSqrtPriceX96(int(seg.TickLower))
		sqrtPriceBX96 := TickToSqrtPriceX96(int(seg.TickUpper))

		// Calculate max liquidity L that fits into both budgets
		liquidity := GetLiquidityForAmounts(pool.SqrtPriceX96, sqrtPriceAX96, sqrtPriceBX96, budget0, budget1)

		// Calculate actual amounts needed for this L
		amt0 := GetAmount0ForLiquidity(pool.SqrtPriceX96, sqrtPriceAX96, sqrtPriceBX96, liquidity)
		amt1 := GetAmount1ForLiquidity(pool.SqrtPriceX96, sqrtPriceAX96, sqrtPriceBX96, liquidity)

		positions[i] = PositionPlan{
			TickLower: int(seg.TickLower),
			TickUpper: int(seg.TickUpper),
			Liquidity: liquidity,
			Amount0:   amt0,
			Amount1:   amt1,
			Weight:    weights[i],
		}

		totalAmount0.Add(totalAmount0, amt0)
		totalAmount1.Add(totalAmount1, amt1)
	}

	return &AllocationResult{
		Positions:     positions,
		TotalAmount0:  totalAmount0,
		TotalAmount1:  totalAmount1,
		SwapAmount:    big.NewInt(0), // No swap allowed
		SwapToken0To1: false,
	}, nil
}

// normalizeWeights converts LiquidityAdded to normalized weights (sum = 1)
func normalizeWeights(segments []coverage.Segment) []float64 {
	var total float64
	liquidityFloats := make([]float64, len(segments))

	for i, seg := range segments {
		if seg.LiquidityAdded != nil {
			f, _ := new(big.Float).SetInt(seg.LiquidityAdded).Float64()
			liquidityFloats[i] = f
			total += f
		}
	}

	weights := make([]float64, len(segments))
	if total == 0 {
		// Equal weights if no liquidity info
		for i := range weights {
			weights[i] = 1.0 / float64(len(segments))
		}
		return weights
	}

	for i := range segments {
		weights[i] = liquidityFloats[i] / total
	}
	return weights
}

// calculateTotalValue computes total value in token1 units
// totalValue = amount0 * price + amount1
func calculateTotalValue(funds UserFunds, pool PoolState) *big.Int {
	// value = amount0 * price + amount1
	// price = sqrtPriceX96^2 / Q192
	// value = amount0 * sqrtPriceX96^2 / Q192 + amount1

	sqrtPriceSquared := new(big.Int).Mul(pool.SqrtPriceX96, pool.SqrtPriceX96)
	amount0Value := new(big.Int).Mul(funds.Amount0, sqrtPriceSquared)
	amount0Value.Div(amount0Value, Q192)

	totalValue := new(big.Int).Add(amount0Value, funds.Amount1)
	return totalValue
}

// calculatePosition computes L, amount0, amount1 for a single position
// given allocated value and price range
func calculatePosition(allocatedValue *big.Int, tickLower, tickUpper int, weight float64, pool PoolState) *PositionPlan {
	sqrtPriceAX96 := TickToSqrtPriceX96(tickLower)
	sqrtPriceBX96 := TickToSqrtPriceX96(tickUpper)

	// Calculate L from allocated value
	liquidity := calculateLiquidityFromValue(
		allocatedValue,
		pool.SqrtPriceX96,
		sqrtPriceAX96,
		sqrtPriceBX96,
		pool,
	)

	// Calculate amount0, amount1 from L
	amount0 := GetAmount0ForLiquidity(pool.SqrtPriceX96, sqrtPriceAX96, sqrtPriceBX96, liquidity)
	amount1 := GetAmount1ForLiquidity(pool.SqrtPriceX96, sqrtPriceAX96, sqrtPriceBX96, liquidity)

	return &PositionPlan{
		TickLower: tickLower,
		TickUpper: tickUpper,
		Liquidity: liquidity,
		Amount0:   amount0,
		Amount1:   amount1,
		Weight:    weight,
	}
}

// calculateLiquidityFromValue solves for L given allocated value (in token1 units)
// value = amount0 * price + amount1, where amounts are functions of L
func calculateLiquidityFromValue(value *big.Int, sqrtPriceX96, sqrtPriceAX96, sqrtPriceBX96 *big.Int, pool PoolState) *big.Int {
	// Ensure sqrtPriceA < sqrtPriceB
	if sqrtPriceAX96.Cmp(sqrtPriceBX96) > 0 {
		sqrtPriceAX96, sqrtPriceBX96 = sqrtPriceBX96, sqrtPriceAX96
	}

	if sqrtPriceX96.Cmp(sqrtPriceAX96) <= 0 {
		// Case 1: price <= priceA (below range, only token0)
		// L = value * sqrtA * sqrtB * Q96 / [price * (sqrtB - sqrtA)]
		// where price = sqrtPriceX96^2 / Q192
		//
		// Rearranged to avoid precision loss:
		// L = value * sqrtA * sqrtB * Q96 * Q192 / [sqrtPriceX96^2 * (sqrtB - sqrtA)]

		diff := new(big.Int).Sub(sqrtPriceBX96, sqrtPriceAX96)
		numerator := new(big.Int).Mul(value, sqrtPriceAX96)
		numerator.Mul(numerator, sqrtPriceBX96)
		numerator.Mul(numerator, Q96)
		numerator.Mul(numerator, Q192)

		sqrtPriceSquared := new(big.Int).Mul(sqrtPriceX96, sqrtPriceX96)
		denominator := new(big.Int).Mul(sqrtPriceSquared, diff)

		return new(big.Int).Div(numerator, denominator)

	} else if sqrtPriceX96.Cmp(sqrtPriceBX96) >= 0 {
		// Case 2: price >= priceB (above range, only token1)
		// L = value * Q96 / (sqrtB - sqrtA)

		diff := new(big.Int).Sub(sqrtPriceBX96, sqrtPriceAX96)
		numerator := new(big.Int).Mul(value, Q96)
		return new(big.Int).Div(numerator, diff)

	} else {
		// Case 3: priceA < price < priceB (in range, both tokens)
		// value = L * [(1/sqrtP - 1/sqrtB) * price + (sqrtP - sqrtA)]
		//
		// Let's compute the coefficient in Q96 format:
		// coef = (1/sqrtP - 1/sqrtB) * price + (sqrtP - sqrtA)
		//      = (sqrtB - sqrtP) * price / (sqrtP * sqrtB) + (sqrtP - sqrtA)

		// Term 1: (sqrtB - sqrtP) * price / (sqrtP * sqrtB)
		// In raw units, then we'll combine
		diffBP := new(big.Int).Sub(sqrtPriceBX96, sqrtPriceX96)
		sqrtPriceSquared := new(big.Int).Mul(sqrtPriceX96, sqrtPriceX96)

		term1Num := new(big.Int).Mul(diffBP, sqrtPriceSquared)
		term1Denom := new(big.Int).Mul(sqrtPriceX96, sqrtPriceBX96)
		term1 := new(big.Int).Div(term1Num, term1Denom)

		// Term 2: (sqrtP - sqrtA) - already in Q96 format
		term2 := new(big.Int).Sub(sqrtPriceX96, sqrtPriceAX96)

		// coef = term1 + term2 (both in Q96-ish units)
		coef := new(big.Int).Add(term1, term2)

		// L = value * Q96 / coef
		numerator := new(big.Int).Mul(value, Q96)
		return new(big.Int).Div(numerator, coef)
	}
}

// calculateSwapNeeded computes the swap amount and direction
// Returns: swapAmount (in source token units), token0To1 (swap direction)
func calculateSwapNeeded(totalNeeded0, totalNeeded1 *big.Int, funds UserFunds, pool PoolState) (swapAmount *big.Int, token0To1 bool) {
	// deficit = needed - have
	deficit0 := new(big.Int).Sub(totalNeeded0, funds.Amount0)
	deficit1 := new(big.Int).Sub(totalNeeded1, funds.Amount1)

	// If deficit0 > 0: we need more token0, swap token1 → token0
	// If deficit1 > 0: we need more token1, swap token0 → token1
	// (Only one can be positive at a time in a valid allocation)

	if deficit0.Sign() > 0 {
		// Need more token0, swap token1 → token0
		// swapAmount = how much token1 to swap
		// deficit0 (token0 needed) * price = token1 amount to swap
		//
		// price = sqrtPriceX96^2 / Q192

		sqrtPriceSquared := new(big.Int).Mul(pool.SqrtPriceX96, pool.SqrtPriceX96)

		// swapAmount (token1) = deficit0 * sqrtPriceX96^2 / Q192
		swapAmount = new(big.Int).Mul(deficit0, sqrtPriceSquared)
		swapAmount.Div(swapAmount, Q192)

		return swapAmount, false // token1 → token0, so token0To1 = false

	} else if deficit1.Sign() > 0 {
		// Need more token1, swap token0 → token1
		// swapAmount = how much token0 to swap
		// deficit1 (token1 needed) / price = token0 amount to swap
		//
		// swapAmount (token0) = deficit1 * Q192 / sqrtPriceX96^2

		sqrtPriceSquared := new(big.Int).Mul(pool.SqrtPriceX96, pool.SqrtPriceX96)

		swapAmount = new(big.Int).Mul(deficit1, Q192)
		swapAmount.Div(swapAmount, sqrtPriceSquared)

		return swapAmount, true // token0 → token1, so token0To1 = true
	}

	// No swap needed
	return big.NewInt(0), false
}
