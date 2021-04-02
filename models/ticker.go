package models

// Ticker details
type Ticker struct {
	Symbol string `json:"symbol"`
	Low    string `json:"low"`
	High   string `json:"high"`
	Bid    string `json:"bid"`
	Last   string `json:"last"`
	Open   string `json:"open"`
	Ask    string `json:"ask"`
}
