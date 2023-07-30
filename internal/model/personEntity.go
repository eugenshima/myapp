// Package model of our entity from database
package model

import (
	"github.com/google/uuid"
)

// Person struct for person entity in the database
type Person struct {
	ID        uuid.UUID `json:"id" bson:"_id"`
	Name      string    `json:"name" bson:"name" validate:"required"`
	Age       int64     `json:"age" bson:"age" validate:"required,min=0,max=140"`
	IsHealthy bool      `json:"ishealthy" bson:"is_healthy"`
}

// PersonRedis struct for Person entity in Redis database
type PersonRedis struct {
	Name      string `json:"name" bson:"name" validate:"required"`
	Age       int    `json:"age" bson:"age" validate:"required,min=0,max=140"`
	IsHealthy bool   `json:"ishealthy" bson:"is_healthy"`
}
