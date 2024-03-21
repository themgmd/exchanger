package currency

import (
	"context"
	"exchanger/internal/currency/types"
	"exchanger/pkg/data"
	"exchanger/pkg/errors"
	"fmt"
	"time"
)

type Repo interface {
	Create(ctx context.Context, pair types.CurrencyPair) error
	CheckExist(ctx context.Context, from, to string) error
	Update(ctx context.Context, id int, update data.Json) error
	GetById(ctx context.Context, id int) (types.CurrencyPair, error)
	Get(ctx context.Context, from, to string) (types.CurrencyPair, error)
	GetRate(ctx context.Context, from, to string) (float64, error)
	List(ctx context.Context, limit, offset int) ([]types.CurrencyPair, int, error)
}

type Cache interface {
	Put(key string, data float64, ttl time.Duration)
	Get(key string) (float64, bool)
	Delete(key string)
}

type Currency struct {
	repo  Repo
	cache Cache
}

func New(repo Repo, cache Cache) *Currency {
	return &Currency{
		repo:  repo,
		cache: cache,
	}
}

func (c Currency) CreatePairs(ctx context.Context, pair types.CurrencyPair) error {
	err := c.repo.CheckExist(ctx, pair.CurrencyFrom, pair.CurrencyTo)
	if err != nil && !errors.Is(err, types.ErrCurrencyPairNotExist) {
		return fmt.Errorf(" uc.repo.CheckExist: %w", err)
	}

	if nil == err {
		return types.ErrCurrencyPairAlreadyExist
	}

	if err = c.repo.Create(ctx, pair); err != nil {
		return fmt.Errorf("uc.repo.Insert: %w", err)
	}

	key := fmt.Sprintf("%s:%s", pair.CurrencyFrom, pair.CurrencyTo)
	c.cache.Put(key, pair.Rate, 1*time.Hour)

	return nil
}

func (c Currency) Exchange(ctx context.Context, from, to string, amount float64) (float64, error) {
	var rate float64
	err := c.repo.CheckExist(ctx, from, to)
	if err != nil && !errors.Is(err, types.ErrCurrencyPairNotExist) {
		return 0, fmt.Errorf("uc.repo.CheckExist: %w", err)
	}

	if nil == err {
		return 0, types.ErrCurrencyPairNotExist
	}

	key := fmt.Sprintf("%s:%s", from, to)
	rate, ok := c.cache.Get(key)
	if !ok {
		rate, err = c.repo.GetRate(ctx, from, to)
		if err != nil {
			return 0, fmt.Errorf("uc.repo.GetRate: %w", err)
		}
	}

	converted := rate * amount
	return converted, nil
}

func (c Currency) UpdateRate(ctx context.Context, id int, rate float64) error {
	err := c.repo.Update(ctx, id, data.Json{
		"rate": rate,
	})
	if err != nil {
		return err
	}

	pair, err := c.repo.GetById(ctx, id)
	if err != nil {
		return err
	}

	key := fmt.Sprintf("%s:%s", pair.CurrencyFrom, pair.CurrencyTo)
	c.cache.Put(key, pair.Rate, 1*time.Hour)

	return nil
}

func (c Currency) GetRate(ctx context.Context, from, to string) (types.CurrencyPair, error) {
	pair, err := c.repo.Get(ctx, from, to)
	if err != nil {
		// TODO: Specify error occurred while not found entity in db
		return types.CurrencyPair{}, fmt.Errorf("uc.repo.Get: %w", err)
	}

	return pair, nil
}

func (c Currency) List(ctx context.Context, limit, offset int) ([]types.CurrencyPair, int, error) {
	pairs, total, err := c.repo.List(ctx, limit, offset)
	if err != nil {
		err = fmt.Errorf("uc.repo.Select: %s", err)
		return pairs, 0, err
	}

	return pairs, total, nil
}
