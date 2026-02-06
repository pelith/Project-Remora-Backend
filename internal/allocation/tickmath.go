package allocation

import (
	"errors"
	"fmt"
	"math/big"
)

// Constants from Uniswap V3 TickMath.
const (
	MinTick = -887272
	MaxTick = 887272
)

var (
	// MinSqrtRatio = 4295128739.
	MinSqrtRatio = big.NewInt(4295128739)
	// MaxSqrtRatio = 1461446703485210103287273052203988822378723970342.
	MaxSqrtRatio = func() *big.Int {
		n, _ := new(big.Int).SetString("1461446703485210103287273052203988822378723970342", 10)
		return n
	}()
)

// GetSqrtRatioAtTick calculates sqrt(1.0001^tick) * 2^96
// Ported from Uniswap V3 TickMath.sol.
func GetSqrtRatioAtTick(tick int) (*big.Int, error) {
	if tick < MinTick || tick > MaxTick {
		return nil, fmt.Errorf("tick %d out of range", tick)
	}

	absTick := tick
	if tick < 0 {
		absTick = -tick
	}

	// ratio starts as 1 << 128
	ratio := new(big.Int).SetBit(new(big.Int), 128, 1)

	if absTick&0x1 != 0 {
		// ratio = ratio * 0xfffcb933bd6fad37aa2d162d1a594001 >> 128
		ratio.Mul(ratio, hexToBig("0xfffcb933bd6fad37aa2d162d1a594001"))
		ratio.Rsh(ratio, 128)
	}

	if absTick&0x2 != 0 {
		ratio.Mul(ratio, hexToBig("0xfff97272373d413259a46990580e213a"))
		ratio.Rsh(ratio, 128)
	}

	if absTick&0x4 != 0 {
		ratio.Mul(ratio, hexToBig("0xfff2e50f5f656932ef12357cf3c7fdcc"))
		ratio.Rsh(ratio, 128)
	}

	if absTick&0x8 != 0 {
		ratio.Mul(ratio, hexToBig("0xffe5caca7e10e4e61c3624eaa0941cd0"))
		ratio.Rsh(ratio, 128)
	}

	if absTick&0x10 != 0 {
		ratio.Mul(ratio, hexToBig("0xffcb9843d60f6159c9db58835c926644"))
		ratio.Rsh(ratio, 128)
	}

	if absTick&0x20 != 0 {
		ratio.Mul(ratio, hexToBig("0xff973b41fa98c081472e6896dfb254c0"))
		ratio.Rsh(ratio, 128)
	}

	if absTick&0x40 != 0 {
		ratio.Mul(ratio, hexToBig("0xff2ea16466c96a3843ec78b326b52861"))
		ratio.Rsh(ratio, 128)
	}

	if absTick&0x80 != 0 {
		ratio.Mul(ratio, hexToBig("0xfe5dee046a99a2a811c461f1969c3053"))
		ratio.Rsh(ratio, 128)
	}

	if absTick&0x100 != 0 {
		ratio.Mul(ratio, hexToBig("0xfcbe86c7900a88aedcffc5d932334409"))
		ratio.Rsh(ratio, 128)
	}

	if absTick&0x200 != 0 {
		ratio.Mul(ratio, hexToBig("0xf987a7253ac413176f2b074cf7815e54"))
		ratio.Rsh(ratio, 128)
	}

	if absTick&0x400 != 0 {
		ratio.Mul(ratio, hexToBig("0xf3392b0822b70005940c7a398e4b70f3"))
		ratio.Rsh(ratio, 128)
	}

	if absTick&0x800 != 0 {
		ratio.Mul(ratio, hexToBig("0xe7159475a2c29b7443b29c7fa6e889d9"))
		ratio.Rsh(ratio, 128)
	}

	if absTick&0x1000 != 0 {
		ratio.Mul(ratio, hexToBig("0xd097f3bdfd2022b8845ad8f792aa5825"))
		ratio.Rsh(ratio, 128)
	}

	if absTick&0x2000 != 0 {
		ratio.Mul(ratio, hexToBig("0xa9f746462d870fdf8a65dc1f90e061e5"))
		ratio.Rsh(ratio, 128)
	}

	if absTick&0x4000 != 0 {
		ratio.Mul(ratio, hexToBig("0x70d869a156d2a1b890bb3df62baf32f7"))
		ratio.Rsh(ratio, 128)
	}

	if absTick&0x8000 != 0 {
		ratio.Mul(ratio, hexToBig("0x31be135f97d08fd981231505542fcfa6"))
		ratio.Rsh(ratio, 128)
	}

	if absTick&0x10000 != 0 {
		ratio.Mul(ratio, hexToBig("0x9aa508b5b7a84e1c677de54f3e99bc9"))
		ratio.Rsh(ratio, 128)
	}

	if absTick&0x20000 != 0 {
		ratio.Mul(ratio, hexToBig("0x5d6af8dedb81196699c329225ee604"))
		ratio.Rsh(ratio, 128)
	}

	if absTick&0x40000 != 0 {
		ratio.Mul(ratio, hexToBig("0x2216e584f5fa1ea926041bedfe98"))
		ratio.Rsh(ratio, 128)
	}

	if absTick&0x80000 != 0 {
		ratio.Mul(ratio, hexToBig("0x48a170391f7dc42444e8fa2"))
		ratio.Rsh(ratio, 128)
	}

	if tick > 0 {
		// ratio = MaxUint256 / ratio
		maxUint256 := new(big.Int).Sub(new(big.Int).Exp(big.NewInt(2), big.NewInt(256), nil), big.NewInt(1))
		ratio.Div(maxUint256, ratio)
	}

	// Round up if necessary (TickMath doesn't explicitly round up here, it truncates, but we must handle shifts carefully)
	// Currently ratio is Q128 (approx). We want Q96.
	// Shift right by 32
	ratio.Rsh(ratio, 32)

	return ratio, nil
}

