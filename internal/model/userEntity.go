// Package model provides a struct for our User entity in database
package model

import "github.com/google/uuid"

// User struct user which represents database entity with the same name
type User struct {
	ID           uuid.UUID `db:"id" bson:"_id"`
	Login        string    `db:"login" bson:"login" validate:"required"`
	Password     []byte    `db:"password" bson:"password" validate:"required"`
	Role         string    `db:"role" bson:"role"`
	RefreshToken []byte    `db:"refreshtoken" bson:"refreshtoken"`
}

// Login struct for user
type Login struct {
	Login    string `db:"login" bson:"login" validate:"required"`
	Password string `db:"password" bson:"password" validate:"required"`
}

// Signup struct for user
type Signup struct {
	Login    string `db:"login" bson:"login" validate:"required"`
	Password string `db:"password" bson:"password" validate:"required"`
	Role     string `json:"role" bson:"role" validate:"required"`
}

// GetUser struct for user
type GetUser struct {
	ID       uuid.UUID `db:"id" bson:"_id"`
	Login    string    `db:"login" bson:"login" validate:"required"`
	Password []byte    `db:"password" bson:"password" validate:"required"`
	Role     string    `json:"role" bson:"role" validate:"required"`
}

// Tokens struct for tokens
type Tokens struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}
