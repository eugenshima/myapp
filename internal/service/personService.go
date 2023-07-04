// Package service provides a set of functions, which include business-logic in it
package service

import (
	"context"

	"github.com/eugenshima/myapp/internal/model"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

// PersonServiceImpl is a struct that contains a reference to the repository interface
type PersonServiceImpl struct {
	rps PersonRepositoryPsql
}

// NewPersonService is a constructor for the PersonServiceImpl struct
func NewPersonService(rps PersonRepositoryPsql) *PersonServiceImpl {
	return &PersonServiceImpl{rps: rps}
}

// PersonRepositoryPsql interface, which contains repository methods
type PersonRepositoryPsql interface {
	GetByID(ctx context.Context, id uuid.UUID) (*model.Person, error)
	GetAll(ctx context.Context) ([]model.Person, error)
	Delete(ctx context.Context, uuidString uuid.UUID) error
	Create(ctx context.Context, entity *model.Person) error
	Update(ctx context.Context, uuidString uuid.UUID, entity *model.Person) error
}

// GetByID is a service function which interacts with PostgreSQL in repository level
func (db *PersonServiceImpl) GetByID(c echo.Context, id uuid.UUID) (*model.Person, error) {
	return db.rps.GetByID(c.Request().Context(), id)
}

// GetAll is a service function which interacts with repository level
func (db *PersonServiceImpl) GetAll(c echo.Context) ([]model.Person, error) {
	return db.rps.GetAll(c.Request().Context())
}

// Delete is a service function which interacts with repository level
func (db *PersonServiceImpl) Delete(c echo.Context, uuidString uuid.UUID) error {
	return db.rps.Delete(c.Request().Context(), uuidString)
}

// Create is a service function which interacts with repository level
func (db *PersonServiceImpl) Create(c echo.Context, entity *model.Person) error {
	return db.rps.Create(c.Request().Context(), entity)
}

// Update is a service function which interacts with repository level
func (db *PersonServiceImpl) Update(c echo.Context, uuidString uuid.UUID, entity *model.Person) error {
	return db.rps.Update(c.Request().Context(), uuidString, entity)
}
