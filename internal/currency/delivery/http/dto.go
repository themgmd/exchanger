package http

type CurrencyPairsDTO struct {
	CurrencyFrom string `json:"currencyFrom" validate:"required,min=3,max=3"`
	CurrencyTo   string `json:"currencyTo" validate:"required,min=3,max=3"`
}

type ExchangeDTO struct {
	CurrencyPairsDTO
	Amount float64 `json:"amount" validate:"required,min=1"`
}

type AggregateDTO struct {
	Limit  int `json:"limit" validate:"max=100"`
	Offset int `json:"offset" validate:"max=100"`
}
