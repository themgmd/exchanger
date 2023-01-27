package usecase

import (
	"context"
	"errors"
	"fmt"
	"onemgvv/exchanger/internal/models"
)

func (uc useCase) UpdateRate(ctx context.Context, params models.CurrencyParams, rate float64) error {
	exists, err := uc.repo.CheckExist(ctx, params)
	if err != nil {
		return fmt.Errorf("uc.repo.CheckExist: %w", err)
	}

	if exists {
		return errors.New("current currency pair already exist")
	}

	updatedPair := models.NewCurrencyPair(params.CurrencyFrom, params.CurrencyTo, rate)
	if err = uc.repo.Update(ctx, *updatedPair); err != nil {
		return fmt.Errorf("uc.repo.Update: %w", err)
	}

	return nil
}
