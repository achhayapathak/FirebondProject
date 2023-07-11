package router

import (
	"backendProject/dbqueries"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
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

func GetExchangeRate(startTime time.Time, filter primitive.D, w http.ResponseWriter) {
	// Connect to the database
	client, err := dbqueries.ConnectMongoDB()
	if err != nil {
		log.Printf("Failed to connect to MongoDB: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
	defer client.Disconnect(context.TODO())

	// Access the MongoDB collection
	collection := client.Database("Currency_Exchange").Collection("exchange_rates")

	// Find the exchange rate documents
	cursor, err := collection.Find(context.TODO(), filter)
	if err != nil {
		fmt.Printf("failed to get exchange rates: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
	defer cursor.Close(context.TODO())

	// Iterate over the cursor, decode the documents and store in response
	var response []ExchangeRateResponse
	for cursor.Next(context.TODO()) {
		var exchangeRate ExchangeRateDB
		if err := cursor.Decode(&exchangeRate); err != nil {
			fmt.Printf("failed to decode exchange rate: %v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}

		singleresponse := ExchangeRateResponse{
			Cryptocurrency: exchangeRate.Cryptocurrency,
			FiatCurrency:   exchangeRate.FiatCurrency,
			Rate:           exchangeRate.Rate,
		}

		response = append(response, singleresponse)
	}

	if err := cursor.Err(); err != nil {
		fmt.Printf("cursor error: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}

	// Writing back response to the client
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
