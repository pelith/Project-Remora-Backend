package liquidity

import "errors"

var (
	// ErrInvalidPoolKey is returned when pool key is invalid.
	ErrInvalidPoolKey = errors.New("invalid pool key")

	// ErrInvalidBinSize is returned when bin size is invalid.
	ErrInvalidBinSize = errors.New("bin size must be positive")

	// ErrInvalidTickRange is returned when tick range is invalid.
	ErrInvalidTickRange = errors.New("invalid tick range")

	// ErrContractCall is returned when contract call fails.
	ErrContractCall = errors.New("contract call failed")

	ErrNoTickRanges       = errors.New("no tick ranges")
	ErrInvalidWeight      = errors.New("invalid weight")
	ErrInvalidTotalAmount = errors.New("invalid total amount")
	ErrInvalidSlippage    = errors.New("invalid slippage tolerance")
	ErrZeroTotalWeight    = errors.New("zero total weight")
)
