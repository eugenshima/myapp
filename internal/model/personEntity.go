// Package model of our entity from database
package model

import (
	"github.com/google/uuid"
)

// Person struct for person entity in the database
type Person struct {
	ID        uuid.UUID `db:"id" bson:"_id"`
	Name      string    `db:"name" bson:"name" validate:"required"`
	Age       int       `db:"age" bson:"age" validate:"required,min=0,max=140"`
	IsHealthy bool      `db:"ishealthy" bson:"ishealthy" validate:"required"`
}
