package coverage

import (
	"log/slog"
	"math"
	"sort"
)

// Run executes the greedy coverage algorithm (default).
func Run(bins []Bin, cfg Config) Result {
	if len(bins) == 0 {
		return Result{}
	}

	// Convert to internal format
	internalBins := toInternalBins(bins)

	// Run LookAhead greedy algorithm
	return runLookAhead(internalBins, cfg)
}

// runLookAhead implements the new look-ahead expansion algorithm.
//
//nolint:cyclop // algorithm complexity
func runLookAhead(bins []internalBin, cfg Config) Result {
	n := len(bins)
	target := make([]float64, n)

	// Initialize target array
	for i, bin := range bins {
		target[i] = bin.liquidity
	}

	// Phase 1: Greedy expansion with look-ahead
	remainingGaps := make([]float64, n)
	copy(remainingGaps, target)

	var segments []internalSegment

	for len(segments) < cfg.N {
		// Find seed: bin with max remaining gap
		seed := -1
		maxGap := 0.0
		for i, g := range remainingGaps {
			if g > maxGap {
				maxGap = g
				seed = i
			}
		}

		if cfg.Debug {
			slog.Debug("[DEBUG] Round", slog.Int("round", len(segments)+1), slog.Int("seed", seed), slog.Float64("maxGap", maxGap))
		}

		if seed < 0 || maxGap <= 0 {
			break // No more gaps
		}

		// Expand from seed using look-ahead
		seg := expandWithLookAhead(remainingGaps, bins, seed, cfg, n)

		if cfg.Debug {
			slog.Debug("[DEBUG] Segment", slog.Int("l", seg.l), slog.Int("r", seg.r), slog.Float64("h", seg.h), slog.Float64("liq", seg.liquidityAdded))
		}

		if seg.h <= 0 {
			break
		}

		segments = append(segments, seg)

		// Update remaining gaps
		for i := seg.l; i <= seg.r; i++ {
			remainingGaps[i] = math.Max(0, remainingGaps[i]-seg.h)
		}
	}

	// Phase 2: Iterative convergence (min_liquidity constraint)
	// threshold = max_liquidity / 2^N
	if cfg.EnableMinLiq {
		segments = enforceMinLiquidity(segments, cfg.N)
	}

	// Convert to output format
	return toSegments(bins, segments, target)
}

// expandWithLookAhead expands a segment from seed using look-ahead strategy.
func expandWithLookAhead(gaps []float64, bins []internalBin, seed int, cfg Config, totalBins int) internalSegment {
	n := len(gaps)
	l, r := seed, seed

	// Calculate initial net score
	h := calcH(gaps, l, r, cfg.Quantile)
	currentScore := calcNetScore(gaps, bins, l, r, h, cfg.Beta, cfg.Lambda, cfg.CurrentBonus, totalBins, cfg.N)

	for {
		bestNewL, bestNewR := l, r
		bestScore := currentScore

		// Try expanding left with look-ahead
		for steps := 1; steps <= cfg.LookAhead && l-steps >= 0; steps++ {
			newL := l - steps
			newH := calcH(gaps, newL, r, cfg.Quantile)
			newScore := calcNetScore(gaps, bins, newL, r, newH, cfg.Beta, cfg.Lambda, cfg.CurrentBonus, totalBins, cfg.N)
			if newScore > bestScore {
				bestScore = newScore
				bestNewL = newL
				bestNewR = r
			}
		}

		// Try expanding right with look-ahead
		for steps := 1; steps <= cfg.LookAhead && r+steps < n; steps++ {
			newR := r + steps
			newH := calcH(gaps, l, newR, cfg.Quantile)
			newScore := calcNetScore(gaps, bins, l, newR, newH, cfg.Beta, cfg.Lambda, cfg.CurrentBonus, totalBins, cfg.N)
			if newScore > bestScore {
				bestScore = newScore
				bestNewL = l
				bestNewR = newR
			}
		}

		// If no improvement, stop
		if bestNewL == l && bestNewR == r {
			break
		}

		// Apply best expansion
		l, r = bestNewL, bestNewR
		currentScore = bestScore
	}

	finalH := calcH(gaps, l, r, cfg.Quantile)
	// internalSegment.liquidityAdded should store Height (h) to be consistent with enforceMinLiquidity
	return internalSegment{l: l, r: r, h: finalH, liquidityAdded: finalH}
}

