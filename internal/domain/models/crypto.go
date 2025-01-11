// Package models contains our domain entities
package models

import (
	"errors"
	"fmt"
	"time"
)

// CryptoPrice represents cryptocurrency price data
// This is our main domain entity that follows DDD principles
type CryptoPrice struct {
	ID           string  `json:"id"`
	Symbol       string  `json:"symbol"`
	Name         string  `json:"name"`
	CurrentPrice float64 `json:"current_price"`
	LastUpdated  string  `json:"last_updated"`
}

// CryptoBatch represents a collection of CryptoPrice
// We'll use this to demonstrate working with slices and concurrent processing
type CryptoBatch struct {
	Prices []CryptoPrice
}

// Validate ensures that the CryptoPrice entity is valid
func (c *CryptoPrice) Validate() error {
	if c.ID == "" {
		return errors.New("crypto ID cannot be empty")
	}
	if c.Symbol == "" {
		return errors.New("crypto symbol cannot be empty")
	}
	if c.Name == "" {
		return errors.New("crypto name cannot be empty")
	}
	if c.CurrentPrice < 0 {
		return errors.New("crypto price cannot be negative")
	}
	return nil
}

// UpdatePrice updates the current price and last updated timestamp
func (c *CryptoPrice) UpdatePrice(newPrice float64) error {
	if newPrice < 0 {
		return errors.New("price cannot be negative")
	}
	c.CurrentPrice = newPrice
	c.LastUpdated = time.Now().UTC().Format(time.RFC3339)
	return nil
}

// AddCrypto adds a new cryptocurrency to the batch
func (b *CryptoBatch) AddCrypto(crypto CryptoPrice) {
	b.Prices = append(b.Prices, crypto)
}

// GetBySymbol finds a cryptocurrency by its symbol
func (b *CryptoBatch) GetBySymbol(symbol string) (CryptoPrice, bool) {
	for _, crypto := range b.Prices {
		if crypto.Symbol == symbol {
			return crypto, true
		}
	}
	return CryptoPrice{}, false
}

// TotalValue calculates the total value of all cryptocurrencies in the batch
func (b *CryptoBatch) TotalValue() float64 {
	total := 0.0
	for _, crypto := range b.Prices {
		total += crypto.CurrentPrice
	}
	return total
}

// MustUpdatePrice updates the price and panics if the price is invalid
// This demonstrates how to test panic scenarios
func (c *CryptoPrice) MustUpdatePrice(newPrice float64) {
	if newPrice < 0 {
		panic(fmt.Sprintf("price cannot be negative: %f", newPrice))
	}
	c.CurrentPrice = newPrice
	c.LastUpdated = time.Now().UTC().Format(time.RFC3339)
}

// GetPriceAt returns the price at a specific index in the batch
// This demonstrates another panic scenario with index out of bounds
func (b *CryptoBatch) GetPriceAt(index int) CryptoPrice {
	if index < 0 || index >= len(b.Prices) {
		panic(fmt.Sprintf("index out of bounds: %d", index))
	}
	return b.Prices[index]
}
