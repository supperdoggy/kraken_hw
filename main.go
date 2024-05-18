package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"
)

type TickerResponse struct {
	Result map[string]TickerInfo `json:"result"`
}

type TickerInfo struct {
	C []string `json:"c"` // c[0] is the last trade closed price
}

type LTPResponse struct {
	Pair   string `json:"pair"`
	Amount string `json:"amount"`
}

type APIResponse struct {
	LTP []LTPResponse `json:"ltp"`
}

var pairs = []string{"BTC/USD", "BTC/CHF", "BTC/EUR"}

func getTicker(pair string) (string, error) {
	url := fmt.Sprintf("https://api.kraken.com/0/public/Ticker?pair=%s", pair)
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var tickerResponse TickerResponse
	if err := json.NewDecoder(resp.Body).Decode(&tickerResponse); err != nil {
		return "", err
	}

	tickerInfo, exists := tickerResponse.Result[pair]
	if !exists {
		return "", fmt.Errorf("pair %s not found in response", pair)
	}

	if len(tickerInfo.C) == 0 {
		return "", fmt.Errorf("no last trade closed price found for pair %s", pair)
	}

	return tickerInfo.C[0], nil
}

func ltpHandler(w http.ResponseWriter, r *http.Request) {
	var wg sync.WaitGroup
	ltpResponses := make([]LTPResponse, len(pairs))
	errors := make(chan error, len(pairs))

	for i, pair := range pairs {
		wg.Add(1)
		go func(i int, pair string) {
			defer wg.Done()
			price, err := getTicker(pair)
			if err != nil {
				errors <- err
				return
			}
			ltpResponses[i] = LTPResponse{
				Pair:   pair,
				Amount: price,
			}
		}(i, pair)
	}

	wg.Wait()
	close(errors)

	if len(errors) > 0 {
		http.Error(w, "Failed to fetch LTP data", http.StatusInternalServerError)
		return
	}

	response := APIResponse{LTP: ltpResponses}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func main() {
	http.HandleFunc("/api/v1/ltp", ltpHandler)
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("Server starting on port %s\n", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), nil))
}
