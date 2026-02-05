package service

import (
	"context"
	"fmt"
	"time"

	"remora/internal/coverage"
	"remora/internal/liquidity"
	"remora/internal/strategy"
)

// Service implements strategy.Service.
type Service struct {
	liquiditySvc liquidity.Service
}

// New creates a new strategy service.
func New(liquiditySvc liquidity.Service) *Service {
	return &Service{
		liquiditySvc: liquiditySvc,
	}
}

// Ensure Service implements strategy.Service.
var _ strategy.Service = (*Service)(nil)

// ComputeTargetPositions computes optimal LP positions based on market liquidity.
func (s *Service) ComputeTargetPositions(ctx context.Context, params *strategy.ComputeParams) (*strategy.ComputeResult, error) {
	// Step 1: Get market liquidity distribution
	dist, err := s.liquiditySvc.GetDistribution(ctx, &liquidity.DistributionParams{
		PoolKey:      params.PoolKey,
		BinSizeTicks: params.BinSizeTicks,
		TickRange:    params.TickRange,
	})
	if err != nil {
		return nil, fmt.Errorf("get distribution: %w", err)
	}

	// Step 2: Convert liquidity bins to allocation bins
	allocationBins := toAllocationBins(dist.Bins, dist.CurrentTick)

	if len(allocationBins) == 0 {
		return &strategy.ComputeResult{
			CurrentTick: dist.CurrentTick,
			Segments:    nil,
			Metrics:     coverage.Metrics{},
			ComputedAt:  time.Now().UTC(),
		}, nil
	}

	// Step 3: Run coverage algorithm
	result := coverage.Run(allocationBins, params.AlgoConfig)

	return &strategy.ComputeResult{
		CurrentTick: dist.CurrentTick,
		Segments:    result.Segments,
		Metrics:     result.Metrics,
		ComputedAt:  time.Now().UTC(),
	}, nil
}

// toAllocationBins converts liquidity.Bin slice to coverage.Bin slice.
func toAllocationBins(liqBins []liquidity.Bin, currentTick int32) []coverage.Bin {
	if len(liqBins) == 0 {
		return nil
	}

	bins := make([]coverage.Bin, len(liqBins))
	for i, b := range liqBins {
		// Determine if this bin contains the current tick
		isCurrent := currentTick >= b.TickLower && currentTick < b.TickUpper

		bins[i] = coverage.Bin{
			TickLower: b.TickLower,
			TickUpper: b.TickUpper,
			Liquidity: b.ActiveLiquidity,
			IsCurrent: isCurrent,
		}
	}
	return bins
}
