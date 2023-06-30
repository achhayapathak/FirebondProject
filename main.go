package main

import (
    "fmt"
    "github.com/go-resty/resty/v2"
)

func main() {
    apiKey := "eeaef8a22a3a7f5998cbd83ecc2fed292698ed28d7adc154738957c8d269a81d"
    fetchExchangeRates(apiKey)
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
        fmt.Println("Exchange rates fetched successfully!")
        fmt.Println("Response Body:", response.String())
        // Parse and process the response here
    } else {
        fmt.Printf("API request failed with status code: %d\n", response.StatusCode())
        // Handle the error scenario here
    }
}
