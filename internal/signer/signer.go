package signer

import (
	"crypto/ecdsa"
	"fmt"
	"math/big"
	"os"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

// Signer holds the private key and provides transaction signing capabilities.
type Signer struct {
	privateKey *ecdsa.PrivateKey
	address    common.Address
	chainID    *big.Int
}

// New creates a new Signer from a hex-encoded private key.
// The privateKey should be 64 hex characters without 0x prefix.
func New(privateKeyHex string, chainID *big.Int) (*Signer, error) {
	privateKey, err := crypto.HexToECDSA(privateKeyHex)
	if err != nil {
		return nil, fmt.Errorf("invalid private key: %w", err)
	}

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return nil, fmt.Errorf("failed to cast public key to ECDSA")
	}

	address := crypto.PubkeyToAddress(*publicKeyECDSA)

	return &Signer{
		privateKey: privateKey,
		address:    address,
		chainID:    chainID,
	}, nil
}

// Address returns the Ethereum address derived from the private key.
func (s *Signer) Address() common.Address {
	return s.address
}

// ChainID returns the chain ID.
func (s *Signer) ChainID() *big.Int {
	return s.chainID
}

// TransactOpts returns bind.TransactOpts for contract interactions.
func (s *Signer) TransactOpts() (*bind.TransactOpts, error) {
	auth, err := bind.NewKeyedTransactorWithChainID(s.privateKey, s.chainID)
	if err != nil {
		return nil, fmt.Errorf("create transactor: %w", err)
	}
	return auth, nil
}

// NewFromEnv creates a Signer from environment variables.
// Reads AGENT_PRIVATE_KEY and CHAIN_ID from environment.
func NewFromEnv() (*Signer, error) {
	privateKey := os.Getenv("AGENT_PRIVATE_KEY")
	if privateKey == "" {
		return nil, fmt.Errorf("AGENT_PRIVATE_KEY not set")
	}
	privateKey = strings.TrimPrefix(privateKey, "0x")

	chainIDStr := os.Getenv("CHAIN_ID")
	if chainIDStr == "" {
		return nil, fmt.Errorf("CHAIN_ID not set")
	}

	chainID, ok := new(big.Int).SetString(chainIDStr, 10)
	if !ok {
		return nil, fmt.Errorf("invalid CHAIN_ID: %s", chainIDStr)
	}

	return New(privateKey, chainID)
}
