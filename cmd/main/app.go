package main

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"github.com/onemgvv/exchanger/internal/config"
	deliveryHttp "github.com/onemgvv/exchanger/internal/delivery/http"
	"github.com/onemgvv/exchanger/internal/domain/usecase/currency"
	currencyRepo "github.com/onemgvv/exchanger/internal/infrastructure/repository/currency"
	"github.com/onemgvv/exchanger/internal/infrastructure/scheduler"
	"github.com/onemgvv/exchanger/internal/server"
	"github.com/onemgvv/exchanger/pkg/database/postgres"
	"log"
	"os"
	"os/signal"
	"syscall"
)

const configDir = "configs"

func main() {
	mode := config.Mode(os.Args[1:][0])
	if mode != config.DEVELOPMENT && mode != config.PRODUCTION {
		log.Fatalf("[ARGS] || [ERROR]: mode '%s' is unknown", mode)
	}

	if err := godotenv.Load(fmt.Sprintf(".env.%s", mode)); err != nil {
		log.Fatalf("[ENV] || [LOAD ERRR]: %s", err.Error())
	}

	cfg, err := config.New(configDir, mode)
	if err != nil {
		log.Fatalf("[CONFIG] || [ERROR]: %v", err)
	}

	db, err := postgres.New(*cfg)
	if err != nil {
		log.Fatalf("[POSTGRES] || [ERROR]: %v", err)
	}

	fib := fiber.New()
	currencyRepo := currencyRepo.NewRepository(db)
	currencyUC := currency.NewUseCase(currencyRepo)

	handler := deliveryHttp.NewHandler(fib, db, currencyUC, cfg)
	handler.InitRoutes()

	app := server.NewServer(fib, cfg)
	worker := scheduler.NewScheduler(cfg, currencyUC)

	go func() {
		if err := app.Run(); err != nil {
			log.Fatalf("[START SERVER ERROR]: %s", err.Error())
		}
	}()

	go worker.Start()

	println("Server Started")

	/**
	 *	Graceful Shutdown
	 */
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	if err := fib.Shutdown(); err != nil {
		log.Fatalf("Shutdown error")
	}

	if err := db.Close(); err != nil {
		log.Fatalf("Database close conn error")
	}

	log.Println("Server stopped!")
}
