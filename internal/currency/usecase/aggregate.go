package usecase

import (
	"context"
	"fmt"
	"onemgvv/exchanger/internal/models"
)

func (uc useCase) GetRate(ctx context.Context, params models.CurrencyParams) (*models.CurrencyPair, error) {
	pair, err := uc.repo.Get(ctx, params)
	if err != nil {
		// TODO: Specify error occurred while not found entity in db
		return nil, fmt.Errorf("uc.repo.Get: %w", err)
	}

	return pair, nil
}

func (uc useCase) Aggregate(ctx context.Context, limit, offset int) ([]models.CurrencyPair, error) {
	pairs, err := uc.repo.Select(ctx, limit, offset)
	if err != nil {
		err = fmt.Errorf("uc.repo.Select: %s", err)
	}

	return pairs, err
}
