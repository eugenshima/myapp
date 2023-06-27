// Package repository provides functions for interacting with a database
// or other persistent storage system in a web service.
// It includes functions for creating, reading, updating, and deleting data from the storage system.
package repository

//TODO: wrap every error here

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/eugenshima/myapp/internal/model"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/labstack/echo/v4"
)

// PsqlConnection is a struct, which contains Pool variable
type PsqlConnection struct {
	pool *pgxpool.Pool
}

// NewPsqlConnection constructor for PsqlConnection
func NewPsqlConnection(pool *pgxpool.Pool) *PsqlConnection {
	return &PsqlConnection{pool: pool}
}

// NewDatabasePsqlConnection function provides Connection with PostgreSQL database
func NewDatabasePsqlConnection() (*pgxpool.Pool, error) {
	// Initialization a connect configuration for a PostgreSQL using pgx driver
	config, err := pgxpool.ParseConfig("postgres://eugen:ur2qly1ini@localhost:5432/eugen")
	if err != nil {
		return nil, err
	}

	// Establishing a new connection to a PostgreSQL database using the pgx driver
	pool, err := pgxpool.ConnectConfig(context.Background(), config)
	if err != nil {
		return nil, err
	}
	// Output to console
	fmt.Println("Connection to PostgreSQL successful")

	return pool, nil
}

// GetByName function executes SQL request to select all rows, where name=Name
func (db *PsqlConnection) GetByName(c echo.Context, Name string) (*model.Entity, error) {
	var entity model.Entity
	entity.Name = Name

	query := `SELECT id, name, age, ishealthy FROM goschema.person WHERE name=$1`

	// Execute a SQL query on a database
	err := db.pool.QueryRow(c.Request().Context(), query, &entity.Name).Scan(&entity.ID, &entity.Name, &entity.Age, &entity.IsHealthy)
	if err != nil {
		return nil, err // Returning error message
	}

	// Convert the result of query into JSON format
	jsonString, err := json.Marshal(entity)
	if err != nil {
		return nil, err
	}
	jsString := string(jsonString)
	fmt.Println("JSON --> ", jsString)
	return &entity, nil
}

// GetAll function executes SQL request to select all rows from Database
func (db *PsqlConnection) GetAll() ([]model.Entity, error) {
	rows, err := db.pool.Query(context.Background(), "SELECT id, name, age, ishealthy FROM goschema.person")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Create slice to store data from our SQL request
	var results []model.Entity

	// go;) through each line
	for rows.Next() {
		entity := model.Entity{}
		err := rows.Scan(&entity.ID, &entity.Name, &entity.Age, &entity.IsHealthy)
		if err != nil {
			return nil, err
		}
		results = append(results, entity)
	}
	fmt.Println(results)
	return results, rows.Err()
}

// Delete function executes SQL reauest to delete row with certain uuid
func (db *PsqlConnection) Delete(uuidString string) error {
	var entity model.Entity
	parsedUUID, err := uuid.Parse(uuidString)
	if err != nil {
		fmt.Println("Error parsiing(JSON)")
		return err
	}
	entity.ID = parsedUUID

	bd, err := db.pool.Exec(context.Background(), "DELETE FROM goschema.person WHERE id=$1", &entity.ID)
	if err != nil && !bd.Delete() {
		fmt.Println("Error deleting data into table:", err)
		return err
	}
	return nil
}

// Insert function executes SQL request to insert person into database
func (db *PsqlConnection) Insert(entity *model.Entity) error {
	entity.ID = uuid.New()

	bd, err := db.pool.Exec(context.Background(),
		`INSERT INTO goschema.person (id, name, age, ishealthy) 
	VALUES($1,$2,$3,$4)`,
		&entity.ID, &entity.Name, &entity.Age, &entity.IsHealthy)

	if err != nil && !bd.Insert() {
		fmt.Println("Error deleting data into table:", err)
		return err
	}
	return nil
}

// Update function executes SQL request to update person data in database
func (db *PsqlConnection) Update(uuidString string, entity *model.Entity) error {
	parsedUUID, err := uuid.Parse(uuidString)
	if err != nil {
		fmt.Println("Error parsiing(JSON)")
		return err
	}
	entity.ID = parsedUUID
	bd, err := db.pool.Exec(context.Background(), "UPDATE goschema.person SET name=$1, age=$2, ishealthy=$3 WHERE id=$4", &entity.Name, &entity.Age, &entity.IsHealthy, &entity.ID)
	if err != nil && !bd.Update() {
		fmt.Println("Error deleting data into table:", err)
		return err
	}
	return nil
}
