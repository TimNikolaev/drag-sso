package app

import (
	"log/slog"
	"time"

	grpcapp "github.com/TimNikolaev/drag-sso/internal/app/grpc"
	"github.com/TimNikolaev/drag-sso/internal/repository/postgres"
	"github.com/TimNikolaev/drag-sso/internal/services/auth"
)

type App struct {
	GRPCServer *grpcapp.App
}

func New(log *slog.Logger, grpcPort int, tokenTTL time.Duration, dsn string) *App {
	//Init repository
	repository, err := postgres.New(dsn)
	if err != nil {
		panic(err)
	}

	//Init service
	authService := auth.New(log, tokenTTL, repository, repository, repository)

	//Init gRPC server
	grpcApp := grpcapp.New(log, grpcPort, authService)

	return &App{
		GRPCServer: grpcApp,
	}

}
