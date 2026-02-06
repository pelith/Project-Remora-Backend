package allocation

import (
	"math/big"
	"testing"
)

func TestTickToSqrtPriceX96_RoundTrip(t *testing.T) {
	tests := []int{0, 100, -100, 1000, -1000, 50000, -50000}

	for _, tick := range tests {
		sqrtPriceX96 := TickToSqrtPriceX96(tick)
		gotTick := SqrtPriceX96ToTick(sqrtPriceX96)

		// Allow ±1 tick difference due to rounding
		diff := gotTick - tick
		if diff < -1 || diff > 1 {
			t.Errorf("tick %d: got %d after round-trip, diff=%d", tick, gotTick, diff)
		}
	}
}

func TestTickToSqrtPriceX96_KnownValues(t *testing.T) {
	// tick 0 → price = 1 → sqrtPrice = 1 → sqrtPriceX96 = Q96
	sqrtPriceX96 := TickToSqrtPriceX96(0)
	if sqrtPriceX96.Cmp(Q96) != 0 {
		t.Errorf("tick 0: expected %s, got %s", Q96.String(), sqrtPriceX96.String())
	}
}

func TestSqrtPriceX96ToPrice(t *testing.T) {
	// tick 0 → price = 1
	sqrtPriceX96 := TickToSqrtPriceX96(0)
	price := SqrtPriceX96ToPrice(sqrtPriceX96)
	if price < 0.999 || price > 1.001 {
		t.Errorf("tick 0: expected price ~1, got %f", price)
	}

	// tick 10000 → price ≈ 2.718 (e^1 ≈ 1.0001^10000)
	sqrtPriceX96 = TickToSqrtPriceX96(10000)
	price = SqrtPriceX96ToPrice(sqrtPriceX96)
	expected := 2.718
	if price < expected*0.99 || price > expected*1.01 {
		t.Errorf("tick 10000: expected price ~%f, got %f", expected, price)
	}
}

func TestGetLiquidityForAmounts_RoundTrip(t *testing.T) {
	// Setup: current price at tick 0, range from tick -1000 to 1000
	sqrtPriceX96 := TickToSqrtPriceX96(0)
	sqrtPriceAX96 := TickToSqrtPriceX96(-1000)
	sqrtPriceBX96 := TickToSqrtPriceX96(1000)

	// Given amounts
	amount0 := big.NewInt(1e18)   // 1 token0
	amount1 := big.NewInt(1000e6) // 1000 token1

	// Calculate liquidity
	liquidity := GetLiquidityForAmounts(sqrtPriceX96, sqrtPriceAX96, sqrtPriceBX96, amount0, amount1)

	// Recalculate amounts from liquidity
	gotAmount0 := GetAmount0ForLiquidity(sqrtPriceX96, sqrtPriceAX96, sqrtPriceBX96, liquidity)
	gotAmount1 := GetAmount1ForLiquidity(sqrtPriceX96, sqrtPriceAX96, sqrtPriceBX96, liquidity)

	// Should be <= original amounts (liquidity is limited by the smaller ratio)
	if gotAmount0.Cmp(amount0) > 0 {
		t.Errorf("amount0: got %s > original %s", gotAmount0.String(), amount0.String())
	}
	if gotAmount1.Cmp(amount1) > 0 {
		t.Errorf("amount1: got %s > original %s", gotAmount1.String(), amount1.String())
	}

	// At least one should be close to original
	ratio0 := new(big.Float).Quo(new(big.Float).SetInt(gotAmount0), new(big.Float).SetInt(amount0))
	ratio1 := new(big.Float).Quo(new(big.Float).SetInt(gotAmount1), new(big.Float).SetInt(amount1))
	r0, _ := ratio0.Float64()
	r1, _ := ratio1.Float64()

	if r0 < 0.99 && r1 < 0.99 {
		t.Errorf("neither amount is close to original: ratio0=%f, ratio1=%f", r0, r1)
	}
}

func TestGetAmountForLiquidity_BelowRange(t *testing.T) {
	// Current price below range: only token0 needed
	sqrtPriceX96 := TickToSqrtPriceX96(-2000)
	sqrtPriceAX96 := TickToSqrtPriceX96(-1000)
	sqrtPriceBX96 := TickToSqrtPriceX96(1000)
	liquidity := big.NewInt(1e18)

	amount0 := GetAmount0ForLiquidity(sqrtPriceX96, sqrtPriceAX96, sqrtPriceBX96, liquidity)
	amount1 := GetAmount1ForLiquidity(sqrtPriceX96, sqrtPriceAX96, sqrtPriceBX96, liquidity)

	if amount0.Sign() <= 0 {
		t.Errorf("below range: amount0 should be > 0, got %s", amount0.String())
	}
	if amount1.Sign() != 0 {
		t.Errorf("below range: amount1 should be 0, got %s", amount1.String())
	}
}

func TestGetAmountForLiquidity_AboveRange(t *testing.T) {
	// Current price above range: only token1 needed
	sqrtPriceX96 := TickToSqrtPriceX96(2000)
	sqrtPriceAX96 := TickToSqrtPriceX96(-1000)
	sqrtPriceBX96 := TickToSqrtPriceX96(1000)
	liquidity := big.NewInt(1e18)

	amount0 := GetAmount0ForLiquidity(sqrtPriceX96, sqrtPriceAX96, sqrtPriceBX96, liquidity)
	amount1 := GetAmount1ForLiquidity(sqrtPriceX96, sqrtPriceAX96, sqrtPriceBX96, liquidity)

	if amount0.Sign() != 0 {
		t.Errorf("above range: amount0 should be 0, got %s", amount0.String())
	}
	if amount1.Sign() <= 0 {
		t.Errorf("above range: amount1 should be > 0, got %s", amount1.String())
	}
}
