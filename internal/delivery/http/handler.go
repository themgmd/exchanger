package delivery_http

import (
	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
	"github.com/onemgvv/exchanger/internal/config"
	"github.com/onemgvv/exchanger/internal/delivery/http/response"
	v1 "github.com/onemgvv/exchanger/internal/delivery/http/v1"
	"time"
)

type Handler struct {
	fiber *fiber.App
	db    *sqlx.DB
	uc    v1.UseCase
	cfg   *config.Config
}

func NewHandler(app *fiber.App, db *sqlx.DB, uc v1.UseCase, cfg *config.Config) *Handler {
	return &Handler{app, db, uc, cfg}
}

func (h Handler) InitRoutes() {
	api := h.fiber.Group("api")

	h.InitAPIV1(api)
}

func (h Handler) InitAPIV1(router fiber.Router) {
	handlerV1 := v1.NewApiV1Handler(h.uc, h.cfg)

	router.Get("/ping", h.pingPong)
	router.Post("/currency", handlerV1.CreatePairs)
	router.Put("/currency", handlerV1.Exchange)
	router.Get("/currency", handlerV1.Aggregate)
}

func (h Handler) pingPong(ctx *fiber.Ctx) error {
	return ctx.Status(fiber.StatusOK).JSON(
		response.PingPongResponse{Timestamp: time.Now().Format("2006-01-02 15:04:05"), Message: "pong"},
	)
}
