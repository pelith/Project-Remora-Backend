package service

import (
	"context"
	"fmt"
	"math/big"
	"sort"

	"remora/internal/liquidity"
	"remora/internal/liquidity/poolid"
)

// Service implements liquidity.Service.
type Service struct {
	repo liquidity.Repository
}

// New creates a new liquidity service.
func New(repo liquidity.Repository) *Service {
	return &Service{
		repo: repo,
	}
}

// Ensure Service implements liquidity.Service.
var _ liquidity.Service = (*Service)(nil)

// GetDistribution returns the liquidity distribution for a pool.
func (s *Service) GetDistribution(ctx context.Context, params *liquidity.DistributionParams) (*liquidity.Distribution, error) {
	if err := validateParams(params); err != nil {
		return nil, err
	}

	poolKey := &params.PoolKey

	// Step 1: Get current pool state (tick and sqrtPrice)
	slot0, err := s.repo.GetSlot0(ctx, poolKey)
	if err != nil {
		return nil, fmt.Errorf("get slot0: %w", err)
	}

	// Step 2: Read initialized ticks in the specified range
	ticks, err := s.getInitializedTicks(ctx, poolKey, slot0.Tick, params.TickRange, poolKey.TickSpacing)
	if err != nil {
		return nil, fmt.Errorf("get initialized ticks: %w", err)
	}

	// Step 3: Calculate active liquidity using prefix sum
	activeLiquidities := s.calculateActiveLiquidity(ticks)

	// Step 4: Aggregate into discretized bins
	bins := s.aggregateBins(ticks, activeLiquidities, params.BinSizeTicks)

	return &liquidity.Distribution{
		CurrentTick:      slot0.Tick,
		SqrtPriceX96:     slot0.SqrtPriceX96.String(),
		InitializedTicks: ticks,
		Bins:             bins,
	}, nil
}

// validateParams validates distribution parameters.
func validateParams(params *liquidity.DistributionParams) error {
	if params.BinSizeTicks <= 0 {
		return liquidity.ErrInvalidBinSize
	}

	if params.TickRange <= 0 {
		return liquidity.ErrInvalidTickRange
	}

	if err := poolid.ValidatePoolKey(&params.PoolKey); err != nil {
		return fmt.Errorf("validate pool key: %w", err)
	}

	return nil
}

// getInitializedTicks retrieves all initialized ticks in the range [currentTick - tickRange, currentTick + tickRange].
func (s *Service) getInitializedTicks(ctx context.Context, poolKey *poolid.PoolKey, currentTick, tickRange, tickSpacing int32) ([]liquidity.TickInfo, error) {
	// Calculate the range of word positions to scan
	tickLower := currentTick - tickRange
	tickUpper := currentTick + tickRange

	// Calculate word positions (each word covers 256 ticks)
	wordPosLower := s.getWordPos(tickLower, tickSpacing)
	wordPosUpper := s.getWordPos(tickUpper, tickSpacing)

	var ticks []liquidity.TickInfo

	// Scan each word position
	for wordPos := wordPosLower; wordPos <= wordPosUpper; wordPos++ {
		bitmap, err := s.repo.GetTickBitmap(ctx, poolKey, wordPos)
		if err != nil {
			return nil, fmt.Errorf("get tick bitmap for wordPos %d: %w", wordPos, err)
		}

		// Parse bitmap to find initialized ticks
		initializedTicks := s.parseTickBitmap(bitmap, wordPos, tickSpacing)

		// Filter ticks within range
		for _, tick := range initializedTicks {
			if tick >= tickLower && tick <= tickUpper {
				tickInfo, err := s.repo.GetTickInfo(ctx, poolKey, tick)
				if err != nil {
					return nil, fmt.Errorf("get tick info for tick %d: %w", tick, err)
				}

				ticks = append(ticks, *tickInfo)
			}
		}
	}

	// Sort ticks in ascending order
	sort.Slice(ticks, func(i, j int) bool {
		return ticks[i].Tick < ticks[j].Tick
	})

	return ticks, nil
}

// getWordPos calculates the word position for a tick.
func (s *Service) getWordPos(tick, tickSpacing int32) int16 {
	// A word covers 256 ticks (2^8)
	const wordShift = 8

	compressed := tick / tickSpacing

	//nolint:gosec // G115: Overflow is expected and safe for word position calculation
	return int16(compressed >> wordShift)
}

// parseTickBitmap parses a bitmap to extract initialized tick positions.
func (s *Service) parseTickBitmap(bitmap *big.Int, wordPos int16, tickSpacing int32) []int32 {
	var ticks []int32

	const (
		bitsPerWord = 256
		wordShift   = 8
	)

	// Each bit in the bitmap represents a tick
	for bitPos := range bitsPerWord {
		if bitmap.Bit(bitPos) == 1 {
			// Calculate the actual tick from word position and bit position
			//nolint:gosec // G115: Overflow is expected and safe for tick calculation
			compressed := (int32(wordPos) << wordShift) + int32(bitPos)
			tick := compressed * tickSpacing
			ticks = append(ticks, tick)
		}
	}

	return ticks
}

// calculateActiveLiquidity calculates active liquidity using prefix sum.
func (s *Service) calculateActiveLiquidity(ticks []liquidity.TickInfo) []*big.Int {
	activeLiquidities := make([]*big.Int, len(ticks))
	currentLiquidity := big.NewInt(0)

	for i, tick := range ticks {
		// Update current liquidity by adding liquidityNet
		currentLiquidity = new(big.Int).Add(currentLiquidity, tick.LiquidityNet)

		// Store the active liquidity at this tick
		activeLiquidities[i] = new(big.Int).Set(currentLiquidity)
	}

	return activeLiquidities
}

// aggregateBins aggregates ticks into discretized bins.
func (s *Service) aggregateBins(ticks []liquidity.TickInfo, activeLiquidities []*big.Int, binSizeTicks int32) []liquidity.Bin {
	if len(ticks) == 0 {
		return []liquidity.Bin{}
	}

	var bins []liquidity.Bin

	// Find the min and max ticks
	minTick := ticks[0].Tick
	maxTick := ticks[len(ticks)-1].Tick

	const halfDivisor = 2

	// Create bins from minTick to maxTick
	for tickLower := minTick; tickLower <= maxTick; tickLower += binSizeTicks {
		tickUpper := tickLower + binSizeTicks

		// Find the representative active liquidity for this bin (use middle point)
		binMiddle := tickLower + binSizeTicks/halfDivisor
		activeLiquidity := s.getActiveLiquidityAtTick(ticks, activeLiquidities, binMiddle)

		bins = append(bins, liquidity.Bin{
			TickLower:       tickLower,
			TickUpper:       tickUpper,
			ActiveLiquidity: activeLiquidity,
		})
	}

	return bins
}

// getActiveLiquidityAtTick finds the active liquidity at a specific tick.
func (s *Service) getActiveLiquidityAtTick(ticks []liquidity.TickInfo, activeLiquidities []*big.Int, targetTick int32) *big.Int {
	// Binary search to find the largest tick <= targetTick
	idx := sort.Search(len(ticks), func(i int) bool {
		return ticks[i].Tick > targetTick
	})

	// If no tick found before targetTick, return 0
	if idx == 0 {
		return big.NewInt(0)
	}

	// Return the active liquidity at the previous tick
	return new(big.Int).Set(activeLiquidities[idx-1])
}
