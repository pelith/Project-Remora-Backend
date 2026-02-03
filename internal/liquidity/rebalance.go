package liquidity

import "github.com/shopspring/decimal"

func BuildRebalanceAllocations(
	ranges []TickRangeWeight,
	totalAmount decimal.Decimal,
	slippageTolerance decimal.Decimal,
) ([]RebalanceAllocation, error) {
	if len(ranges) == 0 {
		return nil, ErrNoTickRanges
	}

	if totalAmount.IsNegative() {
		return nil, ErrInvalidTotalAmount
	}

	if slippageTolerance.IsNegative() || slippageTolerance.GreaterThan(decimal.NewFromInt(1)) {
		return nil, ErrInvalidSlippage
	}

	sumWeight := decimal.NewFromInt(0)
	for _, r := range ranges {
		if r.TickLower >= r.TickUpper {
			return nil, ErrInvalidTickRange
		}

		if r.Weight.LessThan(decimal.NewFromInt(0)) {
			return nil, ErrInvalidWeight
		}

		sumWeight = sumWeight.Add(r.Weight)
	}

	if sumWeight.IsZero() {
		return nil, ErrZeroTotalWeight
	}

	allocations := make([]RebalanceAllocation, 0, len(ranges))
	remaining := totalAmount
	minFactor := decimal.NewFromInt(1).Sub(slippageTolerance)

	for i, r := range ranges {
		var amount decimal.Decimal
		if i == len(ranges)-1 {
			amount = remaining
		} else {
			amount = totalAmount.Mul(r.Weight).DivRound(sumWeight, 18)
			remaining = remaining.Sub(amount)
		}

		allocations = append(allocations, RebalanceAllocation{
			TickLower: r.TickLower,
			TickUpper: r.TickUpper,
			Weight:    r.Weight,
			Amount:    amount,
			AmountMin: amount.Mul(minFactor).Round(18),
		})
	}

	return allocations, nil
}
