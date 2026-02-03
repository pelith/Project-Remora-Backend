package liquidity

import "errors"

var (
	// ErrInvalidPoolKey is returned when pool key is invalid.
	ErrInvalidPoolKey = errors.New("invalid pool key")

	// ErrInvalidBinSize is returned when bin size is invalid.
	ErrInvalidBinSize = errors.New("bin size must be positive")

	// ErrInvalidTickRange is returned when tick range is invalid.
	ErrInvalidTickRange = errors.New("tick range must be positive")

	// ErrContractCall is returned when contract call fails.
	ErrContractCall = errors.New("contract call failed")
)
