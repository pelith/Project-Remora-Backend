package allocation

import (
	"math"
	"math/big"
)

var (
	// Q96 = 2^96
	Q96 = new(big.Int).Exp(big.NewInt(2), big.NewInt(96), nil)
	// Q192 = 2^192
	Q192 = new(big.Int).Exp(big.NewInt(2), big.NewInt(192), nil)
)

// TickToSqrtPriceX96 converts a tick to sqrtPriceX96
// sqrtPriceX96 = sqrt(1.0001^tick) * 2^96
func TickToSqrtPriceX96(tick int) *big.Int {
	sqrtPriceX96, err := GetSqrtRatioAtTick(tick)
	if err != nil {
		return big.NewInt(0)
	}
	return sqrtPriceX96
}

// SqrtPriceX96ToTick converts sqrtPriceX96 to tick
// tick = log(sqrtPriceX96^2 / 2^192) / log(1.0001)
func SqrtPriceX96ToTick(sqrtPriceX96 *big.Int) int {
	tick, err := GetTickAtSqrtRatio(sqrtPriceX96)
	if err != nil {
		return 0
	}
	return tick
}

// GetAmount0ForLiquidity calculates amount0 given liquidity and price range
// When current price is above the range: amount0 = 0
// When current price is below the range: amount0 = L * (1/sqrtPriceA - 1/sqrtPriceB)
// When current price is in range: amount0 = L * (1/sqrtPrice - 1/sqrtPriceB)
func GetAmount0ForLiquidity(sqrtPriceX96, sqrtPriceAX96, sqrtPriceBX96, liquidity *big.Int) *big.Int {
	// Ensure sqrtPriceA < sqrtPriceB
	if sqrtPriceAX96.Cmp(sqrtPriceBX96) > 0 {
		sqrtPriceAX96, sqrtPriceBX96 = sqrtPriceBX96, sqrtPriceAX96
	}

	if sqrtPriceX96.Cmp(sqrtPriceAX96) <= 0 {
		// Current price below range
		return calcAmount0(sqrtPriceAX96, sqrtPriceBX96, liquidity)
	} else if sqrtPriceX96.Cmp(sqrtPriceBX96) >= 0 {
		// Current price above range
		return big.NewInt(0)
	} else {
		// Current price in range
		return calcAmount0(sqrtPriceX96, sqrtPriceBX96, liquidity)
	}
}

// GetAmount1ForLiquidity calculates amount1 given liquidity and price range
// When current price is below the range: amount1 = 0
// When current price is above the range: amount1 = L * (sqrtPriceB - sqrtPriceA)
// When current price is in range: amount1 = L * (sqrtPrice - sqrtPriceA)
func GetAmount1ForLiquidity(sqrtPriceX96, sqrtPriceAX96, sqrtPriceBX96, liquidity *big.Int) *big.Int {
	// Ensure sqrtPriceA < sqrtPriceB
	if sqrtPriceAX96.Cmp(sqrtPriceBX96) > 0 {
		sqrtPriceAX96, sqrtPriceBX96 = sqrtPriceBX96, sqrtPriceAX96
	}

	if sqrtPriceX96.Cmp(sqrtPriceAX96) <= 0 {
		// Current price below range
		return big.NewInt(0)
	} else if sqrtPriceX96.Cmp(sqrtPriceBX96) >= 0 {
		// Current price above range
		return calcAmount1(sqrtPriceAX96, sqrtPriceBX96, liquidity)
	} else {
		// Current price in range
		return calcAmount1(sqrtPriceAX96, sqrtPriceX96, liquidity)
	}
}

// calcAmount0 calculates amount0 = L * Q96 * (sqrtPriceB - sqrtPriceA) / (sqrtPriceA * sqrtPriceB)
// Rearranged to avoid precision loss: amount0 = L * Q96 * (1/sqrtPriceA - 1/sqrtPriceB)
func calcAmount0(sqrtPriceAX96, sqrtPriceBX96, liquidity *big.Int) *big.Int {
	// amount0 = L * (sqrtPriceB - sqrtPriceA) * Q96 / (sqrtPriceA * sqrtPriceB / Q96)
	// Simplified: amount0 = L * (sqrtPriceB - sqrtPriceA) * Q96^2 / (sqrtPriceA * sqrtPriceB)

	diff := new(big.Int).Sub(sqrtPriceBX96, sqrtPriceAX96)
	numerator := new(big.Int).Mul(liquidity, diff)
	numerator.Mul(numerator, Q96)

	denominator := new(big.Int).Mul(sqrtPriceAX96, sqrtPriceBX96)
	// denominator.Div(denominator, Q96) <-- Removed this line

	return new(big.Int).Div(numerator, denominator)
}

