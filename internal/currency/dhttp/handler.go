package dhttp

import (
	"context"
	"exchanger/internal/currency/types"
)

type Service interface {
	CreatePairs(ctx context.Context, pair types.CurrencyPair) error
	Exchange(ctx context.Context, from, to string, amount float64) (float64, error)
	UpdateRate(ctx context.Context, id int, rate float64) error
	GetRate(ctx context.Context, from, to string) (types.CurrencyPair, error)
	List(ctx context.Context, limit, offset int) ([]types.CurrencyPair, int, error)
}

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}
