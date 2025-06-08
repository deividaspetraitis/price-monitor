package monitor

import (
	"context"
	"math"
	"time"

	"github.com/deividaspetraitis/price-monitor/log"
)

type PriceData struct {
	Pair    Pair
	Service string
	Price   float64
}

type Provider interface {
	GetPrices(ctx context.Context, cryptos Pairs) ([]PriceData, error)
}

// Fetch fetches prices for the given pairs from the given providers.
func Fetch(ctx context.Context, providers []Provider, pairs []Pair, timeout time.Duration, logger log.Logger) []PriceData {
	var prices []PriceData
	for _, provider := range providers {
		providerCtx, cancel := context.WithTimeout(ctx, timeout)
		defer cancel()

		p, err := provider.GetPrices(providerCtx, pairs)
		if err != nil {
			logger.Printf("Error fetching prices: %s", err)
			continue
		}

		for crypto, data := range p {
			logger.Printf("%v: %v", crypto, data)
		}

		prices = append(prices, p...)
	}

	return prices
}

// PriceDifference holds details about a detected price mismatch.
type PriceDifference struct {
	Pair       Pair
	PriceA     float64
	PriceB     float64
	Difference float64
	Threshold  float64
}

// Compare checks price differences for each Pair and returns mismatches above the threshold.
func Compare(prices []PriceData, threshold float64) []PriceDifference {
	pairPrices := make(map[Pair][]float64)
	var diffs []PriceDifference

	// Group prices by Pair
	for _, data := range prices {
		pairPrices[data.Pair] = append(pairPrices[data.Pair], data.Price)
	}

	// Compare prices for each Pair
	for pair, ps := range pairPrices {
		for i := 0; i < len(ps); i++ {
			for j := i + 1; j < len(ps); j++ {
				diff := math.Abs(ps[i] - ps[j])
				if diff > threshold {
					diffs = append(diffs, PriceDifference{
						Pair:       pair,
						PriceA:     ps[i],
						PriceB:     ps[j],
						Difference: diff,
						Threshold:  threshold,
					})
				}
			}
		}
	}

	return diffs
}
