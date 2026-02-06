package api

import (
	"encoding/hex"
	"log/slog"
	"math/big"
	"net/http"

	"github.com/ethereum/go-ethereum/common"
	"github.com/go-chi/chi/v5"

	"remora/internal/allocation"
	"remora/internal/httpwrap"
	"remora/internal/liquidity"
	"remora/internal/liquidity/poolid"
	"remora/internal/vault"
)

// VaultFactory creates a vault client for the given address (read-only).
// When nil, vault endpoints return 503.
type VaultFactory func(address common.Address) (vault.Vault, error)

// AddRoutes registers vault-related routes on the provided router.
func AddRoutes(r chi.Router, factory VaultFactory, liquiditySvc liquidity.Service) {
	if factory == nil {
		return
	}

	r.Get("/vaults/{address}/state", httpwrap.Handler(getState(factory)))
	r.Get("/vaults/{address}/positions", httpwrap.Handler(getPositions(factory, liquiditySvc)))
}

// StateResponse is the API response for vault state.
type StateResponse struct {
	Agent            string          `json:"agent"`
	AgentPaused      bool            `json:"agentPaused"`
	SwapAllowed      bool            `json:"swapAllowed"`
	AllowedTickLower int32           `json:"allowedTickLower"`
	AllowedTickUpper int32           `json:"allowedTickUpper"`
	MaxPositionsK    string          `json:"maxPositionsK"`
	PoolKey          PoolKeyResponse `json:"poolKey"`
	PoolID           string          `json:"poolId"`
	Posm             string          `json:"posm"`
	PositionsLength  string          `json:"positionsLength"`
}

// PoolKeyResponse is the pool key in API response.
type PoolKeyResponse struct {
	Currency0   string `json:"currency0"`
	Currency1   string `json:"currency1"`
	Fee         string `json:"fee"`
	TickSpacing string `json:"tickSpacing"`
	Hooks       string `json:"hooks"`
}

// PositionsResponse is the API response for vault positions (with vault totals).
type PositionsResponse struct {
	Amount0   string             `json:"amount0"`
	Amount1   string             `json:"amount1"`
	Positions []PositionResponse `json:"positions"`
}

// PositionResponse is a single position in the API response.
type PositionResponse struct {
	TokenID   string `json:"tokenId"`
	TickLower int32  `json:"tickLower"`
	TickUpper int32  `json:"tickUpper"`
	Liquidity string `json:"liquidity"`
	Amount0   string `json:"amount0"`
	Amount1   string `json:"amount1"`
}

func getState(factory VaultFactory) func(*http.Request) (*httpwrap.Response, *httpwrap.ErrorResponse) {
	return func(r *http.Request) (*httpwrap.Response, *httpwrap.ErrorResponse) {
		addrHex := chi.URLParam(r, "address")
		if addrHex == "" {
			return nil, &httpwrap.ErrorResponse{
				StatusCode: http.StatusBadRequest,
				ErrorMsg:   "missing address",
			}
		}

		if !common.IsHexAddress(addrHex) {
			return nil, &httpwrap.ErrorResponse{
				StatusCode: http.StatusBadRequest,
				ErrorMsg:   "invalid address",
			}
		}

		addr := common.HexToAddress(addrHex)

		v, err := factory(addr)
		if err != nil {
			slog.ErrorContext(r.Context(), "vault factory failed", slog.String("address", addrHex), slog.String("error", err.Error()))

			return nil, &httpwrap.ErrorResponse{
				StatusCode: http.StatusInternalServerError,
				ErrorMsg:   err.Error(),
				Err:        err,
			}
		}

		state, err := v.GetState(r.Context())
		if err != nil {
			slog.ErrorContext(r.Context(), "get vault state failed", slog.String("address", addrHex), slog.String("error", err.Error()))

			return nil, &httpwrap.ErrorResponse{
				StatusCode: http.StatusInternalServerError,
				ErrorMsg:   err.Error(),
				Err:        err,
			}
		}

		return &httpwrap.Response{
			StatusCode: http.StatusOK,
			Body: &StateResponse{
				Agent:            state.Agent.Hex(),
				AgentPaused:      state.AgentPaused,
				SwapAllowed:      state.SwapAllowed,
				AllowedTickLower: state.AllowedTickLower,
				AllowedTickUpper: state.AllowedTickUpper,
				MaxPositionsK:    state.MaxPositionsK.String(),
				PoolKey: PoolKeyResponse{
					Currency0:   state.PoolKey.Currency0.Hex(),
					Currency1:   state.PoolKey.Currency1.Hex(),
					Fee:         state.PoolKey.Fee.String(),
					TickSpacing: state.PoolKey.TickSpacing.String(),
					Hooks:       state.PoolKey.Hooks.Hex(),
				},
				PoolID:          hex.EncodeToString(state.PoolID[:]),
				Posm:            state.Posm.Hex(),
				PositionsLength: state.PositionsLength.String(),
			},
		}, nil
	}
}

