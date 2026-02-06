package api

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"

	"remora/internal/httpwrap"
	"remora/internal/liquidity"
	"remora/internal/liquidity/poolid"
)

// AddRoutes registers liquidity-related routes on the provided router.
func AddRoutes(r chi.Router, svc liquidity.Service) {
	r.Post("/liquidity/distribution", httpwrap.Handler(getDistribution(svc)))
}

// PoolKeyRequest represents the pool key in the API request.
type PoolKeyRequest struct {
	Currency0   string `json:"currency0"`   // Token 0 address (must be < currency1)
	Currency1   string `json:"currency1"`   // Token 1 address
	Fee         uint32 `json:"fee"`         // Fee in hundredths of a bip (e.g., 3000 = 0.3%)
	TickSpacing int32  `json:"tickSpacing"` // Tick spacing
	Hooks       string `json:"hooks"`       // Hooks contract address (0x0 if none)
}

// DistributionRequest is the API request for liquidity distribution.
type DistributionRequest struct {
	PoolKey      PoolKeyRequest `json:"poolKey"`      // Uniswap v4 pool key (PoolId computed server-side)
	BinSizeTicks int32          `json:"binSizeTicks"` // Size of each bin in ticks
	TickRange    int32          `json:"tickRange"`    // Range of ticks to scan (±tickRange from current tick)
}

// TickInfoResponse represents tick information in the API response.
type TickInfoResponse struct {
	Tick           int32  `json:"tick"`
	LiquidityGross string `json:"liquidityGross"`
	LiquidityNet   string `json:"liquidityNet"`
}

// BinResponse represents a liquidity bin in the API response.
type BinResponse struct {
	TickLower       int32  `json:"tickLower"`
	TickUpper       int32  `json:"tickUpper"`
	ActiveLiquidity string `json:"activeLiquidity"`
	PriceLower      string `json:"priceLower,omitempty"`
	PriceUpper      string `json:"priceUpper,omitempty"`
}

// DistributionResponse is the API response for liquidity distribution.
type DistributionResponse struct {
	CurrentTick      int32              `json:"currentTick"`
	SqrtPriceX96     string             `json:"sqrtPriceX96"`
	Liquidity        string             `json:"liquidity"` // Pool total liquidity L
	InitializedTicks []TickInfoResponse `json:"initializedTicks"`
	Bins             []BinResponse      `json:"bins"`
}

// getDistribution returns a handler that fetches liquidity distribution for a pool.
func getDistribution(svc liquidity.Service) func(*http.Request) (*httpwrap.Response, *httpwrap.ErrorResponse) {
	return func(r *http.Request) (*httpwrap.Response, *httpwrap.ErrorResponse) {
		var req DistributionRequest

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			return nil, &httpwrap.ErrorResponse{
				StatusCode: http.StatusBadRequest,
				ErrorMsg:   "invalid request body: " + err.Error(),
				Err:        err,
			}
		}

		ctx := r.Context()

		params := &liquidity.DistributionParams{
			PoolKey: poolid.PoolKey{
				Currency0:   req.PoolKey.Currency0,
				Currency1:   req.PoolKey.Currency1,
				Fee:         req.PoolKey.Fee,
				TickSpacing: req.PoolKey.TickSpacing,
				Hooks:       req.PoolKey.Hooks,
			},
			BinSizeTicks: req.BinSizeTicks,
			TickRange:    req.TickRange,
		}

		dist, err := svc.GetDistribution(ctx, params)
		if err != nil {
			slog.ErrorContext(ctx, "failed to get distribution", //nolint:sloglint // handler error logging, logger not injected in API layer
				slog.String("error", err.Error()),
				slog.String("currency0", req.PoolKey.Currency0),
				slog.String("currency1", req.PoolKey.Currency1),
			)

			return nil, &httpwrap.ErrorResponse{
				StatusCode: http.StatusInternalServerError,
				ErrorMsg:   err.Error(),
				Err:        err,
			}
		}

		// Convert domain model to API response
		ticks := make([]TickInfoResponse, len(dist.InitializedTicks))
		for i, tick := range dist.InitializedTicks {
			ticks[i] = TickInfoResponse{
				Tick:           tick.Tick,
				LiquidityGross: tick.LiquidityGross.String(),
				LiquidityNet:   tick.LiquidityNet.String(),
			}
		}

		bins := make([]BinResponse, len(dist.Bins))
		for i, bin := range dist.Bins {
			bins[i] = BinResponse{
				TickLower:       bin.TickLower,
				TickUpper:       bin.TickUpper,
				ActiveLiquidity: bin.ActiveLiquidity.String(),
				PriceLower:      bin.PriceLower,
				PriceUpper:      bin.PriceUpper,
			}
		}

		return &httpwrap.Response{
			StatusCode: http.StatusOK,
			Body: &DistributionResponse{
				CurrentTick:      dist.CurrentTick,
				SqrtPriceX96:     dist.SqrtPriceX96,
				Liquidity:        dist.Liquidity,
				InitializedTicks: ticks,
				Bins:             bins,
			},
		}, nil
	}
}
