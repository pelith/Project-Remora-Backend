package repository

import (
	"context"
	"fmt"
	"math/big"
	"strings"
	"sync"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"

	"remora/internal/liquidity"
	"remora/internal/liquidity/poolid"
)

const multicallChunkSize = 30 // ticks per multicall batch

// Multicall3 is deployed at a deterministic address on all EVM chains.
var multicall3Addr = common.HexToAddress("0xcA11bde05977b3631167028862bE2a173976CA11")

// Minimal ABIs for encoding/decoding.
var (
	multicall3ABI abi.ABI
	stateViewABI  abi.ABI
)

func init() {
	var err error

	multicall3ABI, err = abi.JSON(strings.NewReader(`[{
		"name":"aggregate3",
		"type":"function",
		"inputs":[{"name":"calls","type":"tuple[]","components":[
			{"name":"target","type":"address"},
			{"name":"allowFailure","type":"bool"},
			{"name":"callData","type":"bytes"}
		]}],
		"outputs":[{"name":"returnData","type":"tuple[]","components":[
			{"name":"success","type":"bool"},
			{"name":"returnData","type":"bytes"}
		]}]
	}]`))
	if err != nil {
		panic("parse multicall3 abi: " + err.Error())
	}

	stateViewABI, err = abi.JSON(strings.NewReader(`[{
		"name":"getTickInfo",
		"type":"function",
		"inputs":[
			{"name":"poolId","type":"bytes32"},
			{"name":"tick","type":"int24"}
		],
		"outputs":[
			{"name":"liquidityGross","type":"uint128"},
			{"name":"liquidityNet","type":"int128"},
			{"name":"feeGrowthOutside0X128","type":"uint256"},
			{"name":"feeGrowthOutside1X128","type":"uint256"}
		]
	}]`))
	if err != nil {
		panic("parse stateview abi: " + err.Error())
	}
}

// multicall3Call matches the Multicall3.Call3 struct layout for ABI encoding.
type multicall3Call struct {
	Target       common.Address
	AllowFailure bool
	CallData     []byte
}

// GetTickInfoBatch fetches tick info for multiple ticks using Multicall3.
// Ticks are split into chunks and fetched in parallel to avoid slow single-call
// execution on forked nodes where each storage read hits the remote RPC.
func (r *Repository) GetTickInfoBatch(ctx context.Context, poolKey *poolid.PoolKey, ticks []int32) ([]liquidity.TickInfo, error) {
	if len(ticks) == 0 {
		return nil, nil
	}

	// Mock mode for testing
	if r.contract == nil {
		result := make([]liquidity.TickInfo, len(ticks))
		for i, t := range ticks {
			result[i] = liquidity.TickInfo{
				Tick:           t,
				LiquidityGross: big.NewInt(0),
				LiquidityNet:   big.NewInt(0),
			}
		}

		return result, nil
	}

	poolID := poolid.CalculatePoolID(poolKey)

	// Split into chunks and fetch in parallel.
	type chunkResult struct {
		index int
		infos []liquidity.TickInfo
		err   error
	}

	chunks := chunkSlice(ticks, multicallChunkSize)
	results := make(chan chunkResult, len(chunks))

	var wg sync.WaitGroup
	for i, chunk := range chunks {
		wg.Add(1)

		go func(idx int, tickChunk []int32) {
			defer wg.Done()

			infos, err := r.fetchTickInfoChunk(ctx, poolID, tickChunk)
			results <- chunkResult{index: idx, infos: infos, err: err}
		}(i, chunk)
	}

	wg.Wait()
	close(results)

	// Reassemble in order.
	ordered := make([][]liquidity.TickInfo, len(chunks))

	for res := range results {
		if res.err != nil {
			return nil, res.err
		}

		ordered[res.index] = res.infos
	}

	tickInfos := make([]liquidity.TickInfo, 0, len(ticks))
	for _, infos := range ordered {
		tickInfos = append(tickInfos, infos...)
	}

	return tickInfos, nil
}

// fetchTickInfoChunk executes a single multicall3 batch for a chunk of ticks.
func (r *Repository) fetchTickInfoChunk(ctx context.Context, poolID [32]byte, ticks []int32) ([]liquidity.TickInfo, error) {
	calls := make([]multicall3Call, len(ticks))
	for i, tick := range ticks {
		callData, err := stateViewABI.Pack("getTickInfo", poolID, big.NewInt(int64(tick)))
		if err != nil {
			return nil, fmt.Errorf("pack getTickInfo for tick %d: %w", tick, err)
		}

		calls[i] = multicall3Call{
			Target:       r.stateViewAddr,
			AllowFailure: false,
			CallData:     callData,
		}
	}

	input, err := multicall3ABI.Pack("aggregate3", calls)
	if err != nil {
		return nil, fmt.Errorf("pack aggregate3: %w", err)
	}

	output, err := r.client.CallContract(ctx, ethereum.CallMsg{
		To:   &multicall3Addr,
		Data: input,
	}, nil)
	if err != nil {
		return nil, fmt.Errorf("multicall3 aggregate3: %w", err)
	}

	decoded, err := multicall3ABI.Unpack("aggregate3", output)
	if err != nil {
		return nil, fmt.Errorf("unpack aggregate3: %w", err)
	}

	resultsRaw, ok := decoded[0].([]struct {
		Success    bool   `json:"success"`
		ReturnData []byte `json:"returnData"`
	})
	if !ok {
		return nil, fmt.Errorf("unexpected aggregate3 result type: %T", decoded[0])
	}

	if len(resultsRaw) != len(ticks) {
		return nil, fmt.Errorf("multicall3 returned %d results, expected %d", len(resultsRaw), len(ticks))
	}

	tickInfos := make([]liquidity.TickInfo, len(ticks))

	for i, res := range resultsRaw {
		if !res.Success {
			return nil, fmt.Errorf("multicall3 call failed for tick %d", ticks[i])
		}

		values, err := stateViewABI.Methods["getTickInfo"].Outputs.Unpack(res.ReturnData)
		if err != nil {
			return nil, fmt.Errorf("unpack getTickInfo result for tick %d: %w", ticks[i], err)
		}

		liquidityGross, _ := values[0].(*big.Int)
		liquidityNet, _ := values[1].(*big.Int)

		tickInfos[i] = liquidity.TickInfo{
			Tick:           ticks[i],
			LiquidityGross: liquidityGross,
			LiquidityNet:   liquidityNet,
		}
	}

	return tickInfos, nil
}

// chunkSlice splits a slice into chunks of the given size.
func chunkSlice(ticks []int32, size int) [][]int32 {
	var chunks [][]int32

	for i := 0; i < len(ticks); i += size {
		end := i + size
		if end > len(ticks) {
			end = len(ticks)
		}

		chunks = append(chunks, ticks[i:end])
	}

	return chunks
}
