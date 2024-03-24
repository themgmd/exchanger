package scheduler

import (
	"context"
	"exchanger/internal/config"
	"exchanger/internal/currency/types"
	currencyApi "exchanger/pkg/currencyapi"
	"github.com/robfig/cron"
	"log"
	"log/slog"
	"time"
)

type Currency interface {
	UpdateRate(ctx context.Context, id int, rate float64) error
	List(ctx context.Context, limit, offset int) ([]types.CurrencyPair, int, error)
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

func (s *Scheduler) Start(ctx context.Context) {
	if err := s.scheduler.AddFunc("@hourly", func() {
		s.UpdatePairs(ctx)
	}); err != nil {
		log.Fatalf("[WHILE CRON START ERROR]: %s", err.Error())
	}
}

func (s *Scheduler) UpdatePairs(ctx context.Context) {
	pairs, _, err := s.currency.List(ctx, 0, 0)
	if err != nil {
		log.Printf("Error occured while aggregate all pairs")
		return
	}

	conf := config.Get()

	currApi := currencyApi.New(
		currencyApi.Key(conf.CurrencyApi.Key),
		currencyApi.Link(conf.CurrencyApi.Link),
	)

	for _, pair := range pairs {
		resp, err := currApi.Fetch(pair.CurrencyFrom, pair.CurrencyTo)
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
