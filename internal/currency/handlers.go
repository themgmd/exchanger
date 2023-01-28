package currency

import "github.com/gofiber/fiber/v2"

type Handlers interface {
	CreatePairs(ctx *fiber.Ctx) error
	Exchange(ctx *fiber.Ctx) error
	GetRate(ctx *fiber.Ctx) error
	Aggregate(ctx *fiber.Ctx) error
}
