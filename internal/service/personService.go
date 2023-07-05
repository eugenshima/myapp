// Package service provides a set of functions, which include business-logic in it
package service

import (
	"context"
	"fmt"

	"github.com/eugenshima/myapp/internal/model"

	"github.com/google/uuid"
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
func (db *PersonServiceImpl) GetByID(ctx context.Context, id uuid.UUID) (*model.Person, error) {
	return db.rps.GetByID(ctx, id)
}

// GetAll is a service function which interacts with repository level
func (db *PersonServiceImpl) GetAll(ctx context.Context) ([]model.Person, error) {
	return db.rps.GetAll(ctx)
}

// Delete is a service function which interacts with repository level
func (db *PersonServiceImpl) Delete(ctx context.Context, uuidString uuid.UUID, accessToken string) error {
	id, role, err := GetPayloadFromToken(accessToken)
	fmt.Println(id)
	if err != nil {
		return err
	}
	if role != "admin" {
		return fmt.Errorf("invalid role: %v", err)
	}
	return db.rps.Delete(ctx, uuidString)
}

// Create is a service function which interacts with repository level
func (db *PersonServiceImpl) Create(ctx context.Context, entity *model.Person, accessToken string) error {
	id, role, err := GetPayloadFromToken(accessToken)
	fmt.Println(id)
	if err != nil {
		return err
	}
	if role != "admin" {
		return fmt.Errorf("invalid role: %v", err)
	}
	return db.rps.Create(ctx, entity)
}

// Update is a service function which interacts with repository level
func (db *PersonServiceImpl) Update(ctx context.Context, uuidString uuid.UUID, entity *model.Person, accessToken string) error {
	id, role, err := GetPayloadFromToken(accessToken)
	fmt.Println(id)
	if err != nil {
		return fmt.Errorf("error getting data from token")
	}
	if role != "admin" {
		str := fmt.Errorf("invalid role: need admin")
		return str
	}
	return db.rps.Update(ctx, uuidString, entity)
}
