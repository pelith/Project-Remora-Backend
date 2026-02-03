package repository

import "errors"

var (
	// ErrConnectionFailed is returned when RPC connection fails.
	ErrConnectionFailed = errors.New("rpc connection failed")

	// ErrContractNotFound is returned when contract is not found.
	ErrContractNotFound = errors.New("contract not found")
)
