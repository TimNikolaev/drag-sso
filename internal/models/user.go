package models

type User struct {
	ID       uint64 `db:"id"`
	Email    string `db:"email"`
	PassHash []byte `db:"pass_hash"`
}
