package model

import "github.com/google/uuid"

type User struct {
	ID       uuid.UUID `db:"id"`
	Login    string    `db:"login"`
	Password string    `db:"password"`
	Role     string    `db:"role"`
}
