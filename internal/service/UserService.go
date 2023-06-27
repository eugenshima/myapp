package service

import "github.com/eugenshima/myapp/internal/model"

type UserService interface {
	GetAllUsers()
}

func (db *Service) GetAllUsers() ([]model.User, error) {
	entity, err := db.DB.GetAllUsers()
	if err != nil {
		return nil, err
	}
	return entity, nil
}
