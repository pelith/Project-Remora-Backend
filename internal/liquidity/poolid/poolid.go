package poolid

import (
	"errors"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

// ErrInvalidPoolKey is returned when pool key is invalid.
var ErrInvalidPoolKey = errors.New("invalid pool key")

// PoolKey identifies a Uniswap v4 pool.
// Must have currency0 < currency1 (by address).
type PoolKey struct {
	Currency0   string `json:"currency0"`
	Currency1   string `json:"currency1"`
	Fee         uint32 `json:"fee"`
	TickSpacing int32  `json:"tickSpacing"`
	Hooks       string `json:"hooks"`
}

// ValidatePoolKey checks that currency0 < currency1 (required by Uniswap v4).
func ValidatePoolKey(key *PoolKey) error {
	addr0 := common.HexToAddress(key.Currency0)
	addr1 := common.HexToAddress(key.Currency1)

	if addr0.Cmp(addr1) >= 0 {
		return ErrInvalidPoolKey
	}

	return nil
}

// CalculatePoolID computes the PoolId from a PoolKey per Uniswap v4 spec.
// PoolId = keccak256(abi.encode(poolKey)).
// See: https://github.com/Uniswap/v4-core/blob/main/src/types/PoolId.sol
func CalculatePoolID(key *PoolKey) [32]byte {
	const (
		abiSlotSize    = 32
		addressLen     = 20
		addressPadding = abiSlotSize - addressLen
		poolKeySlotCnt = 5  // PoolKey has 5 fields: currency0, currency1, fee, tickSpacing, hooks
		feeShiftHigh   = 16 // Bit shift for high byte of uint24
		feeShiftMid    = 8  // Bit shift for middle byte of uint24
		tickShiftHigh  = 16 // Bit shift for high byte of int24
		tickShiftMid   = 8  // Bit shift for middle byte of int24
	)

	data := make([]byte, 0, abiSlotSize*poolKeySlotCnt)

	// currency0 (address) - right-aligned in 32 bytes
	addr0 := common.HexToAddress(key.Currency0)
	slot0 := make([]byte, abiSlotSize)
	copy(slot0[addressPadding:], addr0.Bytes())
	data = append(data, slot0...)

	// currency1 (address) - right-aligned in 32 bytes
	addr1 := common.HexToAddress(key.Currency1)
	slot1 := make([]byte, abiSlotSize)
	copy(slot1[addressPadding:], addr1.Bytes())
	data = append(data, slot1...)

	// fee (uint24) - right-aligned in 32 bytes
	const uint24Offset = 29

	slot2 := make([]byte, abiSlotSize)
	slot2[uint24Offset] = byte(key.Fee >> feeShiftHigh)
	slot2[uint24Offset+1] = byte(key.Fee >> feeShiftMid)
	slot2[uint24Offset+2] = byte(key.Fee)
	data = append(data, slot2...)

	// tickSpacing (int24) - sign-extended, right-aligned in 32 bytes
	const signExtendBytes = 29

	slot3 := make([]byte, abiSlotSize)

	if key.TickSpacing < 0 {
		for i := range signExtendBytes {
			slot3[i] = 0xFF
		}
	}

	slot3[uint24Offset] = byte(key.TickSpacing >> tickShiftHigh)
	slot3[uint24Offset+1] = byte(key.TickSpacing >> tickShiftMid)
	slot3[uint24Offset+2] = byte(key.TickSpacing)
	data = append(data, slot3...)

	// hooks (address) - right-aligned in 32 bytes
	hooks := common.HexToAddress(key.Hooks)
	slot4 := make([]byte, abiSlotSize)
	copy(slot4[addressPadding:], hooks.Bytes())
	data = append(data, slot4...)

	hash := crypto.Keccak256Hash(data)

	return [32]byte(hash)
}
