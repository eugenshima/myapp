package service

import (
	"github.com/eugenshima/myapp/internal/model"
	"github.com/eugenshima/myapp/internal/repository"
)

// Service is a struct that contains a reference to the repository interface
type Service struct {
	DB *repository.PsqlConnection
}

// Constructor for the Service struct
func NewService(DB *repository.PsqlConnection) *Service {
	return &Service{DB: DB}
}

// Interface, which contains repository methods
type PersonService interface {
	GetByName(Name string)
	GetAll()
	Delete(uuidString string)
	Insert(entity *model.Entity)
	Update(uuidString string, entity *model.Entity)
}

func (db *Service) GetByName(Name string) (*model.Entity, error) {
	entity, err := db.DB.GetByName(Name)
	if err != nil {
		return nil, err
	}
	return entity, nil

}

func (db *Service) GetAll() ([]model.Entity, error) {
	entity, err := db.DB.GetAll()
	if err != nil {
		return nil, err
	}
	return entity, nil
}

func (db *Service) Delete(uuidString string) error {
	err := db.DB.Delete(uuidString)
	if err != nil {
		return err
	}
	return nil
}

func (db *Service) Insert(entity *model.Entity) error {
	err := db.DB.Insert(entity)
	if err != nil {
		return err
	}
	return nil
}

func (db *Service) Update(uuidString string, entity *model.Entity) error {
	err := db.DB.Update(uuidString, entity)
	if err != nil {
		return err
	}
	return nil
}
