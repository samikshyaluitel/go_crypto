package api

import (
	"encoding/json"
	"go_crypto/models"
	"go_crypto/ticker"
	"net/http"
	"sync"

	"github.com/julienschmidt/httprouter"
)

// Currency store for API
type API struct {
	CurrencyStore map[string]*models.CurrencyReponse
}

// NewAPi creates a new API struct object
func NewAPi(supportedSymbols []string) (*API, error) {
	store := make(map[string]*models.CurrencyReponse)
	tCLient, err := ticker.NewTickerClient()
	if err != nil {
		return nil, err
	}
	var m sync.RWMutex
	for _, item := range supportedSymbols {
		go func(s string) {
			resp, err := constructCurrency(s)
			if err != nil {
				panic(err)
			}
			m.Lock()
			store[s] = resp
			m.Unlock()
			ch, err := tCLient.SubscribeTicker(s)
			for i := range ch {
				m.Lock()
				resp.Bid = i.Bid
				resp.Ask = i.Ask
				resp.Last = i.Last
				resp.Open = i.Open
				resp.Low = i.Low
				resp.High = i.High
				m.Unlock()
			}
		}(item)
	}
	return &API{CurrencyStore: store}, nil
}

func (ap *API) CurrencyHandler(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	symbol := params.ByName("symbol")

	if symbol != "all" {
		c, err := constructCurrency(symbol)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"err": err.Error()})
			return
		}
		json.NewEncoder(w).Encode(c)
	} else {
		json.NewEncoder(w).Encode(ap.CurrencyStore)
	}

}

func constructCurrency(symbol string) (*models.CurrencyReponse, error) {
	// Fetch symbol details
	s, err := fetchSymbol(symbol)
	if err != nil {
		return nil, err
	}
	c, err := fetchCurrency(s.BaseCurrency)
	if err != nil {
		return nil, err
	}
	t, err := fetchTicker(symbol)
	if err != nil {
		return nil, err
	}
	resp := currencyResponseGenerator(*s, *c, *t)

	return &resp, err
}
