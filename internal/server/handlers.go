package server

import (
	"exchanger/internal/common/response"
	"exchanger/internal/currency"
	"exchanger/internal/currency/delivery/http"
	"github.com/gofiber/fiber/v2"
	"time"
)

type Handlers struct {
	CurrencyHandlers currency.Handlers
}

func (s Server) MapHandlers(h Handlers) {
	s.fiber.Get("/health_check", healthCheck)
	apiGroup := s.fiber.Group("/api")

	currencyGroup := apiGroup.Group("/currency")
	http.MapCurrencyHandlers(currencyGroup, h.CurrencyHandlers)
}

func healthCheck(ctx *fiber.Ctx) error {
	return ctx.JSON(response.PingPongResponse{
		Timestamp: time.Now().Format("2006-01-02 15:04:05"),
		Message:   "service alive",
	})
}
