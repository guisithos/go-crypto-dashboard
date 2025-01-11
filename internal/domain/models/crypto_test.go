package models

import (
	"testing"
	"time"
)

func TestCryptoPrice_Validation(t *testing.T) {
	tests := []struct {
		name    string
		crypto  CryptoPrice
		wantErr bool
	}{
		{
			name: "valid crypto price",
			crypto: CryptoPrice{
				ID:           "bitcoin",
				Symbol:       "btc",
				Name:         "Bitcoin",
				CurrentPrice: 50000.0,
				LastUpdated:  time.Now().UTC().Format(time.RFC3339),
			},
			wantErr: false,
		},
		{
			name: "invalid - empty ID",
			crypto: CryptoPrice{
				ID:           "",
				Symbol:       "btc",
				Name:         "Bitcoin",
				CurrentPrice: 50000.0,
				LastUpdated:  time.Now().UTC().Format(time.RFC3339),
			},
			wantErr: true,
		},
		{
			name: "invalid - negative price",
			crypto: CryptoPrice{
				ID:           "bitcoin",
				Symbol:       "btc",
				Name:         "Bitcoin",
				CurrentPrice: -100.0,
				LastUpdated:  time.Now().UTC().Format(time.RFC3339),
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.crypto.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("CryptoPrice.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCryptoBatch_Operations(t *testing.T) {
	// Create sample data
	btc := CryptoPrice{
		ID:           "bitcoin",
		Symbol:       "btc",
		Name:         "Bitcoin",
		CurrentPrice: 50000.0,
		LastUpdated:  time.Now().UTC().Format(time.RFC3339),
	}
	eth := CryptoPrice{
		ID:           "ethereum",
		Symbol:       "eth",
		Name:         "Ethereum",
		CurrentPrice: 3000.0,
		LastUpdated:  time.Now().UTC().Format(time.RFC3339),
	}

	t.Run("add crypto to batch", func(t *testing.T) {
		batch := CryptoBatch{}
		batch.AddCrypto(btc)
		batch.AddCrypto(eth)

		if len(batch.Prices) != 2 {
			t.Errorf("Expected batch size of 2, got %d", len(batch.Prices))
		}
	})

	t.Run("get crypto by symbol", func(t *testing.T) {
		batch := CryptoBatch{Prices: []CryptoPrice{btc, eth}}

		crypto, found := batch.GetBySymbol("btc")
		if !found {
			t.Error("Expected to find crypto with symbol 'btc'")
		}
		if crypto.ID != "bitcoin" {
			t.Errorf("Expected ID 'bitcoin', got %s", crypto.ID)
		}

		_, found = batch.GetBySymbol("invalid")
		if found {
			t.Error("Expected not to find crypto with invalid symbol")
		}
	})

	t.Run("calculate total value", func(t *testing.T) {
		batch := CryptoBatch{Prices: []CryptoPrice{btc, eth}}
		expected := 53000.0 // 50000 + 3000

		total := batch.TotalValue()
		if total != expected {
			t.Errorf("Expected total value of %f, got %f", expected, total)
		}
	})
}

func TestCryptoPrice_PriceUpdate(t *testing.T) {
	crypto := CryptoPrice{
		ID:           "bitcoin",
		Symbol:       "btc",
		Name:         "Bitcoin",
		CurrentPrice: 50000.0,
		LastUpdated:  time.Now().UTC().Format(time.RFC3339),
	}

	t.Run("valid price update", func(t *testing.T) {
		newPrice := 51000.0
		err := crypto.UpdatePrice(newPrice)
		if err != nil {
			t.Errorf("Unexpected error updating price: %v", err)
		}
		if crypto.CurrentPrice != newPrice {
			t.Errorf("Expected price %f, got %f", newPrice, crypto.CurrentPrice)
		}
	})

	t.Run("invalid price update", func(t *testing.T) {
		newPrice := -1000.0
		err := crypto.UpdatePrice(newPrice)
		if err == nil {
			t.Error("Expected error updating to negative price, got nil")
		}
	})
}

func TestCryptoPrice_MustUpdatePrice_Panic(t *testing.T) {
	crypto := CryptoPrice{
		ID:           "bitcoin",
		Symbol:       "btc",
		Name:         "Bitcoin",
		CurrentPrice: 50000.0,
		LastUpdated:  time.Now().UTC().Format(time.RFC3339),
	}

	t.Run("should panic with negative price", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("Expected panic but got none")
			} else {
				// Check if panic message is as expected
				expected := "price cannot be negative: -100.000000"
				if r.(string) != expected {
					t.Errorf("Expected panic message '%s', got '%v'", expected, r)
				}
			}
		}()

		crypto.MustUpdatePrice(-100.0)
	})

	t.Run("should not panic with valid price", func(t *testing.T) {
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("Unexpected panic: %v", r)
			}
		}()

		crypto.MustUpdatePrice(55000.0)
		if crypto.CurrentPrice != 55000.0 {
			t.Errorf("Expected price 55000.0, got %f", crypto.CurrentPrice)
		}
	})
}

func TestCryptoBatch_GetPriceAt_Panic(t *testing.T) {
	batch := CryptoBatch{
		Prices: []CryptoPrice{
			{
				ID:           "bitcoin",
				Symbol:       "btc",
				Name:         "Bitcoin",
				CurrentPrice: 50000.0,
			},
		},
	}

	tests := []struct {
		name          string
		index         int
		shouldPanic   bool
		expectedPrice float64
	}{
		{
			name:          "valid index",
			index:         0,
			shouldPanic:   false,
			expectedPrice: 50000.0,
		},
		{
			name:        "panic on negative index",
			index:       -1,
			shouldPanic: true,
		},
		{
			name:        "panic on out of bounds index",
			index:       1,
			shouldPanic: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				r := recover()
				if tt.shouldPanic && r == nil {
					t.Error("Expected panic but got none")
				}
				if !tt.shouldPanic && r != nil {
					t.Errorf("Unexpected panic: %v", r)
				}
			}()

			result := batch.GetPriceAt(tt.index)
			if !tt.shouldPanic && result.CurrentPrice != tt.expectedPrice {
				t.Errorf("Expected price %f, got %f", tt.expectedPrice, result.CurrentPrice)
			}
		})
	}
}
