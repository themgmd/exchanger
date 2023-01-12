package server

import (
	"github.com/gofiber/fiber/v2"
	"github.com/onemgvv/exchanger/internal/config"
	"log"
	"os"
	"os/signal"
	"syscall"
)

type Server struct {
	fiber *fiber.App
	cfg   *config.Config
}

func NewServer(app *fiber.App, cfg *config.Config) *Server {
	return &Server{fiber: app, cfg: cfg}
}

func (s *Server) Run() error {
	go func() {
		if err := s.fiber.Listen(":" + s.cfg.HTTP.Port); err != nil {
			log.Fatalf("while fiber listen occured: %s", err.Error())
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, os.Interrupt, syscall.SIGINT)
	<-quit

	return s.fiber.Shutdown()
}
