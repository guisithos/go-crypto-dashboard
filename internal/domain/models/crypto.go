// Package models contains our domain entities
package models

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
