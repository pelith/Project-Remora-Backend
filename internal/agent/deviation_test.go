package agent

import (
	"math"
	"math/big"
	"testing"

	"remora/internal/coverage"
	"remora/internal/strategy"
	"remora/internal/vault"
)

func almostEqual(a, b, eps float64) bool {
	return math.Abs(a-b) < eps
}

func newService() *Service {
	return &Service{}
}

// ─── Edge cases ─────────────────────────────────────────────────────────────

func TestDeviation_NoBinsNoPositions(t *testing.T) {
	s := newService()
	d := s.calculateDeviation(nil, &strategy.ComputeResult{})

	if d != 0.0 {
		t.Errorf("expected 0.0 for empty bins + empty positions, got %f", d)
	}
}

func TestDeviation_NoBinsWithPositions(t *testing.T) {
	s := newService()
	positions := []vault.Position{
		{TickLower: 0, TickUpper: 100, Liquidity: big.NewInt(1000)},
	}
	d := s.calculateDeviation(positions, &strategy.ComputeResult{})

	if d != 1.0 {
		t.Errorf("expected 1.0 for empty bins + existing positions, got %f", d)
	}
}

func TestDeviation_HasBinsNoPositions(t *testing.T) {
	s := newService()
	target := &strategy.ComputeResult{
		Bins: []coverage.Bin{
			{TickLower: 0, TickUpper: 100},
		},
		Segments: []coverage.Segment{
			{TickLower: 0, TickUpper: 100, LiquidityAdded: big.NewInt(500)},
		},
	}

	d := s.calculateDeviation(nil, target)

	// Target has liquidity but current has none → 1.0
	if d != 1.0 {
		t.Errorf("expected 1.0 for target with liquidity + no positions, got %f", d)
	}
}

func TestDeviation_BothEmpty_WithBins(t *testing.T) {
	s := newService()
	target := &strategy.ComputeResult{
		Bins: []coverage.Bin{
			{TickLower: 0, TickUpper: 100},
		},
		Segments: nil, // no target segments
	}

	d := s.calculateDeviation(nil, target)

	// sumTarget=0, sumCurrent=0 → 0.0
	if d != 0.0 {
		t.Errorf("expected 0.0 for both empty distributions, got %f", d)
	}
}

// ─── Identical distributions ────────────────────────────────────────────────

func TestDeviation_IdenticalSingleSegment(t *testing.T) {
	s := newService()
	target := &strategy.ComputeResult{
		Bins: []coverage.Bin{
			{TickLower: 0, TickUpper: 100},
			{TickLower: 100, TickUpper: 200},
		},
		Segments: []coverage.Segment{
			{TickLower: 0, TickUpper: 200, LiquidityAdded: big.NewInt(1000)},
		},
	}
	// Current positions exactly match target: one position covering [0,200) with same liquidity
	positions := []vault.Position{
		{TickLower: 0, TickUpper: 200, Liquidity: big.NewInt(1000)},
	}

	d := s.calculateDeviation(positions, target)

	if !almostEqual(d, 0.0, 0.01) {
		t.Errorf("expected ~0.0 for identical distributions, got %f", d)
	}
}

func TestDeviation_IdenticalProportional(t *testing.T) {
	s := newService()
	target := &strategy.ComputeResult{
		Bins: []coverage.Bin{
			{TickLower: 0, TickUpper: 100},
			{TickLower: 100, TickUpper: 200},
		},
		Segments: []coverage.Segment{
			{TickLower: 0, TickUpper: 200, LiquidityAdded: big.NewInt(500)},
		},
	}
	// Current has different absolute amount but same weight distribution → deviation = 0
	positions := []vault.Position{
		{TickLower: 0, TickUpper: 200, Liquidity: big.NewInt(2000)},
	}

	d := s.calculateDeviation(positions, target)

	if !almostEqual(d, 0.0, 0.01) {
		t.Errorf("expected ~0.0 for proportional distributions, got %f", d)
	}
}

// ─── Completely different distributions ─────────────────────────────────────

func TestDeviation_CompletelyDifferent(t *testing.T) {
	s := newService()
	target := &strategy.ComputeResult{
		Bins: []coverage.Bin{
			{TickLower: 0, TickUpper: 100},
			{TickLower: 100, TickUpper: 200},
		},
		Segments: []coverage.Segment{
			// Target only covers first bin
			{TickLower: 0, TickUpper: 100, LiquidityAdded: big.NewInt(1000)},
		},
	}
	// Current only covers second bin
	positions := []vault.Position{
		{TickLower: 100, TickUpper: 200, Liquidity: big.NewInt(1000)},
	}

	d := s.calculateDeviation(positions, target)

	// Target weights: [1.0, 0.0], Current weights: [0.0, 1.0]
	// L1 = |1-0| + |0-1| = 2, deviation = 2/2 = 1.0
	if !almostEqual(d, 1.0, 0.01) {
		t.Errorf("expected 1.0 for completely different distributions, got %f", d)
	}
}

// ─── Partial overlap ────────────────────────────────────────────────────────

