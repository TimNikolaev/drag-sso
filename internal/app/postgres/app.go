package pgapp

import (
	"fmt"
	"log/slog"

	"github.com/jmoiron/sqlx"
)

type App struct {
	log *slog.Logger
	dsn string
	db  *sqlx.DB
}

func New(log *slog.Logger, dsn string) *App {
	return &App{
		log: log,
		dsn: dsn,
	}
}

func (a *App) Connect() (*sqlx.DB, error) {
	const op = "pgapp.Connect"

	log := a.log.With(slog.String("op", op))

	log.Info("opening the pull connection to postgres")

	db, err := sqlx.Open("postgres", a.dsn)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	log.Info("pull connection is open")

	a.db = db

	return db, nil
}

func (a *App) Close() error {
	const op = "pgapp.Close"

	log := a.log.With(slog.String("op", op))

	log.Info("closing the pull connection of postgres")

	if err := a.db.Close(); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
