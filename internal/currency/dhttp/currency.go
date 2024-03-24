package dhttp

import (
	"github.com/go-chi/chi/v5"
)

type Currency struct {
	handler *Handler
}

func NewCurrency(service Service) *Currency {
	return &Currency{
		handler: NewHandler(service),
	}
}

func (c Currency) SetupRoutes(router chi.Router) {
	router.Post("/currency", c.handler.CreatePair)
	router.Get("/currency", c.handler.List)
	router.Post("/currency/rate", c.handler.Exchange)
	router.Get("/currency/{currency}", c.handler.GetRate)
}
