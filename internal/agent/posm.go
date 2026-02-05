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

// Minimal ABI for POSM
const posmABIJSON = `[
	{"constant":true,"inputs":[{"name":"tokenId","type":"uint256"}],"name":"getPositionLiquidity","outputs":[{"name":"liquidity","type":"uint128"}],"type":"function"}
]`

var posmABI abi.ABI

func init() {
	var err error
	posmABI, err = abi.JSON(strings.NewReader(posmABIJSON))
	if err != nil {
		panic(fmt.Sprintf("failed to parse POSM ABI: %v", err))
	}
}

// getPositionLiquidity fetches the liquidity of a position via POSM's getPositionLiquidity(uint256).
func (s *Service) getPositionLiquidity(ctx context.Context, posmAddr common.Address, tokenID *big.Int) (*big.Int, error) {
	data, err := posmABI.Pack("getPositionLiquidity", tokenID)
	if err != nil {
		return nil, fmt.Errorf("pack getPositionLiquidity: %w", err)
	}

	msg := ethereum.CallMsg{To: &posmAddr, Data: data}
	output, err := s.ethClient.CallContract(ctx, msg, nil)
	if err != nil {
		return nil, fmt.Errorf("call getPositionLiquidity: %w", err)
	}

	results, err := posmABI.Unpack("getPositionLiquidity", output)
	if err != nil {
		return nil, fmt.Errorf("unpack getPositionLiquidity: %w", err)
	}

	return results[0].(*big.Int), nil
}
