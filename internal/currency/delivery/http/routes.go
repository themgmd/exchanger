package http

import (
	"exchanger/internal/currency"
	"github.com/gofiber/fiber/v2"
)

func MapCurrencyHandlers(currencyHandlers fiber.Router, h currency.Handlers) {
	currencyHandlers.Post("/create_pairs", h.CreatePairs)
	currencyHandlers.Post("/exchange", h.Exchange)
	currencyHandlers.Post("/get_rate", h.GetRate)
	currencyHandlers.Post("/aggregate", h.Aggregate)
}
