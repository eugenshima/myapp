// Package service provides a set of functions, which include business-logic in it
package service

import (
	"context"

	"github.com/eugenshima/myapp/internal/model"

	"github.com/google/uuid"
)

// PersonServiceImpl is a struct that contains a reference to the repository interface
type PersonServiceImpl struct {
	rps PersonRepositoryPsql
	rdb PersonRepositoryRedis
}

// NewPersonService is a constructor for the PersonServiceImpl struct
func NewPersonService(rps PersonRepositoryPsql, rdb PersonRepositoryRedis) *PersonServiceImpl {
	return &PersonServiceImpl{rps: rps, rdb: rdb}
}

// PersonRepositoryPsql interface, which contains repository methods
type PersonRepositoryPsql interface {
	GetByID(ctx context.Context, id uuid.UUID) (*model.Person, error)
	GetAll(ctx context.Context) ([]model.Person, error)
	Delete(ctx context.Context, uuidString uuid.UUID) error
	Create(ctx context.Context, entity *model.Person) (uuid.UUID, error)
	Update(ctx context.Context, uuidString uuid.UUID, entity *model.Person) error
}

// PersonRepositoryPsql interface, which contains repository methods
type PersonRepositoryRedis interface {
	RedisGetByID(ctx context.Context, id uuid.UUID) *model.Person
	RedisSetByID(ctx context.Context, entity *model.Person) error
}

// GetByID is a service function which interacts with PostgreSQL in repository level
func (db *PersonServiceImpl) GetByID(ctx context.Context, id uuid.UUID) (*model.Person, error) {
	res := db.rdb.RedisGetByID(ctx, id)
	if res != nil {
		return res, nil
	}
	return db.rps.GetByID(ctx, id)
}

// GetAll is a service function which interacts with repository level
func (db *PersonServiceImpl) GetAll(ctx context.Context) ([]model.Person, error) {
	return db.rps.GetAll(ctx)
}

// Delete is a service function which interacts with repository level
func (db *PersonServiceImpl) Delete(ctx context.Context, uuidString uuid.UUID) error {
	return db.rps.Delete(ctx, uuidString)
}

// Create is a service function which interacts with repository level
func (db *PersonServiceImpl) Create(ctx context.Context, entity *model.Person) (uuid.UUID, error) {
	id, err := db.rps.Create(ctx, entity)
	db.rdb.RedisSetByID(ctx, entity)
	return id, err
}

// Update is a service function which interacts with repository level
func (db *PersonServiceImpl) Update(ctx context.Context, uuidString uuid.UUID, entity *model.Person) error {
	return db.rps.Update(ctx, uuidString, entity)
}
