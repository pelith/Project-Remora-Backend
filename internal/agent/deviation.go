package agent

import (
	"math/big"

	"remora/internal/strategy"
	"remora/internal/vault"
)

// calculateDeviation computes the distance between current positions and target segments.
// Returns a value between 0.0 (identical) and 1.0 (completely different).
// Uses Weight Distribution Distance: Σ |W_curr - W_target| / 2
func (s *Service) calculateDeviation(current []vault.Position, target *strategy.ComputeResult) float64 {
	if len(target.Bins) == 0 {
		if len(current) == 0 {
			return 0.0
		}
		return 1.0
	}

	// 1. Project distributions onto strategy bins
	n := len(target.Bins)
	targetL := make([]float64, n)
	currentL := make([]float64, n)

	// Target distribution (from strategy segments)
	for i, bin := range target.Bins {
		for _, seg := range target.Segments {
			// If bin is within segment range
			if bin.TickLower >= seg.TickLower && bin.TickUpper <= seg.TickUpper {
				if seg.LiquidityAdded != nil {
					f, _ := new(big.Float).SetInt(seg.LiquidityAdded).Float64()
					targetL[i] += f
				}
			}
		}
	}

	// Current distribution (from old positions)
	for i, bin := range target.Bins {
		for _, pos := range current {
			if pos.Liquidity == nil || pos.Liquidity.Sign() == 0 {
				continue
			}
			// If bin is within position range
			if bin.TickLower >= pos.TickLower && bin.TickUpper <= pos.TickUpper {
				f, _ := new(big.Float).SetInt(pos.Liquidity).Float64()
				currentL[i] += f
			}
		}
	}

	// 2. Normalize to weights
	var sumTarget, sumCurrent float64
	for i := 0; i < n; i++ {
		sumTarget += targetL[i]
		sumCurrent += currentL[i]
	}

	// If one side is empty, they are completely different
	if sumTarget == 0 || sumCurrent == 0 {
		if sumTarget == 0 && sumCurrent == 0 {
			return 0.0
		}
		return 1.0
	}

	// 3. Calculate L1 Distance
	var l1Dist float64
	for i := 0; i < n; i++ {
		wTarget := targetL[i] / sumTarget
		wCurrent := currentL[i] / sumCurrent
		
		diff := wCurrent - wTarget
		if diff < 0 {
			diff = -diff
		}
		l1Dist += diff
	}

	return l1Dist / 2.0
}
