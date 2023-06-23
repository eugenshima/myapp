package repository

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/eugenshima/myapp/internal/model"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

// Create a connection to a PostgreSQL database
func createConn() *pgx.Conn {

	//initialization a connect configuration for a PostgreSQL using pgx driver
	connConfig, err := pgx.ParseConfig("postgres://eugen:ur2qly1ini@localhost:5432/eugen")
	if err != nil {
		fmt.Printf("failed to create connection to PostgreSQL: %v", err)
		return nil
	}

	//establishing a new connection to a PostgreSQL database using the pgx driver
	conn, err := pgx.ConnectConfig(context.Background(), connConfig)
	if err != nil {
		fmt.Printf("failed to create connection to PostgreSQL: %v", err)
		return nil
	}

	//Output to console
	fmt.Println("Connection to PostgreSQL successful")
	return conn
}

// Add a person to database table
func CreatePerson() error {
	var entity model.Entity
	entity.ID = uuid.New()
	entity.Name = "eugen"
	entity.Age = 20
	entity.IsHealthy = true
	conn := createConn()

	bd, err := conn.Exec(context.Background(), "INSERT INTO goschema.newtable (id, name, age, ishealthy) VALUES ($1, $2, $3, $4)", entity.ID, entity.Name, entity.Age, entity.IsHealthy)
	if err != nil {
		fmt.Println("Error inserting data into table:", err)
		return err
	}
	fmt.Println(bd, " <-- result of the request")
	fmt.Println("Data successfully inserted into table yauhenishymanski.newtable")
	return nil
}

// Find person by name in database table
func GetByName(Name string) (string, error) {

	var entity model.Entity
	entity.Name = Name
	conn := createConn()

	query := `SELECT id, name, age, ishealthy FROM goschema.newtable WHERE name=$1`
	//execute a SQL query on a database
	err := conn.QueryRow(context.Background(), query, entity.Name).Scan(entity.ID, entity.Name, entity.Age, entity.IsHealthy)
	if err != nil {
		return "nil\n", err //returning error message
	}

	// Convert the result of query into JSON format
	jsonString, err := json.Marshal(entity)
	if err != nil {
		return "nil\n", err
	}
	jsString := string(jsonString)
	fmt.Println("JSON --> ", jsString)
	return jsString, nil
}

// Get all Rows from the table
func GetAll() ([]model.Entity, error) {

	conn := createConn()
	rows, err := conn.Query(context.Background(), "SELECT * FROM goschema.newtable")
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

//Delete row from table by name
func Delete(Name string) error {
	var entity model.Entity
	entity.Name = Name
	conn := createConn()
	bd, err := conn.Exec(context.Background(), "DELETE FROM goschema.newtable WHERE name=$1", entity.Name)
	if err != nil {
		fmt.Println("Error deleting data into table:", err)
		return err
	}
	fmt.Println(bd, " <-- result of the request")
	fmt.Println("Data successfully deleted from table yauhenishymanski.newtable")
	return nil
}

func Update(entity model.Entity) error {
	conn := createConn()
	bd, err := conn.Exec(context.Background(), "UPDATE goschema.newtable SET name=$1, age=$2, ishealthy=$3", entity.Name, entity.Age, entity.IsHealthy)
	if err != nil {
		fmt.Println("Error updating data into table:", err)
		return err
	}
	fmt.Println(bd, " <-- result of the request")
	fmt.Println("Data successfully updated from table yauhenishymanski.newtable")
	return nil
}
