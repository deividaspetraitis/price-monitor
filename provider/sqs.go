package provider

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/deividaspetraitis/price-monitor"
)

// SQSClient represents the client to interact with the SQS API.
type SQSClient struct {
	BaseURL    string
	HTTPClient *http.Client
}

// NewSQSClient creates a new instance of the SQSClient.
func NewSQSClient() *SQSClient {
	return &SQSClient{
		BaseURL:    "http://localhost:9092",
		HTTPClient: &http.Client{},
	}
}

// SQSCoin represents a cryptocurrency in the SQS API.
type SQSCoin struct {
	Coin monitor.Coin
}

// String returns the string representation of the SQSCoin.
func (c SQSCoin) String() string {
	switch c.Coin {
	case monitor.OSMO:
		return "uosmo"
	case monitor.USD:
		return defaultQuote
	}
	return ""
}

// NewSQSCoin creates a new instance of the SQSCoin.
func NewSQSCoin(coin monitor.Coin) SQSCoin {
	return SQSCoin{Coin: coin}
}

const defaultQuote = "ibc/498A0751C798A0D9A389AA3691123DADA57DAA4FE165D5C75894505B876BA6E4"

// GetPrices fetches the prices of cryptocurrencies from the SQS API.
func (s *SQSClient) GetPrices(cryptos monitor.Pairs) ([]monitor.PriceData, error) {
	baseCoins := make([]string, len(cryptos))
	for i, pair := range cryptos {
		baseCoins[i] = NewSQSCoin(pair.Base).String()
	}

	url := fmt.Sprintf("%s/tokens/prices?base=%s", s.BaseURL, strings.Join(baseCoins, ","))
	resp, err := s.HTTPClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch prices: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var rawPrices map[string]map[string]string
	if err := json.NewDecoder(resp.Body).Decode(&rawPrices); err != nil {
		return nil, err
	}

	var pricesData []monitor.PriceData
	for _, pair := range cryptos {
		baseCoin := NewSQSCoin(pair.Base).String()
		if priceStr, ok := rawPrices[baseCoin][defaultQuote]; ok {
			price, err := strconv.ParseFloat(priceStr, 64)
			if err != nil {
				return nil, fmt.Errorf("failed to parse price: %w", err)
			}
			pricesData = append(pricesData, monitor.PriceData{
				Pair:    pair,
				Service: "SQS",
				Price:   price,
			})
		}
	}

	return pricesData, nil
}
