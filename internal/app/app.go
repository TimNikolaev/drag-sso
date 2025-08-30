package app

import (
	"log/slog"
	"time"

	grpcapp "github.com/TimNikolaev/drag-sso/internal/app/grpc"
)

type App struct {
	GRPCServer *grpcapp.App
}

func New(log *slog.Logger, grpcPort int, tokenTTL time.Duration) *App {
	//Init repository

	//Init service

	//Init gRPC server
	grpcApp := grpcapp.New(log, grpcPort)

	return &App{
		GRPCServer: grpcApp,
	}

}
