package scheduler

import (
	"context"
	"exchanger/internal/config"
	"exchanger/internal/currency/types"
	currencyApi "exchanger/pkg/currencyapi"
	"exchanger/pkg/errors"
	"exchanger/pkg/pagination"
	"github.com/robfig/cron"
	"log/slog"
	"math"
	"time"
)

type Currency interface {
	UpdateRate(ctx context.Context, id int, rate float64) error
	List(ctx context.Context, pag pagination.Pagination) ([]types.CurrencyPair, int, error)
}

type Scheduler struct {
	currency  Currency
	scheduler *cron.Cron
}

func NewScheduler(currency Currency) *Scheduler {
	localTime, _ := time.LoadLocation("Europe/Moscow")
	scheduler := cron.NewWithLocation(localTime)

	return &Scheduler{
		currency:  currency,
		scheduler: scheduler,
	}
}

func (s *Scheduler) Start(ctx context.Context) error {
	err := s.scheduler.AddFunc("@hourly", func() {
		s.UpdatePairs(ctx)
	})
	if err != nil {
		return errors.Wrap(err, "cron start error")
	}

	return nil
}

func (s *Scheduler) UpdatePairs(ctx context.Context) {
	pairs, _, err := s.currency.List(ctx, pagination.Pagination{
		Limit:  math.MaxInt,
		Offset: 0,
	})
	if err != nil {
		slog.Info("list pairs error", "error", err.Error())
		return
	}

	conf := config.Get()

	currApi := currencyApi.New(
		currencyApi.Key(conf.CurrencyApi.Key),
		currencyApi.Link(conf.CurrencyApi.Link),
	)

	for _, pair := range pairs {
		resp, err := currApi.Latest(pair.CurrencyFrom, pair.CurrencyTo)
		if err != nil {
			slog.Error("fetch new currency pair rate",
				"error", err.Error(),
				"pair_id", pair.ID,
				"currency_from", pair.CurrencyFrom,
				"currency_to", pair.CurrencyTo,
				"old_rate", pair.Rate,
			)
			continue
		}

		if err = s.currency.UpdateRate(ctx, pair.ID, resp.Data[pair.CurrencyTo]); err != nil {
			slog.Error("fetch new currency pair rate",
				"error", err.Error(),
				"pair_id", pair.ID,
				"currency_from", pair.CurrencyFrom,
				"currency_to", pair.CurrencyTo,
				"old_rate", pair.Rate,
				"response_data", resp.Data[pair.CurrencyTo],
			)
			continue
		}
	}
}
