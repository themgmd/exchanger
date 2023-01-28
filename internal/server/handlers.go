package server

import (
	"exchanger/internal/currency"
	"exchanger/internal/currency/delivery/http"
)

type Handlers struct {
	CurrencyHandlers currency.Handlers
}

func (s Server) MapHandlers(h Handlers) {
	apiGroup := s.fiber.Group("/api")

	currencyGroup := apiGroup.Group("/currency")
	http.MapCurrencyHandlers(currencyGroup, h.CurrencyHandlers)
}
