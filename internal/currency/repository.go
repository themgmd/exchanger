package currency

import (
	"context"
	"exchanger/internal/models"
)

type Repository interface {
	Insert(ctx context.Context, pair models.CurrencyPair) error
	Update(ctx context.Context, pair models.CurrencyPair) error
	CheckExist(ctx context.Context, params models.CurrencyParams) (bool, error)
	GetRate(ctx context.Context, params models.CurrencyParams) (float64, error)
	Get(ctx context.Context, params models.CurrencyParams) (*models.CurrencyPair, error)
	Select(ctx context.Context, limit, offset int) ([]models.CurrencyPair, error)
}
