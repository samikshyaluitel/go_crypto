package api

import (
	"encoding/json"
	"fmt"
	"go_crypto/models"
	"io/ioutil"
	"net/http"
)

// fetchSymbol fetches symbols from symbol api
func fetchSymbol(symbol string) (*models.Symbols, error) {
	url := fmt.Sprintf("https://api.hitbtc.com/api/2/public/symbol/%s", symbol)
	var obj = &models.Symbols{}
	err := fetcher(url, obj)
	return obj, err
}

// fetchCurrency fetches currency from currency api
func fetchCurrency(currency string) (*models.Currency, error) {
	url := fmt.Sprintf("https://api.hitbtc.com/api/2/public/currency/%s", currency)
	var obj = &models.Currency{}
	err := fetcher(url, obj)
	return obj, err
}

// fetchTicker fetches ticker from ticker api
func fetchTicker(symbol string) (*models.Ticker, error) {
	url := fmt.Sprintf("https://api.hitbtc.com/api/2/public/ticker/%s", symbol)
	var obj = &models.Ticker{}
	err := fetcher(url, obj)
	return obj, err
}

// fetcher is the method to make api calls
func fetcher(url string, obj interface{}) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	err = json.Unmarshal(body, obj)
	return err
}

// currencyResponseGenerator generates response from symbol,currency and ticker
func currencyResponseGenerator(s models.Symbols, c models.Currency, t models.Ticker) models.CurrencyReponse {
	return models.CurrencyReponse{
		ID:           c.ID,
		FullName:     c.FullName,
		BaseCurrency: s.BaseCurrency,
		FeeCurrency:  s.FeeCurrency,
		Low:          t.Low,
		High:         t.High,
		Bid:          t.Bid,
		Last:         t.Last,
		Open:         t.Open,
		Ask:          t.Ask,
	}
}
