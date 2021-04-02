package models

// Symbols details
type Symbols struct {
	ID           string `json:"id"`
	BaseCurrency string `json:"baseCurrency"`
	FeeCurrency  string `json:"feeCurrency"`
}
