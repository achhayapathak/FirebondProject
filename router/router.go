package router

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

type ExchangeRateResponse struct {
	Cryptocurrency string    `json:"cryptocurrency"`
	FiatCurrency   string    `json:"fiat_currency"`
	Rate           float64   `json:"rate"`
	Timestamp      time.Time `json:"timestamp"`
}

type ExchangeRateDB struct {
	ID             primitive.ObjectID `bson:"_id,omitempty"`
	Cryptocurrency string             `bson:"cryptocurrency"`
	FiatCurrency   string             `bson:"fiat_currency"`
	Rate           float64            `bson:"rate"`
	Timestamp      time.Time          `bson:"timestamp"`
}

func connectMongoDB() (*mongo.Client, error) {
	mongoURI := os.Getenv("MONGO_URI")
	
	// MongoDB connection options
	clientOptions := options.Client().ApplyURI(mongoURI)

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

func getCurrentExchangeRate(client *mongo.Client, cryptocurrency, fiat string) (*ExchangeRateDB, error) {
	// Access the MongoDB collection
	collection := client.Database("Currency_Exchange").Collection("exchange_rates")

	// Build the query filter
	filter := bson.D{
		{Key: "cryptocurrency", Value: cryptocurrency},
		{Key: "fiat_currency", Value: fiat},
	}

	// Find the exchange rate document
	var exchangeRate ExchangeRateDB
	err := collection.FindOne(context.TODO(), filter).Decode(&exchangeRate)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("exchange rate not found for %s-%s", cryptocurrency, fiat)
		}
		return nil, fmt.Errorf("failed to get exchange rate: %v", err)
	}

	return &exchangeRate, nil
}

func getExchangeRateHistory(client *mongo.Client, cryptocurrency, fiat string) ([]ExchangeRateDB, error) {
	// Access the MongoDB collection
	collection := client.Database("Currency_Exchange").Collection("exchange_rates")

	// Calculate the start time for the past 24 hours
	startTime := time.Now().Add(-24 * time.Hour)

	// Build the query filter
	filter := bson.D{
		{Key: "cryptocurrency", Value: cryptocurrency},
		{Key: "fiat_currency", Value: fiat},
		{Key: "timestamp", Value: bson.D{{Key: "$gte", Value: startTime}}},
	}

	// Find the exchange rate documents within the past 24 hours
	cursor, err := collection.Find(context.TODO(), filter)
	if err != nil {
		return nil, fmt.Errorf("failed to get exchange rate history: %v", err)
	}
	defer cursor.Close(context.TODO())

	// Iterate over the cursor and decode the documents
	var exchangeRates []ExchangeRateDB
	for cursor.Next(context.TODO()) {
		var exchangeRate ExchangeRateDB
		if err := cursor.Decode(&exchangeRate); err != nil {
			return nil, fmt.Errorf("failed to decode exchange rate: %v", err)
		}
		exchangeRates = append(exchangeRates, exchangeRate)
	}

	if err := cursor.Err(); err != nil {
		return nil, fmt.Errorf("cursor error: %v", err)
	}

	return exchangeRates, nil
}

func getExchangeRatesByCryptocurrency(client *mongo.Client, cryptocurrency string) ([]ExchangeRateDB, error) {
	// Access the MongoDB collection
	collection := client.Database("Currency_Exchange").Collection("exchange_rates")

	// Calculate the start time for the past 5 minutes
	startTime := time.Now().Add(-5 * time.Minute)

	// Build the query filter
	filter := bson.D{
		{Key: "cryptocurrency", Value: cryptocurrency},
		{Key: "timestamp", Value: bson.D{{Key: "$gte", Value: startTime}}},
	}

	// Find the exchange rate documents
	cursor, err := collection.Find(context.TODO(), filter)
	if err != nil {
		return nil, fmt.Errorf("failed to get exchange rates: %v", err)
	}
	defer cursor.Close(context.TODO())

	// Iterate over the cursor and decode the documents
	var exchangeRates []ExchangeRateDB
	for cursor.Next(context.TODO()) {
		var exchangeRate ExchangeRateDB
		if err := cursor.Decode(&exchangeRate); err != nil {
			return nil, fmt.Errorf("failed to decode exchange rate: %v", err)
		}
		exchangeRates = append(exchangeRates, exchangeRate)
	}

	if err := cursor.Err(); err != nil {
		return nil, fmt.Errorf("cursor error: %v", err)
	}

	return exchangeRates, nil
}

