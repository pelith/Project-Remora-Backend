package allocation

import (
	"math/big"
	"testing"

	"remora/internal/coverage"
)

func TestNormalizeWeights(t *testing.T) {
	segments := []coverage.Segment{
		{LiquidityAdded: big.NewInt(100)},
		{LiquidityAdded: big.NewInt(200)},
		{LiquidityAdded: big.NewInt(300)},
	}

	weights := normalizeWeights(segments)

	if len(weights) != 3 {
		t.Fatalf("expected 3 weights, got %d", len(weights))
	}

	// Check sum = 1
	sum := weights[0] + weights[1] + weights[2]
	if sum < 0.999 || sum > 1.001 {
		t.Errorf("weights sum: expected 1, got %f", sum)
	}

	// Check ratios
	if weights[0] < 0.166 || weights[0] > 0.167 { // 100/600
		t.Errorf("weight[0]: expected ~0.167, got %f", weights[0])
	}
	if weights[1] < 0.333 || weights[1] > 0.334 { // 200/600
		t.Errorf("weight[1]: expected ~0.333, got %f", weights[1])
	}
	if weights[2] < 0.499 || weights[2] > 0.501 { // 300/600
		t.Errorf("weight[2]: expected ~0.5, got %f", weights[2])
	}
}

func TestNormalizeWeights_ZeroTotal(t *testing.T) {
	segments := []coverage.Segment{
		{LiquidityAdded: big.NewInt(0)},
		{LiquidityAdded: big.NewInt(0)},
	}

	weights := normalizeWeights(segments)

	// Should be equal weights
	if weights[0] != 0.5 || weights[1] != 0.5 {
		t.Errorf("zero total: expected equal weights [0.5, 0.5], got %v", weights)
	}
}

// Use simple setup: both tokens 18 decimals, price = 1 (tick 0)
func simplePool() PoolState {
	return PoolState{
		SqrtPriceX96:   TickToSqrtPriceX96(0), // price = 1
		CurrentTick:    0,
		Token0Decimals: 18,
		Token1Decimals: 18,
	}
}

func TestCalculateTotalValue_SimplePool(t *testing.T) {
	pool := simplePool()

	funds := UserFunds{
		Amount0: big.NewInt(1e18), // 1 token0
		Amount1: big.NewInt(2e18), // 2 token1
	}

	totalValue := calculateTotalValue(funds, pool)

	// price = 1, so totalValue = 1 + 2 = 3 token1 = 3e18
	expected := big.NewInt(3e18)

	// Allow 1% tolerance
	diff := new(big.Int).Sub(totalValue, expected)
	tolerance := big.NewInt(3e16) // 1%
	if diff.CmpAbs(tolerance) > 0 {
		t.Errorf("totalValue: expected ~%s, got %s (diff=%s)", expected.String(), totalValue.String(), diff.String())
	}
}

func TestCalculateSwapNeeded_NeedToken0(t *testing.T) {
	pool := simplePool()

	funds := UserFunds{
		Amount0: big.NewInt(0),    // no token0
		Amount1: big.NewInt(2e18), // 2 token1
	}

	totalNeeded0 := big.NewInt(1e18) // need 1 token0
	totalNeeded1 := big.NewInt(0)    // need 0 token1

	swapAmount, token0To1 := calculateSwapNeeded(totalNeeded0, totalNeeded1, funds, pool)

	// Should swap token1 → token0, so token0To1 = false
	if token0To1 {
		t.Error("expected token0To1=false (swap token1→token0)")
	}

	// swapAmount should be ~1e18 (price = 1)
	expected := big.NewInt(1e18)
	tolerance := big.NewInt(1e16) // 1%
	diff := new(big.Int).Sub(swapAmount, expected)
	if diff.CmpAbs(tolerance) > 0 {
		t.Errorf("swapAmount: expected ~%s, got %s", expected.String(), swapAmount.String())
	}
}

func TestCalculateSwapNeeded_NeedToken1(t *testing.T) {
	pool := simplePool()

	funds := UserFunds{
		Amount0: big.NewInt(2e18), // 2 token0
		Amount1: big.NewInt(0),    // no token1
	}

	totalNeeded0 := big.NewInt(0)    // need 0 token0
	totalNeeded1 := big.NewInt(1e18) // need 1 token1

	swapAmount, token0To1 := calculateSwapNeeded(totalNeeded0, totalNeeded1, funds, pool)

	// Should swap token0 → token1, so token0To1 = true
	if !token0To1 {
		t.Error("expected token0To1=true (swap token0→token1)")
	}

	// swapAmount should be ~1e18 (price = 1)
	expected := big.NewInt(1e18)
	tolerance := big.NewInt(1e16) // 1%
	diff := new(big.Int).Sub(swapAmount, expected)
	if diff.CmpAbs(tolerance) > 0 {
		t.Errorf("swapAmount: expected ~%s, got %s", expected.String(), swapAmount.String())
	}
}

