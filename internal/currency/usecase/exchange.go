package usecase

import (
	"context"
	"errors"
	"exchanger/internal/currency"
	"exchanger/internal/models"
	"fmt"
)

func (uc useCase) Exchange(ctx context.Context, params models.CurrencyParams, amount float64) (float64, error) {
	var rate float64
	exists, err := uc.repo.CheckExist(ctx, params)
	if err != nil && !errors.Is(err, currency.ErrCurrencyPairNotExist) {
		return 0, fmt.Errorf("uc.repo.CheckExist: %w", err)
	}

	if !exists {
		return 0, currency.ErrCurrencyPairNotExist
	}

	rate, err = uc.GetFromInMemory(ctx, params)
	if err != nil {
		rate, err = uc.repo.GetRate(ctx, params)
		if err != nil {
			return 0, fmt.Errorf("uc.repo.GetRate: %w", err)
		}
	}

	converted := rate * amount
	return converted, nil
}
