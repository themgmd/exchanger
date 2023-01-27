package usecase

import (
	"context"
	"fmt"
	"onemgvv/exchanger/internal/models"
	"time"
)

func (uc useCase) SaveInMemory(ctx context.Context, pair models.CurrencyPair) {
	key := fmt.Sprint(pair.CurrencyFrom, pair.CurrencyTo)
	uc.inMemory.Put(key, pair.Rate, 1*time.Hour)
}
