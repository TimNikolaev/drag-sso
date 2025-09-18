package grpcapp

import (
	"fmt"
	"log/slog"
	"net"

	authgrpc "github.com/TimNikolaev/drag-sso/internal/grpc/auth"
	"google.golang.org/grpc"
)

type App struct {
	log        *slog.Logger
	port       int
	gRPCServer *grpc.Server
}

func New(log *slog.Logger, port int, authService authgrpc.Auth) *App {
	gRPCServer := grpc.NewServer()

	authgrpc.Register(gRPCServer, authService)

	return &App{
		log:        log,
		port:       port,
		gRPCServer: gRPCServer,
	}
}

func (a *App) MustRun() {
	if err := a.Run(); err != nil {
		panic(err)
	}
}

func (a *App) Run() error {
	const op = "grpcapp.Run"

	log := a.log.With(slog.String("op", op), slog.Int("port", a.port))

	log.Info("starting gRPC server")

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", a.port))
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	log.Info("gRPC server is running", slog.String("addr", lis.Addr().Network()))

	if err := a.gRPCServer.Serve(lis); err != nil {
		return fmt.Errorf("%s: %d", op, err)
	}

	return nil
}

func (a *App) Stop() {
	const op = "grpcapp.Stop"

	log := a.log.With(slog.String("op", op))

	log.Info("stopping gRPC server", slog.Int("port", a.port))

	a.gRPCServer.GracefulStop()
}
