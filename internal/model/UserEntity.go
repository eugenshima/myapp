package model

import "github.com/google/uuid"

type User struct {
	ID       uuid.UUID `db:"id"`
	login    string    `db:"login"`
	password string    `db:"password"`
}
