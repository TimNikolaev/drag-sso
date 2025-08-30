package main

import (
	"log/slog"
	"os"

	"github.com/TimNikolaev/drag-sso/internal/app"
	"github.com/TimNikolaev/drag-sso/internal/config"
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

	log.Info("starting application", slog.Any("cfg", cfg))

	// Init app
	application := app.New(log, cfg.GRPC.Port, cfg.TokenTTL)

	// Run app's gRPC-server
	application.GRPCServer.MustRun()
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