// calcH calculates the height (h) for a segment using quantile.
func calcH(gaps []float64, l, r int, q float64) float64 {
	segmentGaps := make([]float64, 0, r-l+1)
	for i := l; i <= r; i++ {
		if gaps[i] > 0 {
			segmentGaps = append(segmentGaps, gaps[i])
		}
	}
	if len(segmentGaps) == 0 {
		return 0
	}
	return quantile(segmentGaps, q)
}

// calcNetScore calculates the net score for a segment.
// score = Σ min(gap[i], h)  (captured gap)
// loss  = Σ max(0, gap[i]-h) + β×Σmax(0, h-gap[i]) + λ×widthPenalty
// widthPenalty = max(0, ratio-1) × avgGap, where ratio = numBins / idealWidth
// net_score = score - loss, with currentBonus if segment contains current price.
func calcNetScore(gaps []float64, bins []internalBin, l, r int, h, beta, lambda, currentBonus float64, totalBins, k int) float64 {
	var captured, underCover, waste, sumGap float64
	numBins := float64(r - l + 1)
	containsCurrent := false

	for i := l; i <= r; i++ {
		captured += math.Min(gaps[i], h)     // gap captured by this segment
		underCover += math.Max(0, gaps[i]-h) // gap not covered
		waste += math.Max(0, h-gaps[i])      // liquidity wasted
		sumGap += gaps[i]
		if bins[i].isCurrent {
			containsCurrent = true
		}
	}

	// Width penalty: only penalize if exceeding ideal width
	idealWidth := float64(totalBins) / float64(k)
	ratio := numBins / idealWidth
	excess := math.Max(0, ratio-1)
	avgGap := sumGap / numBins
	widthPenalty := lambda * excess * avgGap

	// Apply current price bonus to captured only
	if containsCurrent && currentBonus > 0 {
		captured *= (1 + currentBonus)
	}

	loss := underCover + beta*waste + widthPenalty
	return captured - loss
}

// enforceMinLiquidity applies the min_liquidity constraint.
// threshold = max_total_liquidity / (N×2).
// Total Liquidity = h * width.
func enforceMinLiquidity(segments []internalSegment, n int) []internalSegment {
	if len(segments) == 0 {
		return segments
	}

	// Find max total liquidity (amount)
	maxAmount := 0.0
	for _, seg := range segments {
		width := float64(seg.r - seg.l + 1)
		amount := seg.liquidityAdded * width
		if amount > maxAmount {
			maxAmount = amount
		}
	}

	// threshold = maxAmount / (N×2)
	threshold := maxAmount / float64(n*2) //nolint:mnd // threshold factor

	// Filter segments below threshold
	var validSegments []internalSegment

	for _, seg := range segments {
		width := float64(seg.r - seg.l + 1)
		amount := seg.liquidityAdded * width

		if amount >= threshold {
			validSegments = append(validSegments, seg)
		}
	}

	return validSegments
}

// quantile calculates the q-th quantile of a slice.
func quantile(data []float64, q float64) float64 {
	if len(data) == 0 {
		return 0
	}

	sorted := make([]float64, len(data))
	copy(sorted, data)
	sort.Float64s(sorted)

	if q <= 0 {
		return sorted[0]
	}
	if q >= 1 {
		return sorted[len(sorted)-1]
	}

	index := q * float64(len(sorted)-1)
	lower := int(math.Floor(index))
	upper := int(math.Ceil(index))

	if lower == upper {
		return sorted[lower]
	}

	// Linear interpolation
	weight := index - float64(lower)
	return sorted[lower]*(1-weight) + sorted[upper]*weight
}
