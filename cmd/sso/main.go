package main

import (
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/TimNikolaev/drag-sso/internal/app"
	"github.com/TimNikolaev/drag-sso/internal/config"
	_ "github.com/lib/pq"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
	// Init config
	cfg := config.MustLoad()

	// Init logger
	log := setupLogger(cfg.Env)

	if cfg.Env == envLocal {
		log.Info("starting application", slog.Any("cfg", cfg))
	}

	// Init app
	application := app.New(log, cfg.GRPC.Port, cfg.TokenTTL, cfg.DSN)

	// Run app's gRPC-server
	go application.GRPCServer.MustRun()

	// Realization Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	signal := <-quit

	log.Info("graceful stopping", slog.String("signal", signal.String()))

	application.GRPCServer.Stop()

	time.Sleep(2 * time.Second)

	if err := application.PostgresDB.Close(); err != nil {
		log.Error("error occurred on db connection close", slog.String("err", err.Error()))
	} else {
		log.Info("db connection closed successfully")
	}

	log.Info("graceful stopped")

}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envLocal:
		log = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envDev:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envProd:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	}

	return log
}
