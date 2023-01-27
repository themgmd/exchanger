package core

import (
	"context"
	"encoding/json"
	"github.com/gofiber/fiber/v2"
	"onemgvv/exchanger/internal/common/mappers"
	"onemgvv/exchanger/internal/config"
	"onemgvv/exchanger/internal/currency"
	currencyHttp "onemgvv/exchanger/internal/currency/delivery/http"
	"onemgvv/exchanger/internal/currency/repository"
	currencyUseCase "onemgvv/exchanger/internal/currency/usecase"
	"onemgvv/exchanger/internal/logger/zaplog"
	"onemgvv/exchanger/internal/scheduler"
	"onemgvv/exchanger/internal/server"
	"onemgvv/exchanger/pkg/database/inmemory"
	"onemgvv/exchanger/pkg/database/postgres"
	"os/signal"
	"syscall"
	"time"
)

type Core struct {
	cfg    *config.Config
	server *server.Server
}

func initPostgres(ctx context.Context, pgConfig config.PostgresConfig) postgres.Client {
	cfg := mappers.MapPostgresConfig(pgConfig)
	db, err := postgres.New(ctx, 5, 5*time.Second, cfg)
	if err != nil {
		zaplog.AppLogger.Fatalf("Error occured while connecting to DB: %v", err.Error())
	}

	if err = db.Ping(ctx); err != nil {
		zaplog.AppLogger.Fatalf("Error occured while ping DB: %v", err.Error())
	}

	return db
}

func initScheduler(ctx context.Context, cfg *config.Config, uc currency.UseCase) *scheduler.Scheduler {
	worker := scheduler.NewScheduler(cfg, uc)
	go worker.Start(ctx)
	return worker
}

func initApiV1(currHandler currency.Handlers) server.Handlers {
	return server.Handlers{CurrencyHandlers: currHandler}
}

func New(configPath string) *Core {
	cfg, err := config.New(configPath)
	if err != nil {
		zaplog.AppLogger.Fatalf("Error occured while init configuration")
	}

	fiberConfig := fiber.Config{
		ReadTimeout:  cfg.HTTP.Timeout.Read,
		WriteTimeout: cfg.HTTP.Timeout.Write,
		JSONEncoder:  json.Marshal,
		JSONDecoder:  json.Unmarshal,
	}

	app := fiber.New(fiberConfig)
	return &Core{
		cfg:    cfg,
		server: server.New(app),
	}
}

func (c *Core) Serve() error {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)
	defer cancel()

	pool := initPostgres(ctx, c.cfg.Database.Postgres)

	currencyRepo := repository.New(pool)
	inMemory := inmemory.New()
	currencyUC := currencyUseCase.New(currencyRepo, inMemory)
	currencyHandlers := currencyHttp.New(ctx, c.cfg, currencyUC)

	initScheduler(ctx, c.cfg, currencyUC)

	handlersV1 := initApiV1(currencyHandlers)
	c.server.MapHandlers(handlersV1)
	return c.server.Run(ctx, c.cfg.HTTP.Port)
}
