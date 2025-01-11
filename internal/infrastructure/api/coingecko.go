package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"crypto-dashboard/internal/domain/models"
)

// CoinGeckoClient handles communication with the CoinGecko API
type CoinGeckoClient struct {
	baseURL    string
	httpClient *http.Client
}

// NewCoinGeckoClient creates a new API client with timeout
func NewCoinGeckoClient() *CoinGeckoClient {
	return &CoinGeckoClient{
		baseURL: "https://api.coingecko.com/api/v3",
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// FetchCryptoPrices demonstrates concurrent API calls and error handling
func (c *CoinGeckoClient) FetchCryptoPrices(cryptoIDs []string) ([]models.CryptoPrice, error) {
	// Create a channel to receive results from goroutines
	results := make(chan models.CryptoPrice, len(cryptoIDs))
	errors := make(chan error, len(cryptoIDs))

	// Launch a goroutine for each crypto ID
	// This demonstrates Go's concurrent execution model
	for _, id := range cryptoIDs {
		go func(cryptoID string) {
			// Dealing with panic
			defer func() {
				if r := recover(); r != nil {
					errors <- fmt.Errorf("panic occurred: %v", r)
				}
			}()

			url := fmt.Sprintf("%s/simple/price?ids=%s&vs_currencies=usd", c.baseURL, cryptoID)
			resp, err := c.httpClient.Get(url)
			if err != nil {
				errors <- err
				return
			}
			defer resp.Body.Close()

			// If status code is not 200, we'll panic to handle the panic
			if resp.StatusCode != http.StatusOK {
				panic(fmt.Sprintf("API returned status code: %d", resp.StatusCode))
			}

			var data map[string]map[string]float64
			if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
				errors <- err
				return
			}

			price := models.CryptoPrice{
				ID:           cryptoID,
				CurrentPrice: data[cryptoID]["usd"],
				LastUpdated:  time.Now().UTC().Format(time.RFC3339),
			}

			results <- price
		}(id)
	}

	// Collect results using a slice
	var prices []models.CryptoPrice
	for i := 0; i < len(cryptoIDs); i++ {
		select {
		case price := <-results:
			prices = append(prices, price)
		case err := <-errors:
			return nil, err
		}
	}

	return prices, nil
}

// MarketData represents the market data for a cryptocurrency
type MarketData struct {
	ID     string  `json:"id"`
	Symbol string  `json:"symbol"`
	Name   string  `json:"name"`
	Price  float64 `json:"current_price"`
}

// GetTopNCryptos fetches the top N cryptocurrencies by market cap
func (c *CoinGeckoClient) GetTopNCryptos(n int) ([]models.CryptoPrice, error) {
	url := fmt.Sprintf("%s/coins/markets?vs_currency=usd&order=market_cap_desc&per_page=%d&page=1", c.baseURL, n)

	resp, err := c.httpClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch top cryptos: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned status code: %d", resp.StatusCode)
	}

	var marketData []MarketData
	if err := json.NewDecoder(resp.Body).Decode(&marketData); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	cryptoPrices := make([]models.CryptoPrice, len(marketData))
	for i, data := range marketData {
		cryptoPrices[i] = models.CryptoPrice{
			ID:           data.ID,
			Symbol:       data.Symbol,
			Name:         data.Name,
			CurrentPrice: data.Price,
			LastUpdated:  time.Now().UTC().Format(time.RFC3339),
		}
	}

	return cryptoPrices, nil
}
