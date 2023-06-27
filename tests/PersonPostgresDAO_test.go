package tests

import (
	"fmt"
	"testing"

	"github.com/eugenshima/myapp/internal/repository"
)

type Test struct {
	DB *repository.PsqlConnection
}

func NewTest(DB *repository.PsqlConnection) *Test {
	return &Test{DB: DB}
}

func TestGetAll(t *testing.T) {

	//arrange

	// Connect to PSQL database
	pool, err := repository.NewDatabasePsqlConnection()
	if err != nil {
		fmt.Println(err)
	}
	dbpool := repository.NewPsqlConnection(pool)
	handlr := NewTest(dbpool)

	//act

	// Make a database call to retrieve all entities
	result, err := handlr.DB.GetAll()

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
