package vault

//go:generate mockgen -destination=mocks/mock_vault.go -package=mocks . Vault

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

// Position represents a managed LP position.
type Position struct {
	TokenID   *big.Int
	TickLower int32
	TickUpper int32
}

// VaultState represents the current state of a vault.
type VaultState struct {
	Agent            common.Address
	AgentPaused      bool
	SwapAllowed      bool
	AllowedTickLower int32
	AllowedTickUpper int32
	MaxPositionsK    *big.Int
	PoolKey          PoolKey
	PoolID           [32]byte
	PositionsLength  *big.Int
}

// Vault defines the interface for interacting with V4AgenticVault contracts.
type Vault interface {
	// Address returns the vault contract address.
	Address() common.Address

	// GetState returns the current vault state.
	GetState(ctx context.Context) (*VaultState, error)

	// GetPositions returns all managed positions.
	GetPositions(ctx context.Context) ([]Position, error)

	// Agent operations
	MintPosition(ctx context.Context, tickLower, tickUpper int32, liquidity *big.Int, amount0Max, amount1Max *big.Int, deadline *big.Int) (*types.Transaction, error)
	IncreaseLiquidity(ctx context.Context, tokenID, liquidity *big.Int, amount0Max, amount1Max *big.Int, deadline *big.Int) (*types.Transaction, error)
	DecreaseLiquidity(ctx context.Context, tokenID, liquidity *big.Int, amount0Min, amount1Min *big.Int, deadline *big.Int) (*types.Transaction, error)
	CollectFees(ctx context.Context, tokenID *big.Int, amount0Min, amount1Min *big.Int, deadline *big.Int) (*types.Transaction, error)
	BurnPosition(ctx context.Context, tokenID *big.Int, amount0Min, amount1Min *big.Int, deadline *big.Int) (*types.Transaction, error)
	Swap(ctx context.Context, zeroForOne bool, amountIn, minAmountOut *big.Int, deadline *big.Int) (*types.Transaction, error)
}

// Client implements Vault interface.
type Client struct {
	address  common.Address
	contract *V4AgenticVault
	auth     *bind.TransactOpts
}

// NewClient creates a new vault client.
func NewClient(address common.Address, backend bind.ContractBackend, auth *bind.TransactOpts) (*Client, error) {
	contract, err := NewV4AgenticVault(address, backend)
	if err != nil {
		return nil, err
	}
	return &Client{
		address:  address,
		contract: contract,
		auth:     auth,
	}, nil
}

// Address returns the vault contract address.
func (c *Client) Address() common.Address {
	return c.address
}

// GetState returns the current vault state.
func (c *Client) GetState(ctx context.Context) (*VaultState, error) {
	opts := &bind.CallOpts{Context: ctx}

	agent, err := c.contract.Agent(opts)
	if err != nil {
		return nil, err
	}

	agentPaused, err := c.contract.AgentPaused(opts)
	if err != nil {
		return nil, err
	}

	swapAllowed, err := c.contract.SwapAllowed(opts)
	if err != nil {
		return nil, err
	}

	allowedTickLower, err := c.contract.AllowedTickLower(opts)
	if err != nil {
		return nil, err
	}

	allowedTickUpper, err := c.contract.AllowedTickUpper(opts)
	if err != nil {
		return nil, err
	}

	maxPositionsK, err := c.contract.MaxPositionsK(opts)
	if err != nil {
		return nil, err
	}

	poolKey, err := c.contract.GetPoolKey(opts)
	if err != nil {
		return nil, err
	}

	poolID, err := c.contract.PoolId(opts)
	if err != nil {
		return nil, err
	}

	positionsLength, err := c.contract.PositionsLength(opts)
	if err != nil {
		return nil, err
	}

	return &VaultState{
		Agent:            agent,
		AgentPaused:      agentPaused,
		SwapAllowed:      swapAllowed,
		AllowedTickLower: int32(allowedTickLower.Int64()),
		AllowedTickUpper: int32(allowedTickUpper.Int64()),
		MaxPositionsK:    maxPositionsK,
		PoolKey:          poolKey,
		PoolID:           poolID,
		PositionsLength:  positionsLength,
	}, nil
}

