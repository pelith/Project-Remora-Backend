package liquidity_test

import (
	"errors"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/shopspring/decimal"

	"remora/internal/liquidity"
)

func TestBuildRebalanceAllocations(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		ranges     []liquidity.TickRangeWeight
		amount     decimal.Decimal
		slippage   decimal.Decimal
		want       []liquidity.RebalanceAllocation
		wantErr    bool
		wantErrIs  error
		wantAmount decimal.Decimal
	}{
		{
			name: "success - weighted allocation with slippage",
			ranges: []liquidity.TickRangeWeight{
				{TickLower: -100, TickUpper: 0, Weight: mustDecimal(t, "1")},
				{TickLower: 0, TickUpper: 100, Weight: mustDecimal(t, "3")},
			},
			amount:   mustDecimal(t, "100"),
			slippage: mustDecimal(t, "0.01"),
			want: []liquidity.RebalanceAllocation{
				{
					TickLower: -100,
					TickUpper: 0,
					Weight:    mustDecimal(t, "1"),
					Amount:    mustDecimal(t, "25"),
					AmountMin: mustDecimal(t, "24.75"),
				},
				{
					TickLower: 0,
					TickUpper: 100,
					Weight:    mustDecimal(t, "3"),
					Amount:    mustDecimal(t, "75"),
					AmountMin: mustDecimal(t, "74.25"),
				},
			},
			wantAmount: mustDecimal(t, "100"),
		},
		{
			name: "success - zero weight keeps entry with zero amounts",
			ranges: []liquidity.TickRangeWeight{
				{TickLower: -200, TickUpper: -100, Weight: mustDecimal(t, "0")},
				{TickLower: -100, TickUpper: 100, Weight: mustDecimal(t, "2")},
			},
			amount:   mustDecimal(t, "10"),
			slippage: mustDecimal(t, "0.05"),
			want: []liquidity.RebalanceAllocation{
				{
					TickLower: -200,
					TickUpper: -100,
					Weight:    mustDecimal(t, "0"),
					Amount:    mustDecimal(t, "0"),
					AmountMin: mustDecimal(t, "0"),
				},
				{
					TickLower: -100,
					TickUpper: 100,
					Weight:    mustDecimal(t, "2"),
					Amount:    mustDecimal(t, "10"),
					AmountMin: mustDecimal(t, "9.5"),
				},
			},
			wantAmount: mustDecimal(t, "10"),
		},
		{
			name:       "error - no ranges",
			ranges:     nil,
			amount:     mustDecimal(t, "1"),
			slippage:   mustDecimal(t, "0"),
			wantErr:    true,
			wantErrIs:  liquidity.ErrNoTickRanges,
			wantAmount: mustDecimal(t, "0"),
		},
		{
			name: "error - invalid tick range",
			ranges: []liquidity.TickRangeWeight{
				{TickLower: 100, TickUpper: 100, Weight: mustDecimal(t, "1")},
			},
			amount:     mustDecimal(t, "1"),
			slippage:   mustDecimal(t, "0"),
			wantErr:    true,
			wantErrIs:  liquidity.ErrInvalidTickRange,
			wantAmount: mustDecimal(t, "0"),
		},
		{
			name: "error - negative weight",
			ranges: []liquidity.TickRangeWeight{
				{TickLower: 0, TickUpper: 10, Weight: mustDecimal(t, "-1")},
			},
			amount:     mustDecimal(t, "1"),
			slippage:   mustDecimal(t, "0"),
			wantErr:    true,
			wantErrIs:  liquidity.ErrInvalidWeight,
			wantAmount: mustDecimal(t, "0"),
		},
		{
			name: "error - negative total amount",
			ranges: []liquidity.TickRangeWeight{
				{TickLower: 0, TickUpper: 10, Weight: mustDecimal(t, "1")},
			},
			amount:     mustDecimal(t, "-1"),
			slippage:   mustDecimal(t, "0"),
			wantErr:    true,
			wantErrIs:  liquidity.ErrInvalidTotalAmount,
			wantAmount: mustDecimal(t, "0"),
		},
		{
			name: "error - slippage > 1",
			ranges: []liquidity.TickRangeWeight{
				{TickLower: 0, TickUpper: 10, Weight: mustDecimal(t, "1")},
			},
			amount:     mustDecimal(t, "1"),
			slippage:   mustDecimal(t, "1.01"),
			wantErr:    true,
			wantErrIs:  liquidity.ErrInvalidSlippage,
			wantAmount: mustDecimal(t, "0"),
		},
		{
			name: "error - total weight zero",
			ranges: []liquidity.TickRangeWeight{
				{TickLower: 0, TickUpper: 10, Weight: mustDecimal(t, "0")},
				{TickLower: 10, TickUpper: 20, Weight: mustDecimal(t, "0")},
			},
			amount:     mustDecimal(t, "1"),
			slippage:   mustDecimal(t, "0"),
			wantErr:    true,
			wantErrIs:  liquidity.ErrZeroTotalWeight,
			wantAmount: mustDecimal(t, "0"),
		},
	}

	decimalComparer := cmp.Comparer(func(a, b decimal.Decimal) bool {
		return a.Equal(b)
	})

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := liquidity.BuildRebalanceAllocations(tt.ranges, tt.amount, tt.slippage)
			if err != nil {
				if !tt.wantErr {
					t.Fatalf("BuildRebalanceAllocations() failed: %v", err)
				}

				if tt.wantErrIs != nil && !errors.Is(err, tt.wantErrIs) {
					t.Fatalf("BuildRebalanceAllocations() error = %v, want %v", err, tt.wantErrIs)
				}

				return
			}

			if tt.wantErr {
				t.Fatalf("BuildRebalanceAllocations() expected error")
			}

			if diff := cmp.Diff(tt.want, got, decimalComparer); diff != "" {
				t.Fatalf("BuildRebalanceAllocations() mismatch (-want +got):\n%s", diff)
			}

			sum := decimal.NewFromInt(0)
			for _, allocation := range got {
				sum = sum.Add(allocation.Amount)
			}

			if !sum.Equal(tt.wantAmount) {
				t.Fatalf("BuildRebalanceAllocations() amount sum = %s, want %s", sum, tt.wantAmount)
			}
		})
	}
}

func mustDecimal(t *testing.T, value string) decimal.Decimal {
	t.Helper()

	dec, err := decimal.NewFromString(value)
	if err != nil {
		t.Fatalf("NewDecimalFromString(%q) failed: %v", value, err)
	}

	return dec
}
