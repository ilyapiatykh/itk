package app

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ilyapiatykh/itk/config"
	"github.com/ilyapiatykh/itk/internal/api"
	"github.com/ilyapiatykh/itk/internal/repo"
	"github.com/ilyapiatykh/itk/internal/service"
	"github.com/ilyapiatykh/itk/pkg/logging"
	_ "github.com/jackc/pgx/v5/stdlib"
)

func Run(cfg *config.Config) {
	URL := fmt.Sprintf(
		"postgres://%s:%s@postgres:5432/%s",
		cfg.User,
		cfg.Password,
		cfg.DBName,
	)

	if err := migrateDB(URL); err != nil {
		logging.Fatal("Failed to migrate db", slog.Any("error", err))
	}

	db, err := sql.Open("pgx", URL)
	if err != nil {
		logging.Fatal("Failed to open db conn", slog.Any("error", err))
	}
	// To reuse conns and not create new ones without PgBouncer
	db.SetMaxIdleConns(10)
	db.SetMaxOpenConns(10)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		logging.Fatal("Failed to ping db", slog.Any("error", err))
	}

	var (
		repo    = repo.NewWallets(db)
		service = service.NewWallets(repo)
		router  = api.NewRouter(&cfg.Server, service)
	)

	var (
		osSignals = make(chan os.Signal, 1)
		errCh     = make(chan error, 1)
	)
	signal.Notify(osSignals, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		err = router.Start()
		if err != nil {
			errCh <- err
		}
	}()

	select {
	case <-osSignals:
		stopService(router)

		slog.Info("Service was stopped")
	case err := <-errCh:
		slog.Error("Service failed", slog.Any("error", err))

		stopService(router)

		slog.Info("Service was stopped")
		os.Exit(1)
	}
}

func stopService(r *api.Router) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	if err := r.Stop(ctx); err != nil {
		slog.Error("Failed to gracefully stop service", slog.Any("error", err))
	}
}
