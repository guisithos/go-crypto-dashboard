package main

import (
	"fmt"
	"log"

	"crypto-dashboard/internal/infrastructure/api"
)

func main() {
	// Create API client
	client := api.NewCoinGeckoClient()

	// Fetch top 20 cryptocurrencies
	prices, err := client.GetTopNCryptos(20)
	if err != nil {
		log.Fatalf("Error fetching top cryptos: %v", err)
	}

	// Print results with more detailed information
	for i, price := range prices {
		fmt.Printf("%2d. %-20s (%s) $%.2f\n",
			i+1,
			price.Name,
			price.Symbol,
			price.CurrentPrice)
	}
}
