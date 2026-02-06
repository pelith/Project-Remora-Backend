package agent

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

// Minimal ABI for factory vault discovery.
const factoryABIJSON = `[
	{"inputs":[],"name":"getAllVaults","outputs":[{"internalType":"address[]","name":"","type":"address[]"}],"stateMutability":"view","type":"function"},
	{"inputs":[],"name":"totalVaults","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"stateMutability":"view","type":"function"},
	{"inputs":[{"internalType":"uint256","name":"","type":"uint256"}],"name":"vaults","outputs":[{"internalType":"address","name":"","type":"address"}],"stateMutability":"view","type":"function"}
]`

var factoryABI abi.ABI

func init() {
	var err error

	factoryABI, err = abi.JSON(strings.NewReader(factoryABIJSON))
	if err != nil {
		panic(fmt.Sprintf("failed to parse factory ABI: %v", err))
	}
}

// FactoryVaultSource reads vaults from the on-chain factory.
type FactoryVaultSource struct {
	client      *ethclient.Client
	factoryAddr common.Address
}

// NewFactoryVaultSource creates a vault source that queries the factory contract.
func NewFactoryVaultSource(client *ethclient.Client, factoryAddr common.Address) *FactoryVaultSource {
	return &FactoryVaultSource{
		client:      client,
		factoryAddr: factoryAddr,
	}
}

// GetVaultAddresses returns all vaults created by the factory.
func (s *FactoryVaultSource) GetVaultAddresses(ctx context.Context) ([]common.Address, error) {
	if s.factoryAddr == (common.Address{}) {
		return nil, errors.New("factory address not set")
	}

	// First try getAllVaults()
	if addrs, err := s.getAllVaults(ctx); err == nil {
		slog.Info("factory getAllVaults ok", slog.String("factory", s.factoryAddr.Hex()), slog.Int("vault_count", len(addrs)))
		return addrs, nil
	} else {
		slog.Warn("factory getAllVaults failed, fallback to totalVaults", slog.String("factory", s.factoryAddr.Hex()), slog.Any("error", err))
	}

	// Fallback: totalVaults + vaults(i)
	total, err := s.getTotalVaults(ctx)
	if err != nil {
		return nil, err
	}

	count := total.Int64()
	if count < 0 {
		return nil, fmt.Errorf("invalid totalVaults: %d", count)
	}

	addrs := make([]common.Address, 0, count)
	for i := range count {
		addr, err := s.getVaultByIndex(ctx, big.NewInt(i))
		if err != nil {
			return nil, err
		}

		addrs = append(addrs, addr)
	}

	slog.Info("factory vaults indexed", slog.String("factory", s.factoryAddr.Hex()), slog.Int64("vault_count", count))

	return addrs, nil
}

func (s *FactoryVaultSource) getAllVaults(ctx context.Context) ([]common.Address, error) {
	data, err := factoryABI.Pack("getAllVaults")
	if err != nil {
		return nil, fmt.Errorf("pack getAllVaults: %w", err)
	}

	msg := ethereum.CallMsg{To: &s.factoryAddr, Data: data}

	output, err := s.client.CallContract(ctx, msg, nil)
	if err != nil {
		return nil, fmt.Errorf("call getAllVaults: %w", err)
	}

	results, err := factoryABI.Unpack("getAllVaults", output)
	if err != nil {
		return nil, fmt.Errorf("unpack getAllVaults: %w", err)
	}

	if len(results) != 1 {
		return nil, fmt.Errorf("unexpected getAllVaults result length: %d", len(results))
	}

	addrs, ok := results[0].([]common.Address)
	if !ok {
		return nil, errors.New("unexpected getAllVaults result type")
	}

	return addrs, nil
}

func (s *FactoryVaultSource) getTotalVaults(ctx context.Context) (*big.Int, error) {
	data, err := factoryABI.Pack("totalVaults")
	if err != nil {
		return nil, fmt.Errorf("pack totalVaults: %w", err)
	}

	msg := ethereum.CallMsg{To: &s.factoryAddr, Data: data}

	output, err := s.client.CallContract(ctx, msg, nil)
	if err != nil {
		return nil, fmt.Errorf("call totalVaults: %w", err)
	}

	results, err := factoryABI.Unpack("totalVaults", output)
	if err != nil {
		return nil, fmt.Errorf("unpack totalVaults: %w", err)
	}

	if len(results) != 1 {
		return nil, fmt.Errorf("unexpected totalVaults result length: %d", len(results))
	}

	total, ok := results[0].(*big.Int)
	if !ok {
		return nil, errors.New("unexpected totalVaults result type")
	}

	return total, nil
}

func (s *FactoryVaultSource) getVaultByIndex(ctx context.Context, index *big.Int) (common.Address, error) {
	data, err := factoryABI.Pack("vaults", index)
	if err != nil {
		return common.Address{}, fmt.Errorf("pack vaults: %w", err)
	}

	msg := ethereum.CallMsg{To: &s.factoryAddr, Data: data}

	output, err := s.client.CallContract(ctx, msg, nil)
	if err != nil {
		return common.Address{}, fmt.Errorf("call vaults: %w", err)
	}

	results, err := factoryABI.Unpack("vaults", output)
	if err != nil {
		return common.Address{}, fmt.Errorf("unpack vaults: %w", err)
	}

	if len(results) != 1 {
		return common.Address{}, fmt.Errorf("unexpected vaults result length: %d", len(results))
	}

	addr, ok := results[0].(common.Address)
	if !ok {
		return common.Address{}, errors.New("unexpected vaults result type")
	}

	return addr, nil
}

// Ensure FactoryVaultSource implements VaultSource.
var _ VaultSource = (*FactoryVaultSource)(nil)