// GetPositions returns all managed positions.
func (c *Client) GetPositions(ctx context.Context) ([]Position, error) {
	opts := &bind.CallOpts{Context: ctx}

	length, err := c.contract.PositionsLength(opts)
	if err != nil {
		return nil, err
	}

	positions := make([]Position, 0, length.Int64())
	for i := int64(0); i < length.Int64(); i++ {
		tokenID, err := c.contract.PositionIds(opts, big.NewInt(i))
		if err != nil {
			return nil, err
		}

		tickLower, err := c.contract.PositionTickLower(opts, tokenID)
		if err != nil {
			return nil, err
		}

		tickUpper, err := c.contract.PositionTickUpper(opts, tokenID)
		if err != nil {
			return nil, err
		}

		positions = append(positions, Position{
			TokenID:   tokenID,
			TickLower: int32(tickLower.Int64()),
			TickUpper: int32(tickUpper.Int64()),
		})
	}

	return positions, nil
}

// MintPosition mints a new LP position.
func (c *Client) MintPosition(ctx context.Context, tickLower, tickUpper int32, liquidity *big.Int, amount0Max, amount1Max *big.Int, deadline *big.Int) (*types.Transaction, error) {
	auth := c.authWithContext(ctx)
	return c.contract.MintPosition(auth, big.NewInt(int64(tickLower)), big.NewInt(int64(tickUpper)), liquidity, amount0Max, amount1Max, deadline)
}

// IncreaseLiquidity increases liquidity for an existing position.
func (c *Client) IncreaseLiquidity(ctx context.Context, tokenID, liquidity *big.Int, amount0Max, amount1Max *big.Int, deadline *big.Int) (*types.Transaction, error) {
	auth := c.authWithContext(ctx)
	return c.contract.IncreaseLiquidity(auth, tokenID, liquidity, amount0Max, amount1Max, deadline)
}

// DecreaseLiquidity decreases liquidity for an existing position.
func (c *Client) DecreaseLiquidity(ctx context.Context, tokenID, liquidity *big.Int, amount0Min, amount1Min *big.Int, deadline *big.Int) (*types.Transaction, error) {
	auth := c.authWithContext(ctx)
	return c.contract.DecreaseLiquidityToVault(auth, tokenID, liquidity, amount0Min, amount1Min, deadline)
}

// CollectFees collects fees from a position.
func (c *Client) CollectFees(ctx context.Context, tokenID *big.Int, amount0Min, amount1Min *big.Int, deadline *big.Int) (*types.Transaction, error) {
	auth := c.authWithContext(ctx)
	return c.contract.CollectFeesToVault(auth, tokenID, amount0Min, amount1Min, deadline)
}

// BurnPosition burns a position and returns tokens to vault.
func (c *Client) BurnPosition(ctx context.Context, tokenID *big.Int, amount0Min, amount1Min *big.Int, deadline *big.Int) (*types.Transaction, error) {
	auth := c.authWithContext(ctx)
	return c.contract.BurnPositionToVault(auth, tokenID, amount0Min, amount1Min, deadline)
}

// Swap executes a single-hop swap.
func (c *Client) Swap(ctx context.Context, zeroForOne bool, amountIn, minAmountOut *big.Int, deadline *big.Int) (*types.Transaction, error) {
	auth := c.authWithContext(ctx)
	return c.contract.SwapExactInputSingle(auth, zeroForOne, amountIn, minAmountOut, deadline)
}

// authWithContext creates a copy of auth with the given context.
func (c *Client) authWithContext(ctx context.Context) *bind.TransactOpts {
	if c.auth == nil {
		return nil
	}
	auth := *c.auth
	auth.Context = ctx
	return &auth
}

// Ensure Client implements Vault.
var _ Vault = (*Client)(nil)
