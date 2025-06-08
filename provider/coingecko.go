package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/deividaspetraitis/price-monitor"
)

// CoinGeckoClient represents the client to interact with the CoinGecko API.
type CoinGeckoClient struct {
	BaseURL    string
	HTTPClient *http.Client
}

// NewCoinGeckoClient creates a new instance of the CoinGeckoClient.
func NewCoinGeckoClient() *CoinGeckoClient {
	return &CoinGeckoClient{
		BaseURL:    "https://api.coingecko.com/api/v3",
		HTTPClient: &http.Client{},
	}
}

// CoinGeckoCoin represents a cryptocurrency in the CoinGecko API.
type CoinGeckoCoin struct {
	Coin monitor.Coin
}

// String returns the string representation of the CoinGeckoCoin.
func (c CoinGeckoCoin) String() string {
	switch c.Coin {
	case monitor.OSMO:
		return "osmosis"
	case monitor.USD:
		return "usd"
	}
	return ""
}

// NewCoinGeckoCoin creates a new instance of the CoinGeckoCoin.
func NewCoinGeckoCoin(coin monitor.Coin) CoinGeckoCoin {
	return CoinGeckoCoin{Coin: coin}
}

// GetPrices fetches the prices of cryptocurrencies in the specified currency.
func (c *CoinGeckoClient) GetPrices(ctx context.Context, cryptos monitor.Pairs) ([]monitor.PriceData, error) {
	baseCoins := make([]string, len(cryptos))
	for i, pair := range cryptos {
		baseCoins[i] = NewCoinGeckoCoin(pair.Base).String()
	}

	quoteCoin := NewCoinGeckoCoin(cryptos[0].Quote).String()
	url := fmt.Sprintf("%s/simple/price?ids=%s&vs_currencies=%s", c.BaseURL, strings.Join(baseCoins, ","), quoteCoin)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch prices: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var rawPrices map[string]map[string]float64
	if err := json.NewDecoder(resp.Body).Decode(&rawPrices); err != nil {
		return nil, err
	}

	var pricesData []monitor.PriceData
	for _, pair := range cryptos {
		baseCoin := NewCoinGeckoCoin(pair.Base).String()
		quoteCoin := NewCoinGeckoCoin(pair.Quote).String()
		if price, ok := rawPrices[baseCoin][quoteCoin]; ok {
			pricesData = append(pricesData, monitor.PriceData{
				Pair:    pair,
				Service: "CoinGecko",
				Price:   price,
			})
		}
	}

	return pricesData, nil
}