func getExchangeRatesByCryptocurrencyHistory(client *mongo.Client, cryptocurrency string) ([]ExchangeRateDB, error) {
	// Access the MongoDB collection
	collection := client.Database("Currency_Exchange").Collection("exchange_rates")

	// Calculate the start time for the past 24 hours
	startTime := time.Now().Add(-24 * time.Hour)

	// Build the query filter
	filter := bson.D{
		{Key: "cryptocurrency", Value: cryptocurrency},
		{Key: "timestamp", Value: bson.D{{Key: "$gte", Value: startTime}}},
	}

	// Find the exchange rate documents within the past 24 hours
	cursor, err := collection.Find(context.TODO(), filter)
	if err != nil {
		return nil, fmt.Errorf("failed to get exchange rate history: %v", err)
	}
	defer cursor.Close(context.TODO())

	// Iterate over the cursor and decode the documents
	var exchangeRates []ExchangeRateDB
	for cursor.Next(context.TODO()) {
		var exchangeRate ExchangeRateDB
		if err := cursor.Decode(&exchangeRate); err != nil {
			return nil, fmt.Errorf("failed to decode exchange rate: %v", err)
		}
		exchangeRates = append(exchangeRates, exchangeRate)
	}

	if err := cursor.Err(); err != nil {
		return nil, fmt.Errorf("cursor error: %v", err)
	}

	return exchangeRates, nil
}

func getExchangeRates(client *mongo.Client) ([]ExchangeRateDB, error) {
	// Access the MongoDB collection
	collection := client.Database("Currency_Exchange").Collection("exchange_rates")

	// Calculate the start time for the past 5 minutes
	startTime := time.Now().Add(-5 * time.Minute)

	// Build the query filter
	filter := bson.D{
		{Key: "timestamp", Value: bson.D{{Key: "$gte", Value: startTime}}},
	}

	// Find all exchange rate documents
	cursor, err := collection.Find(context.TODO(), filter)
	if err != nil {
		return nil, fmt.Errorf("failed to get exchange rates: %v", err)
	}
	defer cursor.Close(context.TODO())

	// Iterate over the cursor and decode the documents
	var exchangeRates []ExchangeRateDB
	for cursor.Next(context.TODO()) {
		var exchangeRate ExchangeRateDB
		if err := cursor.Decode(&exchangeRate); err != nil {
			return nil, fmt.Errorf("failed to decode exchange rate: %v", err)
		}
		exchangeRates = append(exchangeRates, exchangeRate)
	}

	if err := cursor.Err(); err != nil {
		return nil, fmt.Errorf("cursor error: %v", err)
	}

	return exchangeRates, nil
}

func getExchangeRatesHistory(client *mongo.Client) ([]ExchangeRateDB, error) {
	// Access the MongoDB collection
	collection := client.Database("Currency_Exchange").Collection("exchange_rates")

	// Calculate the start time for the past 24 hours
	startTime := time.Now().Add(-24 * time.Hour)

	// Build the query filter
	filter := bson.D{
		{Key: "timestamp", Value: bson.D{{Key: "$gte", Value: startTime}}},
	}

	// Find the exchange rate documents
	cursor, err := collection.Find(context.TODO(), filter)
	if err != nil {
		return nil, fmt.Errorf("failed to get exchange rates: %v", err)
	}
	defer cursor.Close(context.TODO())

	// Iterate over the cursor and decode the documents
	var exchangeRates []ExchangeRateDB
	for cursor.Next(context.TODO()) {
		var exchangeRate ExchangeRateDB
		if err := cursor.Decode(&exchangeRate); err != nil {
			return nil, fmt.Errorf("failed to decode exchange rate: %v", err)
		}
		exchangeRates = append(exchangeRates, exchangeRate)
	}

	if err := cursor.Err(); err != nil {
		return nil, fmt.Errorf("cursor error: %v", err)
	}

	return exchangeRates, nil
}

