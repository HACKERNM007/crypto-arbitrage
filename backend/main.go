// Go package for arbitrage calculation and fetching funding rates from Binance and Delta Exchange

package main

import (
    "encoding/json"
    "fmt"
    "log"
    "net/http"
    "sync"
    "time"
    "github.com/gorilla/handlers"
)

// FundingRate structure to hold the fetched funding rate
type FundingRate struct {
    Symbol      string  `json:"symbol"`
    FundingRate float64 `json:"fundingRate"`
}

// FetchFundingRate fetches the funding rate for a given symbol from Binance
func FetchFundingRate(symbol string) (FundingRate, error) {
    resp, err := http.Get("https://api.binance.com/api/v3/fundingRate?symbol=" + symbol)
    if err != nil {
        return FundingRate{}, err
    }
    defer resp.Body.Close()

    var fundingRates []FundingRate
    if err := json.NewDecoder(resp.Body).Decode(&fundingRates); err != nil {
        return FundingRate{}, err
    }

    return fundingRates[0], nil
}

// Fetch Delta Exchange funding rates (dummy function for demonstration)
func FetchDeltaFundingRate(symbol string) (FundingRate, error) {
    // Placeholder for actual implementation
    return FundingRate{Symbol: symbol, FundingRate: 0.0025}, nil // Dummy rate
}

// NormalizeSymbol normalizes the symbol from Delta Exchange to Binance
func NormalizeSymbol(symbol string) string {
    return symbol + "USDT" // Dummy normalization for demonstration
}

// ArbitrageCalculation performs the arbitrage calculation
func ArbitrageCalculation(binanceRate, deltaRate float64) float64 {
    return (binanceRate - deltaRate) * 100
}

// main function to initialize the server and endpoints
func main() {
    http.HandleFunc("/arbitrage", func(w http.ResponseWriter, r *http.Request) {
        symbols := []string{"BTCUSDT", "ETHUSDT"}
        var wg sync.WaitGroup
        results := make(map[string]float64)

        for _, symbol := range symbols {
            wg.Add(1)
            go func(s string) {
                defer wg.Done()
                binanceRate, err := FetchFundingRate(s)
                if err != nil {
                    log.Println("Error fetching Binance funding rate:", err)
                    return
                }
                deltaRate, err := FetchDeltaFundingRate(NormalizeSymbol(s))
                if err != nil {
                    log.Println("Error fetching Delta funding rate:", err)
                    return
                }
                results[s] = ArbitrageCalculation(binanceRate.FundingRate, deltaRate.FundingRate)
            }(symbol)
        }

        wg.Wait()

        response, err := json.Marshal(results)
        if err != nil {
            http.Error(w, "Unable to marshal response", http.StatusInternalServerError)
            return
        }
        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(http.StatusOK)
        w.Write(response)
    })

    // CORS support
    corsHandler := handlers.CORS()(http.DefaultServeMux)
    log.Println("Starting server on :8080")
    if err := http.ListenAndServe(":8080", corsHandler); err != nil {
        log.Fatal("ListenAndServe: ", err)
    }
}