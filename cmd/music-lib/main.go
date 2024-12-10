package main

import (
	"log/slog"
	"music-lib/internal/config"
	"music-lib/internal/storage/pgsql"
	"os"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
	// define config
	cfg := config.MustLoad()

	// define logger
	log := setupLogger(cfg.AppEnv)
	log.Info("starting server", slog.Any("cfg", cfg))

	// define postgres
	storage := pgsql.New(cfg)
	_ = storage

	// TODO: define router
	// TODO: run server
	// TODO: Graceful shutdown
}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envLocal:
		log = slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envDev:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envProd:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	}

	return log
}
