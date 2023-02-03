package http

import (
	"exchanger/internal/common/response"
	"exchanger/internal/models"
)

type GetRateResponse struct {
	response.DefaultHttpResponse
	Rate models.CurrencyPair `json:"rate"`
}

type ExchangeResponse struct {
	response.DefaultHttpResponse
	Result float64 `json:"result"`
}

type AggregateResponse struct {
	response.DefaultHttpResponse
	Data []models.CurrencyPair `json:"rates"`
}