// GetTickAtSqrtRatio calculates the greatest tick such that
// GetSqrtRatioAtTick(tick) <= sqrtPriceX96.
// Ported from Uniswap V3 TickMath.sol.
func GetTickAtSqrtRatio(sqrtPriceX96 *big.Int) (int, error) {
	if sqrtPriceX96.Cmp(MinSqrtRatio) < 0 || sqrtPriceX96.Cmp(MaxSqrtRatio) >= 0 {
		return 0, errors.New("sqrtPriceX96 out of range")
	}

	// ratio = sqrtPriceX96 << 32
	ratio := new(big.Int).Lsh(sqrtPriceX96, 32)

	// msb = floor(log2(ratio))
	msb := ratio.BitLen() - 1

	// log_2 = (msb - 128) << 64
	log2 := new(big.Int).Lsh(big.NewInt(int64(msb-128)), 64)

	// r = ratio >> (msb - 127) if msb >= 128 else ratio << (127 - msb)
	r := new(big.Int)
	if msb >= 128 {
		r.Rsh(ratio, uint(msb-127))
	} else {
		r.Lsh(ratio, uint(127-msb))
	}

	// for i in 0..13
	for i := range 14 {
		// r = (r * r) >> 127
		r.Mul(r, r)
		r.Rsh(r, 127)

		// f = r >> 128 (either 0 or 1)
		f := 0
		if r.BitLen() > 128 {
			f = 1
		}

		if f == 1 {
			// log2 |= 1 << (63 - i)
			log2.SetBit(log2, 63-i, 1)
			r.Rsh(r, 1)
		}
	}

	// log_sqrt10001 = log_2 * 255738958999603826347141
	logSqrt10001 := new(big.Int).Mul(log2, bigFromDec("255738958999603826347141"))

	// tickLow = (log_sqrt10001 - 3402992956809132418596140100660247210) >> 128
	// tickHigh = (log_sqrt10001 + 291339464771989622907027621153398088495) >> 128
	tickLow := new(big.Int).Sub(logSqrt10001, bigFromDec("3402992956809132418596140100660247210"))
	tickLow.Rsh(tickLow, 128)

	tickHigh := new(big.Int).Add(logSqrt10001, bigFromDec("291339464771989622907027621153398088495"))
	tickHigh.Rsh(tickHigh, 128)

	tLow := int(tickLow.Int64())

	tHigh := int(tickHigh.Int64())
	if tLow == tHigh {
		return tLow, nil
	}

	sqrtAtHigh, err := GetSqrtRatioAtTick(tHigh)
	if err != nil {
		return tLow, err
	}

	if sqrtAtHigh.Cmp(sqrtPriceX96) <= 0 {
		return tHigh, nil
	}

	return tLow, nil
}

func hexToBig(hex string) *big.Int {
	n := new(big.Int)
	n.SetString(hex[2:], 16)

	return n
}

func bigFromDec(dec string) *big.Int {
	n, _ := new(big.Int).SetString(dec, 10)
	return n
}
