package repository

import (
	"context"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"

	"remora/internal/liquidity"
	"remora/internal/liquidity/poolid"
	"remora/internal/liquidity/repository/contracts"
)

// Repository implements liquidity.Repository using StateView contract.
type Repository struct {
	client   *ethclient.Client
	contract *contracts.StateView
}

// Config contains configuration for repository.
type Config struct {
	RPCURL          string
	ContractAddress string
}

// New creates a new liquidity repository with ethereum client.
func New(cfg Config) (*Repository, error) {
	client, err := ethclient.Dial(cfg.RPCURL)
	if err != nil {
		return nil, fmt.Errorf("dial rpc: %w", err)
	}

	contractAddr := common.HexToAddress(cfg.ContractAddress)

	contract, err := contracts.NewStateView(contractAddr, client)
	if err != nil {
		client.Close()

		return nil, fmt.Errorf("create contract: %w", err)
	}

	return &Repository{
		client:   client,
		contract: contract,
	}, nil
}

// NewMock creates a new mock repository for testing (returns empty data).
func NewMock() *Repository {
	return &Repository{
		client:   nil,
		contract: nil,
	}
}

// Close closes the repository and its underlying connections.
func (r *Repository) Close() {
	if r.client != nil {
		r.client.Close()
	}
}

// Ensure Repository implements liquidity.Repository.
var _ liquidity.Repository = (*Repository)(nil)

// GetSlot0 retrieves current pool state (tick and sqrtPrice).
func (r *Repository) GetSlot0(ctx context.Context, poolKey *poolid.PoolKey) (*liquidity.Slot0, error) {
	// Mock mode for testing
	if r.contract == nil {
		return &liquidity.Slot0{
			SqrtPriceX96: big.NewInt(0),
			Tick:         0,
		}, nil
	}

	poolID := poolid.CalculatePoolID(poolKey)

	result, err := r.contract.GetSlot0(&bind.CallOpts{Context: ctx}, poolID)
	if err != nil {
		return nil, fmt.Errorf("get slot0: %w", err)
	}

	//nolint:gosec // G115: Tick is int24 in Solidity, safe to convert to int32
	return &liquidity.Slot0{
		SqrtPriceX96: result.SqrtPriceX96,
		Tick:         int32(result.Tick.Int64()),
	}, nil
}

// GetTickBitmap retrieves the tick bitmap for a word position.
func (r *Repository) GetTickBitmap(ctx context.Context, poolKey *poolid.PoolKey, wordPos int16) (*big.Int, error) {
	// Mock mode for testing
	if r.contract == nil {
		return big.NewInt(0), nil
	}

	poolID := poolid.CalculatePoolID(poolKey)

	bitmap, err := r.contract.GetTickBitmap(&bind.CallOpts{Context: ctx}, poolID, wordPos)
	if err != nil {
		return nil, fmt.Errorf("get tick bitmap: %w", err)
	}

	return bitmap, nil
}

// GetTickInfo retrieves liquidity info for a specific tick.
func (r *Repository) GetTickInfo(ctx context.Context, poolKey *poolid.PoolKey, tick int32) (*liquidity.TickInfo, error) {
	// Mock mode for testing
	if r.contract == nil {
		return &liquidity.TickInfo{
			Tick:           tick,
			LiquidityGross: big.NewInt(0),
			LiquidityNet:   big.NewInt(0),
		}, nil
	}

	poolID := poolid.CalculatePoolID(poolKey)

	result, err := r.contract.GetTickInfo(&bind.CallOpts{Context: ctx}, poolID, big.NewInt(int64(tick)))
	if err != nil {
		return nil, fmt.Errorf("get tick info: %w", err)
	}

	return &liquidity.TickInfo{
		Tick:           tick,
		LiquidityGross: result.LiquidityGross,
		LiquidityNet:   result.LiquidityNet,
	}, nil
}
