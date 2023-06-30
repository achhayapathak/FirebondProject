package main

import (
	"backendProject/router"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/go-resty/resty/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	apiKey := "eeaef8a22a3a7f5998cbd83ecc2fed292698ed28d7adc154738957c8d269a81d"

	// Connect to MongoDB
	client, err := connectMongoDB()
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(context.TODO())

	deleteAllRecords(client)

	fetchExchangeRates(apiKey, client)

	retrieveExchangeRatesFromDB(client)

	router := router.SetupRouter()

	log.Println("Server started")

	// Start the server in a separate goroutine
	go func() {
		log.Fatal(http.ListenAndServe(":8080", router))
	}()

	duration := 5 * time.Minute // Update interval of 5 minutes

	// Set up a ticker to trigger updates at specified intervals
	ticker := time.NewTicker(duration)
	defer ticker.Stop()

	// Run the update process in a separate goroutine
	go func() {
		for range ticker.C {
			fetchExchangeRates(apiKey, client)
		}
	}()

	// Keep the main goroutine running
	select {}

}


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

type ExchangeRateDB struct {
	ID             primitive.ObjectID `bson:"_id,omitempty"`
	Cryptocurrency string             `bson:"cryptocurrency"`
	FiatCurrency   string             `bson:"fiat_currency"`
	Rate           float64            `bson:"rate"`
	Timestamp      time.Time          `bson:"timestamp"`
}

func connectMongoDB() (*mongo.Client, error) {
	// Set up MongoDB connection options
	clientOptions := options.Client().ApplyURI("mongodb+srv://achhayapathak:achhaya@cluster0.syfn4ue.mongodb.net/Currency_Exchange?retryWrites=true&w=majority")

	// Connect to MongoDB
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MongoDB: %v", err)
	}

	// Ping the MongoDB server to ensure the connection is valid
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to ping MongoDB server: %v", err)
	}

	return client, nil
}

func processExchangeRates(responseBody string, client *mongo.Client) {
	var exchangeRateResponse ExchangeRateResponse
	err := json.Unmarshal([]byte(responseBody), &exchangeRateResponse)
	if err != nil {
		fmt.Printf("Error parsing API response: %s\n", err.Error())
		return
	}

	// Insert exchange rate data into MongoDB
	err = insertExchangeRateData(client, exchangeRateResponse)
	if err != nil {
		log.Printf("Error inserting exchange rate data: %s\n", err.Error())
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

func insertExchangeRateData(client *mongo.Client, exchangeRateResponse ExchangeRateResponse) error {
	// Access the MongoDB collection
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

func fetchExchangeRates(apiKey string, client *mongo.Client) {
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

	// Check the response status code
	if response.StatusCode() == 200 {
		// Parse and process the response here
		fmt.Println("Exchange rates fetched successfully!")
		// fmt.Println("Response Body:", response.String())
		processExchangeRates(response.String(), client)
	} else {
		// Handle the error scenario here
		log.Printf("API request failed with status code: %d\n", response.StatusCode())
	}
}

func retrieveExchangeRatesFromDB(client *mongo.Client) {
	// Access the MongoDB collection
	collection := client.Database("Currency_Exchange").Collection("exchange_rates")

	// Define a filter to retrieve all documents
	filter := bson.M{}

	// Execute the find operation
	cursor, err := collection.Find(context.TODO(), filter)
	if err != nil {
		log.Printf("Error retrieving exchange rates from the database: %s\n", err.Error())
		return
	}
	defer cursor.Close(context.TODO())

	// Iterate over the result cursor
	for cursor.Next(context.TODO()) {
		// Define a variable to store each document
		var exchangeRate ExchangeRateDB

		// Decode the document into the exchangeRate variable
		if err := cursor.Decode(&exchangeRate); err != nil {
			log.Printf("Error decoding exchange rate document: %s\n", err.Error())
			continue
		}

		// Display the exchange rate data
		fmt.Println("Exchange Rate:")
		fmt.Printf("Cryptocurrency: %s\n", exchangeRate.Cryptocurrency)
		fmt.Printf("Fiat Currency: %s\n", exchangeRate.FiatCurrency)
		fmt.Printf("Rate: %.2f\n", exchangeRate.Rate)
		fmt.Printf("Timestamp: %s\n", exchangeRate.Timestamp.String())
		fmt.Println("-----------------------------------")
	}

	// Check if any errors occurred during iteration
	if err := cursor.Err(); err != nil {
		log.Printf("Error iterating over exchange rate documents: %s\n", err.Error())
		return
	}
}

func deleteAllRecords(client *mongo.Client) {
	// Access the MongoDB collection
	collection := client.Database("Currency_Exchange").Collection("exchange_rates")

	// Define an empty filter to delete all documents
	filter := bson.M{}

	// Perform the deletion operation
	result, err := collection.DeleteMany(context.TODO(), filter)
	if err != nil {
		log.Printf("Error deleting documents: %s\n", err.Error())
		return
	}

	// Display the number of deleted documents
	fmt.Printf("Deleted %d documents\n", result.DeletedCount)
}


