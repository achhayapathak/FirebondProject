package dbqueries

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type ExchangeRateDB struct {
	ID             primitive.ObjectID `bson:"_id,omitempty"`
	Cryptocurrency string             `bson:"cryptocurrency"`
	FiatCurrency   string             `bson:"fiat_currency"`
	Rate           float64            `bson:"rate"`
	Timestamp      time.Time          `bson:"timestamp"`
}

func RetrieveExchangeRatesFromDB(client *mongo.Client) {
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

func DeleteAllRecords(client *mongo.Client) {
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
