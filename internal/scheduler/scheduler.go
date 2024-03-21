package scheduler

import (
	"context"
	"exchanger/internal/config"
	"exchanger/internal/currency_d"
	"exchanger/internal/models"
	currencyApi "exchanger/pkg/currencyapi"
	"github.com/robfig/cron"
	"log"
	"time"
)

type Scheduler struct {
	cfg       *config.Config
	uc        currency_d.UseCase
	scheduler *cron.Cron
}

func NewScheduler(cfg *config.Config, uc currency_d.UseCase) *Scheduler {
	localTime, _ := time.LoadLocation("Europe/Moscow")
	scheduler := cron.NewWithLocation(localTime)
	return &Scheduler{cfg, uc, scheduler}
}

func (s *Scheduler) Start(ctx context.Context) {
	if err := s.scheduler.AddFunc("@hourly", func() {
		s.UpdatePairs(ctx)
	}); err != nil {
		log.Fatalf("[WHILE CRON START ERROR]: %s", err.Error())
	}
}

func (s *Scheduler) UpdatePairs(ctx context.Context) {
	pairs, err := s.uc.Aggregate(ctx, 0, 0)
	if err != nil {
		log.Printf("Error occured while aggregate all pairs")
		return
	}

	apiConfig := currencyApi.APIConfig{
		Link: s.cfg.API.Link,
		Key:  s.cfg.API.Key,
	}

	for _, pair := range pairs {
		resp, err := currencyApi.FetchCurrency(apiConfig, pair.CurrencyFrom, pair.CurrencyTo)
		if err != nil {
			log.Printf("[Fetch New Pair Value] | [ERROR]: %s", err.Error())
			continue
		}

		params := models.NewCurrencyParams(pair.CurrencyFrom, pair.CurrencyTo)
		if err := s.uc.UpdateRate(ctx, *params, resp.Data[pair.CurrencyTo]); err != nil {
			log.Printf("[Fetch New Pair Value] | [ERROR]: %s", err.Error())
			continue
		}
	}
}
