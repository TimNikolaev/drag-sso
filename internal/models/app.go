package models

type App struct {
	ID        int32  `db:"id"`
	Name      string `db:"name"`
	SecretKey string `db:"secret"`
}
