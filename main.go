package main

import (
	"encoding/json"
    "fmt"
    "github.com/go-resty/resty/v2"
	"time"
)

type ExchangeRateResponse struct {
	BTC ExchangeRate `json:"BTC"`
	ETH ExchangeRate `json:"ETH"`
	LTC ExchangeRate `json:"LTC"`
}

type ExchangeRate struct {
	USD float64 `json:"USD"`
	EUR float64 `json:"EUR"`
	GBP float64 `json:"GBP"`
}


func processExchangeRates(responseBody string) {
	var exchangeRateResponse ExchangeRateResponse
	err := json.Unmarshal([]byte(responseBody), &exchangeRateResponse)
	if err != nil {
		fmt.Printf("Error parsing API response: %s\n", err.Error())
		return
	}

	// Display the exchange rate data
	fmt.Println("Exchange Rates:")
	fmt.Printf("BTC to USD: %.2f\n", exchangeRateResponse.BTC.USD)
	fmt.Printf("BTC to EUR: %.2f\n", exchangeRateResponse.BTC.EUR)
	fmt.Printf("BTC to GBP: %.2f\n", exchangeRateResponse.BTC.GBP)
	fmt.Printf("ETH to USD: %.2f\n", exchangeRateResponse.ETH.USD)
	fmt.Printf("ETH to EUR: %.2f\n", exchangeRateResponse.ETH.EUR)
	fmt.Printf("ETH to GBP: %.2f\n", exchangeRateResponse.ETH.GBP)
	fmt.Printf("LTC to USD: %.2f\n", exchangeRateResponse.LTC.USD)
	fmt.Printf("LTC to EUR: %.2f\n", exchangeRateResponse.LTC.EUR)
	fmt.Printf("LTC to GBP: %.2f\n", exchangeRateResponse.LTC.GBP)
	fmt.Println("-----------------------------------")
}


func fetchExchangeRates(apiKey string) {
    // Create a new HTTP client
    client := resty.New()

    // Set the API key in the request header
    client.SetHeader("Authorization", "Apikey "+apiKey)

    // Make the API request to fetch exchange rates
    response, err := client.R().
        Get("https://min-api.cryptocompare.com/data/pricemulti?fsyms=BTC,ETH,LTC&tsyms=USD,EUR,GBP")

    if err != nil {
        fmt.Printf("Error making API request: %s\n", err.Error())
        return
    }

    // Check the response status code
    if response.StatusCode() == 200 {
        // Parse and process the response here
        fmt.Println("Exchange rates fetched successfully!")
        // fmt.Println("Response Body:", response.String())
		processExchangeRates(response.String())

    } else {
        // Handle the error scenario here
        fmt.Printf("API request failed with status code: %d\n", response.StatusCode())
    }
}

func main() {
	apiKey := "eeaef8a22a3a7f5998cbd83ecc2fed292698ed28d7adc154738957c8d269a81d"
	fetchExchangeRates(apiKey)

	duration := 5 * time.Minute // Update interval of 5 minutes

	// Start the initial data fetch
	fetchExchangeRates(apiKey)

	// Set up a ticker to trigger updates at specified intervals
	ticker := time.NewTicker(duration)
	defer ticker.Stop()

	// Run the update process in a separate goroutine
	go func() {
		for range ticker.C {
			fetchExchangeRates(apiKey)
		}
	}()

	// Keep the main goroutine running
	select {}
}
