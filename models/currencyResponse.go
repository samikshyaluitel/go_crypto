package models

// CurrencyReponse
type CurrencyReponse struct {
	ID           string `json:"id"`
	FullName     string `json:"fullName"`
	BaseCurrency string `json:"baseCurrency"`
	FeeCurrency  string `json:"feeCurrency"`
	Low          string `json:"low"`
	High         string `json:"high"`
	Bid          string `json:"bid"`
	Last         string `json:"last"`
	Open         string `json:"open"`
	Ask          string `json:"ask"`
}
