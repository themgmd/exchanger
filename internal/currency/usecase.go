package currency

import (
	"context"
	"onemgvv/exchanger/internal/models"
)

type UseCase interface {
	CreatePairs(ctx context.Context, params models.CurrencyParams, rate float64) error
	Exchange(ctx context.Context, params models.CurrencyParams, amount float64) (float64, error)
	UpdateRate(ctx context.Context, params models.CurrencyParams, rate float64) error
	GetRate(ctx context.Context, params models.CurrencyParams) (*models.CurrencyPair, error)
	Aggregate(ctx context.Context, limit, offset int) ([]models.CurrencyPair, error)
	SaveInMemory(ctx context.Context, pair models.CurrencyPair)
	GetFromInMemory(ctx context.Context, params models.CurrencyParams) (float64, error)
}
