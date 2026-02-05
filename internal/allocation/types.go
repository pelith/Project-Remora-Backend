package allocation

import (
	"math/big"
)

// UserFunds represents user's input funds
type UserFunds struct {
	Amount0 *big.Int // token0 amount (e.g., ETH in wei, 18 decimals)
	Amount1 *big.Int // token1 amount (e.g., USDC, 6 decimals)
}

// PoolState represents the current state of the pool
type PoolState struct {
	SqrtPriceX96   *big.Int
	CurrentTick    int
	Token0Decimals int
	Token1Decimals int
}

// PositionPlan represents a planned LP position for modifyLiquidities
type PositionPlan struct {
	TickLower int
	TickUpper int
	Liquidity *big.Int // calculated L value
	Amount0   *big.Int // token0 needed (amount0Max)
	Amount1   *big.Int // token1 needed (amount1Max)
	Weight    float64  // original weight
}

// AllocationResult is the output of allocation calculation
type AllocationResult struct {
	Positions     []PositionPlan
	TotalAmount0  *big.Int // sum of all positions' amount0
	TotalAmount1  *big.Int // sum of all positions' amount1
	SwapAmount    *big.Int // amount to swap (in source token units)
	SwapToken0To1 bool     // true = swap token0 to token1, false = opposite

}
