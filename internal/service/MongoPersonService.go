package service

import (
	"context"

	"github.com/eugenshima/myapp/internal/repository"
	"go.mongodb.org/mongo-driver/bson"
)

type MongoService struct {
	MDB *repository.MongoDbConnection
}

func NewMongoService(MDB *repository.MongoDbConnection) *MongoService {
	return &MongoService{MDB: MDB}
}

type MongoPersonService interface {
	FindPersons()
}

func (db *MongoService) FindPersons(ctx context.Context) (bson.M, error) {
	entity, err := db.MDB.FindPersons(ctx)
	if err != nil {
		return nil, err
	}
	return entity, nil
}
