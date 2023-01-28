package usecase

import (
	"context"
	"exchanger/internal/models"
	"fmt"
	"time"
)

func (uc useCase) SaveInMemory(ctx context.Context, pair models.CurrencyPair) {
	key := fmt.Sprint(pair.CurrencyFrom, pair.CurrencyTo)
	uc.inMemory.Put(key, pair.Rate, 1*time.Hour)
}
