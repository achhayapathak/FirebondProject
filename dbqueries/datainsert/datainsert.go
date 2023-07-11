package datainsert

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/go-resty/resty/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type ExchangeRate struct {
	USD float64 `json:"USD"`
	EUR float64 `json:"EUR"`
	GBP float64 `json:"GBP"`
}

type ExchangeRateResponse struct {
	BTC ExchangeRate `json:"BTC"`
	ETH ExchangeRate `json:"ETH"`
	LTC ExchangeRate `json:"LTC"`
}

type ExchangeRateDB struct {
	ID             primitive.ObjectID `bson:"_id,omitempty"`
	Cryptocurrency string             `bson:"cryptocurrency"`
	FiatCurrency   string             `bson:"fiat_currency"`
	Rate           float64            `bson:"rate"`
	Timestamp      time.Time          `bson:"timestamp"`
}

func processExchangeRates(responseBody string, client *mongo.Client) {
	var exchangeRateResponse ExchangeRateResponse
	err := json.Unmarshal([]byte(responseBody), &exchangeRateResponse)
	if err != nil {
		fmt.Printf("Error parsing API response: %s\n", err.Error())
		return
	}

	err = insertExchangeRateData(client, exchangeRateResponse)
	if err != nil {
		log.Printf("Error inserting exchange rate data: %s\n", err.Error())
		return
	}

}

func insertExchangeRateData(client *mongo.Client, exchangeRateResponse ExchangeRateResponse) error {
	collection := client.Database("Currency_Exchange").Collection("exchange_rates")

	// Prepare the exchange rate data for insertion
	exchangeRates := []ExchangeRateDB{
		{
			Cryptocurrency: "BTC",
			FiatCurrency:   "USD",
			Rate:           exchangeRateResponse.BTC.USD,
			Timestamp:      time.Now(),
		},
		{
			Cryptocurrency: "BTC",
			FiatCurrency:   "EUR",
			Rate:           exchangeRateResponse.BTC.EUR,
			Timestamp:      time.Now(),
		},
		{
			Cryptocurrency: "BTC",
			FiatCurrency:   "GBP",
			Rate:           exchangeRateResponse.BTC.GBP,
			Timestamp:      time.Now(),
		},
		{
			Cryptocurrency: "ETH",
			FiatCurrency:   "USD",
			Rate:           exchangeRateResponse.ETH.USD,
			Timestamp:      time.Now(),
		},
		{
			Cryptocurrency: "ETH",
			FiatCurrency:   "EUR",
			Rate:           exchangeRateResponse.ETH.EUR,
			Timestamp:      time.Now(),
		},
		{
			Cryptocurrency: "ETH",
			FiatCurrency:   "GBP",
			Rate:           exchangeRateResponse.ETH.GBP,
			Timestamp:      time.Now(),
		},
		{
			Cryptocurrency: "LTC",
			FiatCurrency:   "USD",
			Rate:           exchangeRateResponse.LTC.USD,
			Timestamp:      time.Now(),
		},
		{
			Cryptocurrency: "LTC",
			FiatCurrency:   "EUR",
			Rate:           exchangeRateResponse.LTC.EUR,
			Timestamp:      time.Now(),
		},
		{
			Cryptocurrency: "LTC",
			FiatCurrency:   "GBP",
			Rate:           exchangeRateResponse.LTC.GBP,
			Timestamp:      time.Now(),
		},
	}

	// Insert the exchange rate data into the collection
	var exchangeRatesInterface []interface{}
	for _, rate := range exchangeRates {
		exchangeRatesInterface = append(exchangeRatesInterface, rate)
	}
	_, err := collection.InsertMany(context.TODO(), exchangeRatesInterface)

	if err != nil {
		return fmt.Errorf("failed to insert exchange rate data: %v", err)
	}

	return nil
}

func FetchExchangeRates(apiKey string, client *mongo.Client) {
	// Create a new HTTP client
	client1 := resty.New()

	// Set the API key in the request header
	client1.SetHeader("Authorization", "Apikey "+apiKey)

	// Make the API request to fetch exchange rates
	response, err := client1.R().
		Get("https://min-api.cryptocompare.com/data/pricemulti?fsyms=BTC,ETH,LTC&tsyms=USD,EUR,GBP")

	if err != nil {
		log.Printf("Error making API request: %s\n", err.Error())
		return
	}

	if response.StatusCode() == 200 {
		fmt.Println("Exchange rates fetched successfully!")
		processExchangeRates(response.String(), client)
	} else {
		log.Printf("API request failed with status code: %d\n", response.StatusCode())
	}
}

