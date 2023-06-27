package repository

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/eugenshima/myapp/internal/model"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v4/pgxpool"
)

type PsqlConnection struct {
	pool *pgxpool.Pool
}

// Struct constructor
func NewPsqlConnection(pool *pgxpool.Pool) *PsqlConnection {
	return &PsqlConnection{pool: pool}
}

func NewDatabasePsqlConnection() (*pgxpool.Pool, error) {
	//initialization a connect configuration for a PostgreSQL using pgx driver
	config, err := pgxpool.ParseConfig("postgres://eugen:ur2qly1ini@localhost:5432/eugen")
	if err != nil {
		return nil, err
	}

	//establishing a new connection to a PostgreSQL database using the pgx driver
	pool, err := pgxpool.ConnectConfig(context.Background(), config)
	if err != nil {
		return nil, err
	}
	//Output to console
	fmt.Println("Connection to PostgreSQL successful")

	return pool, nil
}

func (db *PsqlConnection) GetByName(Name string) (*model.Entity, error) {
	var entity model.Entity
	entity.Name = Name

	query := `SELECT id, name, age, ishealthy FROM goschema.person WHERE name=$1`

	//execute a SQL query on a database
	err := db.pool.QueryRow(context.Background(), query, &entity.Name).Scan(&entity.ID, &entity.Name, &entity.Age, &entity.IsHealthy)
	if err != nil {
		return nil, err //returning error message
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

func (db *PsqlConnection) Insert(entity *model.Entity) error {
	entity.ID = uuid.New()
	bd, err := db.pool.Exec(context.Background(), "INSERT INTO goschema.person (id, name, age, ishealthy) VALUES($1,$2,$3,$4)", &entity.ID, &entity.Name, &entity.Age, &entity.IsHealthy)
	if err != nil && !bd.Insert() {
		fmt.Println("Error deleting data into table:", err)
		return err
	}
	return nil
}

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
