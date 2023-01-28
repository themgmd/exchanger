package utils

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

// ReadRequest read request body and validate
func ReadRequest(ctx *fiber.Ctx, request interface{}) error {
	if err := ctx.BodyParser(&request); err != nil {
		return err
	}

	validate := validator.New()
	return validate.Struct(request)
}
