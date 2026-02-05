package coverage

import (
	"context"
	"math"
	"math/big"
	"testing"
)

// ─── toInternalBins ─────────────────────────────────────────────────────────

func TestToInternalBins_Empty(t *testing.T) {
	result := toInternalBins(nil)
	if len(result) != 0 {
		t.Errorf("expected empty, got %d", len(result))
	}
}

func TestToInternalBins_NilLiquidity(t *testing.T) {
	bins := []Bin{
		{TickLower: 0, TickUpper: 100, Liquidity: nil, IsCurrent: true},
	}
	result := toInternalBins(bins)

	if len(result) != 1 {
		t.Fatalf("expected 1, got %d", len(result))
	}
	if result[0].liquidity != 0 {
		t.Errorf("expected 0 for nil liquidity, got %f", result[0].liquidity)
	}
	if !result[0].isCurrent {
		t.Error("expected isCurrent=true")
	}
}

func TestToInternalBins_ConvertsBigInt(t *testing.T) {
	bins := []Bin{
		{
			TickLower:  -100,
			TickUpper:  200,
			PriceLower: 1.5,
			PriceUpper: 2.5,
			Liquidity:  big.NewInt(999),
			IsCurrent:  false,
		},
	}
	result := toInternalBins(bins)

	if result[0].tickLower != -100 {
		t.Errorf("tickLower: expected -100, got %d", result[0].tickLower)
	}
	if result[0].tickUpper != 200 {
		t.Errorf("tickUpper: expected 200, got %d", result[0].tickUpper)
	}
	if result[0].priceLower != 1.5 {
		t.Errorf("priceLower: expected 1.5, got %f", result[0].priceLower)
	}
	if result[0].priceUpper != 2.5 {
		t.Errorf("priceUpper: expected 2.5, got %f", result[0].priceUpper)
	}
	if !almostEqual(result[0].liquidity, 999, 0.01) {
		t.Errorf("liquidity: expected 999, got %f", result[0].liquidity)
	}
}

func TestToInternalBins_MultipleBins(t *testing.T) {
	bins := []Bin{
		{TickLower: 0, TickUpper: 100, Liquidity: big.NewInt(100), IsCurrent: true},
		{TickLower: 100, TickUpper: 200, Liquidity: big.NewInt(200), IsCurrent: false},
		{TickLower: 200, TickUpper: 300, Liquidity: big.NewInt(300), IsCurrent: false},
	}
	result := toInternalBins(bins)

	if len(result) != 3 {
		t.Fatalf("expected 3 bins, got %d", len(result))
	}
	if !result[0].isCurrent {
		t.Error("bin 0 should be current")
	}
	if result[1].isCurrent || result[2].isCurrent {
		t.Error("bins 1,2 should not be current")
	}
}

// ─── toSegments ─────────────────────────────────────────────────────────────

func TestToSegments_Empty(t *testing.T) {
	bins := []internalBin{{tickLower: 0, tickUpper: 100}}
	result := toSegments(bins, nil, []float64{100})

	if len(result.Segments) != 0 {
		t.Errorf("expected 0 segments, got %d", len(result.Segments))
	}
}

func TestToSegments_SingleSegment(t *testing.T) {
	bins := []internalBin{
		{tickLower: 0, tickUpper: 100, priceLower: 1.0, priceUpper: 1.5},
		{tickLower: 100, tickUpper: 200, priceLower: 1.5, priceUpper: 2.0},
		{tickLower: 200, tickUpper: 300, priceLower: 2.0, priceUpper: 2.5},
	}
	segments := []internalSegment{
		{l: 0, r: 2, h: 50, liquidityAdded: 50},
	}
	target := []float64{100, 100, 100}

	result := toSegments(bins, segments, target)

	if len(result.Segments) != 1 {
		t.Fatalf("expected 1 segment, got %d", len(result.Segments))
	}

	seg := result.Segments[0]
	if seg.TickLower != 0 {
		t.Errorf("TickLower: expected 0, got %d", seg.TickLower)
	}
	if seg.TickUpper != 300 {
		t.Errorf("TickUpper: expected 300, got %d", seg.TickUpper)
	}
	if seg.PriceLower != 1.0 {
		t.Errorf("PriceLower: expected 1.0, got %f", seg.PriceLower)
	}
	if seg.PriceUpper != 2.5 {
		t.Errorf("PriceUpper: expected 2.5, got %f", seg.PriceUpper)
	}
	if seg.LiquidityAdded.Cmp(big.NewInt(50)) != 0 {
		t.Errorf("LiquidityAdded: expected 50, got %s", seg.LiquidityAdded)
	}
}

