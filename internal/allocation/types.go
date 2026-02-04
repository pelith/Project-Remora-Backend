package allocation

import "math/big"

// Bin represents a single tick bin with liquidity data.
// Uses blockchain-native types for seamless integration with Uniswap v4.
type Bin struct {
	TickLower  int32    `json:"tickLower"`
	TickUpper  int32    `json:"tickUpper"`
	PriceLower float64  `json:"priceLower"`
	PriceUpper float64  `json:"priceUpper"`
	Liquidity  *big.Int `json:"liquidity"`
	IsCurrent  bool     `json:"isCurrent"`
}

// Segment represents an LP position segment.
// Uses blockchain-native types for seamless integration with Uniswap v4.
type Segment struct {
	TickLower      int32    `json:"tickLower"`
	TickUpper      int32    `json:"tickUpper"`
	PriceLower     float64  `json:"priceLower"`
	PriceUpper     float64  `json:"priceUpper"`
	LiquidityAdded *big.Int `json:"liquidityAdded"`
}

// Config holds the algorithm configuration.
type Config struct {
	Algo         string  // algorithm: greedy, dp
	N            int     // max number of segments (positions)
	MinWidth     int     // minimum width in number of bins
	MaxWidth     int     // maximum width in number of bins (0 = unlimited)
	Lambda       float64 // width penalty coefficient
	Beta         float64 // waste penalty coefficient (h > gap)
	CurrentBonus float64 // score bonus for segments containing current price (e.g., 0.2 = +20%)
	EnableMinLiq bool    // enable min liquidity filter (threshold = max / 2^N)
	WeightMode   string  // "avg" or "quantile"
	Quantile     float64 // quantile value (used when WeightMode = "quantile")
	LookAhead    int     // look-ahead steps for expansion (0 = use old algorithm)
	Debug        bool    // enable debug output
}

// DefaultConfig returns a default configuration.
func DefaultConfig() Config {
	return Config{
		N:          5,    //nolint:mnd // default segment count
		MinWidth:   1,
		MaxWidth:   0,
		Lambda:     50.0, //nolint:mnd // default width penalty coefficient
		Beta:       0.5,  //nolint:mnd // default waste penalty coefficient
		WeightMode: "quantile",
		Quantile:   0.6, //nolint:mnd // default quantile value
		LookAhead:  3,   //nolint:mnd // default look-ahead steps
	}
}

// Result holds the algorithm output and metrics.
type Result struct {
	Segments []Segment
	Metrics  Metrics
}

// Metrics holds coverage evaluation metrics.
type Metrics struct {
	Covered float64 // Σ min(target, pred)
	Gap     float64 // Σ max(0, target - pred)
	Over    float64 // Σ max(0, pred - target)
}
