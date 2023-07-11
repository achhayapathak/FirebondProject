package router

import (
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
)

func handleHome(w http.ResponseWriter, r *http.Request) {
	htmlContent := `
		<!DOCTYPE html>
		<html>
		<head>	<title>API Instructions</title>	</head>
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

func handleGetExchangeRate(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	cryptocurrency := vars["cryptocurrency"]
	fiat := vars["fiat"]

	// Starting time to query
	startTime := time.Now().Add(-5 * time.Minute)

	// Build the query filter
	filter := bson.D{
		{Key: "cryptocurrency", Value: cryptocurrency},
		{Key: "fiat_currency", Value: fiat},
		{Key: "timestamp", Value: bson.D{{Key: "$gte", Value: startTime}}},
	}

	GetExchangeRate(startTime, filter, w)
}

func handleGetExchangeRateHistory(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	cryptocurrency := vars["cryptocurrency"]
	fiat := vars["fiat"]

	startTime := time.Now().Add(-24 * time.Hour)

	filter := bson.D{
		{Key: "cryptocurrency", Value: cryptocurrency},
		{Key: "fiat_currency", Value: fiat},
		{Key: "timestamp", Value: bson.D{{Key: "$gte", Value: startTime}}},
	}

	GetExchangeRate(startTime, filter, w)
}

func handleGetExchangeRatesByCryptocurrency(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	cryptocurrency := vars["cryptocurrency"]

	startTime := time.Now().Add(-5 * time.Minute)

	filter := bson.D{
		{Key: "cryptocurrency", Value: cryptocurrency},
		{Key: "timestamp", Value: bson.D{{Key: "$gte", Value: startTime}}},
	}

	GetExchangeRate(startTime, filter, w)
}

func handleGetExchangeRatesByCryptocurrencyHistory(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	cryptocurrency := vars["cryptocurrency"]

	startTime := time.Now().Add(-24 * time.Hour)

	filter := bson.D{
		{Key: "cryptocurrency", Value: cryptocurrency},
		{Key: "timestamp", Value: bson.D{{Key: "$gte", Value: startTime}}},
	}

	GetExchangeRate(startTime, filter, w)
}

func handleGetExchangeRates(w http.ResponseWriter, r *http.Request) {
	startTime := time.Now().Add(-5 * time.Minute)

	filter := bson.D{
		{Key: "timestamp", Value: bson.D{{Key: "$gte", Value: startTime}}},
	}

	GetExchangeRate(startTime, filter, w)
}

func handleGetExchangeRatesHistory(w http.ResponseWriter, r *http.Request) {
	startTime := time.Now().Add(-24 * time.Hour)

	filter := bson.D{
		{Key: "timestamp", Value: bson.D{{Key: "$gte", Value: startTime}}},
	}

	GetExchangeRate(startTime, filter, w)
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
