package main

import (
	"go_crypto/api"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/julienschmidt/httprouter"
)

func main() {
	var supportedItems []string
	args := os.Args[1]
	if len(args) == 0 {
		supportedItems = []string{"BTCUSD", "ETHBTC"}
	} else {
		supportedItems = strings.Split(args, ",")
	}
	currencyAPI, err := api.NewAPi(supportedItems)
	if err != nil {
		panic(err)
	}
	router := httprouter.New()
	router.GET("/currency/:symbol", currencyAPI.CurrencyHandler)
	log.Fatal(http.ListenAndServe(":8080", router))
}
