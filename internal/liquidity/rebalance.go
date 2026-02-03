package liquidity

import "github.com/shopspring/decimal"

func BuildRebalanceAllocations(
	ranges []TickRangeWeight,
	totalAmount decimal.Decimal,
) ([]RebalanceAllocation, error) {
	if len(ranges) == 0 {
		return nil, ErrNoTickRanges
	}

	if totalAmount.IsNegative() {
		return nil, ErrInvalidTotalAmount
	}

	sumWeight := decimal.Zero
	for _, r := range ranges {
		if r.TickLower >= r.TickUpper {
			return nil, ErrInvalidTickRange
		}

		if r.Weight.IsNegative() {
			return nil, ErrInvalidWeight
		}

		sumWeight = sumWeight.Add(r.Weight)
	}

	if sumWeight.IsZero() {
		return nil, ErrZeroTotalWeight
	}

	allocations := make([]RebalanceAllocation, len(ranges))
	remaining := totalAmount

	for i, r := range ranges {
		var amount decimal.Decimal
		if i == len(ranges)-1 {
			amount = remaining
		} else {
			amount = totalAmount.Mul(r.Weight).Div(sumWeight)
			remaining = remaining.Sub(amount)
		}

		allocations[i] = RebalanceAllocation{
			TickLower: r.TickLower,
			TickUpper: r.TickUpper,
			Weight:    r.Weight,
			Amount:    amount,
		}
	}

	return allocations, nil
}
