package entity

type CurrencyPair struct {
	CurrencyFrom string  `json:"currencyFrom" db:"currency_from"`
	CurrencyTo   string  `json:"currencyTo" db:"currency_to"`
	Well         float64 `json:"value" db:"well"`
}

func NewCurrencyPair(from, to string, well float64) *CurrencyPair {
	return &CurrencyPair{from, to, well}
}

func (c *CurrencyPair) ChangeWell(val float64) {
	if val < 0 {
		return
	}
	c.Well = val
}

type CurrencyPairParams struct {
	CurrencyFrom string `json:"currencyFrom"`
	CurrencyTo   string `json:"currencyTo"`
}
