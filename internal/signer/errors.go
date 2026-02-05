package signer

import "errors"

var (
	ErrPublicKeyNotECDSA     = errors.New("failed to cast public key to ECDSA")
	ErrAgentPrivateKeyNotSet = errors.New("AGENT_PRIVATE_KEY not set")
	ErrChainIDNotSet         = errors.New("CHAIN_ID not set")
	ErrInvalidChainID        = errors.New("invalid CHAIN_ID")
)
