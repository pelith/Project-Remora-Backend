// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package vault

import (
	"errors"
	"math/big"
	"strings"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"
)

// Reference imports to suppress errors if they are not otherwise used.
var (
	_ = errors.New
	_ = big.NewInt
	_ = strings.NewReader
	_ = ethereum.NotFound
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
	_ = abi.ConvertType
)

// PoolKey is an auto generated low-level Go binding around an user-defined struct.
type PoolKey struct {
	Currency0   common.Address
	Currency1   common.Address
	Fee         *big.Int
	TickSpacing *big.Int
	Hooks       common.Address
}

// V4AgenticVaultMetaData contains all meta data concerning the V4AgenticVault contract.
var V4AgenticVaultMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"constructor\",\"inputs\":[{\"name\":\"_owner\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"_agent\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"_posm\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"_universalRouter\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"_permit2\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"_poolKey\",\"type\":\"tuple\",\"internalType\":\"structPoolKey\",\"components\":[{\"name\":\"currency0\",\"type\":\"address\",\"internalType\":\"Currency\"},{\"name\":\"currency1\",\"type\":\"address\",\"internalType\":\"Currency\"},{\"name\":\"fee\",\"type\":\"uint24\",\"internalType\":\"uint24\"},{\"name\":\"tickSpacing\",\"type\":\"int24\",\"internalType\":\"int24\"},{\"name\":\"hooks\",\"type\":\"address\",\"internalType\":\"contractIHooks\"}]},{\"name\":\"_initialAllowedTickLower\",\"type\":\"int24\",\"internalType\":\"int24\"},{\"name\":\"_initialAllowedTickUpper\",\"type\":\"int24\",\"internalType\":\"int24\"},{\"name\":\"_swapAllowed\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"_maxPositionsK\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"receive\",\"stateMutability\":\"payable\"},{\"type\":\"function\",\"name\":\"agent\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"agentPaused\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"allowedTickLower\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"int24\",\"internalType\":\"int24\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"allowedTickUpper\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"int24\",\"internalType\":\"int24\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"approveTokenWithPermit2\",\"inputs\":[{\"name\":\"currency\",\"type\":\"address\",\"internalType\":\"Currency\"},{\"name\":\"spender\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint160\",\"internalType\":\"uint160\"},{\"name\":\"expiration\",\"type\":\"uint48\",\"internalType\":\"uint48\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"burnPositionToVault\",\"inputs\":[{\"name\":\"tokenId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"amount0Min\",\"type\":\"uint128\",\"internalType\":\"uint128\"},{\"name\":\"amount1Min\",\"type\":\"uint128\",\"internalType\":\"uint128\"},{\"name\":\"deadline\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"collectFeesToVault\",\"inputs\":[{\"name\":\"tokenId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"amount0Min\",\"type\":\"uint128\",\"internalType\":\"uint128\"},{\"name\":\"amount1Min\",\"type\":\"uint128\",\"internalType\":\"uint128\"},{\"name\":\"deadline\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"currency0\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"Currency\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"currency1\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"Currency\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"decreaseLiquidityToVault\",\"inputs\":[{\"name\":\"tokenId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"liquidity\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"amount0Min\",\"type\":\"uint128\",\"internalType\":\"uint128\"},{\"name\":\"amount1Min\",\"type\":\"uint128\",\"internalType\":\"uint128\"},{\"name\":\"deadline\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"fee\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint24\",\"internalType\":\"uint24\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getPoolKey\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"tuple\",\"internalType\":\"structPoolKey\",\"components\":[{\"name\":\"currency0\",\"type\":\"address\",\"internalType\":\"Currency\"},{\"name\":\"currency1\",\"type\":\"address\",\"internalType\":\"Currency\"},{\"name\":\"fee\",\"type\":\"uint24\",\"internalType\":\"uint24\"},{\"name\":\"tickSpacing\",\"type\":\"int24\",\"internalType\":\"int24\"},{\"name\":\"hooks\",\"type\":\"address\",\"internalType\":\"contractIHooks\"}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"hooks\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"increaseLiquidity\",\"inputs\":[{\"name\":\"tokenId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"liquidity\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"amount0Max\",\"type\":\"uint128\",\"internalType\":\"uint128\"},{\"name\":\"amount1Max\",\"type\":\"uint128\",\"internalType\":\"uint128\"},{\"name\":\"deadline\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"isManagedPosition\",\"inputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"maxPositionsK\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"mintPosition\",\"inputs\":[{\"name\":\"tickLower\",\"type\":\"int24\",\"internalType\":\"int24\"},{\"name\":\"tickUpper\",\"type\":\"int24\",\"internalType\":\"int24\"},{\"name\":\"liquidity\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"amount0Max\",\"type\":\"uint128\",\"internalType\":\"uint128\"},{\"name\":\"amount1Max\",\"type\":\"uint128\",\"internalType\":\"uint128\"},{\"name\":\"deadline\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"tokenId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"onERC721Received\",\"inputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bytes4\",\"internalType\":\"bytes4\"}],\"stateMutability\":\"pure\"},{\"type\":\"function\",\"name\":\"owner\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"pauseAndExitAll\",\"inputs\":[{\"name\":\"deadline\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"permit2\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractIPermit2\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"poolId\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"PoolId\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"positionIds\",\"inputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"positionTickLower\",\"inputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"int24\",\"internalType\":\"int24\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"positionTickUpper\",\"inputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"int24\",\"internalType\":\"int24\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"positionsLength\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"posm\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractIPositionManager\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"renounceOwnership\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setAgent\",\"inputs\":[{\"name\":\"newAgent\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setAgentPaused\",\"inputs\":[{\"name\":\"paused\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setAllowedTickRange\",\"inputs\":[{\"name\":\"tickLower\",\"type\":\"int24\",\"internalType\":\"int24\"},{\"name\":\"tickUpper\",\"type\":\"int24\",\"internalType\":\"int24\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setMaxPositionsK\",\"inputs\":[{\"name\":\"k\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setSwapAllowed\",\"inputs\":[{\"name\":\"allowed\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"swapAllowed\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"swapExactInputSingle\",\"inputs\":[{\"name\":\"zeroForOne\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"amountIn\",\"type\":\"uint128\",\"internalType\":\"uint128\"},{\"name\":\"minAmountOut\",\"type\":\"uint128\",\"internalType\":\"uint128\"},{\"name\":\"deadline\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"amountOut\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"tickSpacing\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"int24\",\"internalType\":\"int24\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"transferOwnership\",\"inputs\":[{\"name\":\"newOwner\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"universalRouter\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractIUniversalRouter\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"withdraw\",\"inputs\":[{\"name\":\"currency\",\"type\":\"address\",\"internalType\":\"Currency\"},{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"to\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"event\",\"name\":\"AgentPaused\",\"inputs\":[{\"name\":\"paused\",\"type\":\"bool\",\"indexed\":false,\"internalType\":\"bool\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"AgentUpdated\",\"inputs\":[{\"name\":\"newAgent\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"AllowedTickRangeUpdated\",\"inputs\":[{\"name\":\"tickLower\",\"type\":\"int24\",\"indexed\":false,\"internalType\":\"int24\"},{\"name\":\"tickUpper\",\"type\":\"int24\",\"indexed\":false,\"internalType\":\"int24\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"MaxPositionsKUpdated\",\"inputs\":[{\"name\":\"k\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"OwnershipTransferred\",\"inputs\":[{\"name\":\"previousOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"newOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"PositionAdded\",\"inputs\":[{\"name\":\"tokenId\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"tickLower\",\"type\":\"int24\",\"indexed\":false,\"internalType\":\"int24\"},{\"name\":\"tickUpper\",\"type\":\"int24\",\"indexed\":false,\"internalType\":\"int24\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"PositionRemoved\",\"inputs\":[{\"name\":\"tokenId\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"SwapAllowed\",\"inputs\":[{\"name\":\"allowed\",\"type\":\"bool\",\"indexed\":false,\"internalType\":\"bool\"}],\"anonymous\":false},{\"type\":\"error\",\"name\":\"OwnableInvalidOwner\",\"inputs\":[{\"name\":\"owner\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"OwnableUnauthorizedAccount\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ReentrancyGuardReentrantCall\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"SafeERC20FailedOperation\",\"inputs\":[{\"name\":\"token\",\"type\":\"address\",\"internalType\":\"address\"}]}]",
}

// V4AgenticVaultABI is the input ABI used to generate the binding from.
// Deprecated: Use V4AgenticVaultMetaData.ABI instead.
var V4AgenticVaultABI = V4AgenticVaultMetaData.ABI

// V4AgenticVault is an auto generated Go binding around an Ethereum contract.
type V4AgenticVault struct {
	V4AgenticVaultCaller     // Read-only binding to the contract
	V4AgenticVaultTransactor // Write-only binding to the contract
	V4AgenticVaultFilterer   // Log filterer for contract events
}

// V4AgenticVaultCaller is an auto generated read-only Go binding around an Ethereum contract.
type V4AgenticVaultCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// V4AgenticVaultTransactor is an auto generated write-only Go binding around an Ethereum contract.
type V4AgenticVaultTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// V4AgenticVaultFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type V4AgenticVaultFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// V4AgenticVaultSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type V4AgenticVaultSession struct {
	Contract     *V4AgenticVault   // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// V4AgenticVaultCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type V4AgenticVaultCallerSession struct {
	Contract *V4AgenticVaultCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts         // Call options to use throughout this session
}

// V4AgenticVaultTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type V4AgenticVaultTransactorSession struct {
	Contract     *V4AgenticVaultTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts         // Transaction auth options to use throughout this session
}

// V4AgenticVaultRaw is an auto generated low-level Go binding around an Ethereum contract.
type V4AgenticVaultRaw struct {
	Contract *V4AgenticVault // Generic contract binding to access the raw methods on
}

// V4AgenticVaultCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type V4AgenticVaultCallerRaw struct {
	Contract *V4AgenticVaultCaller // Generic read-only contract binding to access the raw methods on
}

// V4AgenticVaultTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type V4AgenticVaultTransactorRaw struct {
	Contract *V4AgenticVaultTransactor // Generic write-only contract binding to access the raw methods on
}

