package strategy

//go:generate mockgen -destination=mocks/mock_service.go -package=mocks . Service

import (
	"context"
	"math/big"
	"time"

	"remora/internal/allocation"
	"remora/internal/liquidity"
)

// Service defines the strategy orchestration use cases.
type Service interface {
	// ComputeTargetPositions computes optimal LP positions based on market liquidity.
	ComputeTargetPositions(ctx context.Context, params *ComputeParams) (*ComputeResult, error)
}

// ComputeParams contains parameters for computing target positions.
type ComputeParams struct {
	PoolKey      liquidity.PoolKey // Uniswap v4 pool key
	BinSizeTicks int32             // Size of each bin in ticks
	TickRange    int32             // Range of ticks to scan (±tickRange from current tick)
	AlgoConfig   allocation.Config // Algorithm configuration
}

// ComputeResult contains the computed target positions.
type ComputeResult struct {
	CurrentTick int32                // Current pool tick
	Segments    []allocation.Segment // Target LP segments
	Metrics     allocation.Metrics   // Coverage metrics
	ComputedAt  time.Time            // Timestamp when computation was performed
}

// Position represents an existing LP position (for future use with gap calculation).
type Position struct {
	TickLower int32
	TickUpper int32
	Liquidity *big.Int
}
