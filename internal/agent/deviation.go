package agent

import (
	"math/big"
	"sort"

	"remora/internal/allocation"
	"remora/internal/vault"
)

// calculateDeviation computes the distance between current positions and planned positions.
// Returns a value between 0.0 (identical) and 1.0 (completely different).
// Projects both distributions onto a common bin grid derived from tick boundaries,
// then computes width-weighted L1 distance: Σ |W_curr - W_plan| / 2
func (s *Service) calculateDeviation(current []vault.Position, planned []allocation.PositionPlan) float64 {
	if len(current) == 0 && len(planned) == 0 {
		return 0.0
	}
	if len(current) == 0 || len(planned) == 0 {
		return 1.0
	}

	// Collect all unique tick boundaries from positions with actual liquidity
	tickSet := map[int32]struct{}{}
	for _, pos := range current {
		if pos.Liquidity != nil && pos.Liquidity.Sign() > 0 {
			tickSet[pos.TickLower] = struct{}{}
			tickSet[pos.TickUpper] = struct{}{}
		}
	}
	for _, p := range planned {
		if p.Liquidity != nil && p.Liquidity.Sign() > 0 {
			tickSet[int32(p.TickLower)] = struct{}{}
			tickSet[int32(p.TickUpper)] = struct{}{}
		}
	}

	ticks := make([]int32, 0, len(tickSet))
	for t := range tickSet {
		ticks = append(ticks, t)
	}
	sort.Slice(ticks, func(i, j int) bool { return ticks[i] < ticks[j] })

	if len(ticks) < 2 {
		// No valid tick range from either side
		return 1.0
	}

	// Build bins from consecutive tick pairs and project both distributions
	n := len(ticks) - 1
	currentL := make([]float64, n)
	plannedL := make([]float64, n)

	for i := 0; i < n; i++ {
		lower := ticks[i]
		upper := ticks[i+1]
		width := float64(upper - lower)

		for _, pos := range current {
			if pos.Liquidity == nil || pos.Liquidity.Sign() == 0 {
				continue
			}
			if lower >= pos.TickLower && upper <= pos.TickUpper {
				f, _ := new(big.Float).SetInt(pos.Liquidity).Float64()
				currentL[i] += f * width
			}
		}
		for _, p := range planned {
			if p.Liquidity == nil || p.Liquidity.Sign() == 0 {
				continue
			}
			if lower >= int32(p.TickLower) && upper <= int32(p.TickUpper) {
				f, _ := new(big.Float).SetInt(p.Liquidity).Float64()
				plannedL[i] += f * width
			}
		}
	}

	// Normalize to weights
	var sumCurrent, sumPlanned float64
	for i := 0; i < n; i++ {
		sumCurrent += currentL[i]
		sumPlanned += plannedL[i]
	}

	if sumCurrent == 0 || sumPlanned == 0 {
		if sumCurrent == 0 && sumPlanned == 0 {
			return 0.0
		}
		return 1.0
	}

	// L1 distance / 2
	var l1 float64
	for i := 0; i < n; i++ {
		wCurrent := currentL[i] / sumCurrent
		wPlanned := plannedL[i] / sumPlanned
		diff := wCurrent - wPlanned
		if diff < 0 {
			diff = -diff
		}
		l1 += diff
	}

	return l1 / 2.0
}
