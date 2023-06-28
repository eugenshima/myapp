// Package repository provides functions for interacting with a database
// or other persistent storage system in a web service.
// It includes functions for creating, reading, updating, and deleting data from the storage system.
package repository

import (
	"context"

	"github.com/eugenshima/myapp/internal/model"
)

// GetAllUsers function executes a query, which select all rows from user table
func (db *PsqlConnection) GetAllUsers() ([]model.User, error) {
	rows, err := db.pool.Query(context.Background(), "SELECT id, login, password, role FROM goschema.user")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Create slice to store data from our SQL request
	var results []model.User

	// go;) through each line
	for rows.Next() {
		entity := model.User{}
		err := rows.Scan(&entity.ID, &entity.Login, &entity.Password, &entity.Role)
		if err != nil {
			return nil, err
		}
		results = append(results, entity)
	}
	return results, rows.Err()
}
