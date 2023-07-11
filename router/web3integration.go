package router

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/gorilla/mux"
)

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
