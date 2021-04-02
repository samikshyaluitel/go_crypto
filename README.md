# GO_CRYPTO

Online assignment to fetch crypto rates using `https://api.hitbtc.com/`

## PROJECT STRUCTURE

    ├── README.md
    ├── api
    │   ├── currencyHandler.go
    │   └── helpers.go
    ├── go.mod
    ├── go.sum
    ├── main
    ├── main.go
    ├── models
    │   ├── currency.go
    │   ├── currencyResponse.go
    │   ├── symbol.go
    │   └── ticker.go
    └── ticker
        └── ticker_listener.go

### Api

Includes the handlers and helper methods to fetch currency/symbols/tickers from the api

### Models

Includes go files and models for currency,symbol,ticker and api response

### Ticker

Methods to listen to ticker socket api and update responses

### Main.go

This is the entry point of the program and takes configurable symbols in the args and fetches currency,symbol and ticker details from the respective apis.

## Building and running the program

`GO111MODULE="on" go build main.go`

This will create a `main` binary which can be directly executed.

The binary takes symbols as comma separated args.
If no args are provided then the default symbols are picked up