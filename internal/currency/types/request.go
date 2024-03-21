package types

type CurrencyPairsRequest struct {
	CurrencyFrom string `json:"currencyFrom"`
	CurrencyTo   string `json:"currencyTo"`
}

type ExchangeRequest struct {
	CurrencyPairsRequest
	Amount float64 `json:"amount"`
}

type Pagination struct {
	Limit  int `json:"limit"`
	Offset int `json:"offset"`
}
