package service

import (
	"context"

	"github.com/eugenshima/myapp/internal/repository"
	"go.mongodb.org/mongo-driver/bson"
)

// MongoService struct with *repository.MongoDBConnection variable
type MongoService struct {
	MDB *repository.MongoDBConnection
}

// NewMongoService is a constructor for MongoService
func NewMongoService(MDB *repository.MongoDBConnection) *MongoService {
	return &MongoService{MDB: MDB}
}

// MongoPersonService interface which contains repository methods
type MongoPersonService interface {
	FindPersons()
}

// FindPersons function
func (db *MongoService) FindPersons(ctx context.Context) (bson.M, error) {
	entity, err := db.MDB.FindPersons(ctx)
	if err != nil {
		return nil, err
	}
	return entity, nil
}
