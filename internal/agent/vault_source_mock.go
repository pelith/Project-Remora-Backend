package agent

import (
	"context"

	"github.com/ethereum/go-ethereum/common"
)

// MockVaultSource is a mock implementation of VaultSource for testing.
type MockVaultSource struct {
	addresses []common.Address
}

// NewMockVaultSource creates a mock vault source with predefined addresses.
func NewMockVaultSource(addresses []common.Address) *MockVaultSource {
	return &MockVaultSource{addresses: addresses}
}

// GetVaultAddresses returns the list of vault addresses.
func (m *MockVaultSource) GetVaultAddresses(ctx context.Context) ([]common.Address, error) {
	return m.addresses, nil
}

// Ensure MockVaultSource implements VaultSource.
var _ VaultSource = (*MockVaultSource)(nil)
