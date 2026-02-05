package strategy

//go:generate mockgen -destination=mocks/mock_service.go -package=mocks . Service

import (
	"context"
	"math/big"
	"remora/internal/coverage"
	"remora/internal/liquidity/poolid"
	"time"
)

// Service defines the strategy orchestration use cases.
type Service interface {
	// ComputeTargetPositions computes optimal LP positions based on market liquidity.
	ComputeTargetPositions(ctx context.Context, params *ComputeParams) (*ComputeResult, error)
}

// ComputeParams contains parameters for computing target positions.
type ComputeParams struct {
	PoolKey          poolid.PoolKey  // Uniswap v4 pool key
	BinSizeTicks     int32           // Size of each bin in ticks
	TickRange        int32           // Range of ticks to scan (±tickRange from current tick)
	AlgoConfig       coverage.Config // Algorithm configuration
	AllowedTickLower int32           // Vault's allowed lower tick bound
	AllowedTickUpper int32           // Vault's allowed upper tick bound
}

// ComputeResult contains the computed target positions.
type ComputeResult struct {
	CurrentTick  int32              // Current pool tick
	SqrtPriceX96 *big.Int           // Current pool sqrtPriceX96
	Segments     []coverage.Segment // Target LP segments
	Bins         []coverage.Bin    // Original market liquidity bins
	Metrics      coverage.Metrics   // Coverage metrics
	ComputedAt   time.Time          // Timestamp when computation was performed
}

// Position represents an existing LP position (for future use with gap calculation).
type Position struct {
	TickLower int32
	TickUpper int32
	Liquidity *big.Int
}
