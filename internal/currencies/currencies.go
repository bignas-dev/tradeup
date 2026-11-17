// Package currencies syncs and serves exchange-rate data for supported currency codes.
package currencies

import (
	"encoding/json"
	"fmt"
	"net/http"
	"slices"
	"strings"
)

const currencyRatesURL = "https://cdn.jsdelivr.net/npm/@fawazahmed0/currency-api@latest/v1/currencies/eur.json"

var supportedCodes = []string{
	"AUD", "CAD", "CHF", "CNY",
	"EUR", "GBP", "HKD", "INR",
	"JPY", "KRW", "NOK", "NZD",
	"SEK", "SGD", "USD",
}

type ratesResponse struct {
	Data map[string]float64 `json:"eur"`
}

// GetCurrencies fetches current EUR-based rates and upserts USD-relative rates for the supported codes.
func GetCurrencies() (map[string]float64, error) {
	resp, err := http.Get(currencyRatesURL)
	if err != nil {
		return nil, fmt.Errorf("http.Get failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("status code %d", resp.StatusCode)
	}

	var decResp ratesResponse
	decoder := json.NewDecoder(resp.Body)
	if err := decoder.Decode(&decResp); err != nil {
		return nil, fmt.Errorf("json decoding failed: %w", err)
	}

	answer := make(map[string]float64)
	for key, val := range decResp.Data {
		upperKey := strings.ToUpper(key)
		if slices.Contains(supportedCodes, upperKey) {
			answer[upperKey] = val
		}
	}

	return answer, nil
}
