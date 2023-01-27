package usecase

import (
	"onemgvv/exchanger/internal/currency"
)

type useCase struct {
	repo     currency.Repository
	inMemory currency.InMemory
}

func New(repo currency.Repository, inMemory currency.InMemory) currency.UseCase {
	return &useCase{repo, inMemory}
}
