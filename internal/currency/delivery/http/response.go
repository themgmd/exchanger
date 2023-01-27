package http

import "onemgvv/exchanger/internal/models"

type PingPongResponse struct {
	Timestamp string `json:"timestamp"`
	Message   string `json:"message"`
}

type DefaultHttpResponse struct {
	Success bool   `json:"success"`
	Comment string `json:"comment"`
}

type GetRateResponse struct {
	DefaultHttpResponse
	Rate models.CurrencyPair `json:"rate"`
}

type ExchangeResponse struct {
	DefaultHttpResponse
	Result float64 `json:"result"`
}

type AggregateResponse struct {
	DefaultHttpResponse
	Data []models.CurrencyPair `json:"rates"`
}
