package usecase

import (
	"context"
	"errors"
	"fmt"
	"onemgvv/exchanger/internal/models"
)

func (uc useCase) CreatePairs(ctx context.Context, params models.CurrencyParams, rate float64) error {
	exists, err := uc.repo.CheckExist(ctx, params)
	if err != nil {
		return fmt.Errorf(" uc.repo.CheckExist: %w", err)
	}

	if exists {
		return errors.New("current currency pair already exist")
	}

	newPair := models.NewCurrencyPair(params.CurrencyFrom, params.CurrencyTo, rate)
	if err = uc.repo.Insert(ctx, *newPair); err != nil {
		return fmt.Errorf("uc.repo.Insert: %w", err)
	}

	uc.SaveInMemory(ctx, *newPair)
	return nil
}
