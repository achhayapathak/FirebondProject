package main

import (
	"backendProject/dbqueries/datainsert"
	"backendProject/dbqueries/dataretrieve"
	"backendProject/dbqueries/dbconnection"
	"backendProject/router"
	"context"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
)

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// CryptoCompare API KEY
	apiKey := os.Getenv("API_KEY")

	// Connect to MongoDB
	client, err := dbconnection.ConnectMongoDB()
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(context.TODO())

	// The fetchExchangeRates function is responsible for making the API request
	// to fetch the exchange rates and handling the response and any potential errors.
	datainsert.FetchExchangeRates(apiKey, client)

	// The retrieveExchangeRatesFromDB function retrieves all the exchange rate documents
	// from the MongoDB collection and displays the same.
	dataretrieve.RetrieveExchangeRatesFromDB(client)

	router := router.SetupRouter()

	// Start the server in a separate goroutine
	go func() {
		log.Fatal(http.ListenAndServe(":8080", router))
		log.Println("Server started")
	}()

	// this code sets up a ticker to trigger updates at a specified interval
	duration := 5 * time.Minute // Update interval of 5 minutes

	// Set up a ticker to trigger updates at specified intervals
	ticker := time.NewTicker(duration)
	defer ticker.Stop()

	// Run the update process in a separate goroutine
	go func() {
		for range ticker.C {
			datainsert.FetchExchangeRates(apiKey, client)
		}
	}()

	// Keep the main goroutine running
	select {}

}

