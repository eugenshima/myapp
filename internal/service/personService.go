// Package service provides a set of functions, which include business-logic in it
package service

import (
	"github.com/eugenshima/myapp/internal/model"
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

// PersonService interface, which contains repository methods
type PersonRepositoryPsql interface {
	GetByName(c echo.Context, Name string) (*model.Person, error)
	GetAll(c echo.Context) ([]model.Person, error)
	Delete(c echo.Context, uuidString string) error
	Create(c echo.Context, entity *model.Person) error
	Update(c echo.Context, uuidString string, entity *model.Person) error
}

// GetByName is a service function which interacts with PostgreSQL in repository level
func (db *PersonServiceImpl) GetByName(c echo.Context, Name string) (*model.Person, error) {
	return db.rps.GetByName(c, Name)
}

// GetAll is a service function which interacts with repository level
func (db *PersonServiceImpl) GetAll(c echo.Context) ([]model.Person, error) {
	return db.rps.GetAll(c)
}

// Delete is a service function which interacts with repository level
func (db *PersonServiceImpl) Delete(c echo.Context, uuidString string) error {
	return db.rps.Delete(c, uuidString)
}

// Create is a service function which interacts with repository level
func (db *PersonServiceImpl) Create(c echo.Context, entity *model.Person) error {
	return db.rps.Create(c, entity)
}

// Update is a service function which interacts with repository level
func (db *PersonServiceImpl) Update(c echo.Context, uuidString string, entity *model.Person) error {
	return db.rps.Update(c, uuidString, entity)
}