func TestToSegments_MultipleSegments(t *testing.T) {
	bins := []internalBin{
		{tickLower: 0, tickUpper: 100, priceLower: 1.0, priceUpper: 1.1},
		{tickLower: 100, tickUpper: 200, priceLower: 1.1, priceUpper: 1.2},
		{tickLower: 200, tickUpper: 300, priceLower: 1.2, priceUpper: 1.3},
		{tickLower: 300, tickUpper: 400, priceLower: 1.3, priceUpper: 1.4},
	}
	segments := []internalSegment{
		{l: 0, r: 1, h: 80, liquidityAdded: 80},
		{l: 2, r: 3, h: 60, liquidityAdded: 60},
	}
	target := []float64{100, 100, 100, 100}

	result := toSegments(bins, segments, target)

	if len(result.Segments) != 2 {
		t.Fatalf("expected 2 segments, got %d", len(result.Segments))
	}

	// First segment: bins 0-1
	if result.Segments[0].TickLower != 0 || result.Segments[0].TickUpper != 200 {
		t.Errorf("seg0 ticks: expected [0,200], got [%d,%d]",
			result.Segments[0].TickLower, result.Segments[0].TickUpper)
	}
	// Second segment: bins 2-3
	if result.Segments[1].TickLower != 200 || result.Segments[1].TickUpper != 400 {
		t.Errorf("seg1 ticks: expected [200,400], got [%d,%d]",
			result.Segments[1].TickLower, result.Segments[1].TickUpper)
	}
}

// ─── calcMetrics ────────────────────────────────────────────────────────────

func TestCalcMetrics_PerfectCoverage(t *testing.T) {
	target := []float64{100, 200, 300}
	pred := []float64{100, 200, 300}
	m := calcMetrics(target, pred)

	if !almostEqual(m.Covered, 600, 0.01) {
		t.Errorf("Covered: expected 600, got %f", m.Covered)
	}
	if !almostEqual(m.Gap, 0, 0.01) {
		t.Errorf("Gap: expected 0, got %f", m.Gap)
	}
	if !almostEqual(m.Over, 0, 0.01) {
		t.Errorf("Over: expected 0, got %f", m.Over)
	}
}

func TestCalcMetrics_UnderCoverage(t *testing.T) {
	target := []float64{100, 200}
	pred := []float64{50, 100}
	m := calcMetrics(target, pred)

	// covered = min(100,50) + min(200,100) = 50 + 100 = 150
	if !almostEqual(m.Covered, 150, 0.01) {
		t.Errorf("Covered: expected 150, got %f", m.Covered)
	}
	// gap = (100-50) + (200-100) = 50 + 100 = 150
	if !almostEqual(m.Gap, 150, 0.01) {
		t.Errorf("Gap: expected 150, got %f", m.Gap)
	}
	if !almostEqual(m.Over, 0, 0.01) {
		t.Errorf("Over: expected 0, got %f", m.Over)
	}
}

func TestCalcMetrics_OverCoverage(t *testing.T) {
	target := []float64{100, 200}
	pred := []float64{150, 300}
	m := calcMetrics(target, pred)

	// covered = 100 + 200 = 300
	if !almostEqual(m.Covered, 300, 0.01) {
		t.Errorf("Covered: expected 300, got %f", m.Covered)
	}
	if !almostEqual(m.Gap, 0, 0.01) {
		t.Errorf("Gap: expected 0, got %f", m.Gap)
	}
	// over = (150-100) + (300-200) = 50 + 100 = 150
	if !almostEqual(m.Over, 150, 0.01) {
		t.Errorf("Over: expected 150, got %f", m.Over)
	}
}

