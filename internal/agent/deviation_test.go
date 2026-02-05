package agent

import (
	"math"
	"math/big"
	"testing"

	"remora/internal/allocation"
	"remora/internal/vault"
)

func almostEqual(a, b, eps float64) bool {
	return math.Abs(a-b) < eps
}

func newService() *Service {
	return &Service{}
}

// ─── Edge cases ─────────────────────────────────────────────────────────────

func TestDeviation_NoCurrentNoPlanned(t *testing.T) {
	s := newService()
	d := s.calculateDeviation(nil, nil)

	if d != 0.0 {
		t.Errorf("expected 0.0 for empty current + empty planned, got %f", d)
	}
}

func TestDeviation_NoCurrentWithPlanned(t *testing.T) {
	s := newService()
	planned := []allocation.PositionPlan{
		{TickLower: 0, TickUpper: 100, Liquidity: big.NewInt(500)},
	}

	d := s.calculateDeviation(nil, planned)

	if d != 1.0 {
		t.Errorf("expected 1.0 for no current + planned positions, got %f", d)
	}
}

func TestDeviation_HasCurrentNoPlanned(t *testing.T) {
	s := newService()
	positions := []vault.Position{
		{TickLower: 0, TickUpper: 100, Liquidity: big.NewInt(1000)},
	}

	d := s.calculateDeviation(positions, nil)

	if d != 1.0 {
		t.Errorf("expected 1.0 for current positions + no planned, got %f", d)
	}
}

func TestDeviation_BothEmpty(t *testing.T) {
	s := newService()
	d := s.calculateDeviation([]vault.Position{}, []allocation.PositionPlan{})

	if d != 0.0 {
		t.Errorf("expected 0.0 for both empty, got %f", d)
	}
}

// ─── Identical distributions ────────────────────────────────────────────────

func TestDeviation_IdenticalSinglePosition(t *testing.T) {
	s := newService()
	positions := []vault.Position{
		{TickLower: 0, TickUpper: 200, Liquidity: big.NewInt(1000)},
	}
	planned := []allocation.PositionPlan{
		{TickLower: 0, TickUpper: 200, Liquidity: big.NewInt(1000)},
	}

	d := s.calculateDeviation(positions, planned)

	if !almostEqual(d, 0.0, 0.01) {
		t.Errorf("expected ~0.0 for identical positions, got %f", d)
	}
}

func TestDeviation_IdenticalProportional(t *testing.T) {
	s := newService()
	positions := []vault.Position{
		{TickLower: 0, TickUpper: 200, Liquidity: big.NewInt(2000)},
	}
	planned := []allocation.PositionPlan{
		{TickLower: 0, TickUpper: 200, Liquidity: big.NewInt(500)},
	}

	d := s.calculateDeviation(positions, planned)

	if !almostEqual(d, 0.0, 0.01) {
		t.Errorf("expected ~0.0 for proportional distributions, got %f", d)
	}
}

// ─── Completely different distributions ─────────────────────────────────────

func TestDeviation_CompletelyDifferent(t *testing.T) {
	s := newService()
	positions := []vault.Position{
		{TickLower: 100, TickUpper: 200, Liquidity: big.NewInt(1000)},
	}
	planned := []allocation.PositionPlan{
		{TickLower: 0, TickUpper: 100, Liquidity: big.NewInt(1000)},
	}

	d := s.calculateDeviation(positions, planned)

	// No overlap → 1.0
	if !almostEqual(d, 1.0, 0.01) {
		t.Errorf("expected 1.0 for completely different distributions, got %f", d)
	}
}

// ─── Partial overlap ────────────────────────────────────────────────────────