func TestDeviation_PartialOverlap(t *testing.T) {
	s := newService()
	target := &strategy.ComputeResult{
		Bins: []coverage.Bin{
			{TickLower: 0, TickUpper: 100},
			{TickLower: 100, TickUpper: 200},
		},
		Segments: []coverage.Segment{
			{TickLower: 0, TickUpper: 200, LiquidityAdded: big.NewInt(1000)},
		},
	}
	// Current only covers first bin
	positions := []vault.Position{
		{TickLower: 0, TickUpper: 100, Liquidity: big.NewInt(1000)},
	}

	d := s.calculateDeviation(positions, target)

	// Target weights: [0.5, 0.5], Current weights: [1.0, 0.0]
	// L1 = |1-0.5| + |0-0.5| = 1.0, deviation = 1.0/2 = 0.5
	if !almostEqual(d, 0.5, 0.01) {
		t.Errorf("expected 0.5, got %f", d)
	}
}

// ─── Multiple segments and positions ────────────────────────────────────────

func TestDeviation_MultipleSegments(t *testing.T) {
	s := newService()
	target := &strategy.ComputeResult{
		Bins: []coverage.Bin{
			{TickLower: 0, TickUpper: 100},
			{TickLower: 100, TickUpper: 200},
			{TickLower: 200, TickUpper: 300},
			{TickLower: 300, TickUpper: 400},
		},
		Segments: []coverage.Segment{
			{TickLower: 0, TickUpper: 200, LiquidityAdded: big.NewInt(100)},
			{TickLower: 200, TickUpper: 400, LiquidityAdded: big.NewInt(100)},
		},
	}
	// Current matches exactly
	positions := []vault.Position{
		{TickLower: 0, TickUpper: 200, Liquidity: big.NewInt(100)},
		{TickLower: 200, TickUpper: 400, Liquidity: big.NewInt(100)},
	}

	d := s.calculateDeviation(positions, target)

	if !almostEqual(d, 0.0, 0.01) {
		t.Errorf("expected ~0.0, got %f", d)
	}
}

// ─── Null/zero liquidity positions ──────────────────────────────────────────

func TestDeviation_NilLiquidityPosition(t *testing.T) {
	s := newService()
	target := &strategy.ComputeResult{
		Bins: []coverage.Bin{
			{TickLower: 0, TickUpper: 100},
		},
		Segments: []coverage.Segment{
			{TickLower: 0, TickUpper: 100, LiquidityAdded: big.NewInt(1000)},
		},
	}
	// Position with nil liquidity should be skipped
	positions := []vault.Position{
		{TickLower: 0, TickUpper: 100, Liquidity: nil},
	}

	d := s.calculateDeviation(positions, target)

	// Current sums to 0 → 1.0
	if d != 1.0 {
		t.Errorf("expected 1.0 for nil liquidity position, got %f", d)
	}
}

func TestDeviation_ZeroLiquidityPosition(t *testing.T) {
	s := newService()
	target := &strategy.ComputeResult{
		Bins: []coverage.Bin{
			{TickLower: 0, TickUpper: 100},
		},
		Segments: []coverage.Segment{
			{TickLower: 0, TickUpper: 100, LiquidityAdded: big.NewInt(1000)},
		},
	}
	positions := []vault.Position{
		{TickLower: 0, TickUpper: 100, Liquidity: big.NewInt(0)},
	}

	d := s.calculateDeviation(positions, target)

	if d != 1.0 {
		t.Errorf("expected 1.0 for zero liquidity position, got %f", d)
	}
}

// ─── Range check ────────────────────────────────────────────────────────────

func TestDeviation_AlwaysBetweenZeroAndOne(t *testing.T) {
	s := newService()

	cases := []struct {
		name      string
		positions []vault.Position
		segments  []coverage.Segment
	}{
		{
			name:      "shifted",
			positions: []vault.Position{{TickLower: 100, TickUpper: 300, Liquidity: big.NewInt(500)}},
			segments:  []coverage.Segment{{TickLower: 0, TickUpper: 200, LiquidityAdded: big.NewInt(1000)}},
		},
		{
			name:      "uneven",
			positions: []vault.Position{{TickLower: 0, TickUpper: 400, Liquidity: big.NewInt(1)}},
			segments: []coverage.Segment{
				{TickLower: 0, TickUpper: 200, LiquidityAdded: big.NewInt(9999)},
				{TickLower: 200, TickUpper: 400, LiquidityAdded: big.NewInt(1)},
			},
		},
	}

	bins := []coverage.Bin{
		{TickLower: 0, TickUpper: 100},
		{TickLower: 100, TickUpper: 200},
		{TickLower: 200, TickUpper: 300},
		{TickLower: 300, TickUpper: 400},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			target := &strategy.ComputeResult{
				Bins:     bins,
				Segments: tc.segments,
			}

			d := s.calculateDeviation(tc.positions, target)

			if d < 0.0 || d > 1.0 {
				t.Errorf("deviation %f out of [0,1] range", d)
			}
		})
	}
}

// ─── Nil segment LiquidityAdded ─────────────────────────────────────────────

func TestDeviation_NilSegmentLiquidity(t *testing.T) {
	s := newService()
	target := &strategy.ComputeResult{
		Bins: []coverage.Bin{
			{TickLower: 0, TickUpper: 100},
		},
		Segments: []coverage.Segment{
			{TickLower: 0, TickUpper: 100, LiquidityAdded: nil},
		},
	}
	positions := []vault.Position{
		{TickLower: 0, TickUpper: 100, Liquidity: big.NewInt(1000)},
	}

	d := s.calculateDeviation(positions, target)

	// Target sums to 0, current > 0 → 1.0
	if d != 1.0 {
		t.Errorf("expected 1.0 for nil segment liquidity, got %f", d)
	}
}
