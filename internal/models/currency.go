package models

type CurrencyPair struct {
	CurrencyFrom string  `json:"currencyFrom" db:"currency_from"`
	CurrencyTo   string  `json:"currencyTo" db:"currency_to"`
	Rate         float64 `json:"rate" db:"rate"`
}

func NewCurrencyPair(from, to string, rate float64) *CurrencyPair {
	return &CurrencyPair{from, to, rate}
}

type CurrencyParams struct {
	CurrencyFrom string `json:"currencyFrom"`
	CurrencyTo   string `json:"currencyTo"`
}

func NewCurrencyParams(from, to string) *CurrencyParams {
	return &CurrencyParams{from, to}
}
