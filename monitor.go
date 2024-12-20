package monitor

import (
	"log"
	"math"
)

type PriceData struct {
	Pair    Pair
	Service string
	Price   float64
}

type Provider interface {
	GetPrices(cryptos Pairs) ([]PriceData, error)
}

// Fetch fetches prices for the given pairs from the given providers.
func Fetch(providers []Provider, pairs []Pair, logger *log.Logger) []PriceData {
	var prices []PriceData
	for _, provider := range providers {
		p, err := provider.GetPrices(pairs)
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

// Compare compares prices for each Pair and logs an error if the price difference exceeds the threshold.
func Compare(prices []PriceData, threshold float64, logger *log.Logger) {
	pairPrices := make(map[Pair][]float64)

	// Group prices by Pair
	for _, data := range prices {
		pairPrices[data.Pair] = append(pairPrices[data.Pair], data.Price)
	}

	// Compare prices for each Pair
	for pair, prices := range pairPrices {
		for i := 0; i < len(prices); i++ {
			for j := i + 1; j < len(prices); j++ {
				priceDiff := math.Abs(prices[i] - prices[j])
				if priceDiff > threshold {
					logger.Printf(
						"Error: Price difference for pair %s/%s exceeds threshold: %.4f > %.4f\n",
						pair.Base,
						pair.Quote,
						priceDiff,
						threshold,
					)
					PricingErrorCounter.Inc() // Increment error counter
				}
			}
		}
	}

	PricingHeartbeatCounter.Inc() // Send heartbeat signal
}
