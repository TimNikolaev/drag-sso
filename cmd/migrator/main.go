package main

import (
	"errors"
	"flag"
	"fmt"

	"github.com/TimNikolaev/drag-sso/internal/config"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func main() {
	var migrationsPath, migrationsDB string

	flag.StringVar(&migrationsPath, "migrations-path", "", "path to migrations files")
	flag.StringVar(&migrationsDB, "migrations-db", "", "db name for migrations")

	cfg := config.MustLoad()

	if migrationsPath == "" {
		panic("migrations-path is required")
	}

	if migrationsDB == "" {
		panic("migrations-db is required")
	}

	db := MustDBConnect(cfg.DSN)
	defer db.Close()

	driver, err := postgres.WithInstance(db.DB, &postgres.Config{})
	if err != nil {
		panic(err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://"+migrationsPath,
		migrationsDB,
		driver,
	)
	if err != nil {
		panic(err)
	}

	if err := m.Up(); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			fmt.Println("no migrations to apply")

			return
		}

		panic(err)
	}

	fmt.Println("migrations applied successfully")

}

func MustDBConnect(dsn string) *sqlx.DB {
	db, err := DBConnect(dsn)
	if err != nil {
		panic(err)
	}

	return db
}

func DBConnect(dsn string) (*sqlx.DB, error) {

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
