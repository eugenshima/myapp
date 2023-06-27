// Package service provides a set of functions, which include business-logic in it
package service

import "github.com/eugenshima/myapp/internal/model"

// UserService interface, which contains repository methods
type UserService interface {
	GetAllUsers()
}

// GetAllUsers is a service function which interacts with repository level
func (db *Service) GetAllUsers() ([]model.User, error) {
	entity, err := db.DB.GetAllUsers()
	if err != nil {
		return nil, err
	}
	return entity, nil
}