func handleGetExchangeRate(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	cryptocurrency := vars["cryptocurrency"]
	fiat := vars["fiat"]

	client, err := connectMongoDB()
	if err != nil {
		log.Printf("Failed to connect to MongoDB: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	defer client.Disconnect(context.TODO())

	exchangeRate, err := getCurrentExchangeRate(client, cryptocurrency, fiat)
	if err != nil {
		log.Printf("Failed to get exchange rate: %v", err)
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	response := struct {
		Cryptocurrency string  `json:"cryptocurrency"`
		FiatCurrency   string  `json:"fiat_currency"`
		Rate           float64 `json:"rate"`
	}{
		Cryptocurrency: exchangeRate.Cryptocurrency,
		FiatCurrency:   exchangeRate.FiatCurrency,
		Rate:           exchangeRate.Rate,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func handleGetExchangeRateHistory(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	cryptocurrency := vars["cryptocurrency"]
	fiat := vars["fiat"]

	client, err := connectMongoDB()
	if err != nil {
		log.Printf("Failed to connect to MongoDB: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	defer client.Disconnect(context.TODO())

	exchangeRates, err := getExchangeRateHistory(client, cryptocurrency, fiat)
	if err != nil {
		log.Printf("Failed to get exchange rate history: %v", err)
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	response := make([]ExchangeRateResponse, len(exchangeRates))
	for i, exchangeRate := range exchangeRates {
		response[i] = ExchangeRateResponse{
			Cryptocurrency: exchangeRate.Cryptocurrency,
			FiatCurrency:   exchangeRate.FiatCurrency,
			Rate:           exchangeRate.Rate,
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func handleGetExchangeRatesByCryptocurrency(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	cryptocurrency := vars["cryptocurrency"]

	client, err := connectMongoDB()
	if err != nil {
		log.Printf("Failed to connect to MongoDB: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	defer client.Disconnect(context.TODO())

	exchangeRates, err := getExchangeRatesByCryptocurrency(client, cryptocurrency)
	if err != nil {
		log.Printf("Failed to get exchange rates: %v", err)
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	response := make([]ExchangeRateResponse, len(exchangeRates))
	for i, exchangeRate := range exchangeRates {
		response[i] = ExchangeRateResponse{
			Cryptocurrency: exchangeRate.Cryptocurrency,
			FiatCurrency:   exchangeRate.FiatCurrency,
			Rate:           exchangeRate.Rate,
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func handleGetExchangeRatesByCryptocurrencyHistory(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	cryptocurrency := vars["cryptocurrency"]

	client, err := connectMongoDB()
	if err != nil {
		log.Printf("Failed to connect to MongoDB: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	defer client.Disconnect(context.TODO())

	exchangeRates, err := getExchangeRatesByCryptocurrencyHistory(client, cryptocurrency)
	if err != nil {
		log.Printf("Failed to get exchange rate history: %v", err)
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	response := make([]ExchangeRateResponse, len(exchangeRates))
	for i, exchangeRate := range exchangeRates {
		response[i] = ExchangeRateResponse{
			Cryptocurrency: exchangeRate.Cryptocurrency,
			FiatCurrency:   exchangeRate.FiatCurrency,
			Rate:           exchangeRate.Rate,
			Timestamp:      exchangeRate.Timestamp,
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func handleGetExchangeRates(w http.ResponseWriter, r *http.Request) {
	client, err := connectMongoDB()
	if err != nil {
		log.Printf("Failed to connect to MongoDB: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	defer client.Disconnect(context.TODO())

	exchangeRates, err := getExchangeRates(client)
	if err != nil {
		log.Printf("Failed to get exchange rates: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	response := make([]ExchangeRateResponse, len(exchangeRates))
	for i, exchangeRate := range exchangeRates {
		response[i] = ExchangeRateResponse{
			Cryptocurrency: exchangeRate.Cryptocurrency,
			FiatCurrency:   exchangeRate.FiatCurrency,
			Rate:           exchangeRate.Rate,
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func handleGetExchangeRatesHistory(w http.ResponseWriter, r *http.Request) {
	client, err := connectMongoDB()
	if err != nil {
		log.Printf("Failed to connect to MongoDB: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	defer client.Disconnect(context.TODO())

	exchangeRates, err := getExchangeRatesHistory(client)
	if err != nil {
		log.Printf("Failed to get exchange rates: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	response := make([]ExchangeRateResponse, len(exchangeRates))
	for i, exchangeRate := range exchangeRates {
		response[i] = ExchangeRateResponse{
			Cryptocurrency: exchangeRate.Cryptocurrency,
			FiatCurrency:   exchangeRate.FiatCurrency,
			Rate:           exchangeRate.Rate,
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func GetAddressBalance(address string) (string, error) {

	// Using Infura's API to extract balance from an address
	uri := os.Getenv("INFURA_URI")

	client, err := ethclient.Dial(uri)
	if err != nil {
		return "", err
	}

	account := common.HexToAddress(address)
	balance, err := client.BalanceAt(context.Background(), account, nil)
	if err != nil {
		return "", err
	}

	return balance.String(), nil
}

func GetBalanceHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	address := params["address"]

	balance, err := GetAddressBalance(address)
	if err != nil {
		log.Println(err)
		http.Error(w, "Failed to retrieve balance", http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "Balance of address %s: %s", address, balance)
}

func handleHome(w http.ResponseWriter, r *http.Request) {
	htmlContent := `
		<!DOCTYPE html>
		<html>
		<head>
			<title>API Instructions</title>
		</head>
		<body>
			<h1>API Instructions</h1>
			<p>Use the following endpoints to interact with the API:</p>
			<ul>
				<li><strong>GET /rates/{cryptocurrency}/{fiat}</strong>: Returns the current exchange rate between the specified cryptocurrency and fiat currency.</li><br>
				<li><strong>GET /rates/history/{cryptocurrency}/{fiat}</strong>: Returns the exchange rate history between the specified cryptocurrency and fiat currency for the past 24 hours.</li><br>
				<li><strong>GET /rates/{cryptocurrency}</strong>: Returns the current exchange rates between the specified cryptocurrency and all supported fiat currencies.</li><br>
				<li><strong>GET /history/rates/{cryptocurrency}</strong>: Returns the exchange rate history between the specified cryptocurrency and all supported fiat currencies for the past 24 hours.</li><br>
				<li><strong>GET /rates</strong>: Returns the current exchange rates for all supported cryptocurrency-fiat pairs.</li><br>
				<li><strong>GET /history/rates</strong>: Returns the exchange rate for all supported cryptocurrency-fiat pairs for the past 24 hours.</li><br>
				<li><strong>GET /balance/{address}</strong>: Returns the balance of the provided Ethereum Address.</li>
			</ul>
		</body>
		</html>
	`

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write([]byte(htmlContent))
}

func SetupRouter() *mux.Router {
	r := mux.NewRouter()

	r.HandleFunc("/", handleHome).Methods("GET")

	r.HandleFunc("/rates/{cryptocurrency}/{fiat}", handleGetExchangeRate).Methods("GET")
	r.HandleFunc("/rates/history/{cryptocurrency}/{fiat}", handleGetExchangeRateHistory).Methods("GET")

	r.HandleFunc("/rates/{cryptocurrency}", handleGetExchangeRatesByCryptocurrency).Methods("GET")
	r.HandleFunc("/history/rates/{cryptocurrency}", handleGetExchangeRatesByCryptocurrencyHistory).Methods("GET")
	
	r.HandleFunc("/rates", handleGetExchangeRates).Methods("GET")
	r.HandleFunc("/history/rates", handleGetExchangeRatesHistory).Methods("GET")

	r.HandleFunc("/balance/{address}", GetBalanceHandler).Methods("GET")

	return r
}
