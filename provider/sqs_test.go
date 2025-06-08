package provider

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/deividaspetraitis/price-monitor"
	"github.com/stretchr/testify/assert"
)

func TestSQSClient_GetPrices(t *testing.T) {
	tests := []struct {
		name           string
		pairs          monitor.Pairs
		mockResponse   map[string]map[string]string
		expectedPrices []monitor.PriceData
		expectedError  string
		cancelContext  bool
	}{
		{
			name:  "successful request",
			pairs: monitor.Pairs{{Base: monitor.OSMO, Quote: monitor.USD}},
			mockResponse: map[string]map[string]string{
				"uosmo": {defaultQuote: "1.23"},
			},
			expectedPrices: []monitor.PriceData{
				{
					Pair:    monitor.Pair{Base: monitor.OSMO, Quote: monitor.USD},
					Service: "SQS",
					Price:   1.23,
				},
			},
		},
		{
			name:          "server error",
			pairs:         monitor.Pairs{{Base: monitor.OSMO, Quote: monitor.USD}},
			mockResponse:  nil,
			expectedError: "unexpected status code: 500",
		},
		{
			name:  "missing price",
			pairs: monitor.Pairs{{Base: monitor.OSMO, Quote: monitor.USD}},
			mockResponse: map[string]map[string]string{
				"uatom": {defaultQuote: "10.5"},
			},
			expectedPrices: nil,
		},
		{
			name:          "context cancelled",
			pairs:         monitor.Pairs{{Base: monitor.OSMO, Quote: monitor.USD}},
			mockResponse:  map[string]map[string]string{},
			expectedError: "context canceled",
			cancelContext: true,
		},
		{
			name:  "invalid price format",
			pairs: monitor.Pairs{{Base: monitor.OSMO, Quote: monitor.USD}},
			mockResponse: map[string]map[string]string{
				"uosmo": {defaultQuote: "invalid"},
			},
			expectedError: "failed to parse price",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if tt.mockResponse == nil {
					w.WriteHeader(http.StatusInternalServerError)
					return
				}
				w.WriteHeader(http.StatusOK)
				json.NewEncoder(w).Encode(tt.mockResponse)
			}))
			defer server.Close()

			client := &SQSClient{
				BaseURL:    server.URL,
				HTTPClient: server.Client(),
			}

			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			if tt.cancelContext {
				cancel()
			}

			prices, err := client.GetPrices(ctx, tt.pairs)

			if tt.expectedError != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedPrices, prices)
			}
		})
	}
}
