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

// Minimal ABIs for POSM and StateView
const (
	posmABIJSON = `[
		{"constant":true,"inputs":[{"name":"tokenId","type":"uint256"}],"name":"positions","outputs":[
			{"name":"poolKey","type":"tuple","components":[
				{"name":"currency0","type":"address"},
				{"name":"currency1","type":"address"},
				{"name":"fee","type":"uint24"},
				{"name":"tickSpacing","type":"int24"},
				{"name":"hooks","type":"address"}
			]},
			{"name":"tickLower","type":"int24"},
			{"name":"tickUpper","type":"int24"},
			{"name":"salt","type":"uint256"}
		],"type":"function"}
	]`

	stateViewABIJSON = `[
		{"constant":true,"inputs":[{"name":"poolId","type":"bytes32"},{"name":"owner","type":"address"},{"name":"tickLower","type":"int24"},{"name":"tickUpper","type":"int24"},{"name":"salt","type":"bytes32"}],"name":"getPositionInfo","outputs":[
			{"name":"liquidity","type":"uint128"},
			{"name":"feeGrowthInside0LastX128","type":"uint256"},
			{"name":"feeGrowthInside1LastX128","type":"uint256"}
		],"type":"function"}
	]`
)

var (
	posmABI        abi.ABI
	stateViewABI   abi.ABI
)

func init() {
	var err error
	posmABI, err = abi.JSON(strings.NewReader(posmABIJSON))
	if err != nil {
		panic(fmt.Sprintf("failed to parse POSM ABI: %v", err))
	}
	stateViewABI, err = abi.JSON(strings.NewReader(stateViewABIJSON))
	if err != nil {
		panic(fmt.Sprintf("failed to parse StateView ABI: %v", err))
	}
}

// getPositionLiquidity fetches the liquidity of a position using POSM and StateView.
func (s *Service) getPositionLiquidity(ctx context.Context, posmAddr common.Address, poolId [32]byte, tokenID *big.Int) (*big.Int, error) {
	// 1. Get salt and tick range from POSM
	data, err := posmABI.Pack("positions", tokenID)
	if err != nil {
		return nil, fmt.Errorf("pack positions: %w", err)
	}

	msg := ethereum.CallMsg{To: &posmAddr, Data: data}
	output, err := s.ethClient.CallContract(ctx, msg, nil)
	if err != nil {
		return nil, fmt.Errorf("call positions: %w", err)
	}

	results, err := posmABI.Unpack("positions", output)
	if err != nil {
		return nil, fmt.Errorf("unpack positions: %w", err)
	}

	// results: [PoolKey, tickLower, tickUpper, salt]
	// tickLower is results[1], tickUpper is results[2], salt is results[3]
	tickLower := results[1].(*big.Int)
	tickUpper := results[2].(*big.Int)
	salt := results[3].(*big.Int)

	// 2. Convert salt to [32]byte for StateView
	saltBytes := [32]byte{}
	salt.FillBytes(saltBytes[:])

	// 3. Get liquidity from StateView
	// owner is the POSM address
	if s.stateViewAddr == (common.Address{}) {
		return nil, fmt.Errorf("stateView address not configured")
	}

	svData, err := stateViewABI.Pack("getPositionInfo", poolId, posmAddr, tickLower, tickUpper, saltBytes)
	if err != nil {
		return nil, fmt.Errorf("pack getPositionInfo: %w", err)
	}

	svMsg := ethereum.CallMsg{To: &s.stateViewAddr, Data: svData}
	svOutput, err := s.ethClient.CallContract(ctx, svMsg, nil)
	if err != nil {
		return nil, fmt.Errorf("call getPositionInfo: %w", err)
	}

	svResults, err := stateViewABI.Unpack("getPositionInfo", svOutput)
	if err != nil {
		return nil, fmt.Errorf("unpack getPositionInfo: %w", err)
	}

	// svResults[0] is liquidity (uint128)
	return svResults[0].(*big.Int), nil
}