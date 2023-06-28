// Package model of our entity from database
package model

import (
	"github.com/google/uuid"
)

// Person for person entity in database
type Person struct {
	ID        uuid.UUID `db:"id" bson:"_id"`
	Name      string    `db:"name" bson:"name"`
	Age       int       `db:"age" bson:"age"`
	IsHealthy bool      `db:"ishealthy" bson:"ishealthy"`
}