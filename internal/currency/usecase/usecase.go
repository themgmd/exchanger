package usecase

import (
	"exchanger/internal/currency"
	"exchanger/pkg/database/inmemory"
)

type useCase struct {
	repo     currency.Repository
	inMemory inmemory.InMemory
}

func New(repo currency.Repository, inMemory inmemory.InMemory) currency.UseCase {
	return &useCase{repo, inMemory}
}