// NewV4AgenticVault creates a new instance of V4AgenticVault, bound to a specific deployed contract.
func NewV4AgenticVault(address common.Address, backend bind.ContractBackend) (*V4AgenticVault, error) {
	contract, err := bindV4AgenticVault(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &V4AgenticVault{V4AgenticVaultCaller: V4AgenticVaultCaller{contract: contract}, V4AgenticVaultTransactor: V4AgenticVaultTransactor{contract: contract}, V4AgenticVaultFilterer: V4AgenticVaultFilterer{contract: contract}}, nil
}

// NewV4AgenticVaultCaller creates a new read-only instance of V4AgenticVault, bound to a specific deployed contract.
func NewV4AgenticVaultCaller(address common.Address, caller bind.ContractCaller) (*V4AgenticVaultCaller, error) {
	contract, err := bindV4AgenticVault(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &V4AgenticVaultCaller{contract: contract}, nil
}

// NewV4AgenticVaultTransactor creates a new write-only instance of V4AgenticVault, bound to a specific deployed contract.
func NewV4AgenticVaultTransactor(address common.Address, transactor bind.ContractTransactor) (*V4AgenticVaultTransactor, error) {
	contract, err := bindV4AgenticVault(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &V4AgenticVaultTransactor{contract: contract}, nil
}

// NewV4AgenticVaultFilterer creates a new log filterer instance of V4AgenticVault, bound to a specific deployed contract.
func NewV4AgenticVaultFilterer(address common.Address, filterer bind.ContractFilterer) (*V4AgenticVaultFilterer, error) {
	contract, err := bindV4AgenticVault(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &V4AgenticVaultFilterer{contract: contract}, nil
}

// bindV4AgenticVault binds a generic wrapper to an already deployed contract.
func bindV4AgenticVault(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := V4AgenticVaultMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_V4AgenticVault *V4AgenticVaultRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _V4AgenticVault.Contract.V4AgenticVaultCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_V4AgenticVault *V4AgenticVaultRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _V4AgenticVault.Contract.V4AgenticVaultTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_V4AgenticVault *V4AgenticVaultRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _V4AgenticVault.Contract.V4AgenticVaultTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_V4AgenticVault *V4AgenticVaultCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _V4AgenticVault.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_V4AgenticVault *V4AgenticVaultTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _V4AgenticVault.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_V4AgenticVault *V4AgenticVaultTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _V4AgenticVault.Contract.contract.Transact(opts, method, params...)
}

// Agent is a free data retrieval call binding the contract method 0xf5ff5c76.
//
// Solidity: function agent() view returns(address)
func (_V4AgenticVault *V4AgenticVaultCaller) Agent(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _V4AgenticVault.contract.Call(opts, &out, "agent")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Agent is a free data retrieval call binding the contract method 0xf5ff5c76.
//
// Solidity: function agent() view returns(address)
func (_V4AgenticVault *V4AgenticVaultSession) Agent() (common.Address, error) {
	return _V4AgenticVault.Contract.Agent(&_V4AgenticVault.CallOpts)
}

// Agent is a free data retrieval call binding the contract method 0xf5ff5c76.
//
// Solidity: function agent() view returns(address)
func (_V4AgenticVault *V4AgenticVaultCallerSession) Agent() (common.Address, error) {
	return _V4AgenticVault.Contract.Agent(&_V4AgenticVault.CallOpts)
}

// AgentPaused is a free data retrieval call binding the contract method 0x0ba0d3de.
//
// Solidity: function agentPaused() view returns(bool)
func (_V4AgenticVault *V4AgenticVaultCaller) AgentPaused(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _V4AgenticVault.contract.Call(opts, &out, "agentPaused")

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// AgentPaused is a free data retrieval call binding the contract method 0x0ba0d3de.
//
// Solidity: function agentPaused() view returns(bool)
func (_V4AgenticVault *V4AgenticVaultSession) AgentPaused() (bool, error) {
	return _V4AgenticVault.Contract.AgentPaused(&_V4AgenticVault.CallOpts)
}

// AgentPaused is a free data retrieval call binding the contract method 0x0ba0d3de.
//
// Solidity: function agentPaused() view returns(bool)
func (_V4AgenticVault *V4AgenticVaultCallerSession) AgentPaused() (bool, error) {
	return _V4AgenticVault.Contract.AgentPaused(&_V4AgenticVault.CallOpts)
}

// AllowedTickLower is a free data retrieval call binding the contract method 0x01d32bf5.
//
// Solidity: function allowedTickLower() view returns(int24)
func (_V4AgenticVault *V4AgenticVaultCaller) AllowedTickLower(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _V4AgenticVault.contract.Call(opts, &out, "allowedTickLower")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// AllowedTickLower is a free data retrieval call binding the contract method 0x01d32bf5.
//
// Solidity: function allowedTickLower() view returns(int24)
func (_V4AgenticVault *V4AgenticVaultSession) AllowedTickLower() (*big.Int, error) {
	return _V4AgenticVault.Contract.AllowedTickLower(&_V4AgenticVault.CallOpts)
}

// AllowedTickLower is a free data retrieval call binding the contract method 0x01d32bf5.
//
// Solidity: function allowedTickLower() view returns(int24)
func (_V4AgenticVault *V4AgenticVaultCallerSession) AllowedTickLower() (*big.Int, error) {
	return _V4AgenticVault.Contract.AllowedTickLower(&_V4AgenticVault.CallOpts)
}

// AllowedTickUpper is a free data retrieval call binding the contract method 0xc6dc6db0.
//
// Solidity: function allowedTickUpper() view returns(int24)
func (_V4AgenticVault *V4AgenticVaultCaller) AllowedTickUpper(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _V4AgenticVault.contract.Call(opts, &out, "allowedTickUpper")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// AllowedTickUpper is a free data retrieval call binding the contract method 0xc6dc6db0.
//
// Solidity: function allowedTickUpper() view returns(int24)
func (_V4AgenticVault *V4AgenticVaultSession) AllowedTickUpper() (*big.Int, error) {
	return _V4AgenticVault.Contract.AllowedTickUpper(&_V4AgenticVault.CallOpts)
}

// AllowedTickUpper is a free data retrieval call binding the contract method 0xc6dc6db0.
//
// Solidity: function allowedTickUpper() view returns(int24)
func (_V4AgenticVault *V4AgenticVaultCallerSession) AllowedTickUpper() (*big.Int, error) {
	return _V4AgenticVault.Contract.AllowedTickUpper(&_V4AgenticVault.CallOpts)
}

// Currency0 is a free data retrieval call binding the contract method 0x79f1232b.
//
// Solidity: function currency0() view returns(address)
func (_V4AgenticVault *V4AgenticVaultCaller) Currency0(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _V4AgenticVault.contract.Call(opts, &out, "currency0")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Currency0 is a free data retrieval call binding the contract method 0x79f1232b.
//
// Solidity: function currency0() view returns(address)
func (_V4AgenticVault *V4AgenticVaultSession) Currency0() (common.Address, error) {
	return _V4AgenticVault.Contract.Currency0(&_V4AgenticVault.CallOpts)
}

// Currency0 is a free data retrieval call binding the contract method 0x79f1232b.
//
// Solidity: function currency0() view returns(address)
func (_V4AgenticVault *V4AgenticVaultCallerSession) Currency0() (common.Address, error) {
	return _V4AgenticVault.Contract.Currency0(&_V4AgenticVault.CallOpts)
}

// Currency1 is a free data retrieval call binding the contract method 0x10d737b8.
//
// Solidity: function currency1() view returns(address)
func (_V4AgenticVault *V4AgenticVaultCaller) Currency1(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _V4AgenticVault.contract.Call(opts, &out, "currency1")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Currency1 is a free data retrieval call binding the contract method 0x10d737b8.
//
// Solidity: function currency1() view returns(address)
func (_V4AgenticVault *V4AgenticVaultSession) Currency1() (common.Address, error) {
	return _V4AgenticVault.Contract.Currency1(&_V4AgenticVault.CallOpts)
}

// Currency1 is a free data retrieval call binding the contract method 0x10d737b8.
//
// Solidity: function currency1() view returns(address)
func (_V4AgenticVault *V4AgenticVaultCallerSession) Currency1() (common.Address, error) {
	return _V4AgenticVault.Contract.Currency1(&_V4AgenticVault.CallOpts)
}

// Fee is a free data retrieval call binding the contract method 0xddca3f43.
//
// Solidity: function fee() view returns(uint24)
func (_V4AgenticVault *V4AgenticVaultCaller) Fee(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _V4AgenticVault.contract.Call(opts, &out, "fee")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Fee is a free data retrieval call binding the contract method 0xddca3f43.
//
// Solidity: function fee() view returns(uint24)
func (_V4AgenticVault *V4AgenticVaultSession) Fee() (*big.Int, error) {
	return _V4AgenticVault.Contract.Fee(&_V4AgenticVault.CallOpts)
}

// Fee is a free data retrieval call binding the contract method 0xddca3f43.
//
// Solidity: function fee() view returns(uint24)
func (_V4AgenticVault *V4AgenticVaultCallerSession) Fee() (*big.Int, error) {
	return _V4AgenticVault.Contract.Fee(&_V4AgenticVault.CallOpts)
}

// GetPoolKey is a free data retrieval call binding the contract method 0x683e76e0.
//
// Solidity: function getPoolKey() view returns((address,address,uint24,int24,address))
func (_V4AgenticVault *V4AgenticVaultCaller) GetPoolKey(opts *bind.CallOpts) (PoolKey, error) {
	var out []interface{}
	err := _V4AgenticVault.contract.Call(opts, &out, "getPoolKey")

	if err != nil {
		return *new(PoolKey), err
	}

	out0 := *abi.ConvertType(out[0], new(PoolKey)).(*PoolKey)

	return out0, err

}

// GetPoolKey is a free data retrieval call binding the contract method 0x683e76e0.
//
// Solidity: function getPoolKey() view returns((address,address,uint24,int24,address))
func (_V4AgenticVault *V4AgenticVaultSession) GetPoolKey() (PoolKey, error) {
	return _V4AgenticVault.Contract.GetPoolKey(&_V4AgenticVault.CallOpts)
}

// GetPoolKey is a free data retrieval call binding the contract method 0x683e76e0.
//
// Solidity: function getPoolKey() view returns((address,address,uint24,int24,address))
func (_V4AgenticVault *V4AgenticVaultCallerSession) GetPoolKey() (PoolKey, error) {
	return _V4AgenticVault.Contract.GetPoolKey(&_V4AgenticVault.CallOpts)
}

// Hooks is a free data retrieval call binding the contract method 0xcd7033c4.
//
// Solidity: function hooks() view returns(address)
func (_V4AgenticVault *V4AgenticVaultCaller) Hooks(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _V4AgenticVault.contract.Call(opts, &out, "hooks")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Hooks is a free data retrieval call binding the contract method 0xcd7033c4.
//
// Solidity: function hooks() view returns(address)
func (_V4AgenticVault *V4AgenticVaultSession) Hooks() (common.Address, error) {
	return _V4AgenticVault.Contract.Hooks(&_V4AgenticVault.CallOpts)
}

// Hooks is a free data retrieval call binding the contract method 0xcd7033c4.
//
// Solidity: function hooks() view returns(address)
func (_V4AgenticVault *V4AgenticVaultCallerSession) Hooks() (common.Address, error) {
	return _V4AgenticVault.Contract.Hooks(&_V4AgenticVault.CallOpts)
}

// IsManagedPosition is a free data retrieval call binding the contract method 0x8041c480.
//
// Solidity: function isManagedPosition(uint256 ) view returns(bool)
func (_V4AgenticVault *V4AgenticVaultCaller) IsManagedPosition(opts *bind.CallOpts, arg0 *big.Int) (bool, error) {
	var out []interface{}
	err := _V4AgenticVault.contract.Call(opts, &out, "isManagedPosition", arg0)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsManagedPosition is a free data retrieval call binding the contract method 0x8041c480.
//
// Solidity: function isManagedPosition(uint256 ) view returns(bool)
func (_V4AgenticVault *V4AgenticVaultSession) IsManagedPosition(arg0 *big.Int) (bool, error) {
	return _V4AgenticVault.Contract.IsManagedPosition(&_V4AgenticVault.CallOpts, arg0)
}

// IsManagedPosition is a free data retrieval call binding the contract method 0x8041c480.
//
// Solidity: function isManagedPosition(uint256 ) view returns(bool)
func (_V4AgenticVault *V4AgenticVaultCallerSession) IsManagedPosition(arg0 *big.Int) (bool, error) {
	return _V4AgenticVault.Contract.IsManagedPosition(&_V4AgenticVault.CallOpts, arg0)
}

// MaxPositionsK is a free data retrieval call binding the contract method 0x140b3d4c.
//
// Solidity: function maxPositionsK() view returns(uint256)
func (_V4AgenticVault *V4AgenticVaultCaller) MaxPositionsK(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _V4AgenticVault.contract.Call(opts, &out, "maxPositionsK")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// MaxPositionsK is a free data retrieval call binding the contract method 0x140b3d4c.
//
// Solidity: function maxPositionsK() view returns(uint256)
func (_V4AgenticVault *V4AgenticVaultSession) MaxPositionsK() (*big.Int, error) {
	return _V4AgenticVault.Contract.MaxPositionsK(&_V4AgenticVault.CallOpts)
}

// MaxPositionsK is a free data retrieval call binding the contract method 0x140b3d4c.
//
// Solidity: function maxPositionsK() view returns(uint256)
func (_V4AgenticVault *V4AgenticVaultCallerSession) MaxPositionsK() (*big.Int, error) {
	return _V4AgenticVault.Contract.MaxPositionsK(&_V4AgenticVault.CallOpts)
}

// OnERC721Received is a free data retrieval call binding the contract method 0x150b7a02.
//
// Solidity: function onERC721Received(address , address , uint256 , bytes ) pure returns(bytes4)
func (_V4AgenticVault *V4AgenticVaultCaller) OnERC721Received(opts *bind.CallOpts, arg0 common.Address, arg1 common.Address, arg2 *big.Int, arg3 []byte) ([4]byte, error) {
	var out []interface{}
	err := _V4AgenticVault.contract.Call(opts, &out, "onERC721Received", arg0, arg1, arg2, arg3)

	if err != nil {
		return *new([4]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([4]byte)).(*[4]byte)

	return out0, err

}

// OnERC721Received is a free data retrieval call binding the contract method 0x150b7a02.
//
// Solidity: function onERC721Received(address , address , uint256 , bytes ) pure returns(bytes4)
func (_V4AgenticVault *V4AgenticVaultSession) OnERC721Received(arg0 common.Address, arg1 common.Address, arg2 *big.Int, arg3 []byte) ([4]byte, error) {
	return _V4AgenticVault.Contract.OnERC721Received(&_V4AgenticVault.CallOpts, arg0, arg1, arg2, arg3)
}

// OnERC721Received is a free data retrieval call binding the contract method 0x150b7a02.
//
// Solidity: function onERC721Received(address , address , uint256 , bytes ) pure returns(bytes4)
func (_V4AgenticVault *V4AgenticVaultCallerSession) OnERC721Received(arg0 common.Address, arg1 common.Address, arg2 *big.Int, arg3 []byte) ([4]byte, error) {
	return _V4AgenticVault.Contract.OnERC721Received(&_V4AgenticVault.CallOpts, arg0, arg1, arg2, arg3)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_V4AgenticVault *V4AgenticVaultCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _V4AgenticVault.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_V4AgenticVault *V4AgenticVaultSession) Owner() (common.Address, error) {
	return _V4AgenticVault.Contract.Owner(&_V4AgenticVault.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_V4AgenticVault *V4AgenticVaultCallerSession) Owner() (common.Address, error) {
	return _V4AgenticVault.Contract.Owner(&_V4AgenticVault.CallOpts)
}

// Permit2 is a free data retrieval call binding the contract method 0x12261ee7.
//
// Solidity: function permit2() view returns(address)
func (_V4AgenticVault *V4AgenticVaultCaller) Permit2(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _V4AgenticVault.contract.Call(opts, &out, "permit2")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Permit2 is a free data retrieval call binding the contract method 0x12261ee7.
//
// Solidity: function permit2() view returns(address)
func (_V4AgenticVault *V4AgenticVaultSession) Permit2() (common.Address, error) {
	return _V4AgenticVault.Contract.Permit2(&_V4AgenticVault.CallOpts)
}

// Permit2 is a free data retrieval call binding the contract method 0x12261ee7.
//
// Solidity: function permit2() view returns(address)
func (_V4AgenticVault *V4AgenticVaultCallerSession) Permit2() (common.Address, error) {
	return _V4AgenticVault.Contract.Permit2(&_V4AgenticVault.CallOpts)
}

// PoolId is a free data retrieval call binding the contract method 0x3e0dc34e.
//
// Solidity: function poolId() view returns(bytes32)
func (_V4AgenticVault *V4AgenticVaultCaller) PoolId(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _V4AgenticVault.contract.Call(opts, &out, "poolId")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// PoolId is a free data retrieval call binding the contract method 0x3e0dc34e.
//
// Solidity: function poolId() view returns(bytes32)
func (_V4AgenticVault *V4AgenticVaultSession) PoolId() ([32]byte, error) {
	return _V4AgenticVault.Contract.PoolId(&_V4AgenticVault.CallOpts)
}

// PoolId is a free data retrieval call binding the contract method 0x3e0dc34e.
//
// Solidity: function poolId() view returns(bytes32)
func (_V4AgenticVault *V4AgenticVaultCallerSession) PoolId() ([32]byte, error) {
	return _V4AgenticVault.Contract.PoolId(&_V4AgenticVault.CallOpts)
}

// PositionIds is a free data retrieval call binding the contract method 0x939c5f7a.
//
// Solidity: function positionIds(uint256 ) view returns(uint256)
func (_V4AgenticVault *V4AgenticVaultCaller) PositionIds(opts *bind.CallOpts, arg0 *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _V4AgenticVault.contract.Call(opts, &out, "positionIds", arg0)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// PositionIds is a free data retrieval call binding the contract method 0x939c5f7a.
//
// Solidity: function positionIds(uint256 ) view returns(uint256)
func (_V4AgenticVault *V4AgenticVaultSession) PositionIds(arg0 *big.Int) (*big.Int, error) {
	return _V4AgenticVault.Contract.PositionIds(&_V4AgenticVault.CallOpts, arg0)
}

// PositionIds is a free data retrieval call binding the contract method 0x939c5f7a.
//
// Solidity: function positionIds(uint256 ) view returns(uint256)
func (_V4AgenticVault *V4AgenticVaultCallerSession) PositionIds(arg0 *big.Int) (*big.Int, error) {
	return _V4AgenticVault.Contract.PositionIds(&_V4AgenticVault.CallOpts, arg0)
}

// PositionTickLower is a free data retrieval call binding the contract method 0x2e3f4461.
//
// Solidity: function positionTickLower(uint256 ) view returns(int24)
func (_V4AgenticVault *V4AgenticVaultCaller) PositionTickLower(opts *bind.CallOpts, arg0 *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _V4AgenticVault.contract.Call(opts, &out, "positionTickLower", arg0)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// PositionTickLower is a free data retrieval call binding the contract method 0x2e3f4461.
//
// Solidity: function positionTickLower(uint256 ) view returns(int24)
func (_V4AgenticVault *V4AgenticVaultSession) PositionTickLower(arg0 *big.Int) (*big.Int, error) {
	return _V4AgenticVault.Contract.PositionTickLower(&_V4AgenticVault.CallOpts, arg0)
}

// PositionTickLower is a free data retrieval call binding the contract method 0x2e3f4461.
//
// Solidity: function positionTickLower(uint256 ) view returns(int24)
func (_V4AgenticVault *V4AgenticVaultCallerSession) PositionTickLower(arg0 *big.Int) (*big.Int, error) {
	return _V4AgenticVault.Contract.PositionTickLower(&_V4AgenticVault.CallOpts, arg0)
}

// PositionTickUpper is a free data retrieval call binding the contract method 0x36ced346.
//
// Solidity: function positionTickUpper(uint256 ) view returns(int24)
func (_V4AgenticVault *V4AgenticVaultCaller) PositionTickUpper(opts *bind.CallOpts, arg0 *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _V4AgenticVault.contract.Call(opts, &out, "positionTickUpper", arg0)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// PositionTickUpper is a free data retrieval call binding the contract method 0x36ced346.
//
// Solidity: function positionTickUpper(uint256 ) view returns(int24)
func (_V4AgenticVault *V4AgenticVaultSession) PositionTickUpper(arg0 *big.Int) (*big.Int, error) {
	return _V4AgenticVault.Contract.PositionTickUpper(&_V4AgenticVault.CallOpts, arg0)
}

// PositionTickUpper is a free data retrieval call binding the contract method 0x36ced346.
//
// Solidity: function positionTickUpper(uint256 ) view returns(int24)
func (_V4AgenticVault *V4AgenticVaultCallerSession) PositionTickUpper(arg0 *big.Int) (*big.Int, error) {
	return _V4AgenticVault.Contract.PositionTickUpper(&_V4AgenticVault.CallOpts, arg0)
}

// PositionsLength is a free data retrieval call binding the contract method 0xd6887bfa.
//
// Solidity: function positionsLength() view returns(uint256)
func (_V4AgenticVault *V4AgenticVaultCaller) PositionsLength(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _V4AgenticVault.contract.Call(opts, &out, "positionsLength")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// PositionsLength is a free data retrieval call binding the contract method 0xd6887bfa.
//
// Solidity: function positionsLength() view returns(uint256)
func (_V4AgenticVault *V4AgenticVaultSession) PositionsLength() (*big.Int, error) {
	return _V4AgenticVault.Contract.PositionsLength(&_V4AgenticVault.CallOpts)
}

// PositionsLength is a free data retrieval call binding the contract method 0xd6887bfa.
//
// Solidity: function positionsLength() view returns(uint256)
func (_V4AgenticVault *V4AgenticVaultCallerSession) PositionsLength() (*big.Int, error) {
	return _V4AgenticVault.Contract.PositionsLength(&_V4AgenticVault.CallOpts)
}

// Posm is a free data retrieval call binding the contract method 0x6f70a7fa.
//
// Solidity: function posm() view returns(address)
func (_V4AgenticVault *V4AgenticVaultCaller) Posm(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _V4AgenticVault.contract.Call(opts, &out, "posm")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Posm is a free data retrieval call binding the contract method 0x6f70a7fa.
//
// Solidity: function posm() view returns(address)
func (_V4AgenticVault *V4AgenticVaultSession) Posm() (common.Address, error) {
	return _V4AgenticVault.Contract.Posm(&_V4AgenticVault.CallOpts)
}

// Posm is a free data retrieval call binding the contract method 0x6f70a7fa.
//
// Solidity: function posm() view returns(address)
func (_V4AgenticVault *V4AgenticVaultCallerSession) Posm() (common.Address, error) {
	return _V4AgenticVault.Contract.Posm(&_V4AgenticVault.CallOpts)
}

// SwapAllowed is a free data retrieval call binding the contract method 0x172869c4.
//
// Solidity: function swapAllowed() view returns(bool)
func (_V4AgenticVault *V4AgenticVaultCaller) SwapAllowed(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _V4AgenticVault.contract.Call(opts, &out, "swapAllowed")

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// SwapAllowed is a free data retrieval call binding the contract method 0x172869c4.
//
// Solidity: function swapAllowed() view returns(bool)
func (_V4AgenticVault *V4AgenticVaultSession) SwapAllowed() (bool, error) {
	return _V4AgenticVault.Contract.SwapAllowed(&_V4AgenticVault.CallOpts)
}

// SwapAllowed is a free data retrieval call binding the contract method 0x172869c4.
//
// Solidity: function swapAllowed() view returns(bool)
func (_V4AgenticVault *V4AgenticVaultCallerSession) SwapAllowed() (bool, error) {
	return _V4AgenticVault.Contract.SwapAllowed(&_V4AgenticVault.CallOpts)
}

// TickSpacing is a free data retrieval call binding the contract method 0xd0c93a7c.
//
// Solidity: function tickSpacing() view returns(int24)
func (_V4AgenticVault *V4AgenticVaultCaller) TickSpacing(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _V4AgenticVault.contract.Call(opts, &out, "tickSpacing")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// TickSpacing is a free data retrieval call binding the contract method 0xd0c93a7c.
//
// Solidity: function tickSpacing() view returns(int24)
func (_V4AgenticVault *V4AgenticVaultSession) TickSpacing() (*big.Int, error) {
	return _V4AgenticVault.Contract.TickSpacing(&_V4AgenticVault.CallOpts)
}

// TickSpacing is a free data retrieval call binding the contract method 0xd0c93a7c.
//
// Solidity: function tickSpacing() view returns(int24)
func (_V4AgenticVault *V4AgenticVaultCallerSession) TickSpacing() (*big.Int, error) {
	return _V4AgenticVault.Contract.TickSpacing(&_V4AgenticVault.CallOpts)
}

// UniversalRouter is a free data retrieval call binding the contract method 0x35a9e4df.
//
// Solidity: function universalRouter() view returns(address)
func (_V4AgenticVault *V4AgenticVaultCaller) UniversalRouter(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _V4AgenticVault.contract.Call(opts, &out, "universalRouter")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// UniversalRouter is a free data retrieval call binding the contract method 0x35a9e4df.
//
// Solidity: function universalRouter() view returns(address)
func (_V4AgenticVault *V4AgenticVaultSession) UniversalRouter() (common.Address, error) {
	return _V4AgenticVault.Contract.UniversalRouter(&_V4AgenticVault.CallOpts)
}

// UniversalRouter is a free data retrieval call binding the contract method 0x35a9e4df.
//
// Solidity: function universalRouter() view returns(address)
func (_V4AgenticVault *V4AgenticVaultCallerSession) UniversalRouter() (common.Address, error) {
	return _V4AgenticVault.Contract.UniversalRouter(&_V4AgenticVault.CallOpts)
}

// ApproveTokenWithPermit2 is a paid mutator transaction binding the contract method 0x7c9684ef.
//
// Solidity: function approveTokenWithPermit2(address currency, address spender, uint160 amount, uint48 expiration) returns()
func (_V4AgenticVault *V4AgenticVaultTransactor) ApproveTokenWithPermit2(opts *bind.TransactOpts, currency common.Address, spender common.Address, amount *big.Int, expiration *big.Int) (*types.Transaction, error) {
	return _V4AgenticVault.contract.Transact(opts, "approveTokenWithPermit2", currency, spender, amount, expiration)
}

// ApproveTokenWithPermit2 is a paid mutator transaction binding the contract method 0x7c9684ef.
//
// Solidity: function approveTokenWithPermit2(address currency, address spender, uint160 amount, uint48 expiration) returns()
func (_V4AgenticVault *V4AgenticVaultSession) ApproveTokenWithPermit2(currency common.Address, spender common.Address, amount *big.Int, expiration *big.Int) (*types.Transaction, error) {
	return _V4AgenticVault.Contract.ApproveTokenWithPermit2(&_V4AgenticVault.TransactOpts, currency, spender, amount, expiration)
}

// ApproveTokenWithPermit2 is a paid mutator transaction binding the contract method 0x7c9684ef.
//
// Solidity: function approveTokenWithPermit2(address currency, address spender, uint160 amount, uint48 expiration) returns()
func (_V4AgenticVault *V4AgenticVaultTransactorSession) ApproveTokenWithPermit2(currency common.Address, spender common.Address, amount *big.Int, expiration *big.Int) (*types.Transaction, error) {
	return _V4AgenticVault.Contract.ApproveTokenWithPermit2(&_V4AgenticVault.TransactOpts, currency, spender, amount, expiration)
}

// BurnPositionToVault is a paid mutator transaction binding the contract method 0x18c1653f.
//
// Solidity: function burnPositionToVault(uint256 tokenId, uint128 amount0Min, uint128 amount1Min, uint256 deadline) returns()
func (_V4AgenticVault *V4AgenticVaultTransactor) BurnPositionToVault(opts *bind.TransactOpts, tokenId *big.Int, amount0Min *big.Int, amount1Min *big.Int, deadline *big.Int) (*types.Transaction, error) {
	return _V4AgenticVault.contract.Transact(opts, "burnPositionToVault", tokenId, amount0Min, amount1Min, deadline)
}

// BurnPositionToVault is a paid mutator transaction binding the contract method 0x18c1653f.
//
// Solidity: function burnPositionToVault(uint256 tokenId, uint128 amount0Min, uint128 amount1Min, uint256 deadline) returns()
func (_V4AgenticVault *V4AgenticVaultSession) BurnPositionToVault(tokenId *big.Int, amount0Min *big.Int, amount1Min *big.Int, deadline *big.Int) (*types.Transaction, error) {
	return _V4AgenticVault.Contract.BurnPositionToVault(&_V4AgenticVault.TransactOpts, tokenId, amount0Min, amount1Min, deadline)
}

// BurnPositionToVault is a paid mutator transaction binding the contract method 0x18c1653f.
//
// Solidity: function burnPositionToVault(uint256 tokenId, uint128 amount0Min, uint128 amount1Min, uint256 deadline) returns()
func (_V4AgenticVault *V4AgenticVaultTransactorSession) BurnPositionToVault(tokenId *big.Int, amount0Min *big.Int, amount1Min *big.Int, deadline *big.Int) (*types.Transaction, error) {
	return _V4AgenticVault.Contract.BurnPositionToVault(&_V4AgenticVault.TransactOpts, tokenId, amount0Min, amount1Min, deadline)
}

// CollectFeesToVault is a paid mutator transaction binding the contract method 0x8696caae.
//
// Solidity: function collectFeesToVault(uint256 tokenId, uint128 amount0Min, uint128 amount1Min, uint256 deadline) returns()
func (_V4AgenticVault *V4AgenticVaultTransactor) CollectFeesToVault(opts *bind.TransactOpts, tokenId *big.Int, amount0Min *big.Int, amount1Min *big.Int, deadline *big.Int) (*types.Transaction, error) {
	return _V4AgenticVault.contract.Transact(opts, "collectFeesToVault", tokenId, amount0Min, amount1Min, deadline)
}

// CollectFeesToVault is a paid mutator transaction binding the contract method 0x8696caae.
//
// Solidity: function collectFeesToVault(uint256 tokenId, uint128 amount0Min, uint128 amount1Min, uint256 deadline) returns()
func (_V4AgenticVault *V4AgenticVaultSession) CollectFeesToVault(tokenId *big.Int, amount0Min *big.Int, amount1Min *big.Int, deadline *big.Int) (*types.Transaction, error) {
	return _V4AgenticVault.Contract.CollectFeesToVault(&_V4AgenticVault.TransactOpts, tokenId, amount0Min, amount1Min, deadline)
}

// CollectFeesToVault is a paid mutator transaction binding the contract method 0x8696caae.
//
// Solidity: function collectFeesToVault(uint256 tokenId, uint128 amount0Min, uint128 amount1Min, uint256 deadline) returns()
func (_V4AgenticVault *V4AgenticVaultTransactorSession) CollectFeesToVault(tokenId *big.Int, amount0Min *big.Int, amount1Min *big.Int, deadline *big.Int) (*types.Transaction, error) {
	return _V4AgenticVault.Contract.CollectFeesToVault(&_V4AgenticVault.TransactOpts, tokenId, amount0Min, amount1Min, deadline)
}

// DecreaseLiquidityToVault is a paid mutator transaction binding the contract method 0xcabc348d.
//
// Solidity: function decreaseLiquidityToVault(uint256 tokenId, uint256 liquidity, uint128 amount0Min, uint128 amount1Min, uint256 deadline) returns()
func (_V4AgenticVault *V4AgenticVaultTransactor) DecreaseLiquidityToVault(opts *bind.TransactOpts, tokenId *big.Int, liquidity *big.Int, amount0Min *big.Int, amount1Min *big.Int, deadline *big.Int) (*types.Transaction, error) {
	return _V4AgenticVault.contract.Transact(opts, "decreaseLiquidityToVault", tokenId, liquidity, amount0Min, amount1Min, deadline)
}

// DecreaseLiquidityToVault is a paid mutator transaction binding the contract method 0xcabc348d.
//
// Solidity: function decreaseLiquidityToVault(uint256 tokenId, uint256 liquidity, uint128 amount0Min, uint128 amount1Min, uint256 deadline) returns()
func (_V4AgenticVault *V4AgenticVaultSession) DecreaseLiquidityToVault(tokenId *big.Int, liquidity *big.Int, amount0Min *big.Int, amount1Min *big.Int, deadline *big.Int) (*types.Transaction, error) {
	return _V4AgenticVault.Contract.DecreaseLiquidityToVault(&_V4AgenticVault.TransactOpts, tokenId, liquidity, amount0Min, amount1Min, deadline)
}

// DecreaseLiquidityToVault is a paid mutator transaction binding the contract method 0xcabc348d.
//
// Solidity: function decreaseLiquidityToVault(uint256 tokenId, uint256 liquidity, uint128 amount0Min, uint128 amount1Min, uint256 deadline) returns()
func (_V4AgenticVault *V4AgenticVaultTransactorSession) DecreaseLiquidityToVault(tokenId *big.Int, liquidity *big.Int, amount0Min *big.Int, amount1Min *big.Int, deadline *big.Int) (*types.Transaction, error) {
	return _V4AgenticVault.Contract.DecreaseLiquidityToVault(&_V4AgenticVault.TransactOpts, tokenId, liquidity, amount0Min, amount1Min, deadline)
}

// IncreaseLiquidity is a paid mutator transaction binding the contract method 0x61f88c73.
//
// Solidity: function increaseLiquidity(uint256 tokenId, uint256 liquidity, uint128 amount0Max, uint128 amount1Max, uint256 deadline) returns()
func (_V4AgenticVault *V4AgenticVaultTransactor) IncreaseLiquidity(opts *bind.TransactOpts, tokenId *big.Int, liquidity *big.Int, amount0Max *big.Int, amount1Max *big.Int, deadline *big.Int) (*types.Transaction, error) {
	return _V4AgenticVault.contract.Transact(opts, "increaseLiquidity", tokenId, liquidity, amount0Max, amount1Max, deadline)
}

// IncreaseLiquidity is a paid mutator transaction binding the contract method 0x61f88c73.
//
// Solidity: function increaseLiquidity(uint256 tokenId, uint256 liquidity, uint128 amount0Max, uint128 amount1Max, uint256 deadline) returns()
func (_V4AgenticVault *V4AgenticVaultSession) IncreaseLiquidity(tokenId *big.Int, liquidity *big.Int, amount0Max *big.Int, amount1Max *big.Int, deadline *big.Int) (*types.Transaction, error) {
	return _V4AgenticVault.Contract.IncreaseLiquidity(&_V4AgenticVault.TransactOpts, tokenId, liquidity, amount0Max, amount1Max, deadline)
}

// IncreaseLiquidity is a paid mutator transaction binding the contract method 0x61f88c73.
//
// Solidity: function increaseLiquidity(uint256 tokenId, uint256 liquidity, uint128 amount0Max, uint128 amount1Max, uint256 deadline) returns()
func (_V4AgenticVault *V4AgenticVaultTransactorSession) IncreaseLiquidity(tokenId *big.Int, liquidity *big.Int, amount0Max *big.Int, amount1Max *big.Int, deadline *big.Int) (*types.Transaction, error) {
	return _V4AgenticVault.Contract.IncreaseLiquidity(&_V4AgenticVault.TransactOpts, tokenId, liquidity, amount0Max, amount1Max, deadline)
}

// MintPosition is a paid mutator transaction binding the contract method 0xd7364b09.
//
// Solidity: function mintPosition(int24 tickLower, int24 tickUpper, uint256 liquidity, uint128 amount0Max, uint128 amount1Max, uint256 deadline) returns(uint256 tokenId)
func (_V4AgenticVault *V4AgenticVaultTransactor) MintPosition(opts *bind.TransactOpts, tickLower *big.Int, tickUpper *big.Int, liquidity *big.Int, amount0Max *big.Int, amount1Max *big.Int, deadline *big.Int) (*types.Transaction, error) {
	return _V4AgenticVault.contract.Transact(opts, "mintPosition", tickLower, tickUpper, liquidity, amount0Max, amount1Max, deadline)
}

// MintPosition is a paid mutator transaction binding the contract method 0xd7364b09.
//
// Solidity: function mintPosition(int24 tickLower, int24 tickUpper, uint256 liquidity, uint128 amount0Max, uint128 amount1Max, uint256 deadline) returns(uint256 tokenId)
func (_V4AgenticVault *V4AgenticVaultSession) MintPosition(tickLower *big.Int, tickUpper *big.Int, liquidity *big.Int, amount0Max *big.Int, amount1Max *big.Int, deadline *big.Int) (*types.Transaction, error) {
	return _V4AgenticVault.Contract.MintPosition(&_V4AgenticVault.TransactOpts, tickLower, tickUpper, liquidity, amount0Max, amount1Max, deadline)
}

// MintPosition is a paid mutator transaction binding the contract method 0xd7364b09.
//
// Solidity: function mintPosition(int24 tickLower, int24 tickUpper, uint256 liquidity, uint128 amount0Max, uint128 amount1Max, uint256 deadline) returns(uint256 tokenId)
func (_V4AgenticVault *V4AgenticVaultTransactorSession) MintPosition(tickLower *big.Int, tickUpper *big.Int, liquidity *big.Int, amount0Max *big.Int, amount1Max *big.Int, deadline *big.Int) (*types.Transaction, error) {
	return _V4AgenticVault.Contract.MintPosition(&_V4AgenticVault.TransactOpts, tickLower, tickUpper, liquidity, amount0Max, amount1Max, deadline)
}

// PauseAndExitAll is a paid mutator transaction binding the contract method 0x38785dcd.
//
// Solidity: function pauseAndExitAll(uint256 deadline) returns()
func (_V4AgenticVault *V4AgenticVaultTransactor) PauseAndExitAll(opts *bind.TransactOpts, deadline *big.Int) (*types.Transaction, error) {
	return _V4AgenticVault.contract.Transact(opts, "pauseAndExitAll", deadline)
}

// PauseAndExitAll is a paid mutator transaction binding the contract method 0x38785dcd.
//
// Solidity: function pauseAndExitAll(uint256 deadline) returns()
func (_V4AgenticVault *V4AgenticVaultSession) PauseAndExitAll(deadline *big.Int) (*types.Transaction, error) {
	return _V4AgenticVault.Contract.PauseAndExitAll(&_V4AgenticVault.TransactOpts, deadline)
}

// PauseAndExitAll is a paid mutator transaction binding the contract method 0x38785dcd.
//
// Solidity: function pauseAndExitAll(uint256 deadline) returns()
func (_V4AgenticVault *V4AgenticVaultTransactorSession) PauseAndExitAll(deadline *big.Int) (*types.Transaction, error) {
	return _V4AgenticVault.Contract.PauseAndExitAll(&_V4AgenticVault.TransactOpts, deadline)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_V4AgenticVault *V4AgenticVaultTransactor) RenounceOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _V4AgenticVault.contract.Transact(opts, "renounceOwnership")
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_V4AgenticVault *V4AgenticVaultSession) RenounceOwnership() (*types.Transaction, error) {
	return _V4AgenticVault.Contract.RenounceOwnership(&_V4AgenticVault.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_V4AgenticVault *V4AgenticVaultTransactorSession) RenounceOwnership() (*types.Transaction, error) {
	return _V4AgenticVault.Contract.RenounceOwnership(&_V4AgenticVault.TransactOpts)
}

// SetAgent is a paid mutator transaction binding the contract method 0xbcf685ed.
//
// Solidity: function setAgent(address newAgent) returns()
func (_V4AgenticVault *V4AgenticVaultTransactor) SetAgent(opts *bind.TransactOpts, newAgent common.Address) (*types.Transaction, error) {
	return _V4AgenticVault.contract.Transact(opts, "setAgent", newAgent)
}

// SetAgent is a paid mutator transaction binding the contract method 0xbcf685ed.
//
// Solidity: function setAgent(address newAgent) returns()
func (_V4AgenticVault *V4AgenticVaultSession) SetAgent(newAgent common.Address) (*types.Transaction, error) {
	return _V4AgenticVault.Contract.SetAgent(&_V4AgenticVault.TransactOpts, newAgent)
}

// SetAgent is a paid mutator transaction binding the contract method 0xbcf685ed.
//
// Solidity: function setAgent(address newAgent) returns()
func (_V4AgenticVault *V4AgenticVaultTransactorSession) SetAgent(newAgent common.Address) (*types.Transaction, error) {
	return _V4AgenticVault.Contract.SetAgent(&_V4AgenticVault.TransactOpts, newAgent)
}

// SetAgentPaused is a paid mutator transaction binding the contract method 0x7c75005b.
//
// Solidity: function setAgentPaused(bool paused) returns()
func (_V4AgenticVault *V4AgenticVaultTransactor) SetAgentPaused(opts *bind.TransactOpts, paused bool) (*types.Transaction, error) {
	return _V4AgenticVault.contract.Transact(opts, "setAgentPaused", paused)
}

// SetAgentPaused is a paid mutator transaction binding the contract method 0x7c75005b.
//
// Solidity: function setAgentPaused(bool paused) returns()
func (_V4AgenticVault *V4AgenticVaultSession) SetAgentPaused(paused bool) (*types.Transaction, error) {
	return _V4AgenticVault.Contract.SetAgentPaused(&_V4AgenticVault.TransactOpts, paused)
}

// SetAgentPaused is a paid mutator transaction binding the contract method 0x7c75005b.
//
// Solidity: function setAgentPaused(bool paused) returns()
func (_V4AgenticVault *V4AgenticVaultTransactorSession) SetAgentPaused(paused bool) (*types.Transaction, error) {
	return _V4AgenticVault.Contract.SetAgentPaused(&_V4AgenticVault.TransactOpts, paused)
}

// SetAllowedTickRange is a paid mutator transaction binding the contract method 0x2b1d2125.
//
// Solidity: function setAllowedTickRange(int24 tickLower, int24 tickUpper) returns()
func (_V4AgenticVault *V4AgenticVaultTransactor) SetAllowedTickRange(opts *bind.TransactOpts, tickLower *big.Int, tickUpper *big.Int) (*types.Transaction, error) {
	return _V4AgenticVault.contract.Transact(opts, "setAllowedTickRange", tickLower, tickUpper)
}

// SetAllowedTickRange is a paid mutator transaction binding the contract method 0x2b1d2125.
//
// Solidity: function setAllowedTickRange(int24 tickLower, int24 tickUpper) returns()
func (_V4AgenticVault *V4AgenticVaultSession) SetAllowedTickRange(tickLower *big.Int, tickUpper *big.Int) (*types.Transaction, error) {
	return _V4AgenticVault.Contract.SetAllowedTickRange(&_V4AgenticVault.TransactOpts, tickLower, tickUpper)
}

// SetAllowedTickRange is a paid mutator transaction binding the contract method 0x2b1d2125.
//
// Solidity: function setAllowedTickRange(int24 tickLower, int24 tickUpper) returns()
func (_V4AgenticVault *V4AgenticVaultTransactorSession) SetAllowedTickRange(tickLower *big.Int, tickUpper *big.Int) (*types.Transaction, error) {
	return _V4AgenticVault.Contract.SetAllowedTickRange(&_V4AgenticVault.TransactOpts, tickLower, tickUpper)
}

// SetMaxPositionsK is a paid mutator transaction binding the contract method 0x39ca413e.
//
// Solidity: function setMaxPositionsK(uint256 k) returns()
func (_V4AgenticVault *V4AgenticVaultTransactor) SetMaxPositionsK(opts *bind.TransactOpts, k *big.Int) (*types.Transaction, error) {
	return _V4AgenticVault.contract.Transact(opts, "setMaxPositionsK", k)
}

// SetMaxPositionsK is a paid mutator transaction binding the contract method 0x39ca413e.
//
// Solidity: function setMaxPositionsK(uint256 k) returns()
func (_V4AgenticVault *V4AgenticVaultSession) SetMaxPositionsK(k *big.Int) (*types.Transaction, error) {
	return _V4AgenticVault.Contract.SetMaxPositionsK(&_V4AgenticVault.TransactOpts, k)
}

// SetMaxPositionsK is a paid mutator transaction binding the contract method 0x39ca413e.
//
// Solidity: function setMaxPositionsK(uint256 k) returns()
func (_V4AgenticVault *V4AgenticVaultTransactorSession) SetMaxPositionsK(k *big.Int) (*types.Transaction, error) {
	return _V4AgenticVault.Contract.SetMaxPositionsK(&_V4AgenticVault.TransactOpts, k)
}

// SetSwapAllowed is a paid mutator transaction binding the contract method 0xd3f829ed.
//
// Solidity: function setSwapAllowed(bool allowed) returns()
func (_V4AgenticVault *V4AgenticVaultTransactor) SetSwapAllowed(opts *bind.TransactOpts, allowed bool) (*types.Transaction, error) {
	return _V4AgenticVault.contract.Transact(opts, "setSwapAllowed", allowed)
}

// SetSwapAllowed is a paid mutator transaction binding the contract method 0xd3f829ed.
//
// Solidity: function setSwapAllowed(bool allowed) returns()
func (_V4AgenticVault *V4AgenticVaultSession) SetSwapAllowed(allowed bool) (*types.Transaction, error) {
	return _V4AgenticVault.Contract.SetSwapAllowed(&_V4AgenticVault.TransactOpts, allowed)
}

// SetSwapAllowed is a paid mutator transaction binding the contract method 0xd3f829ed.
//
// Solidity: function setSwapAllowed(bool allowed) returns()
func (_V4AgenticVault *V4AgenticVaultTransactorSession) SetSwapAllowed(allowed bool) (*types.Transaction, error) {
	return _V4AgenticVault.Contract.SetSwapAllowed(&_V4AgenticVault.TransactOpts, allowed)
}

// SwapExactInputSingle is a paid mutator transaction binding the contract method 0xa224ca66.
//
// Solidity: function swapExactInputSingle(bool zeroForOne, uint128 amountIn, uint128 minAmountOut, uint256 deadline) returns(uint256 amountOut)
func (_V4AgenticVault *V4AgenticVaultTransactor) SwapExactInputSingle(opts *bind.TransactOpts, zeroForOne bool, amountIn *big.Int, minAmountOut *big.Int, deadline *big.Int) (*types.Transaction, error) {
	return _V4AgenticVault.contract.Transact(opts, "swapExactInputSingle", zeroForOne, amountIn, minAmountOut, deadline)
}

// SwapExactInputSingle is a paid mutator transaction binding the contract method 0xa224ca66.
//
// Solidity: function swapExactInputSingle(bool zeroForOne, uint128 amountIn, uint128 minAmountOut, uint256 deadline) returns(uint256 amountOut)
func (_V4AgenticVault *V4AgenticVaultSession) SwapExactInputSingle(zeroForOne bool, amountIn *big.Int, minAmountOut *big.Int, deadline *big.Int) (*types.Transaction, error) {
	return _V4AgenticVault.Contract.SwapExactInputSingle(&_V4AgenticVault.TransactOpts, zeroForOne, amountIn, minAmountOut, deadline)
}

// SwapExactInputSingle is a paid mutator transaction binding the contract method 0xa224ca66.
//
// Solidity: function swapExactInputSingle(bool zeroForOne, uint128 amountIn, uint128 minAmountOut, uint256 deadline) returns(uint256 amountOut)
func (_V4AgenticVault *V4AgenticVaultTransactorSession) SwapExactInputSingle(zeroForOne bool, amountIn *big.Int, minAmountOut *big.Int, deadline *big.Int) (*types.Transaction, error) {
	return _V4AgenticVault.Contract.SwapExactInputSingle(&_V4AgenticVault.TransactOpts, zeroForOne, amountIn, minAmountOut, deadline)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_V4AgenticVault *V4AgenticVaultTransactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _V4AgenticVault.contract.Transact(opts, "transferOwnership", newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_V4AgenticVault *V4AgenticVaultSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _V4AgenticVault.Contract.TransferOwnership(&_V4AgenticVault.TransactOpts, newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_V4AgenticVault *V4AgenticVaultTransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _V4AgenticVault.Contract.TransferOwnership(&_V4AgenticVault.TransactOpts, newOwner)
}

// Withdraw is a paid mutator transaction binding the contract method 0x69328dec.
//
// Solidity: function withdraw(address currency, uint256 amount, address to) returns()
func (_V4AgenticVault *V4AgenticVaultTransactor) Withdraw(opts *bind.TransactOpts, currency common.Address, amount *big.Int, to common.Address) (*types.Transaction, error) {
	return _V4AgenticVault.contract.Transact(opts, "withdraw", currency, amount, to)
}

// Withdraw is a paid mutator transaction binding the contract method 0x69328dec.
//
// Solidity: function withdraw(address currency, uint256 amount, address to) returns()
func (_V4AgenticVault *V4AgenticVaultSession) Withdraw(currency common.Address, amount *big.Int, to common.Address) (*types.Transaction, error) {
	return _V4AgenticVault.Contract.Withdraw(&_V4AgenticVault.TransactOpts, currency, amount, to)
}

// Withdraw is a paid mutator transaction binding the contract method 0x69328dec.
//
// Solidity: function withdraw(address currency, uint256 amount, address to) returns()
func (_V4AgenticVault *V4AgenticVaultTransactorSession) Withdraw(currency common.Address, amount *big.Int, to common.Address) (*types.Transaction, error) {
	return _V4AgenticVault.Contract.Withdraw(&_V4AgenticVault.TransactOpts, currency, amount, to)
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_V4AgenticVault *V4AgenticVaultTransactor) Receive(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _V4AgenticVault.contract.RawTransact(opts, nil) // calldata is disallowed for receive function
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_V4AgenticVault *V4AgenticVaultSession) Receive() (*types.Transaction, error) {
	return _V4AgenticVault.Contract.Receive(&_V4AgenticVault.TransactOpts)
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_V4AgenticVault *V4AgenticVaultTransactorSession) Receive() (*types.Transaction, error) {
	return _V4AgenticVault.Contract.Receive(&_V4AgenticVault.TransactOpts)
}

// V4AgenticVaultAgentPausedIterator is returned from FilterAgentPaused and is used to iterate over the raw logs and unpacked data for AgentPaused events raised by the V4AgenticVault contract.
type V4AgenticVaultAgentPausedIterator struct {
	Event *V4AgenticVaultAgentPaused // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *V4AgenticVaultAgentPausedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(V4AgenticVaultAgentPaused)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(V4AgenticVaultAgentPaused)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *V4AgenticVaultAgentPausedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *V4AgenticVaultAgentPausedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// V4AgenticVaultAgentPaused represents a AgentPaused event raised by the V4AgenticVault contract.
type V4AgenticVaultAgentPaused struct {
	Paused bool
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterAgentPaused is a free log retrieval operation binding the contract event 0x5a563da584fcc9c1e5a7134b461ee58c183cf7ddb624d6cabbc7641a5de2b689.
//
// Solidity: event AgentPaused(bool paused)
func (_V4AgenticVault *V4AgenticVaultFilterer) FilterAgentPaused(opts *bind.FilterOpts) (*V4AgenticVaultAgentPausedIterator, error) {

	logs, sub, err := _V4AgenticVault.contract.FilterLogs(opts, "AgentPaused")
	if err != nil {
		return nil, err
	}
	return &V4AgenticVaultAgentPausedIterator{contract: _V4AgenticVault.contract, event: "AgentPaused", logs: logs, sub: sub}, nil
}

// WatchAgentPaused is a free log subscription operation binding the contract event 0x5a563da584fcc9c1e5a7134b461ee58c183cf7ddb624d6cabbc7641a5de2b689.
//
// Solidity: event AgentPaused(bool paused)
func (_V4AgenticVault *V4AgenticVaultFilterer) WatchAgentPaused(opts *bind.WatchOpts, sink chan<- *V4AgenticVaultAgentPaused) (event.Subscription, error) {

	logs, sub, err := _V4AgenticVault.contract.WatchLogs(opts, "AgentPaused")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(V4AgenticVaultAgentPaused)
				if err := _V4AgenticVault.contract.UnpackLog(event, "AgentPaused", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseAgentPaused is a log parse operation binding the contract event 0x5a563da584fcc9c1e5a7134b461ee58c183cf7ddb624d6cabbc7641a5de2b689.
//
// Solidity: event AgentPaused(bool paused)
func (_V4AgenticVault *V4AgenticVaultFilterer) ParseAgentPaused(log types.Log) (*V4AgenticVaultAgentPaused, error) {
	event := new(V4AgenticVaultAgentPaused)
	if err := _V4AgenticVault.contract.UnpackLog(event, "AgentPaused", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// V4AgenticVaultAgentUpdatedIterator is returned from FilterAgentUpdated and is used to iterate over the raw logs and unpacked data for AgentUpdated events raised by the V4AgenticVault contract.
type V4AgenticVaultAgentUpdatedIterator struct {
	Event *V4AgenticVaultAgentUpdated // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *V4AgenticVaultAgentUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(V4AgenticVaultAgentUpdated)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(V4AgenticVaultAgentUpdated)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *V4AgenticVaultAgentUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *V4AgenticVaultAgentUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// V4AgenticVaultAgentUpdated represents a AgentUpdated event raised by the V4AgenticVault contract.
type V4AgenticVaultAgentUpdated struct {
	NewAgent common.Address
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterAgentUpdated is a free log retrieval operation binding the contract event 0xe9f337c154e801e0f86b6bc993df9a2cd349bb210385592c7a52e38ea726334f.
//
// Solidity: event AgentUpdated(address indexed newAgent)
func (_V4AgenticVault *V4AgenticVaultFilterer) FilterAgentUpdated(opts *bind.FilterOpts, newAgent []common.Address) (*V4AgenticVaultAgentUpdatedIterator, error) {

	var newAgentRule []interface{}
	for _, newAgentItem := range newAgent {
		newAgentRule = append(newAgentRule, newAgentItem)
	}

	logs, sub, err := _V4AgenticVault.contract.FilterLogs(opts, "AgentUpdated", newAgentRule)
	if err != nil {
		return nil, err
	}
	return &V4AgenticVaultAgentUpdatedIterator{contract: _V4AgenticVault.contract, event: "AgentUpdated", logs: logs, sub: sub}, nil
}

// WatchAgentUpdated is a free log subscription operation binding the contract event 0xe9f337c154e801e0f86b6bc993df9a2cd349bb210385592c7a52e38ea726334f.
//
// Solidity: event AgentUpdated(address indexed newAgent)
func (_V4AgenticVault *V4AgenticVaultFilterer) WatchAgentUpdated(opts *bind.WatchOpts, sink chan<- *V4AgenticVaultAgentUpdated, newAgent []common.Address) (event.Subscription, error) {

	var newAgentRule []interface{}
	for _, newAgentItem := range newAgent {
		newAgentRule = append(newAgentRule, newAgentItem)
	}

	logs, sub, err := _V4AgenticVault.contract.WatchLogs(opts, "AgentUpdated", newAgentRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(V4AgenticVaultAgentUpdated)
				if err := _V4AgenticVault.contract.UnpackLog(event, "AgentUpdated", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseAgentUpdated is a log parse operation binding the contract event 0xe9f337c154e801e0f86b6bc993df9a2cd349bb210385592c7a52e38ea726334f.
//
// Solidity: event AgentUpdated(address indexed newAgent)
func (_V4AgenticVault *V4AgenticVaultFilterer) ParseAgentUpdated(log types.Log) (*V4AgenticVaultAgentUpdated, error) {
	event := new(V4AgenticVaultAgentUpdated)
	if err := _V4AgenticVault.contract.UnpackLog(event, "AgentUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// V4AgenticVaultAllowedTickRangeUpdatedIterator is returned from FilterAllowedTickRangeUpdated and is used to iterate over the raw logs and unpacked data for AllowedTickRangeUpdated events raised by the V4AgenticVault contract.
type V4AgenticVaultAllowedTickRangeUpdatedIterator struct {
	Event *V4AgenticVaultAllowedTickRangeUpdated // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *V4AgenticVaultAllowedTickRangeUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(V4AgenticVaultAllowedTickRangeUpdated)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(V4AgenticVaultAllowedTickRangeUpdated)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *V4AgenticVaultAllowedTickRangeUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *V4AgenticVaultAllowedTickRangeUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// V4AgenticVaultAllowedTickRangeUpdated represents a AllowedTickRangeUpdated event raised by the V4AgenticVault contract.
type V4AgenticVaultAllowedTickRangeUpdated struct {
	TickLower *big.Int
	TickUpper *big.Int
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterAllowedTickRangeUpdated is a free log retrieval operation binding the contract event 0x0dd7e9ab7f58c004f7e120690d3c6fe686e999de01c07ffbf5f5af7aa2af363d.
//
// Solidity: event AllowedTickRangeUpdated(int24 tickLower, int24 tickUpper)
func (_V4AgenticVault *V4AgenticVaultFilterer) FilterAllowedTickRangeUpdated(opts *bind.FilterOpts) (*V4AgenticVaultAllowedTickRangeUpdatedIterator, error) {

	logs, sub, err := _V4AgenticVault.contract.FilterLogs(opts, "AllowedTickRangeUpdated")
	if err != nil {
		return nil, err
	}
	return &V4AgenticVaultAllowedTickRangeUpdatedIterator{contract: _V4AgenticVault.contract, event: "AllowedTickRangeUpdated", logs: logs, sub: sub}, nil
}

// WatchAllowedTickRangeUpdated is a free log subscription operation binding the contract event 0x0dd7e9ab7f58c004f7e120690d3c6fe686e999de01c07ffbf5f5af7aa2af363d.
//
// Solidity: event AllowedTickRangeUpdated(int24 tickLower, int24 tickUpper)
func (_V4AgenticVault *V4AgenticVaultFilterer) WatchAllowedTickRangeUpdated(opts *bind.WatchOpts, sink chan<- *V4AgenticVaultAllowedTickRangeUpdated) (event.Subscription, error) {

	logs, sub, err := _V4AgenticVault.contract.WatchLogs(opts, "AllowedTickRangeUpdated")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(V4AgenticVaultAllowedTickRangeUpdated)
				if err := _V4AgenticVault.contract.UnpackLog(event, "AllowedTickRangeUpdated", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseAllowedTickRangeUpdated is a log parse operation binding the contract event 0x0dd7e9ab7f58c004f7e120690d3c6fe686e999de01c07ffbf5f5af7aa2af363d.
//
// Solidity: event AllowedTickRangeUpdated(int24 tickLower, int24 tickUpper)
func (_V4AgenticVault *V4AgenticVaultFilterer) ParseAllowedTickRangeUpdated(log types.Log) (*V4AgenticVaultAllowedTickRangeUpdated, error) {
	event := new(V4AgenticVaultAllowedTickRangeUpdated)
	if err := _V4AgenticVault.contract.UnpackLog(event, "AllowedTickRangeUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// V4AgenticVaultMaxPositionsKUpdatedIterator is returned from FilterMaxPositionsKUpdated and is used to iterate over the raw logs and unpacked data for MaxPositionsKUpdated events raised by the V4AgenticVault contract.
type V4AgenticVaultMaxPositionsKUpdatedIterator struct {
	Event *V4AgenticVaultMaxPositionsKUpdated // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *V4AgenticVaultMaxPositionsKUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(V4AgenticVaultMaxPositionsKUpdated)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(V4AgenticVaultMaxPositionsKUpdated)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *V4AgenticVaultMaxPositionsKUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *V4AgenticVaultMaxPositionsKUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// V4AgenticVaultMaxPositionsKUpdated represents a MaxPositionsKUpdated event raised by the V4AgenticVault contract.
type V4AgenticVaultMaxPositionsKUpdated struct {
	K   *big.Int
	Raw types.Log // Blockchain specific contextual infos
}

// FilterMaxPositionsKUpdated is a free log retrieval operation binding the contract event 0x73a6aedd92b89778061d6e6a15db15d453064446aebd75aa07b46819d4a127b2.
//
// Solidity: event MaxPositionsKUpdated(uint256 k)
func (_V4AgenticVault *V4AgenticVaultFilterer) FilterMaxPositionsKUpdated(opts *bind.FilterOpts) (*V4AgenticVaultMaxPositionsKUpdatedIterator, error) {

	logs, sub, err := _V4AgenticVault.contract.FilterLogs(opts, "MaxPositionsKUpdated")
	if err != nil {
		return nil, err
	}
	return &V4AgenticVaultMaxPositionsKUpdatedIterator{contract: _V4AgenticVault.contract, event: "MaxPositionsKUpdated", logs: logs, sub: sub}, nil
}

// WatchMaxPositionsKUpdated is a free log subscription operation binding the contract event 0x73a6aedd92b89778061d6e6a15db15d453064446aebd75aa07b46819d4a127b2.
//
// Solidity: event MaxPositionsKUpdated(uint256 k)
func (_V4AgenticVault *V4AgenticVaultFilterer) WatchMaxPositionsKUpdated(opts *bind.WatchOpts, sink chan<- *V4AgenticVaultMaxPositionsKUpdated) (event.Subscription, error) {

	logs, sub, err := _V4AgenticVault.contract.WatchLogs(opts, "MaxPositionsKUpdated")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(V4AgenticVaultMaxPositionsKUpdated)
				if err := _V4AgenticVault.contract.UnpackLog(event, "MaxPositionsKUpdated", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseMaxPositionsKUpdated is a log parse operation binding the contract event 0x73a6aedd92b89778061d6e6a15db15d453064446aebd75aa07b46819d4a127b2.
//
// Solidity: event MaxPositionsKUpdated(uint256 k)
func (_V4AgenticVault *V4AgenticVaultFilterer) ParseMaxPositionsKUpdated(log types.Log) (*V4AgenticVaultMaxPositionsKUpdated, error) {
	event := new(V4AgenticVaultMaxPositionsKUpdated)
	if err := _V4AgenticVault.contract.UnpackLog(event, "MaxPositionsKUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// V4AgenticVaultOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the V4AgenticVault contract.
type V4AgenticVaultOwnershipTransferredIterator struct {
	Event *V4AgenticVaultOwnershipTransferred // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *V4AgenticVaultOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(V4AgenticVaultOwnershipTransferred)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(V4AgenticVaultOwnershipTransferred)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *V4AgenticVaultOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *V4AgenticVaultOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// V4AgenticVaultOwnershipTransferred represents a OwnershipTransferred event raised by the V4AgenticVault contract.
type V4AgenticVaultOwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_V4AgenticVault *V4AgenticVaultFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*V4AgenticVaultOwnershipTransferredIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _V4AgenticVault.contract.FilterLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &V4AgenticVaultOwnershipTransferredIterator{contract: _V4AgenticVault.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_V4AgenticVault *V4AgenticVaultFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *V4AgenticVaultOwnershipTransferred, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _V4AgenticVault.contract.WatchLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(V4AgenticVaultOwnershipTransferred)
				if err := _V4AgenticVault.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseOwnershipTransferred is a log parse operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_V4AgenticVault *V4AgenticVaultFilterer) ParseOwnershipTransferred(log types.Log) (*V4AgenticVaultOwnershipTransferred, error) {
	event := new(V4AgenticVaultOwnershipTransferred)
	if err := _V4AgenticVault.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// V4AgenticVaultPositionAddedIterator is returned from FilterPositionAdded and is used to iterate over the raw logs and unpacked data for PositionAdded events raised by the V4AgenticVault contract.
type V4AgenticVaultPositionAddedIterator struct {
	Event *V4AgenticVaultPositionAdded // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *V4AgenticVaultPositionAddedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(V4AgenticVaultPositionAdded)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(V4AgenticVaultPositionAdded)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *V4AgenticVaultPositionAddedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *V4AgenticVaultPositionAddedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// V4AgenticVaultPositionAdded represents a PositionAdded event raised by the V4AgenticVault contract.
type V4AgenticVaultPositionAdded struct {
	TokenId   *big.Int
	TickLower *big.Int
	TickUpper *big.Int
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterPositionAdded is a free log retrieval operation binding the contract event 0xcb7876fb20f38b331b8a25936883f72a96d078e25f941e7cfd88eed8ef4ef4b9.
//
// Solidity: event PositionAdded(uint256 indexed tokenId, int24 tickLower, int24 tickUpper)
func (_V4AgenticVault *V4AgenticVaultFilterer) FilterPositionAdded(opts *bind.FilterOpts, tokenId []*big.Int) (*V4AgenticVaultPositionAddedIterator, error) {

	var tokenIdRule []interface{}
	for _, tokenIdItem := range tokenId {
		tokenIdRule = append(tokenIdRule, tokenIdItem)
	}

	logs, sub, err := _V4AgenticVault.contract.FilterLogs(opts, "PositionAdded", tokenIdRule)
	if err != nil {
		return nil, err
	}
	return &V4AgenticVaultPositionAddedIterator{contract: _V4AgenticVault.contract, event: "PositionAdded", logs: logs, sub: sub}, nil
}

// WatchPositionAdded is a free log subscription operation binding the contract event 0xcb7876fb20f38b331b8a25936883f72a96d078e25f941e7cfd88eed8ef4ef4b9.
//
// Solidity: event PositionAdded(uint256 indexed tokenId, int24 tickLower, int24 tickUpper)
func (_V4AgenticVault *V4AgenticVaultFilterer) WatchPositionAdded(opts *bind.WatchOpts, sink chan<- *V4AgenticVaultPositionAdded, tokenId []*big.Int) (event.Subscription, error) {

	var tokenIdRule []interface{}
	for _, tokenIdItem := range tokenId {
		tokenIdRule = append(tokenIdRule, tokenIdItem)
	}

	logs, sub, err := _V4AgenticVault.contract.WatchLogs(opts, "PositionAdded", tokenIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(V4AgenticVaultPositionAdded)
				if err := _V4AgenticVault.contract.UnpackLog(event, "PositionAdded", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParsePositionAdded is a log parse operation binding the contract event 0xcb7876fb20f38b331b8a25936883f72a96d078e25f941e7cfd88eed8ef4ef4b9.
//
// Solidity: event PositionAdded(uint256 indexed tokenId, int24 tickLower, int24 tickUpper)
func (_V4AgenticVault *V4AgenticVaultFilterer) ParsePositionAdded(log types.Log) (*V4AgenticVaultPositionAdded, error) {
	event := new(V4AgenticVaultPositionAdded)
	if err := _V4AgenticVault.contract.UnpackLog(event, "PositionAdded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// V4AgenticVaultPositionRemovedIterator is returned from FilterPositionRemoved and is used to iterate over the raw logs and unpacked data for PositionRemoved events raised by the V4AgenticVault contract.
type V4AgenticVaultPositionRemovedIterator struct {
	Event *V4AgenticVaultPositionRemoved // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *V4AgenticVaultPositionRemovedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(V4AgenticVaultPositionRemoved)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(V4AgenticVaultPositionRemoved)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *V4AgenticVaultPositionRemovedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *V4AgenticVaultPositionRemovedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// V4AgenticVaultPositionRemoved represents a PositionRemoved event raised by the V4AgenticVault contract.
type V4AgenticVaultPositionRemoved struct {
	TokenId *big.Int
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterPositionRemoved is a free log retrieval operation binding the contract event 0x8dab17171ceb83acca6e0d86e9c8b6d3f2eec3a1ec6a5a2f9a3f1b7e17ad1035.
//
// Solidity: event PositionRemoved(uint256 indexed tokenId)
func (_V4AgenticVault *V4AgenticVaultFilterer) FilterPositionRemoved(opts *bind.FilterOpts, tokenId []*big.Int) (*V4AgenticVaultPositionRemovedIterator, error) {

	var tokenIdRule []interface{}
	for _, tokenIdItem := range tokenId {
		tokenIdRule = append(tokenIdRule, tokenIdItem)
	}

	logs, sub, err := _V4AgenticVault.contract.FilterLogs(opts, "PositionRemoved", tokenIdRule)
	if err != nil {
		return nil, err
	}
	return &V4AgenticVaultPositionRemovedIterator{contract: _V4AgenticVault.contract, event: "PositionRemoved", logs: logs, sub: sub}, nil
}

// WatchPositionRemoved is a free log subscription operation binding the contract event 0x8dab17171ceb83acca6e0d86e9c8b6d3f2eec3a1ec6a5a2f9a3f1b7e17ad1035.
//
// Solidity: event PositionRemoved(uint256 indexed tokenId)
func (_V4AgenticVault *V4AgenticVaultFilterer) WatchPositionRemoved(opts *bind.WatchOpts, sink chan<- *V4AgenticVaultPositionRemoved, tokenId []*big.Int) (event.Subscription, error) {

	var tokenIdRule []interface{}
	for _, tokenIdItem := range tokenId {
		tokenIdRule = append(tokenIdRule, tokenIdItem)
	}

	logs, sub, err := _V4AgenticVault.contract.WatchLogs(opts, "PositionRemoved", tokenIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(V4AgenticVaultPositionRemoved)
				if err := _V4AgenticVault.contract.UnpackLog(event, "PositionRemoved", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParsePositionRemoved is a log parse operation binding the contract event 0x8dab17171ceb83acca6e0d86e9c8b6d3f2eec3a1ec6a5a2f9a3f1b7e17ad1035.
//
// Solidity: event PositionRemoved(uint256 indexed tokenId)
func (_V4AgenticVault *V4AgenticVaultFilterer) ParsePositionRemoved(log types.Log) (*V4AgenticVaultPositionRemoved, error) {
	event := new(V4AgenticVaultPositionRemoved)
	if err := _V4AgenticVault.contract.UnpackLog(event, "PositionRemoved", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// V4AgenticVaultSwapAllowedIterator is returned from FilterSwapAllowed and is used to iterate over the raw logs and unpacked data for SwapAllowed events raised by the V4AgenticVault contract.
type V4AgenticVaultSwapAllowedIterator struct {
	Event *V4AgenticVaultSwapAllowed // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *V4AgenticVaultSwapAllowedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(V4AgenticVaultSwapAllowed)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(V4AgenticVaultSwapAllowed)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *V4AgenticVaultSwapAllowedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *V4AgenticVaultSwapAllowedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// V4AgenticVaultSwapAllowed represents a SwapAllowed event raised by the V4AgenticVault contract.
type V4AgenticVaultSwapAllowed struct {
	Allowed bool
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterSwapAllowed is a free log retrieval operation binding the contract event 0xe4d01fad3ca9a51a8f50a718e1168e0a0eab75526b6517ff484b31d6468f1429.
//
// Solidity: event SwapAllowed(bool allowed)
func (_V4AgenticVault *V4AgenticVaultFilterer) FilterSwapAllowed(opts *bind.FilterOpts) (*V4AgenticVaultSwapAllowedIterator, error) {

	logs, sub, err := _V4AgenticVault.contract.FilterLogs(opts, "SwapAllowed")
	if err != nil {
		return nil, err
	}
	return &V4AgenticVaultSwapAllowedIterator{contract: _V4AgenticVault.contract, event: "SwapAllowed", logs: logs, sub: sub}, nil
}

// WatchSwapAllowed is a free log subscription operation binding the contract event 0xe4d01fad3ca9a51a8f50a718e1168e0a0eab75526b6517ff484b31d6468f1429.
//
// Solidity: event SwapAllowed(bool allowed)
func (_V4AgenticVault *V4AgenticVaultFilterer) WatchSwapAllowed(opts *bind.WatchOpts, sink chan<- *V4AgenticVaultSwapAllowed) (event.Subscription, error) {

	logs, sub, err := _V4AgenticVault.contract.WatchLogs(opts, "SwapAllowed")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(V4AgenticVaultSwapAllowed)
				if err := _V4AgenticVault.contract.UnpackLog(event, "SwapAllowed", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseSwapAllowed is a log parse operation binding the contract event 0xe4d01fad3ca9a51a8f50a718e1168e0a0eab75526b6517ff484b31d6468f1429.
//
// Solidity: event SwapAllowed(bool allowed)
func (_V4AgenticVault *V4AgenticVaultFilterer) ParseSwapAllowed(log types.Log) (*V4AgenticVaultSwapAllowed, error) {
	event := new(V4AgenticVaultSwapAllowed)
	if err := _V4AgenticVault.contract.UnpackLog(event, "SwapAllowed", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
