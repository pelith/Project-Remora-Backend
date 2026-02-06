package coverage

import (
	"context"
	"math"
	"math/big"
	"testing"
)

// ─── helpers ────────────────────────────────────────────────────────────────

func makeBins(liquidities []float64, tickWidth int32, currentIdx int) []Bin {
	bins := make([]Bin, len(liquidities))
	for i, liq := range liquidities {
		lower := int32(i) * tickWidth
		upper := lower + tickWidth
		bins[i] = Bin{
			TickLower:  lower,
			TickUpper:  upper,
			PriceLower: float64(lower),
			PriceUpper: float64(upper),
			Liquidity:  big.NewInt(int64(liq)),
			IsCurrent:  i == currentIdx,
		}
	}
	return bins
}

func almostEqual(a, b, eps float64) bool {
	return math.Abs(a-b) < eps
}

// ─── Run ────────────────────────────────────────────────────────────────────

func TestRun_EmptyBins(t *testing.T) {
	result := Run(context.Background(), nil, DefaultConfig())
	if len(result.Segments) != 0 {
		t.Fatalf("expected 0 segments, got %d", len(result.Segments))
	}
}

func TestRun_SingleBin(t *testing.T) {
	bins := makeBins([]float64{1000}, 100, 0)
	cfg := DefaultConfig()
	cfg.N = 3

	result := Run(context.Background(), bins, cfg)

	if len(result.Segments) != 1 {
		t.Fatalf("expected 1 segment, got %d", len(result.Segments))
	}

	seg := result.Segments[0]
	if seg.TickLower != 0 || seg.TickUpper != 100 {
		t.Errorf("expected tick range [0,100], got [%d,%d]", seg.TickLower, seg.TickUpper)
	}
	if seg.LiquidityAdded.Cmp(big.NewInt(1000)) != 0 {
		t.Errorf("expected liquidity 1000, got %s", seg.LiquidityAdded)
	}
}

func TestRun_UniformLiquidity(t *testing.T) {
	// 5 bins all with same liquidity — should be covered by 1 wide segment
	bins := makeBins([]float64{500, 500, 500, 500, 500}, 100, 2)
	cfg := DefaultConfig()
	cfg.N = 1

	result := Run(context.Background(), bins, cfg)

	if len(result.Segments) == 0 {
		t.Fatal("expected at least 1 segment")
	}

	// With uniform liquidity and N=1, algorithm should produce a single segment
	if len(result.Segments) != 1 {
		t.Fatalf("expected 1 segment for uniform distribution with N=1, got %d", len(result.Segments))
	}

	seg := result.Segments[0]
	if seg.TickLower != 0 || seg.TickUpper != 500 {
		t.Errorf("expected full range [0,500], got [%d,%d]", seg.TickLower, seg.TickUpper)
	}
}

func TestRun_TwoPeaks(t *testing.T) {
	// Two distinct peaks separated by a gap — should produce 2 segments
	liqs := []float64{0, 1000, 1000, 0, 0, 0, 800, 800, 0, 0}
	bins := makeBins(liqs, 100, 1)
	cfg := DefaultConfig()
	cfg.N = 5
	cfg.LookAhead = 2
	cfg.Lambda = 50

	result := Run(context.Background(), bins, cfg)

	if len(result.Segments) < 2 {
		t.Fatalf("expected at least 2 segments for two-peak distribution, got %d", len(result.Segments))
	}
}

func TestRun_MetricsCoverage(t *testing.T) {
	bins := makeBins([]float64{100, 200, 300}, 100, 1)
	cfg := DefaultConfig()
	cfg.N = 1

	result := Run(context.Background(), bins, cfg)

	// Covered + Gap should equal total target liquidity
	totalTarget := 100.0 + 200.0 + 300.0
	coveredPlusGap := result.Metrics.Covered + result.Metrics.Gap
	if !almostEqual(coveredPlusGap, totalTarget, 1.0) {
		t.Errorf("Covered(%.1f) + Gap(%.1f) = %.1f, expected %.1f",
			result.Metrics.Covered, result.Metrics.Gap, coveredPlusGap, totalTarget)
	}
}

