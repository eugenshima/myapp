package repository_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/eugenshima/myapp/internal/repository"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/labstack/echo/v4"
)

type Test struct {
	DB *repository.PsqlConnection
}

func NewTest(DB *repository.PsqlConnection) *Test {
	return &Test{DB: DB}
}

// NewDBPsql function provides Connection with PostgreSQL database
func NewDBPsql() (*pgxpool.Pool, error) {
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

func TestGetAll(t *testing.T) {

	//arrange

	// Connect to PSQL database
	pool, err := NewDBPsql()
	if err != nil {
		fmt.Println(err)
	}
	dbpool := repository.NewPsqlConnection(pool)
	handlr := NewTest(dbpool)

	//act

	// Make a database call to retrieve all entities
	result, err := handlr.DB.GetAll(echo.New().AcquireContext())

	//result

	if err != nil {
		t.Errorf("Error occurred calling GetAll(): %v", err)
	}

	// Ensure that there are no errors
	if result == nil {
		t.Errorf("Error result is empty")
	}

	// Ensure that number of results returned is greater than 0
	if len(result) == 0 {
		t.Error("Number of results returned is 0")
	}

	// Ensure that each result has the expected fields and values
	for _, entity := range result {
		// Ensure that the ID field is uuid type

		// TODO: write a type check

		// Ensure that the Name field is not empty
		if entity.Name == "" {
			t.Error("Name field is empty")
		}

		// Ensure that the Age field is not negative
		if entity.Age < 0 {
			t.Error("Age field is negative")
		}
	}
}