func TestDeviation_PartialOverlap(t *testing.T) {
	s := newService()
	// Current covers [0,100)
	positions := []vault.Position{
		{TickLower: 0, TickUpper: 100, Liquidity: big.NewInt(1000)},
	}
	// Plan covers [0,200) — same L over twice the range
	planned := []allocation.PositionPlan{
		{TickLower: 0, TickUpper: 200, Liquidity: big.NewInt(1000)},
	}

	d := s.calculateDeviation(positions, planned)

	// Plan weights: [0,100)=0.5, [100,200)=0.5
	// Current weights: [0,100)=1.0, [100,200)=0.0
	// L1 = |1-0.5| + |0-0.5| = 1.0, deviation = 0.5
	if !almostEqual(d, 0.5, 0.01) {
		t.Errorf("expected 0.5, got %f", d)
	}
}

// ─── Multiple segments and positions ────────────────────────────────────────

func TestDeviation_MultiplePositions(t *testing.T) {
	s := newService()
	positions := []vault.Position{
		{TickLower: 0, TickUpper: 200, Liquidity: big.NewInt(100)},
		{TickLower: 200, TickUpper: 400, Liquidity: big.NewInt(100)},
	}
	planned := []allocation.PositionPlan{
		{TickLower: 0, TickUpper: 200, Liquidity: big.NewInt(100)},
		{TickLower: 200, TickUpper: 400, Liquidity: big.NewInt(100)},
	}

	d := s.calculateDeviation(positions, planned)

	if !almostEqual(d, 0.0, 0.01) {
		t.Errorf("expected ~0.0, got %f", d)
	}
}

// ─── Null/zero liquidity positions ──────────────────────────────────────────

func TestDeviation_NilLiquidityPosition(t *testing.T) {
	s := newService()
	positions := []vault.Position{
		{TickLower: 0, TickUpper: 100, Liquidity: nil},
	}
	planned := []allocation.PositionPlan{
		{TickLower: 0, TickUpper: 100, Liquidity: big.NewInt(1000)},
	}

	d := s.calculateDeviation(positions, planned)

	// Current has no effective liquidity → 1.0
	if d != 1.0 {
		t.Errorf("expected 1.0 for nil liquidity position, got %f", d)
	}
}

func TestDeviation_ZeroLiquidityPosition(t *testing.T) {
	s := newService()
	positions := []vault.Position{
		{TickLower: 0, TickUpper: 100, Liquidity: big.NewInt(0)},
	}
	planned := []allocation.PositionPlan{
		{TickLower: 0, TickUpper: 100, Liquidity: big.NewInt(1000)},
	}

	d := s.calculateDeviation(positions, planned)

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
		planned   []allocation.PositionPlan
	}{
		{
			name:      "shifted",
			positions: []vault.Position{{TickLower: 100, TickUpper: 300, Liquidity: big.NewInt(500)}},
			planned:   []allocation.PositionPlan{{TickLower: 0, TickUpper: 200, Liquidity: big.NewInt(1000)}},
		},
		{
			name:      "uneven",
			positions: []vault.Position{{TickLower: 0, TickUpper: 400, Liquidity: big.NewInt(1)}},
			planned: []allocation.PositionPlan{
				{TickLower: 0, TickUpper: 200, Liquidity: big.NewInt(9999)},
				{TickLower: 200, TickUpper: 400, Liquidity: big.NewInt(1)},
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			d := s.calculateDeviation(tc.positions, tc.planned)

			if d < 0.0 || d > 1.0 {
				t.Errorf("deviation %f out of [0,1] range", d)
			}
		})
	}
}

// ─── Nil planned liquidity ─────────────────────────────────────────────────

func TestDeviation_NilPlannedLiquidity(t *testing.T) {
	s := newService()
	positions := []vault.Position{
		{TickLower: 0, TickUpper: 100, Liquidity: big.NewInt(1000)},
	}
	planned := []allocation.PositionPlan{
		{TickLower: 0, TickUpper: 100, Liquidity: nil},
	}

	d := s.calculateDeviation(positions, planned)

	// Planned has no effective liquidity → 1.0
	if d != 1.0 {
		t.Errorf("expected 1.0 for nil planned liquidity, got %f", d)
	}
}
