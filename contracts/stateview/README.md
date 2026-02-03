# StateView Contract

## Contract Information

- **Name**: StateView
- **Purpose**: View-only contract for reading Uniswap v4 pool state
- **Repository**: https://github.com/Uniswap/v4-periphery
- **Documentation**: https://docs.uniswap.org/contracts/v4/reference/periphery/lens/StateView

## Deployed Addresses

### Ethereum Mainnet
- **Address**: `0x7ffe42c4a5deea5b0fec41c94c136cf115597227`
- **Explorer**: https://etherscan.io/address/0x7ffe42c4a5deea5b0fec41c94c136cf115597227

### Unichain (Chain 130)
- **Address**: `0x86e8631a016f9068c3f085faf484ee3f5fdee8f2`

## ABI Generation

This directory contains the StateView contract ABI for generating Go bindings.

### Generate Go Bindings

```bash
make abigen
```

This will generate `stateview.go` in `internal/liquidity/repository/contracts/`

### Manual Generation

```bash
abigen \
  --abi contracts/stateview/StateView.json \
  --pkg contracts \
  --type StateView \
  --out internal/liquidity/repository/contracts/stateview.go
```

## Core Functions

### getSlot0(PoolKey) 
Returns pool's current state:
- `sqrtPriceX96`: Square root price in Q96 format
- `tick`: Current tick
- `protocolFee`: Protocol fee
- `lpFee`: LP swap fee

### getTickBitmap(PoolKey, int16)
Returns tick bitmap for a word position (256 ticks per word)

### getTickInfo(PoolKey, int24)
Returns tick liquidity information:
- `liquidityGross`: Total position liquidity
- `liquidityNet`: Net liquidity change when crossing tick
- `feeGrowthOutside0X128`: Fee growth outside (token0)
- `feeGrowthOutside1X128`: Fee growth outside (token1)

## Notes

- The ABI includes only the essential functions needed for liquidity distribution reading
- Full contract source: https://github.com/Uniswap/v4-periphery/blob/main/src/lens/StateView.sol
