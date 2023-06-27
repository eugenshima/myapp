// Package service provides a set of functions, which include business-logic in it
package service

import (
	"github.com/eugenshima/myapp/internal/model"
	"github.com/eugenshima/myapp/internal/repository"
	"github.com/labstack/echo/v4"
)

// Service is a struct that contains a reference to the repository interface
type Service struct {
	DB *repository.PsqlConnection
}

// NewService is a constructor for the Service struct
func NewService(DB *repository.PsqlConnection) *Service {
	return &Service{DB: DB}
}

// PersonService interface, which contains repository methods
type PersonService interface {
	GetByName(Name string)
	GetAll()
	Delete(uuidString string)
	Insert(entity *model.Entity)
	Update(uuidString string, entity *model.Entity)
}

// GetByName is a service function which interacts with repository level
func (db *Service) GetByName(c echo.Context, Name string) (*model.Entity, error) {
	return db.DB.GetByName(c, Name)
}

// GetAll is a service function which interacts with repository level
func (db *Service) GetAll() ([]model.Entity, error) {
	return db.DB.GetAll()
}

// Delete is a service function which interacts with repository level
func (db *Service) Delete(uuidString string) error {
	return db.DB.Delete(uuidString)
}

// Insert is a service function which interacts with repository level
func (db *Service) Insert(entity *model.Entity) error {
	return db.DB.Insert(entity)
}

// Update is a service function which interacts with repository level
func (db *Service) Update(uuidString string, entity *model.Entity) error {
	return db.DB.Update(uuidString, entity)
}
