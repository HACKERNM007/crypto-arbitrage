package main

import (
    "encoding/json"
    "fmt"
    "log"
    "net/http"
)

// FundingRate represents the funding rate information from an exchange
type FundingRate struct {
    Symbol string  `json:"symbol"`
    Rate   float64 `json:"fundingRate"`
}

// FetchFundingRate fetches funding rates from the specified exchange
func FetchFundingRate(url string) ([]FundingRate, error) {
    resp, err := http.Get(url)
    if err != nil {
        return nil, fmt.Errorf("failed to fetch funding rates: %w", err)
    }
    defer resp.Body.Close()

    var rates []FundingRate
    if err := json.NewDecoder(resp.Body).Decode(&rates); err != nil {
        return nil, fmt.Errorf("failed to decode response: %w", err)
    }
    return rates, nil
}

// NormalizeSymbols converts exchange symbols to a common format
func NormalizeSymbols(symbols []string) map[string]string {
    normalized := make(map[string]string)
    for _, symbol := range symbols {
        // Example normalization logic
        normalized[symbol] = symbol
    }
    return normalized
}

// CalculateArbitrageOpportunity identifies arbitrage opportunities between two exchanges
func CalculateArbitrageOpportunity(binanceRates, deltaRates []FundingRate) {
    for _, binanceRate := range binanceRates {
        for _, deltaRate := range deltaRates {
            if binanceRate.Symbol == deltaRate.Symbol {
                opportunity := binanceRate.Rate - deltaRate.Rate
                if opportunity > 0 {
                    log.Printf("Arbitrage opportunity found for %s: %f\n", binanceRate.Symbol, opportunity)
                }
            }
        }
    }
}

func main() {
    binanceURL := "https://api.binance.com/v3/fundingRate"
    deltaURL := "https://api.delta.exchange/v2/funding_rates"

    binanceRates, err := FetchFundingRate(binanceURL)
    if err != nil {
        log.Fatalf("Error fetching Binance rates: %v", err)
    }

    deltaRates, err := FetchFundingRate(deltaURL)
    if err != nil {
        log.Fatalf("Error fetching Delta rates: %v", err)
    }

    CalculateArbitrageOpportunity(binanceRates, deltaRates)
}