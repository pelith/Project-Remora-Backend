package liquidity

//go:generate mockgen -destination=mocks/mock_repository.go -package=mocks . Repository

import (
	"context"
	"math/big"

	"remora/internal/liquidity/poolid"
)

// Slot0 contains the current state of the pool.
type Slot0 struct {
	SqrtPriceX96 *big.Int `json:"sqrtPriceX96"`
	Tick         int32    `json:"tick"`
}

// TickInfo contains liquidity information for a specific tick.
type TickInfo struct {
	Tick           int32    `json:"tick"`
	LiquidityGross *big.Int `json:"liquidityGross"`
	LiquidityNet   *big.Int `json:"liquidityNet"`
}

// Bin represents a discretized liquidity range.
type Bin struct {
	TickLower       int32    `json:"tickLower"`
	TickUpper       int32    `json:"tickUpper"`
	ActiveLiquidity *big.Int `json:"activeLiquidity"`
	PriceLower      string   `json:"priceLower,omitempty"`
	PriceUpper      string   `json:"priceUpper,omitempty"`
}

// DistributionParams contains parameters for liquidity distribution query.
type DistributionParams struct {
	PoolKey      poolid.PoolKey `json:"poolKey"`      // Uniswap v4 pool key (PoolId is computed from this)
	BinSizeTicks int32          `json:"binSizeTicks"` // Size of each bin in ticks
	TickRange    int32          `json:"tickRange"`    // Range of ticks to scan (±tickRange from current tick)
}

// Distribution contains the liquidity distribution result.
type Distribution struct {
	CurrentTick      int32      `json:"currentTick"`
	SqrtPriceX96     string     `json:"sqrtPriceX96"`
	InitializedTicks []TickInfo `json:"initializedTicks"`
	Bins             []Bin      `json:"bins"`
}

// Service defines the use cases for liquidity distribution.
type Service interface {
	// GetDistribution returns the liquidity distribution for a pool.
	GetDistribution(ctx context.Context, params *DistributionParams) (*Distribution, error)
}

// Repository abstracts blockchain interaction for liquidity data.
type Repository interface {
	// GetSlot0 retrieves current pool state (tick and sqrtPrice).
	GetSlot0(ctx context.Context, poolKey *poolid.PoolKey) (*Slot0, error)

	// GetTickBitmap retrieves the tick bitmap for a word position.
	GetTickBitmap(ctx context.Context, poolKey *poolid.PoolKey, wordPos int16) (*big.Int, error)

	// GetTickInfo retrieves liquidity info for a specific tick.
	GetTickInfo(ctx context.Context, poolKey *poolid.PoolKey, tick int32) (*TickInfo, error)

	// GetTickInfoBatch retrieves liquidity info for multiple ticks in one call.
	GetTickInfoBatch(ctx context.Context, poolKey *poolid.PoolKey, ticks []int32) ([]TickInfo, error)
}
