package types

import "exchanger/pkg/postgre"

type CurrencyPair struct {
	*postgre.BaseModel
	CurrencyFrom string  `json:"currencyFrom" db:"currency_from"`
	CurrencyTo   string  `json:"currencyTo" db:"currency_to"`
	Rate         float64 `json:"rate" db:"rate"`
}

func NewCurrencyPair(from, to string, rate float64) *CurrencyPair {
	return &CurrencyPair{
		CurrencyFrom: from,
		CurrencyTo:   to,
		Rate:         rate,
	}
}