// calcAmount1 calculates amount1 = L * (sqrtPriceB - sqrtPriceA) / Q96
func calcAmount1(sqrtPriceAX96, sqrtPriceBX96, liquidity *big.Int) *big.Int {
	diff := new(big.Int).Sub(sqrtPriceBX96, sqrtPriceAX96)
	result := new(big.Int).Mul(liquidity, diff)
	return result.Div(result, Q96)
}

// GetLiquidityForAmounts calculates liquidity given amounts and price range
// Returns the maximum liquidity that can be minted with the given amounts
func GetLiquidityForAmounts(sqrtPriceX96, sqrtPriceAX96, sqrtPriceBX96, amount0, amount1 *big.Int) *big.Int {
	// Ensure sqrtPriceA < sqrtPriceB
	if sqrtPriceAX96.Cmp(sqrtPriceBX96) > 0 {
		sqrtPriceAX96, sqrtPriceBX96 = sqrtPriceBX96, sqrtPriceAX96
	}

	if sqrtPriceX96.Cmp(sqrtPriceAX96) <= 0 {
		// Current price below range: only amount0 matters
		return getLiquidityForAmount0(sqrtPriceAX96, sqrtPriceBX96, amount0)
	} else if sqrtPriceX96.Cmp(sqrtPriceBX96) >= 0 {
		// Current price above range: only amount1 matters
		return getLiquidityForAmount1(sqrtPriceAX96, sqrtPriceBX96, amount1)
	} else {
		// Current price in range: take the minimum of the two
		l0 := getLiquidityForAmount0(sqrtPriceX96, sqrtPriceBX96, amount0)
		l1 := getLiquidityForAmount1(sqrtPriceAX96, sqrtPriceX96, amount1)
		if l0.Cmp(l1) < 0 {
			return l0
		}
		return l1
	}
}

// getLiquidityForAmount0 calculates L given amount0
// L = amount0 * sqrtPriceA * sqrtPriceB / (Q96 * (sqrtPriceB - sqrtPriceA))
func getLiquidityForAmount0(sqrtPriceAX96, sqrtPriceBX96, amount0 *big.Int) *big.Int {
	if amount0.Sign() == 0 {
		return big.NewInt(0)
	}
	diff := new(big.Int).Sub(sqrtPriceBX96, sqrtPriceAX96)
	if diff.Sign() == 0 {
		return big.NewInt(0)
	}

	// L = amount0 * sqrtPriceA * sqrtPriceB / (Q96 * (sqrtPriceB - sqrtPriceA))
	// Rearranged: L = amount0 * (sqrtPriceA * sqrtPriceB / Q96) / (sqrtPriceB - sqrtPriceA)
	product := new(big.Int).Mul(sqrtPriceAX96, sqrtPriceBX96)
	product.Div(product, Q96)

	numerator := new(big.Int).Mul(amount0, product)
	return numerator.Div(numerator, diff)
}

// getLiquidityForAmount1 calculates L given amount1
// L = amount1 * Q96 / (sqrtPriceB - sqrtPriceA)
func getLiquidityForAmount1(sqrtPriceAX96, sqrtPriceBX96, amount1 *big.Int) *big.Int {
	if amount1.Sign() == 0 {
		return big.NewInt(0)
	}
	diff := new(big.Int).Sub(sqrtPriceBX96, sqrtPriceAX96)
	if diff.Sign() == 0 {
		return big.NewInt(0)
	}

	numerator := new(big.Int).Mul(amount1, Q96)
	return numerator.Div(numerator, diff)
}

// SqrtPriceX96ToPrice converts sqrtPriceX96 to human-readable price (float64)
func SqrtPriceX96ToPrice(sqrtPriceX96 *big.Int) float64 {
	sqrtPriceFloat := new(big.Float).SetInt(sqrtPriceX96)
	q96Float := new(big.Float).SetInt(Q96)
	sqrtPriceFloat.Quo(sqrtPriceFloat, q96Float)

	sqrtPrice, _ := sqrtPriceFloat.Float64()
	return sqrtPrice * sqrtPrice
}

// PriceToSqrtPriceX96 converts human-readable price to sqrtPriceX96
func PriceToSqrtPriceX96(price float64) *big.Int {
	sqrtPrice := math.Sqrt(price)
	sqrtPriceX96 := new(big.Float).SetFloat64(sqrtPrice)
	q96Float := new(big.Float).SetInt(Q96)
	sqrtPriceX96.Mul(sqrtPriceX96, q96Float)

	result := new(big.Int)
	sqrtPriceX96.Int(result)
	return result
}