func TestRun_RespectsMaxSegments(t *testing.T) {
	liqs := []float64{100, 200, 300, 400, 500, 600, 700, 800, 900, 1000}
	bins := makeBins(liqs, 100, 5)
	cfg := DefaultConfig()
	cfg.N = 3

	result := Run(context.Background(), bins, cfg)

	if len(result.Segments) > cfg.N {
		t.Errorf("segments %d exceeds N=%d", len(result.Segments), cfg.N)
	}
}

func TestRun_ZeroLiquidity(t *testing.T) {
	bins := makeBins([]float64{0, 0, 0}, 100, 1)
	cfg := DefaultConfig()

	result := Run(context.Background(), bins, cfg)

	if len(result.Segments) != 0 {
		t.Errorf("expected 0 segments for zero liquidity, got %d", len(result.Segments))
	}
}

// ─── quantile ───────────────────────────────────────────────────────────────

func TestQuantile_Empty(t *testing.T) {
	if q := quantile(nil, 0.5); q != 0 {
		t.Errorf("expected 0, got %f", q)
	}
}

func TestQuantile_Single(t *testing.T) {
	if q := quantile([]float64{42}, 0.5); q != 42 {
		t.Errorf("expected 42, got %f", q)
	}
}

func TestQuantile_Boundaries(t *testing.T) {
	data := []float64{10, 20, 30, 40, 50}

	if q := quantile(data, 0); q != 10 {
		t.Errorf("q=0: expected 10, got %f", q)
	}
	if q := quantile(data, 1); q != 50 {
		t.Errorf("q=1: expected 50, got %f", q)
	}
}

func TestQuantile_Median(t *testing.T) {
	data := []float64{10, 20, 30, 40, 50}
	q := quantile(data, 0.5)
	if !almostEqual(q, 30, 0.01) {
		t.Errorf("median: expected 30, got %f", q)
	}
}

func TestQuantile_Interpolation(t *testing.T) {
	data := []float64{10, 20, 30, 40}
	// q=0.5 → index=1.5 → interpolate between 20 and 30
	q := quantile(data, 0.5)
	if !almostEqual(q, 25, 0.01) {
		t.Errorf("expected 25, got %f", q)
	}
}

func TestQuantile_UnsortedInput(t *testing.T) {
	data := []float64{50, 10, 40, 20, 30}
	q := quantile(data, 0.5)
	if !almostEqual(q, 30, 0.01) {
		t.Errorf("expected 30 (median of sorted), got %f", q)
	}
}

// ─── calcH ──────────────────────────────────────────────────────────────────

func TestCalcH_AllZero(t *testing.T) {
	gaps := []float64{0, 0, 0}
	h := calcH(gaps, 0, 2, 0.6)
	if h != 0 {
		t.Errorf("expected 0 for all-zero gaps, got %f", h)
	}
}

func TestCalcH_SingleBin(t *testing.T) {
	gaps := []float64{100}
	h := calcH(gaps, 0, 0, 0.6)
	if !almostEqual(h, 100, 0.01) {
		t.Errorf("expected 100, got %f", h)
	}
}

func TestCalcH_Subset(t *testing.T) {
	gaps := []float64{10, 20, 30, 40, 50}
	// Only use indices 1..3 → {20, 30, 40}
	h := calcH(gaps, 1, 3, 0.5)
	expected := quantile([]float64{20, 30, 40}, 0.5)
	if !almostEqual(h, expected, 0.01) {
		t.Errorf("expected %f, got %f", expected, h)
	}
}

func TestCalcH_SkipsZeros(t *testing.T) {
	gaps := []float64{0, 100, 0}
	h := calcH(gaps, 0, 2, 0.5)
	// Only non-zero is 100
	if !almostEqual(h, 100, 0.01) {
		t.Errorf("expected 100, got %f", h)
	}
}

// ─── calcNetScore ───────────────────────────────────────────────────────────

func TestCalcNetScore_PerfectCoverage(t *testing.T) {
	gaps := []float64{100, 100, 100}
	bins := []internalBin{
		{isCurrent: false},
		{isCurrent: false},
		{isCurrent: false},
	}
	// h=100 perfectly covers all gaps
	score := calcNetScore(gaps, bins, 0, 2, 100, 0.5, 50, 0, 3, 1)
	if score <= 0 {
		t.Errorf("expected positive score for perfect coverage, got %f", score)
	}
}

