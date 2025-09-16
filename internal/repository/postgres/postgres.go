package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/TimNikolaev/drag-sso/internal/models"
	"github.com/TimNikolaev/drag-sso/internal/repository"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jmoiron/sqlx"
)

type Repository struct {
	db *sqlx.DB
}

func New(dsn string) (*Repository, error) {
	const op = "repository.postgres.New"

	db, err := Connect(dsn)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &Repository{
		db: db,
	}, nil
}

const (
	usersTable = "users"
	appsTable  = "apps"

	userExistsCode = "23505"
)

func (r *Repository) CreateUser(ctx context.Context, email string, passHash []byte) (uint64, error) {
	const op = "repository.postgres.CreateUser"

	var id uint64

	query := fmt.Sprintf("INSERT INTO %s (email, pass_hash) VALUES ($1, $2) RETURNING id", usersTable)

	if err := r.db.QueryRowContext(ctx, query, email, passHash).Scan(&id); err != nil {
		var pgErr *pgconn.PgError

		if errors.As(err, &pgErr) && pgErr.Code == userExistsCode {
			return 0, fmt.Errorf("%s: %w", op, repository.ErrUserExists)
		}

		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return id, nil

}

func (r *Repository) GetUser(ctx context.Context, email string) (*models.User, error) {
	const op = "repository.postgres.GetUser"

	panic("not implemented")
}

func (r *Repository) GetApp(ctx context.Context, appID int32) (*models.App, error) {
	const op = "repository.postgres.GetApp"

	panic("not implemented")
}

func MustConnect(dsn string) *sqlx.DB {
	db, err := Connect(dsn)
	if err != nil {
		panic(err)
	}

	return db
}

func Connect(dsn string) (*sqlx.DB, error) {

	db, err := sqlx.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}
