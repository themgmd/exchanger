package currency

import (
	"errors"
	"github.com/onemgvv/exchanger/internal/domain/entity"
	"log"
)

type Repository interface {
	Insert(entity.CurrencyPair) error
	Update(entity.CurrencyPair) error
	Get(entity.CurrencyPairParams) *entity.CurrencyPair
	Select() []entity.CurrencyPairParams
}

type UseCase struct {
	repo Repository
}

func NewUseCase(repo Repository) *UseCase {
	return &UseCase{repo}
}

func (u UseCase) CreatePair(params entity.CurrencyPairParams, well float64) (*entity.CurrencyPair, error) {
	candidate := u.repo.Get(params)
	if candidate != nil {
		return candidate, errors.New("this pair already exists")
	}

	pair := entity.NewCurrencyPair(params.CurrencyFrom, params.CurrencyTo, well)
	if err := u.repo.Insert(*pair); err != nil {
		log.Printf("[Change Pair | ERROR INSERT]: %s\n", err.Error())
		return nil, err
	}

	return pair, nil
}

func (u UseCase) Exchange(params entity.CurrencyPair) (float64, bool) {
	find := entity.CurrencyPairParams{
		CurrencyFrom: params.CurrencyFrom,
		CurrencyTo:   params.CurrencyTo,
	}
	candidate := u.repo.Get(find)
	if candidate == nil {
		return -1, false
	}

	return params.Well * candidate.Well, true
}

func (u UseCase) UpdateRate(params entity.CurrencyPair) error {
	find := entity.CurrencyPairParams{
		CurrencyFrom: params.CurrencyFrom,
		CurrencyTo:   params.CurrencyTo,
	}

	candidate := u.repo.Get(find)
	if candidate == nil {
		return errors.New("course pair not found")
	}
	if err := u.repo.Update(params); err != nil {
		log.Printf("[Update Rate | ERROR Update]: %s\n", err.Error())
		return err
	}
	return nil
}

func (u UseCase) Aggregate() []entity.CurrencyPairParams {
	return u.repo.Select()
}
