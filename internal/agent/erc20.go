package agent

import (
	"context"
	"fmt"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
)

// ERC20 ABI (minimal) for balanceOf and decimals
const erc20ABIJSON = `[
	{"constant":true,"inputs":[{"name":"_owner","type":"address"}],"name":"balanceOf","outputs":[{"name":"balance","type":"uint256"}],"type":"function"},
	{"constant":true,"inputs":[],"name":"decimals","outputs":[{"name":"","type":"uint8"}],"type":"function"}
]`

var erc20ABI abi.ABI

func init() {
	var err error
	erc20ABI, err = abi.JSON(strings.NewReader(erc20ABIJSON))
	if err != nil {
		panic(fmt.Sprintf("failed to parse ERC20 ABI: %v", err))
	}
}

// getTokenDecimals fetches the decimals for a given token address.
func (s *Service) getTokenDecimals(ctx context.Context, token common.Address) (uint8, error) {
	// Pack the input for the "decimals" call
	data, err := erc20ABI.Pack("decimals")
	if err != nil {
		return 0, fmt.Errorf("pack decimals: %w", err)
	}

	// Create the call message
	msg := ethereum.CallMsg{
		To:   &token,
		Data: data,
	}

	// Execute the call
	output, err := s.ethClient.CallContract(ctx, msg, nil)
	if err != nil {
		return 0, fmt.Errorf("call decimals: %w", err)
	}

	// Unpack the output
	var decimals uint8
	results, err := erc20ABI.Unpack("decimals", output)
	if err != nil {
		return 0, fmt.Errorf("unpack decimals: %w", err)
	}
	
	// Handle flexible unpacking (sometimes returns []interface{})
	if len(results) > 0 {
		decimals = results[0].(uint8)
	} else {
		return 0, fmt.Errorf("unexpected empty result from decimals")
	}

	return decimals, nil
}

// getTokenBalance fetches the balance of a token for a specific owner.
func (s *Service) getTokenBalance(ctx context.Context, token common.Address, owner common.Address) (*big.Int, error) {
	// Pack the input for the "balanceOf" call
	data, err := erc20ABI.Pack("balanceOf", owner)
	if err != nil {
		return nil, fmt.Errorf("pack balanceOf: %w", err)
	}

	// Create the call message
	msg := ethereum.CallMsg{
		To:   &token,
		Data: data,
	}

	// Execute the call
	output, err := s.ethClient.CallContract(ctx, msg, nil)
	if err != nil {
		return nil, fmt.Errorf("call balanceOf: %w", err)
	}

	// Unpack the output
	var balance *big.Int
	results, err := erc20ABI.Unpack("balanceOf", output)
	if err != nil {
		return nil, fmt.Errorf("unpack balanceOf: %w", err)
	}

	if len(results) > 0 {
		balance = results[0].(*big.Int)
	} else {
		return nil, fmt.Errorf("unexpected empty result from balanceOf")
	}

	return balance, nil
}
