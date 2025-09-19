package app

import (
	"log/slog"
	"time"

	grpcapp "github.com/TimNikolaev/drag-sso/internal/app/grpc"
	pgapp "github.com/TimNikolaev/drag-sso/internal/app/postgres"
	"github.com/TimNikolaev/drag-sso/internal/repository/postgres"
	"github.com/TimNikolaev/drag-sso/internal/services/auth"
)

type App struct {
	PostgresDB *pgapp.App
	GRPCServer *grpcapp.App
}

func New(log *slog.Logger, grpcPort int, tokenTTL time.Duration, dsn string) *App {

	pgApp := pgapp.New(log, dsn)

	db, err := pgApp.Connect()
	if err != nil {
		panic(err)
	}

	//Init repository
	repository := postgres.New(db)

	//Init service
	authService := auth.New(log, tokenTTL, repository, repository, repository)

	//Init gRPC server
	grpcApp := grpcapp.New(log, grpcPort, authService)

	return &App{
		PostgresDB: pgApp,
		GRPCServer: grpcApp,
	}

}
