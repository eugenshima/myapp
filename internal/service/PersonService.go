package service

import (
	"github.com/eugenshima/myapp/internal/model"
	"github.com/eugenshima/myapp/internal/repository"
)

// Service is a struct that contains a reference to the repository interface
type Service struct {
	DB *repository.PsqlConnection
}

func NewService(DB *repository.PsqlConnection) *Service {
	return &Service{DB: DB}
}

type PersonService interface {
	GetByName(Name string) (interface{}, error)
	GetAll()
	Delete(uuidString string)
	Insert(entity *model.Entity)
	Update(uuidString string, entity *model.Entity)
}

func GetByName(Name string) {

}
