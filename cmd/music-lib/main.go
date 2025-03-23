package main

import (
	"context"
	"errors"
	"log/slog"
	"music-lib/internal/config"
	"music-lib/internal/http/router"
	"music-lib/internal/storage/pgsql"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
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
	log.Info("starting server", slog.Any("cfg", cfg.AppUrl))

	// define postgres
	storage := pgsql.New(cfg)

	// define router
	routes := router.New(storage, log)

	// run server
	server := http.Server{
		Addr:    cfg.AppUrl + ":" + cfg.AppPort,
		Handler: routes,
	}

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// Запускаем сервер в горутине
	go func() {
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Error("failed to start server", slog.Any("error", err))
		}
	}()

	log.Info("Server started", slog.String("addr", cfg.AppUrl+":"+cfg.AppPort))

	// Ожидаем сигнал для graceful shutdown
	sig := <-quit
	log.Info("shutting down server...", slog.Any("signal", sig))

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Error("server forced to shutdown", slog.Any("error", err))
	} else {
		log.Info("server gracefully stopped")
	}
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
