package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
)

// Funding rate structures
type FundingRate struct {
	Symbol      string  `json:"symbol"`
	FundingRate float64 `json:"lastFundingRate"`
}

// Fetch funding rates from Binance
func fetchBinanceFundingRates(wg *sync.WaitGroup, rates chan<- FundingRate) {
	defer wg.Done()
	response, err := http.Get("https://fapi.binance.com/fapi/v1/premiumIndex")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer response.Body.Close()

	var binanceRates []FundingRate
	if err := json.NewDecoder(response.Body).Decode(&binanceRates); err != nil {
		fmt.Println(err)
		return
	}

	for _, rate := range binanceRates {
		rates <- rate
	}
}

// Fetch funding rates from Delta
func fetchDeltaFundingRates(wg *sync.WaitGroup, rates chan<- FundingRate) {
	defer wg.Done()
	response, err := http.Get("https://api.delta.exchange/v2/tickers")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer response.Body.Close()

	var deltaRates []FundingRate
	if err := json.NewDecoder(response.Body).Decode(&deltaRates); err != nil {
		fmt.Println(err)
		return
	}

	for _, rate := range deltaRates {
		rates <- rate
	}
}

// Normalize symbols for comparison
func normalizeSymbol(symbol string) string {
	switch symbol {
	case "BTCUSDT":
		return "BTC-USD"
	case "SOL-PERP":
		return "SOL-PERP"
	default:
		return symbol
	}
}

// Build arbitrage table
func buildArbitrageTable(binanceRates, deltaRates []FundingRate) {
	// Implementation of matching coins and calculating spreads
	fmt.Println("Building arbitrage table...")
}

// CORS middleware
func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		if r.Method == "OPTIONS" {
			return
		}
		next.ServeHTTP(w, r)
	})
}

// Health check handler
func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "OK")
}

// Arbitrage handler
func arbitrageHandler(w http.ResponseWriter, r *http.Request) {
	// Fetch funding rates
	var wg sync.WaitGroup
	binanceRates := make(chan FundingRate)
	deltaRates := make(chan FundingRate)

	wg.Add(2)
	go fetchBinanceFundingRates(&wg, binanceRates)
	go fetchDeltaFundingRates(&wg, deltaRates)

	// Close channels after goroutines finish
	go func() {
		wg.Wait()
		close(binanceRates)
		close(deltaRates)
	}()

	// Create slices to hold rates
	var binanceSlice, deltaSlice []FundingRate

	// Read from channels
	for rate := range binanceRates {
		binanceSlice = append(binanceSlice, rate)
	}
	for rate := range deltaRates {
		deltaSlice = append(deltaSlice, rate)
	}

	// Build arbitrage table
	buildArbitrageTable(binanceSlice, deltaSlice)
}

func main() {
	http.Handle("/api/arbitrage", corsMiddleware(http.HandlerFunc(arbitrageHandler)))
	http.Handle("/health", corsMiddleware(http.HandlerFunc(healthHandler)))
	
	fmt.Println("Server is running on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Println(err)
	}
}