func TestCalcNetScore_CurrentBonus(t *testing.T) {
	gaps := []float64{100}
	binsNoCurrent := []internalBin{{isCurrent: false}}
	binsCurrent := []internalBin{{isCurrent: true}}

	scoreNoCurrent := calcNetScore(gaps, binsNoCurrent, 0, 0, 100, 0.5, 0, 0.2, 1, 1)
	scoreCurrent := calcNetScore(gaps, binsCurrent, 0, 0, 100, 0.5, 0, 0.2, 1, 1)

	if scoreCurrent <= scoreNoCurrent {
		t.Errorf("current bonus should increase score: current=%f, noCurrent=%f", scoreCurrent, scoreNoCurrent)
	}
}

func TestCalcNetScore_WidthPenalty(t *testing.T) {
	gaps := []float64{100, 100, 100, 100, 100}
	bins := make([]internalBin, 5)

	// Wider segment with same h should get penalized
	scoreNarrow := calcNetScore(gaps, bins, 0, 0, 100, 0.5, 50, 0, 5, 5)
	scoreWide := calcNetScore(gaps, bins, 0, 4, 100, 0.5, 50, 0, 5, 5)

	// Narrow segment covers 1 bin with ideal width=1, no penalty
	// Wide segment covers 5 bins with ideal width=1, heavy penalty
	if scoreWide >= scoreNarrow*5 {
		t.Errorf("width penalty should limit wide segment advantage: narrow=%f, wide=%f", scoreNarrow, scoreWide)
	}
}

// ─── enforceMinLiquidity ────────────────────────────────────────────────────

func TestEnforceMinLiquidity_Empty(t *testing.T) {
	result := enforceMinLiquidity(nil, 3)
	if len(result) != 0 {
		t.Errorf("expected empty, got %d", len(result))
	}
}

func TestEnforceMinLiquidity_FiltersTinySegments(t *testing.T) {
	segments := []internalSegment{
		{l: 0, r: 4, h: 1000, liquidityAdded: 1000}, // width=5, amount=5000
		{l: 6, r: 6, h: 1, liquidityAdded: 1},       // width=1, amount=1 (tiny)
	}

	result := enforceMinLiquidity(segments, 3)

	// threshold = 5000 / (3*2) = 833.3 — second segment (amount=1) should be filtered
	if len(result) != 1 {
		t.Fatalf("expected 1 segment after filtering, got %d", len(result))
	}
	if result[0].l != 0 || result[0].r != 4 {
		t.Errorf("expected big segment to survive, got l=%d r=%d", result[0].l, result[0].r)
	}
}

func TestEnforceMinLiquidity_KeepsAll(t *testing.T) {
	segments := []internalSegment{
		{l: 0, r: 4, h: 100, liquidityAdded: 100}, // amount=500
		{l: 6, r: 9, h: 80, liquidityAdded: 80},   // amount=320
	}

	result := enforceMinLiquidity(segments, 5)

	// threshold = 500 / (5*2) = 50 — both exceed
	if len(result) != 2 {
		t.Errorf("expected 2 segments kept, got %d", len(result))
	}
}

// ─── DefaultConfig ──────────────────────────────────────────────────────────

func TestDefaultConfig(t *testing.T) {
	cfg := DefaultConfig()

	if cfg.N != 5 {
		t.Errorf("N: expected 5, got %d", cfg.N)
	}
	if cfg.MinWidth != 1 {
		t.Errorf("MinWidth: expected 1, got %d", cfg.MinWidth)
	}
	if cfg.MaxWidth != 0 {
		t.Errorf("MaxWidth: expected 0, got %d", cfg.MaxWidth)
	}
	if cfg.Lambda != 50.0 {
		t.Errorf("Lambda: expected 50, got %f", cfg.Lambda)
	}
	if cfg.Beta != 0.5 {
		t.Errorf("Beta: expected 0.5, got %f", cfg.Beta)
	}
	if cfg.WeightMode != "quantile" {
		t.Errorf("WeightMode: expected quantile, got %s", cfg.WeightMode)
	}
	if cfg.Quantile != 0.6 {
		t.Errorf("Quantile: expected 0.6, got %f", cfg.Quantile)
	}
	if cfg.LookAhead != 3 {
		t.Errorf("LookAhead: expected 3, got %d", cfg.LookAhead)
	}
}
