package vault

import (
	"context"
	"fmt"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
)

// Minimal ABI for POSM getPositionLiquidity (shared with agent for consistency).
const posmABIJSON = `[
	{"constant":true,"inputs":[{"name":"tokenId","type":"uint256"}],"name":"getPositionLiquidity","outputs":[{"name":"liquidity","type":"uint128"}],"type":"function"}
]`

var posmABI abi.ABI

func init() {
	var err error

	posmABI, err = abi.JSON(strings.NewReader(posmABIJSON))
	if err != nil {
		panic(fmt.Sprintf("parse POSM ABI: %v", err))
	}
}

// getPositionLiquidity fetches the liquidity of a position via POSM's getPositionLiquidity(uint256).
// Caller can be vault.Client's caller or any bind.ContractCaller (e.g. ethclient).
func getPositionLiquidity(ctx context.Context, caller bind.ContractCaller, posmAddr common.Address, tokenID *big.Int) (*big.Int, error) {
	data, err := posmABI.Pack("getPositionLiquidity", tokenID)
	if err != nil {
		return nil, fmt.Errorf("pack getPositionLiquidity: %w", err)
	}

	msg := ethereum.CallMsg{To: &posmAddr, Data: data}

	output, err := caller.CallContract(ctx, msg, nil)
	if err != nil {
		return nil, fmt.Errorf("call getPositionLiquidity: %w", err)
	}

	results, err := posmABI.Unpack("getPositionLiquidity", output)
	if err != nil {
		return nil, fmt.Errorf("unpack getPositionLiquidity: %w", err)
	}

	liq, ok := results[0].(*big.Int)
	if !ok {
		return nil, fmt.Errorf("unexpected getPositionLiquidity result type: %T", results[0])
	}

	return liq, nil
}
