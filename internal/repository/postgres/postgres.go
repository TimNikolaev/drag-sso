package postgres

import "github.com/jmoiron/sqlx"

func MustConnect(dsn string) *sqlx.DB {
	db, err := New(dsn)
	if err != nil {
		panic(err)
	}

	return db
}

func New(dsn string) (*sqlx.DB, error) {

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
