// Package service provides a set of functions, which include business-logic in it
package service

import (
	"context"
	"fmt"

	"github.com/eugenshima/myapp/internal/model"

	"github.com/google/uuid"
)

// PersonService is a struct that contains a reference to the repository interface
type PersonService struct {
	rps PersonRepositoryPsql
	rdb PersonRepositoryRedis
}

// NewPersonService is a constructor for the PersonServiceImpl struct
func NewPersonService(rps PersonRepositoryPsql, rdb PersonRepositoryRedis) *PersonService {
	return &PersonService{
		rps: rps,
		rdb: rdb,
	}
}

// PersonRepositoryPsql interface, which contains repository methods
type PersonRepositoryPsql interface {
	GetByID(ctx context.Context, id uuid.UUID) (*model.Person, error)
	GetAll(ctx context.Context) ([]*model.Person, error)
	Delete(ctx context.Context, uuidString uuid.UUID) error
	Create(ctx context.Context, entity *model.Person) (uuid.UUID, error)
	Update(ctx context.Context, uuidString uuid.UUID, entity *model.Person) error
}

// PersonRepositoryPsql interface, which contains repository methods
type PersonRepositoryRedis interface {
	RedisGetByID(ctx context.Context, id uuid.UUID) (*model.Person, error)
	RedisSetByID(ctx context.Context, entity *model.Person) error
	RedisDeleteByID(ctx context.Context, id uuid.UUID) error
}

// GetByID is a service function which interacts with PostgreSQL in repository level
func (db *PersonService) GetByID(ctx context.Context, id uuid.UUID) (*model.Person, error) {
	res, err := db.rdb.RedisGetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("RedisGetByID: %w", err)
	}
	if res != nil {
		return res, nil
	}
	res, err = db.rps.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("error getting person: %v", err)
	}
	return res, nil
}

// GetAll is a service function which interacts with repository level
func (db *PersonService) GetAll(ctx context.Context) ([]*model.Person, error) {
	return db.rps.GetAll(ctx)
}

// Delete is a service function which interacts with repository level
func (db *PersonService) Delete(ctx context.Context, uuidString uuid.UUID) error {
	db.rdb.RedisDeleteByID(ctx, uuidString)
	return db.rps.Delete(ctx, uuidString)
}

// Create is a service function which interacts with repository level
func (db *PersonService) Create(ctx context.Context, entity *model.Person) (uuid.UUID, error) {
	id, err := db.rps.Create(ctx, entity)
	db.rdb.RedisSetByID(ctx, entity)
	return id, err
}

// Update is a service function which interacts with repository level
func (db *PersonService) Update(ctx context.Context, id uuid.UUID, entity *model.Person) error {
	// Overwriting cache
	err := db.rdb.RedisDeleteByID(ctx, id)
	if err != nil {
		return fmt.Errorf("Update: %w", err)
	}
	err = db.rdb.RedisSetByID(ctx, entity)
	if err != nil {
		return fmt.Errorf("Update: %w", err)
	}
	return db.rps.Update(ctx, id, entity)
}
