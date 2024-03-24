package currency

import (
	"exchanger/internal/config"
	"exchanger/internal/currency/dhttp"
	"exchanger/internal/currency/repo"
	"exchanger/pkg/cache"
	currencyApi "exchanger/pkg/currencyapi"
	"exchanger/pkg/postgre"
	"github.com/go-chi/chi/v5"
)

func Setup(db *postgre.DB, router chi.Router) *Currency {
	currRepo := repo.NewCurrency(db)
	inMemory := cache.New()
	currApi := currencyApi.New(
		currencyApi.Link(config.Get().CurrencyApi.Link),
		currencyApi.Key(config.Get().CurrencyApi.Key),
	)

	service := New(currRepo, inMemory, currApi)
	dhttp.NewCurrency(service).SetupRoutes(router)
	return service
}
