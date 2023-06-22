package repository

//here will be a model of our entity from database
import (
	"github.com/google/uuid"
)

type Entity struct {
	ID        uuid.UUID `db:"id"`
	Name      string    `db:"name"`
	Age       int       `db:"age"`
	IsHealthy bool      `db:"ishealthy"`
}

func NewEntity(name string, age int, isHealthy bool) *Entity {
	return &Entity{
		ID:        uuid.New(),
		Name:      name,
		Age:       age,
		IsHealthy: isHealthy,
	}
}
