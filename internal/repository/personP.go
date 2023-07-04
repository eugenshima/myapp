// Package repository provides functions for interacting with a database
// or other persistent storage system in a web service.
// It includes functions for creating, reading, updating, and deleting data from the storage system.
package repository

import (
	"context"
	"fmt"

	"github.com/eugenshima/myapp/internal/model"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v4/pgxpool"
)

// PsqlConnection is a struct, which contains Pool variable
type PsqlConnection struct {
	pool *pgxpool.Pool
}

// NewPsqlConnection constructor for PsqlConnection
func NewPsqlConnection(pool *pgxpool.Pool) *PsqlConnection {
	return &PsqlConnection{pool: pool}
}

// GetByID function executes SQL request to select all rows, where id=Id
func (db *PsqlConnection) GetByID(ctx context.Context, ID uuid.UUID) (*model.Person, error) {
	var entity model.Person

	query := `SELECT id, name, age, ishealthy FROM goschema.person WHERE id=$1`

	// Execute a SQL query on a database
	err := db.pool.QueryRow(ctx, query, ID).Scan(&entity.ID, &entity.Name, &entity.Age, &entity.IsHealthy)
	if err != nil {
		return nil, fmt.Errorf("error in PersonP.go GetByname() QueryRow(): %v", err) // Returning error message
	}
	return &entity, nil
}

// GetAll function executes SQL request to select all rows from Database
func (db *PsqlConnection) GetAll(ctx context.Context) ([]model.Person, error) {
	rows, err := db.pool.Query(ctx, "SELECT id, name, age, ishealthy FROM goschema.person")
	if err != nil {
		return nil, fmt.Errorf("error in PersonP.go GetAll() Query(): %v", err) // Returning error message
	}
	defer rows.Close()

	// Create slice to store data from our SQL request
	var results []model.Person

	// go;) through each line
	for rows.Next() {
		entity := model.Person{}
		err := rows.Scan(&entity.ID, &entity.Name, &entity.Age, &entity.IsHealthy)
		if err != nil {
			return nil, fmt.Errorf("error in PersonP.go GetAll() rows.Next(): %v", err) // Returning error message
		}
		results = append(results, entity)
	}
	return results, rows.Err()
}

// Delete function executes SQL reauest to delete row with certain uuid
func (db *PsqlConnection) Delete(ctx context.Context, uuidString uuid.UUID) error {
	bd, err := db.pool.Exec(ctx, "DELETE FROM goschema.person WHERE id=$1", uuidString)
	if err != nil && !bd.Delete() {
		return fmt.Errorf("error deleting data from table: %v", err) // Returning error message
	}
	return nil
}

// Create function executes SQL request to insert person into database
func (db *PsqlConnection) Create(ctx context.Context, entity *model.Person) error {
	entity.ID = uuid.New()

	bd, err := db.pool.Exec(ctx,
		`INSERT INTO goschema.person (id, name, age, ishealthy) 
	VALUES($1,$2,$3,$4)`,
		entity.ID, entity.Name, entity.Age, entity.IsHealthy)

	if err != nil && !bd.Insert() {
		return fmt.Errorf("error deleting data into table: %v", err) // Returning error message
	}
	return nil
}

// Update function executes SQL request to update person data in database
func (db *PsqlConnection) Update(ctx context.Context, uuidString uuid.UUID, entity *model.Person) error {
	bd, err := db.pool.Exec(ctx, "UPDATE goschema.person SET name=$1, age=$2, ishealthy=$3 WHERE id=$4", entity.Name, entity.Age, entity.IsHealthy, uuidString)
	if err != nil && !bd.Update() {
		return fmt.Errorf("error updating data in table: %v", err) // Returning error message
	}
	return nil
}
