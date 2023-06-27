package model

// Here will be a model of our entity from database
import (
	"github.com/google/uuid"
)

// Entity for person entity in database
type Entity struct {
	ID        uuid.UUID `db:"id"`
	Name      string    `db:"name"`
	Age       int       `db:"age"`
	IsHealthy bool      `db:"ishealthy"`
}
