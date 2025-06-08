package monitor

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCompare(t *testing.T) {
	const (
		BTC = iota + 3
		ETH
	)

	tests := []struct {
		name      string
		prices    []PriceData
		threshold float64
		expected  []PriceDifference
	}{
		{
			name: "no differences",
			prices: []PriceData{
				{Pair: Pair{Base: BTC, Quote: USD}, Service: "Provider1", Price: 10},
				{Pair: Pair{Base: BTC, Quote: USD}, Service: "Provider2", Price: 10},
				{Pair: Pair{Base: ETH, Quote: USD}, Service: "Provider1", Price: 30.00},
				{Pair: Pair{Base: ETH, Quote: USD}, Service: "Provider2", Price: 30.04},
			},
			threshold: 0.05,
			expected:  nil,
		},
		{
			name: "one difference above threshold",
			prices: []PriceData{
				{Pair: Pair{Base: BTC, Quote: USD}, Service: "Provider1", Price: 50000},
				{Pair: Pair{Base: BTC, Quote: USD}, Service: "Provider2", Price: 50000},
				{Pair: Pair{Base: OSMO, Quote: USD}, Service: "Provider1", Price: 0.198925},
				{Pair: Pair{Base: OSMO, Quote: USD}, Service: "Provider2", Price: 0.1697239635201354},
			},
			threshold: 0.01,
			expected: []PriceDifference{
				{
					Pair:       Pair{Base: OSMO, Quote: USD},
					PriceA:     0.198925,
					PriceB:     0.1697239635201354,
					Difference: 0.02920103647986458,
					Threshold:  0.01,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Compare(tt.prices, tt.threshold)
			assert.Equal(t, tt.expected, result)
		})
	}
}
