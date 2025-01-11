package api

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestNewCoinGeckoClient(t *testing.T) {
	client := NewCoinGeckoClient()

	if client.baseURL != "https://api.coingecko.com/api/v3" {
		t.Errorf("Expected base URL to be 'https://api.coingecko.com/api/v3', got %s", client.baseURL)
	}

	if client.httpClient.Timeout != 10*time.Second {
		t.Errorf("Expected timeout to be 10 seconds, got %v", client.httpClient.Timeout)
	}
}

func TestFetchCryptoPrices_Success(t *testing.T) {
	// Create a test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"bitcoin":{"usd":50000}}`))
	}))
	defer server.Close()

	// Create client with test server URL
	client := &CoinGeckoClient{
		baseURL:    server.URL,
		httpClient: &http.Client{Timeout: 10 * time.Second},
	}

	prices, err := client.FetchCryptoPrices([]string{"bitcoin"})
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if len(prices) != 1 {
		t.Errorf("Expected 1 price, got %d", len(prices))
	}

	if prices[0].ID != "bitcoin" || prices[0].CurrentPrice != 50000 {
		t.Errorf("Expected bitcoin price to be 50000, got %f", prices[0].CurrentPrice)
	}
}

func TestFetchCryptoPrices_PanicRecovery(t *testing.T) {
	// Create a test server that will cause a panic
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`invalid json that will cause panic`))
	}))
	defer server.Close()

	client := &CoinGeckoClient{
		baseURL:    server.URL,
		httpClient: &http.Client{Timeout: 10 * time.Second},
	}

	_, err := client.FetchCryptoPrices([]string{"bitcoin"})
	if err == nil {
		t.Error("Expected error from panic recovery, got nil")
	}
}

func TestFetchCryptoPrices_Error(t *testing.T) {
	// Create a test server that returns an error
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	client := &CoinGeckoClient{
		baseURL:    server.URL,
		httpClient: &http.Client{Timeout: 10 * time.Second},
	}

	_, err := client.FetchCryptoPrices([]string{"bitcoin"})
	if err == nil {
		t.Error("Expected error from API call, got nil")
	}
}
