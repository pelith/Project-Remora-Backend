package allocation

import (
	"math"
	"math/big"
)

// internalBin is used for float64-based algorithm calculations.
type internalBin struct {
	tickLower  int
	tickUpper  int
	priceLower float64
	priceUpper float64
	liquidity  float64
	isCurrent  bool
}

// internalSegment holds segment data during algorithm execution.
type internalSegment struct {
	l, r           int
	h              float64
	liquidityAdded float64
}

// toInternalBins converts public Bin slice to internal format for algorithm processing.
func toInternalBins(bins []Bin) []internalBin {
	result := make([]internalBin, len(bins))
	for i, b := range bins {
		var liq float64
		if b.Liquidity != nil {
			liq, _ = new(big.Float).SetInt(b.Liquidity).Float64()
		}

		result[i] = internalBin{
			tickLower:  int(b.TickLower),
			tickUpper:  int(b.TickUpper),
			priceLower: b.PriceLower,
			priceUpper: b.PriceUpper,
			liquidity:  liq,
			isCurrent:  b.IsCurrent,
		}
	}

	return result
}

// toSegments converts internal segments to public Segment format.
func toSegments(bins []internalBin, segments []internalSegment, target []float64) Result {
	outputSegments := make([]Segment, 0, len(segments))
	pred := make([]float64, len(bins))

	for _, seg := range segments {
		// Convert float64 liquidity to big.Int
		// We use Floor to avoid fractional liquidity
		liqFloat := big.NewFloat(seg.h)
		liqInt, _ := liqFloat.Int(nil)

		outputSeg := Segment{
			TickLower:      int32(bins[seg.l].tickLower), //nolint:gosec // safe conversion
			TickUpper:      int32(bins[seg.r].tickUpper), //nolint:gosec // safe conversion
			PriceLower:     bins[seg.l].priceLower,
			PriceUpper:     bins[seg.r].priceUpper,
			LiquidityAdded: liqInt,
		}
		outputSegments = append(outputSegments, outputSeg)

		// Update pred for metrics calculation
		for i := seg.l; i <= seg.r; i++ {
			pred[i] += seg.h
		}
	}

	metrics := calcMetrics(target, pred)

	return Result{
		Segments: outputSegments,
		Metrics:  metrics,
	}
}

// calcMetrics calculates coverage evaluation metrics.
func calcMetrics(target, pred []float64) Metrics {
	var covered, gap, over float64

	for i := range target {
		covered += math.Min(target[i], pred[i])
		gap += math.Max(0, target[i]-pred[i])
		over += math.Max(0, pred[i]-target[i])
	}

	return Metrics{
		Covered: covered,
		Gap:     gap,
		Over:    over,
	}
}
