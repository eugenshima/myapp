// Package model provides a struct for our User entity in database
package model

import "github.com/google/uuid"

// User struct user which represents database entity with the same name
type User struct {
	ID       uuid.UUID `db:"id"`
	Login    string    `db:"login" validate:"required"`
	Password []byte    `db:"password" validate:"required"`
	Email    string    `db:"email"`
}
