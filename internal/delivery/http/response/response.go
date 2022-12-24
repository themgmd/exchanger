package response

import "github.com/onemgvv/exchanger/internal/domain/entity"

type PingPongResponse struct {
	Timestamp string `json:"timestamp"`
	Message   string `json:"message"`
}

type DefaultHttpResponse struct {
	Success bool   `json:"success"`
	Comment string `json:"comment"`
}

type ExchangeResponse struct {
	DefaultHttpResponse
	entity.CurrencyPairParams
	Result float64 `json:"result"`
}

type AggregateResponse struct {
	DefaultHttpResponse
	Data []entity.CurrencyPairParams
}
