// Package model provides a struct for our User entity in database
package model

import "github.com/google/uuid"

// User struct user which represents database entity with the same name
type User struct {
	ID           uuid.UUID `db:"id"`
	Login        string    `db:"login" validate:"required"`
	Password     []byte    `db:"password" validate:"required"`
	Role         string    `db:"role"`
	RefreshToken []byte    `db:"refreshtoken"`
}

// Login struct for user
type Login struct {
	Login    string `db:"login" validate:"required"`
	Password string `db:"password" validate:"required"`
}

// Signup struct for user
type Signup struct {
	Login    string `json:"login"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

type Refresh struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}
