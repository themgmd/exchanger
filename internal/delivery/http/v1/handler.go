package v1

import (
	"github.com/gofiber/fiber/v2"
	"github.com/onemgvv/exchanger/internal/config"
	"github.com/onemgvv/exchanger/internal/delivery/http/response"
	"github.com/onemgvv/exchanger/internal/domain/entity"
	"github.com/onemgvv/exchanger/internal/infrastructure/currencies"
	"strings"
)

type UseCase interface {
	CreatePair(params entity.CurrencyPairParams, well float64) (*entity.CurrencyPair, error)
	Exchange(params entity.CurrencyPair) (float64, bool)
	Aggregate() []entity.CurrencyPairParams
}

type ApiV1Handler struct {
	cfg *config.Config
	uc  UseCase
}

func NewApiV1Handler(uc UseCase, cfg *config.Config) *ApiV1Handler {
	return &ApiV1Handler{cfg, uc}
}

func (h ApiV1Handler) CreatePairs(ctx *fiber.Ctx) error {
	var body entity.CurrencyPairParams
	if err := ctx.BodyParser(&body); err != nil {
		return ctx.Status(fiber.StatusUnprocessableEntity).JSON(
			response.DefaultHttpResponse{Success: false, Comment: "Invalid request body"},
		)
	}
	body.CurrencyFrom = strings.ToUpper(body.CurrencyFrom)
	body.CurrencyTo = strings.ToUpper(body.CurrencyTo)

	resp, err := currencies.FetchCurrency(h.cfg, body.CurrencyFrom, body.CurrencyTo)
	if err != nil {
		return ctx.SendStatus(fiber.StatusInternalServerError)
	}

	result, err := h.uc.CreatePair(body, resp.Data[body.CurrencyTo])
	if result != nil && err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(
			response.DefaultHttpResponse{Success: false, Comment: err.Error()},
		)
	}
	if err != nil && result == nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(
			response.DefaultHttpResponse{Success: false, Comment: "Unknown Error."},
		)
	}

	return ctx.Status(fiber.StatusCreated).JSON(result)
}

func (h ApiV1Handler) Exchange(ctx *fiber.Ctx) error {
	var body entity.CurrencyPair
	if err := ctx.BodyParser(&body); err != nil {
		return ctx.Status(fiber.StatusUnprocessableEntity).JSON(
			response.DefaultHttpResponse{Success: false, Comment: "Invalid request body"},
		)
	}
	body.CurrencyFrom = strings.ToUpper(body.CurrencyFrom)
	body.CurrencyTo = strings.ToUpper(body.CurrencyTo)

	result, ok := h.uc.Exchange(body)
	if result == -1 && !ok {
		return ctx.Status(fiber.StatusOK).JSON(
			response.DefaultHttpResponse{
				Success: ok,
				Comment: "Currency pair not found",
			})
	}

	if result == 0 && !ok {
		return ctx.Status(fiber.StatusInternalServerError).JSON(
			response.DefaultHttpResponse{
				Success: ok,
				Comment: "Unknown Error",
			})
	}

	return ctx.Status(fiber.StatusOK).JSON(
		response.ExchangeResponse{
			DefaultHttpResponse: response.DefaultHttpResponse{
				Success: ok,
			},
			CurrencyPairParams: entity.CurrencyPairParams{
				CurrencyTo:   body.CurrencyTo,
				CurrencyFrom: body.CurrencyFrom,
			},
			Result: result,
		})
}

func (h ApiV1Handler) Aggregate(ctx *fiber.Ctx) error {
	var success bool
	result := h.uc.Aggregate()
	if len(result) > 0 {
		success = true
	}

	return ctx.Status(fiber.StatusOK).JSON(
		response.AggregateResponse{
			DefaultHttpResponse: response.DefaultHttpResponse{
				Success: success,
			},
			Data: result,
		},
	)
}
