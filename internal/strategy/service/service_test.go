package service

import (
	"math/big"
	"testing"

	"remora/internal/coverage"
	"remora/internal/liquidity"
)

func TestToAllocationBins_Empty(t *testing.T) {
	result := toAllocationBins(nil, 0, -1000, 1000)
	if len(result) != 0 {
		t.Errorf("expected empty, got %d", len(result))
	}
}

func TestToAllocationBins_AllWithinRange(t *testing.T) {
	bins := []liquidity.Bin{
		{TickLower: 0, TickUpper: 100, ActiveLiquidity: big.NewInt(500)},
		{TickLower: 100, TickUpper: 200, ActiveLiquidity: big.NewInt(600)},
		{TickLower: 200, TickUpper: 300, ActiveLiquidity: big.NewInt(700)},
	}

	result := toAllocationBins(bins, 150, 0, 300)

	if len(result) != 3 {
		t.Fatalf("expected 3 bins, got %d", len(result))
	}
}

func TestToAllocationBins_FiltersLowerBound(t *testing.T) {
	bins := []liquidity.Bin{
		{TickLower: -200, TickUpper: -100, ActiveLiquidity: big.NewInt(100)}, // below allowed
		{TickLower: -100, TickUpper: 0, ActiveLiquidity: big.NewInt(200)},
		{TickLower: 0, TickUpper: 100, ActiveLiquidity: big.NewInt(300)},
	}

	result := toAllocationBins(bins, 50, -100, 100)

	if len(result) != 2 {
		t.Fatalf("expected 2 bins after filtering lower bound, got %d", len(result))
	}
	if result[0].TickLower != -100 {
		t.Errorf("first bin TickLower: expected -100, got %d", result[0].TickLower)
	}
}

func TestToAllocationBins_FiltersUpperBound(t *testing.T) {
	bins := []liquidity.Bin{
		{TickLower: 0, TickUpper: 100, ActiveLiquidity: big.NewInt(100)},
		{TickLower: 100, TickUpper: 200, ActiveLiquidity: big.NewInt(200)},
		{TickLower: 200, TickUpper: 300, ActiveLiquidity: big.NewInt(300)}, // above allowed
	}

	result := toAllocationBins(bins, 50, 0, 200)

	if len(result) != 2 {
		t.Fatalf("expected 2 bins after filtering upper bound, got %d", len(result))
	}
	if result[1].TickUpper != 200 {
		t.Errorf("last bin TickUpper: expected 200, got %d", result[1].TickUpper)
	}
}

func TestToAllocationBins_FiltersBothBounds(t *testing.T) {
	bins := []liquidity.Bin{
		{TickLower: -500, TickUpper: -400, ActiveLiquidity: big.NewInt(10)},  // out
		{TickLower: -100, TickUpper: 0, ActiveLiquidity: big.NewInt(100)},    // in
		{TickLower: 0, TickUpper: 100, ActiveLiquidity: big.NewInt(200)},     // in
		{TickLower: 100, TickUpper: 200, ActiveLiquidity: big.NewInt(300)},   // in
		{TickLower: 400, TickUpper: 500, ActiveLiquidity: big.NewInt(10)},    // out
	}

	result := toAllocationBins(bins, 50, -100, 200)

	if len(result) != 3 {
		t.Fatalf("expected 3 bins, got %d", len(result))
	}
	if result[0].TickLower != -100 || result[2].TickUpper != 200 {
		t.Errorf("unexpected range: [%d, %d]", result[0].TickLower, result[2].TickUpper)
	}
}

func TestToAllocationBins_AllFiltered(t *testing.T) {
	bins := []liquidity.Bin{
		{TickLower: -500, TickUpper: -400, ActiveLiquidity: big.NewInt(100)},
		{TickLower: 400, TickUpper: 500, ActiveLiquidity: big.NewInt(200)},
	}

	result := toAllocationBins(bins, 0, -100, 100)

	if len(result) != 0 {
		t.Errorf("expected 0 bins when all are out of range, got %d", len(result))
	}
}

func TestToAllocationBins_IsCurrent(t *testing.T) {
	bins := []liquidity.Bin{
		{TickLower: 0, TickUpper: 100, ActiveLiquidity: big.NewInt(100)},
		{TickLower: 100, TickUpper: 200, ActiveLiquidity: big.NewInt(200)},
		{TickLower: 200, TickUpper: 300, ActiveLiquidity: big.NewInt(300)},
	}

	result := toAllocationBins(bins, 150, 0, 300)

	if result[0].IsCurrent {
		t.Error("bin [0,100) should not be current for tick 150")
	}
	if !result[1].IsCurrent {
		t.Error("bin [100,200) should be current for tick 150")
	}
	if result[2].IsCurrent {
		t.Error("bin [200,300) should not be current for tick 150")
	}
}

func TestToAllocationBins_LiquidityPreserved(t *testing.T) {
	bins := []liquidity.Bin{
		{TickLower: 0, TickUpper: 100, ActiveLiquidity: big.NewInt(12345)},
	}

	result := toAllocationBins(bins, 50, 0, 100)

	if result[0].Liquidity.Cmp(big.NewInt(12345)) != 0 {
		t.Errorf("liquidity not preserved: expected 12345, got %s", result[0].Liquidity)
	}
}

func TestToAllocationBins_EdgeExactBoundary(t *testing.T) {
	bins := []liquidity.Bin{
		{TickLower: -100, TickUpper: 0, ActiveLiquidity: big.NewInt(100)},   // exactly at lower bound
		{TickLower: 0, TickUpper: 100, ActiveLiquidity: big.NewInt(200)},    // middle
		{TickLower: 100, TickUpper: 200, ActiveLiquidity: big.NewInt(300)},  // exactly at upper bound
	}

	result := toAllocationBins(bins, 50, -100, 200)

	// All bins are exactly within bounds
	if len(result) != 3 {
		t.Fatalf("expected 3 bins at exact boundaries, got %d", len(result))
	}
}

func TestToAllocationBins_PartialOverlapFiltered(t *testing.T) {
	// A bin that straddles the allowed lower bound should be filtered
	bins := []liquidity.Bin{
		{TickLower: -150, TickUpper: -50, ActiveLiquidity: big.NewInt(100)}, // straddles: lower < allowedLower
		{TickLower: -100, TickUpper: 0, ActiveLiquidity: big.NewInt(200)},   // fully inside
	}

	result := toAllocationBins(bins, -50, -100, 100)

	if len(result) != 1 {
		t.Fatalf("expected 1 bin (partial overlap should be filtered), got %d", len(result))
	}
	if result[0].TickLower != -100 {
		t.Errorf("expected bin at -100, got %d", result[0].TickLower)
	}
}

// Verify output type is correct coverage.Bin
func TestToAllocationBins_OutputType(t *testing.T) {
	bins := []liquidity.Bin{
		{TickLower: 0, TickUpper: 100, ActiveLiquidity: big.NewInt(500)},
	}

	result := toAllocationBins(bins, 50, 0, 100)

	// Type assertion - compile-time check
	var _ []coverage.Bin = result
	if len(result) != 1 {
		t.Fatal("expected 1 bin")
	}
}
