package main

import (
	"flag"

	"github.com/TimNikolaev/drag-sso/internal/config"
	postgres_repo "github.com/TimNikolaev/drag-sso/internal/repository/postgres"
	"github.com/golang-migrate/migrate/v4"
  "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

func main() {
	var migrationsPath, migrationsDB string

	flag.StringVar(&migrationsPath, "migrations-path", "", "path to migrations files")
	flag.StringVar(&migrationsDB, "migrations-db", "", "db name for migrations")
	flag.Parse()

	if migrationsPath == "" {
		panic("migrations-path is required")
	}

	if migrationsDB == "" {
		panic("migrations-db is required")
	}

	cfg := config.MustLoad()

	db := postgres_repo.MustConnect(cfg.DSN)

	driver, err := postgres.WithInstance(db, &postgres.Config{})
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

		if err := m.Up(); err != nil{
			if errors.Is(err, migrate.ErrNotChange){
				fmt.Println("no migrations to apply")

				return 
			}

			panic(err)
		}

		fmt.Println("migrations applied successfully")

}
