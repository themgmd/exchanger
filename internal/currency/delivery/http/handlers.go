package http

import (
	"context"
	"github.com/gofiber/fiber/v2"
	"onemgvv/exchanger/internal/config"
	"onemgvv/exchanger/internal/currency"
	"onemgvv/exchanger/internal/models"
	currencyApi "onemgvv/exchanger/pkg/currency_api"
	"onemgvv/exchanger/pkg/utils"
)

type handlers struct {
	ctx context.Context
	cfg *config.Config
	uc  currency.UseCase
}

func New(ctx context.Context, cfg *config.Config, uc currency.UseCase) currency.Handlers {
	return &handlers{ctx, cfg, uc}
}

func (h handlers) CreatePairs(ctx *fiber.Ctx) error {
	var (
		dto      CurrencyPairsDTO
		respBody DefaultHttpResponse
	)

	// Validate Request Body
	if err := utils.ReadRequest(ctx, &dto); err != nil {
		respBody.Comment = err.Error()
		return ctx.Status(fiber.StatusBadRequest).JSON(respBody)
	}

	// Prepare config to fetch currency rates
	apiConfig := currencyApi.APIConfig{Link: h.cfg.API.Link, Key: h.cfg.API.Key}
	// Fetch currency rate
	rateResp, err := currencyApi.FetchCurrency(apiConfig, dto.CurrencyFrom, dto.CurrencyTo)
	if err != nil {
		respBody.Comment = err.Error()
		return ctx.JSON(respBody)
	}

	// map dto to business entity
	params := models.NewCurrencyParams(dto.CurrencyFrom, dto.CurrencyTo)

	// create currency pairs
	err = h.uc.CreatePairs(h.ctx, *params, rateResp.Data[dto.CurrencyTo])
	if err != nil {
		respBody.Comment = err.Error()
		return ctx.JSON(respBody)
	}

	respBody.Success = true
	return ctx.JSON(respBody)
}

func (h handlers) Exchange(ctx *fiber.Ctx) error {
	var (
		dto      ExchangeDTO
		respBody ExchangeResponse
	)

	// Validate Request Body
	if err := utils.ReadRequest(ctx, &dto); err != nil {
		respBody.Comment = err.Error()
		return ctx.Status(fiber.StatusBadRequest).JSON(respBody)
	}

	// Prepare currencies to business layer
	params := models.NewCurrencyParams(dto.CurrencyFrom, dto.CurrencyTo)
	result, err := h.uc.Exchange(h.ctx, *params, dto.Amount)
	if err != nil {
		respBody.Comment = err.Error()
		return ctx.JSON(respBody)
	}

	respBody.Success = true
	respBody.Result = result
	return ctx.JSON(respBody)
}

func (h handlers) GetRate(ctx *fiber.Ctx) error {
	var (
		dto      CurrencyPairsDTO
		respBody GetRateResponse
	)

	// Validate Request Body
	if err := utils.ReadRequest(ctx, &dto); err != nil {
		respBody.Comment = err.Error()
		return ctx.Status(fiber.StatusBadRequest).JSON(respBody)
	}

	// Prepare currencies to business layer
	params := models.NewCurrencyParams(dto.CurrencyFrom, dto.CurrencyTo)
	rate, err := h.uc.GetRate(h.ctx, *params)
	if err != nil {
		respBody.Comment = err.Error()
		return ctx.JSON(respBody)
	}

	respBody.Rate = *rate
	respBody.Success = true
	return ctx.JSON(respBody)
}

func (h handlers) Aggregate(ctx *fiber.Ctx) error {
	var (
		dto      AggregateDTO
		respBody AggregateResponse
	)

	// Validate Request Body
	if err := utils.ReadRequest(ctx, &dto); err != nil {
		respBody.Comment = err.Error()
		return ctx.Status(fiber.StatusBadRequest).JSON(respBody)
	}

	pairs, err := h.uc.Aggregate(h.ctx, dto.Limit, dto.Offset)
	if err != nil {
		respBody.Comment = err.Error()
		return ctx.JSON(respBody)
	}

	respBody.Success = true
	respBody.Data = pairs
	return ctx.JSON(respBody)
}
