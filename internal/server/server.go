package server

import (
	"context"
	"github.com/gofiber/fiber/v2"
	"onemgvv/exchanger/internal/logger/zaplog"
)

type Server struct {
	fiber *fiber.App
}

func New(fiber *fiber.App) *Server {
	return &Server{fiber}
}

func (s Server) Run(ctx context.Context, port string) error {
	go func() {
		if err := s.fiber.Listen(":" + port); err != nil {
			zaplog.HttpLogger.Fatalf("[FIBER | LISTEN ERROR]: %+v", err)
		}
	}()

	<-ctx.Done()
	zaplog.HttpLogger.Info("http server is stopped")
	return s.Stop()
}

func (s Server) Stop() error {
	return s.fiber.Shutdown()
}
