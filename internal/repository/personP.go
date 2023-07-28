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
	var person model.Person
	query := `SELECT id, name, age, is_healthy FROM goschema.person WHERE id=$1`

	// Execute a SQL query on a database
	err := db.pool.QueryRow(ctx, query, ID).Scan(&person.ID, &person.Name, &person.Age, &person.IsHealthy)
	if err != nil {
		return nil, fmt.Errorf("QueryRow(): %w", err)
	}
	return &person, nil
}

// GetAll function executes SQL request to select all rows from Database
func (db *PsqlConnection) GetAll(ctx context.Context) ([]*model.Person, error) {
	rows, err := db.pool.Query(ctx, "SELECT id, name, age, is_healthy FROM goschema.person")
	if err != nil {
		return nil, fmt.Errorf("Query(): %w", err)
	}
	defer rows.Close()

	// Create slice to store data from our SQL request
	var results []*model.Person

	// go;) through each line
	for rows.Next() {
		person := &model.Person{}
		err := rows.Scan(&person.ID, &person.Name, &person.Age, &person.IsHealthy)
		if err != nil {
			return nil, fmt.Errorf("Scan(): %w", err) // Returning error message
		}
		results = append(results, person)
	}
	fmt.Println(results)
	return results, rows.Err()
}

// Delete function executes SQL reauest to delete row with certain uuid
func (db *PsqlConnection) Delete(ctx context.Context, uuidString uuid.UUID) (uuid.UUID, error) {
	// Execute a SQL query on a database
	err := db.pool.QueryRow(ctx, `SELECT id FROM goschema.person WHERE id=$1`, uuidString).Scan(&uuidString)
	if err != nil {
		return uuid.Nil, fmt.Errorf("QueryRow(): %w", err)
	}
	bd, err := db.pool.Exec(ctx, "DELETE FROM goschema.person WHERE id=$1", uuidString)
	if err != nil && !bd.Delete() {
		return uuid.Nil, fmt.Errorf("Exec(): %w", err) // Returning error message
	}
	return uuidString, nil
}

// Create function executes SQL request to insert person into database
func (db *PsqlConnection) Create(ctx context.Context, entity *model.Person) (uuid.UUID, error) {
	entity.ID = uuid.New()

	bd, err := db.pool.Exec(ctx,
		`INSERT INTO goschema.person (id, name, age, is_healthy) 
	VALUES($1,$2,$3,$4)`,
		entity.ID, entity.Name, entity.Age, entity.IsHealthy)
	if err != nil && !bd.Insert() {
		return uuid.Nil, fmt.Errorf("Exec(): %w", err) // Returning error message
	}
	return entity.ID, nil
}

// Update function executes SQL request to update person data in database
func (db *PsqlConnection) Update(ctx context.Context, uuidString uuid.UUID, person *model.Person) (uuid.UUID, error) {
	// Execute a SQL query on a database
	err := db.pool.QueryRow(ctx, `SELECT id FROM goschema.person WHERE id=$1`, uuidString).Scan(&uuidString)
	if err != nil {
		return uuid.Nil, fmt.Errorf("QueryRow(): %w", err)
	}
	bd, err := db.pool.Exec(ctx, "UPDATE goschema.person SET name=$1, age=$2, is_healthy=$3 WHERE id=$4", person.Name, person.Age, person.IsHealthy, uuidString)
	if err != nil && !bd.Update() {
		return uuid.Nil, fmt.Errorf("Exec(): %w", err) // Returning error message
	}
	return uuidString, nil
}