func getPositions(factory VaultFactory, liquiditySvc liquidity.Service) func(*http.Request) (*httpwrap.Response, *httpwrap.ErrorResponse) {
	return func(r *http.Request) (*httpwrap.Response, *httpwrap.ErrorResponse) {
		addrHex := chi.URLParam(r, "address")
		if addrHex == "" {
			return nil, &httpwrap.ErrorResponse{
				StatusCode: http.StatusBadRequest,
				ErrorMsg:   "missing address",
			}
		}

		if !common.IsHexAddress(addrHex) {
			return nil, &httpwrap.ErrorResponse{
				StatusCode: http.StatusBadRequest,
				ErrorMsg:   "invalid address",
			}
		}

		addr := common.HexToAddress(addrHex)

		v, err := factory(addr)
		if err != nil {
			slog.ErrorContext(r.Context(), "vault factory failed", slog.String("address", addrHex), slog.String("error", err.Error()))

			return nil, &httpwrap.ErrorResponse{
				StatusCode: http.StatusInternalServerError,
				ErrorMsg:   err.Error(),
				Err:        err,
			}
		}

		state, err := v.GetState(r.Context())
		if err != nil {
			slog.ErrorContext(r.Context(), "get vault state failed", slog.String("address", addrHex), slog.String("error", err.Error()))

			return nil, &httpwrap.ErrorResponse{
				StatusCode: http.StatusInternalServerError,
				ErrorMsg:   err.Error(),
				Err:        err,
			}
		}

		positions, err := v.GetPositions(r.Context())
		if err != nil {
			slog.ErrorContext(r.Context(), "get vault positions failed", slog.String("address", addrHex), slog.String("error", err.Error()))

			return nil, &httpwrap.ErrorResponse{
				StatusCode: http.StatusInternalServerError,
				ErrorMsg:   err.Error(),
				Err:        err,
			}
		}

		poolKey := vaultPoolKeyToLiquidity(&state.PoolKey)

		var slot0 *liquidity.Slot0
		if liquiditySvc != nil {
			slot0, err = liquiditySvc.GetSlot0(r.Context(), poolKey)
			if err != nil {
				slog.ErrorContext(r.Context(), "get pool slot0 failed", slog.String("address", addrHex), slog.String("error", err.Error()))

				return nil, &httpwrap.ErrorResponse{
					StatusCode: http.StatusInternalServerError,
					ErrorMsg:   err.Error(),
					Err:        err,
				}
			}
		}

		resp := make([]PositionResponse, len(positions))
		total0, total1 := new(big.Int), new(big.Int)

		for i, p := range positions {
			amount0, amount1 := "0", "0"

			if slot0 != nil {
				sqrtPriceA := allocation.TickToSqrtPriceX96(int(p.TickLower))
				sqrtPriceB := allocation.TickToSqrtPriceX96(int(p.TickUpper))
				a0 := allocation.GetAmount0ForLiquidity(slot0.SqrtPriceX96, sqrtPriceA, sqrtPriceB, p.Liquidity)
				a1 := allocation.GetAmount1ForLiquidity(slot0.SqrtPriceX96, sqrtPriceA, sqrtPriceB, p.Liquidity)
				amount0 = a0.String()
				amount1 = a1.String()

				total0.Add(total0, a0)
				total1.Add(total1, a1)
			}

			resp[i] = PositionResponse{
				TokenID:   p.TokenID.String(),
				TickLower: p.TickLower,
				TickUpper: p.TickUpper,
				Liquidity: p.Liquidity.String(),
				Amount0:   amount0,
				Amount1:   amount1,
			}
		}

		return &httpwrap.Response{
			StatusCode: http.StatusOK,
			Body: &PositionsResponse{
				Amount0:   total0.String(),
				Amount1:   total1.String(),
				Positions: resp,
			},
		}, nil
	}
}

// vaultPoolKeyToLiquidity converts vault.PoolKey to poolid.PoolKey for liquidity service.
func vaultPoolKeyToLiquidity(k *vault.PoolKey) *poolid.PoolKey {
	if k == nil {
		return nil
	}

	fee := uint32(0)
	if k.Fee != nil && k.Fee.IsUint64() {
		fee = uint32(k.Fee.Uint64())
	}

	tickSpacing := int32(0)
	if k.TickSpacing != nil && k.TickSpacing.IsInt64() {
		tickSpacing = int32(k.TickSpacing.Int64())
	}

	return &poolid.PoolKey{
		Currency0:   k.Currency0.Hex(),
		Currency1:   k.Currency1.Hex(),
		Fee:         fee,
		TickSpacing: tickSpacing,
		Hooks:       k.Hooks.Hex(),
	}
}
