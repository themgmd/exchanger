package main

import (
	"context"
	"exchanger/internal/config"
	"exchanger/internal/currency"
	"exchanger/internal/scheduler"
	"exchanger/migrations"
	"exchanger/pkg/errors"
	"exchanger/pkg/healthcheck"
	"exchanger/pkg/postgre"
	"fmt"
	"github.com/go-chi/chi/v5"
	"golang.org/x/sync/errgroup"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func runApp() error {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGTERM)
	defer cancel()

	slog.Info("config initialization")
	conf := config.Get()

	slog.Info("postgre open connection")
	dbConn, err := postgre.NewConn(ctx,
		postgre.DSN(conf.Postgre.DSN()),
		postgre.MaxOpenConn(conf.Postgre.MaxOpenConn),
		postgre.MaxIdleConn(conf.Postgre.MaxIdleConn),
	)
	if err != nil {
		return errors.Wrap(err, "connect to database")
	}

	slog.Info("applying migrations")
	applied, err := migrations.Apply(conf.Postgre.DSN())
	if err != nil {
		return errors.Wrap(err, "migrations apply")
	}

	slog.Info("migrations applied", "applied count", applied)

	slog.Info("mux router initialization")
	router := chi.NewRouter()
	apiRouter := chi.NewRouter()

	slog.Info("currency initialization")
	service := currency.Setup(dbConn, apiRouter)
	router.Mount("/api", apiRouter)

	eg, _ := errgroup.WithContext(ctx)

	slog.Info("cron worker initialization")
	cron := scheduler.NewScheduler(service)
	eg.Go(func() error {
		return cron.Start(ctx)
	})

	slog.Info("server initialization")
	server := &http.Server{
		Addr:              conf.HTTP.Host,
		Handler:           router,
		ReadHeaderTimeout: conf.HTTP.ReadHeaderTimeout,
	}

	eg.Go(func() error {
		defer slog.Info("stop listen http server")

		slog.Info("init http server", "host", conf.HTTP.Host)
		healthcheck.Get().MarkAsUp()
		serr := server.ListenAndServe()
		if errors.Is(serr, http.ErrServerClosed) {
			return nil
		}

		return fmt.Errorf("group listed http server: %w", serr)
	})

	eg.Go(func() error {
		defer slog.Info("shutdown http server")

		<-ctx.Done()
		healthcheck.Get().MarkAsDown()
		serr := server.Shutdown(ctx)
		if err != nil {
			return fmt.Errorf("group shutdown http server: %w", serr)
		}
		return nil
	})

	if err = eg.Wait(); err != nil {
		return err
	}

	slog.Info("shutdown app")
	time.Sleep(time.Second * 1)

	return nil
}
