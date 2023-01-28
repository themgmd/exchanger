package usecase

import (
	"context"
	"errors"
	"exchanger/internal/currency"
	"exchanger/internal/models"
	"fmt"
)

func (uc useCase) CreatePairs(ctx context.Context, params models.CurrencyParams, rate float64) error {
	exists, err := uc.repo.CheckExist(ctx, params)
	if err != nil && !errors.Is(err, currency.ErrCurrencyPairNotExist) {
		return fmt.Errorf(" uc.repo.CheckExist: %w", err)
	}

	if exists {
		return currency.ErrCurrencyPairAlreadyExist
	}

	newPair := models.NewCurrencyPair(params.CurrencyFrom, params.CurrencyTo, rate)
	if err = uc.repo.Insert(ctx, *newPair); err != nil {
		return fmt.Errorf("uc.repo.Insert: %w", err)
	}

	uc.SaveInMemory(ctx, *newPair)
	return nil
}