func TestCalcMetrics_ZeroPrediction(t *testing.T) {
	target := []float64{100, 200}
	pred := []float64{0, 0}
	m := calcMetrics(target, pred)

	if !almostEqual(m.Covered, 0, 0.01) {
		t.Errorf("Covered: expected 0, got %f", m.Covered)
	}
	if !almostEqual(m.Gap, 300, 0.01) {
		t.Errorf("Gap: expected 300, got %f", m.Gap)
	}
}

func TestCalcMetrics_Invariant(t *testing.T) {
	// Invariant: covered + gap = sum(target)
	target := []float64{73, 150, 220, 10}
	pred := []float64{50, 200, 100, 0}
	m := calcMetrics(target, pred)

	totalTarget := 0.0
	for _, v := range target {
		totalTarget += v
	}

	// covered + gap should always equal total target
	if !almostEqual(m.Covered+m.Gap, totalTarget, 0.01) {
		t.Errorf("invariant broken: Covered(%f) + Gap(%f) = %f, totalTarget = %f",
			m.Covered, m.Gap, m.Covered+m.Gap, totalTarget)
	}

	// Also: covered + over = sum(pred)
	totalPred := 0.0
	for _, v := range pred {
		totalPred += v
	}
	if !almostEqual(m.Covered+m.Over, totalPred, 0.01) {
		t.Errorf("invariant broken: Covered(%f) + Over(%f) = %f, totalPred = %f",
			m.Covered, m.Over, m.Covered+m.Over, totalPred)
	}
}

// ─── Integration: Run with various configs ──────────────────────────────────

func TestRun_WithMinLiqEnabled(t *testing.T) {
	// Big peak + tiny peak
	liqs := []float64{1000, 1000, 0, 0, 1, 0}
	bins := makeBins(liqs, 100, 0)
	cfg := DefaultConfig()
	cfg.N = 5
	cfg.EnableMinLiq = true

	result := Run(context.Background(), bins, cfg)

	// The tiny segment (liquidity=1) should be filtered out
	for _, seg := range result.Segments {
		if seg.LiquidityAdded.Cmp(big.NewInt(0)) <= 0 {
			t.Errorf("found segment with non-positive liquidity: %s", seg.LiquidityAdded)
		}
	}
}

func TestRun_LookAheadZero(t *testing.T) {
	bins := makeBins([]float64{100, 200, 300}, 100, 1)
	cfg := DefaultConfig()
	cfg.LookAhead = 0

	// Should still work, just without look-ahead expansion
	result := Run(context.Background(), bins, cfg)

	if len(result.Segments) == 0 {
		t.Fatal("expected at least 1 segment even with LookAhead=0")
	}
}

func TestRun_SegmentsHaveValidTickOrder(t *testing.T) {
	liqs := []float64{50, 100, 200, 150, 80, 300, 250, 100}
	bins := makeBins(liqs, 200, 3)
	cfg := DefaultConfig()
	cfg.N = 3

	result := Run(context.Background(), bins, cfg)

	for i, seg := range result.Segments {
		if seg.TickLower >= seg.TickUpper {
			t.Errorf("segment %d: TickLower(%d) >= TickUpper(%d)", i, seg.TickLower, seg.TickUpper)
		}
		if seg.LiquidityAdded == nil || seg.LiquidityAdded.Sign() <= 0 {
			t.Errorf("segment %d: invalid LiquidityAdded %v", i, seg.LiquidityAdded)
		}
		if math.IsNaN(seg.PriceLower) || math.IsNaN(seg.PriceUpper) {
			t.Errorf("segment %d: NaN price values", i)
		}
	}
}
