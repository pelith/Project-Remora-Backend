package liquidity

import "github.com/shopspring/decimal"

type TickRangeWeight struct {
	TickLower int
	TickUpper int
	Weight    decimal.Decimal
}

type RebalanceAllocation struct {
	TickLower int
	TickUpper int
	Weight    decimal.Decimal
	Amount    decimal.Decimal
	AmountMin decimal.Decimal
}
