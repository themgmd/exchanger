package scheduler

import (
	"github.com/onemgvv/exchanger/internal/config"
	"github.com/onemgvv/exchanger/internal/domain/entity"
	"github.com/onemgvv/exchanger/internal/domain/usecase/currency"
	"github.com/onemgvv/exchanger/internal/infrastructure/currencies"
	"github.com/robfig/cron"
	"log"
	"time"
)

type Scheduler struct {
	cfg *config.Config
	*currency.UseCase
	scheduler *cron.Cron
}

func NewScheduler(cfg *config.Config, uc *currency.UseCase) *Scheduler {
	localTime, _ := time.LoadLocation("Europe/Moscow")
	scheduler := cron.NewWithLocation(localTime)
	return &Scheduler{cfg, uc, scheduler}
}

func (s *Scheduler) Start() {
	if err := s.scheduler.AddFunc("@hourly", s.UpdatePairs); err != nil {
		log.Fatalf("[WHILE CRON START ERROR]: %s", err.Error())
	}
}

func (s *Scheduler) UpdatePairs() {
	for _, pair := range s.Aggregate() {
		resp, err := currencies.FetchCurrency(s.cfg, pair.CurrencyFrom, pair.CurrencyTo)
		if err != nil {
			log.Printf("[Fetch New Pair Value] | [ERROR]: %s", err.Error())
			continue
		}

		params := entity.CurrencyPair{
			CurrencyFrom: pair.CurrencyFrom,
			CurrencyTo:   pair.CurrencyTo,
			Well:         resp.Data[pair.CurrencyTo],
		}

		if err := s.UpdateRate(params); err != nil {
			log.Printf("[Fetch New Pair Value] | [ERROR]: %s", err.Error())
			continue
		}
	}
}