func TestCalculateSwapNeeded_NoSwap(t *testing.T) {
	pool := simplePool()

	funds := UserFunds{
		Amount0: big.NewInt(1e18),
		Amount1: big.NewInt(1e18),
	}

	totalNeeded0 := big.NewInt(1e18)
	totalNeeded1 := big.NewInt(1e18)

	swapAmount, _ := calculateSwapNeeded(totalNeeded0, totalNeeded1, funds, pool)

	if swapAmount.Sign() != 0 {
		t.Errorf("expected no swap, got %s", swapAmount.String())
	}
}

func TestAllocate_Integration(t *testing.T) {
	pool := simplePool()

	funds := UserFunds{
		Amount0: big.NewInt(1e18), // 1 token0
		Amount1: big.NewInt(1e18), // 1 token1
	}

	// Three segments around current price (tick 0)
	segments := []coverage.Segment{
		{TickLower: -1000, TickUpper: -500, LiquidityAdded: big.NewInt(100)}, // below
		{TickLower: -500, TickUpper: 500, LiquidityAdded: big.NewInt(200)},   // in range
		{TickLower: 500, TickUpper: 1000, LiquidityAdded: big.NewInt(100)},   // above
	}

	result, err := Allocate(segments, funds, pool, true)
	if err != nil {
		t.Fatalf("Allocate error: %v", err)
	}

	// Should have 3 positions
	if len(result.Positions) != 3 {
		t.Fatalf("expected 3 positions, got %d", len(result.Positions))
	}

	// All positions should have positive liquidity
	for i, pos := range result.Positions {
		if pos.Liquidity.Sign() <= 0 {
			t.Errorf("position %d: liquidity should be > 0, got %s", i, pos.Liquidity.String())
		}
		t.Logf("position %d: L=%s, amt0=%s, amt1=%s, weight=%.3f",
			i, pos.Liquidity.String(), pos.Amount0.String(), pos.Amount1.String(), pos.Weight)
	}

	// Weights should match input ratios (100:200:100 = 0.25:0.5:0.25)
	if result.Positions[0].Weight < 0.24 || result.Positions[0].Weight > 0.26 {
		t.Errorf("position 0 weight: expected ~0.25, got %f", result.Positions[0].Weight)
	}
	if result.Positions[1].Weight < 0.49 || result.Positions[1].Weight > 0.51 {
		t.Errorf("position 1 weight: expected ~0.5, got %f", result.Positions[1].Weight)
	}

	t.Logf("TotalAmount0: %s", result.TotalAmount0.String())
	t.Logf("TotalAmount1: %s", result.TotalAmount1.String())
	t.Logf("SwapAmount: %s, Token0To1: %v", result.SwapAmount.String(), result.SwapToken0To1)
}

func TestCalculateLiquidityFromValue_InRange_Verify(t *testing.T) {
	// Setup pool at Tick 0 (Price = 1 raw)
	sqrtP := TickToSqrtPriceX96(0)
	
	// Range [-887220, 887220] is max, let's use [-600, 600]
	tickLower := -600
	tickUpper := 600
	sqrtA := TickToSqrtPriceX96(tickLower)
	sqrtB := TickToSqrtPriceX96(tickUpper)

	testCases := []struct {
		name string
		d0   int
		d1   int
		val  *big.Int // input value in token1 units
	}{
		{
			name: "Equal Decimals 18/18",
			d0:   18,
			d1:   18,
			val:  new(big.Int).Mul(big.NewInt(10), big.NewInt(1e18)), // 10 token1
		},
		{
			name: "Diff Decimals 18/6 (ETH/USDC style)",
			d0:   18,
			d1:   6,
			val:  big.NewInt(10e6), // 10 USDC (small enough for int64)
		},
		{
			name: "Diff Decimals 6/18",
			d0:   6,
			d1:   18,
			val:  new(big.Int).Mul(big.NewInt(10), big.NewInt(1e18)), // 10 token1
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			pool := PoolState{
				SqrtPriceX96:   sqrtP,
				CurrentTick:    0,
				Token0Decimals: tc.d0,
				Token1Decimals: tc.d1,
			}

			// Calculate L
			L := calculateLiquidityFromValue(tc.val, sqrtP, sqrtA, sqrtB, pool)
			
			// Verify by calculating amounts back
			a0 := GetAmount0ForLiquidity(sqrtP, sqrtA, sqrtB, L)
			a1 := GetAmount1ForLiquidity(sqrtP, sqrtA, sqrtB, L)

			// Recalculate value = a0 * price + a1
			funds := UserFunds{Amount0: a0, Amount1: a1}
			recalcVal := calculateTotalValue(funds, pool)

			// Check difference
			diff := new(big.Int).Sub(recalcVal, tc.val)
			tolerance := new(big.Int).Div(tc.val, big.NewInt(100)) // 1% tolerance (should be much tighter actually)

			if diff.CmpAbs(tolerance) > 0 {
				t.Errorf("Value mismatch: Input: %s, Recalc: %s, Diff: %s", 
					tc.val.String(), recalcVal.String(), diff.String())
				t.Logf("L: %s", L.String())
				t.Logf("a0: %s", a0.String())
				t.Logf("a1: %s", a1.String())
			}
		})
	}
}
