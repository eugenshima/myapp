// Package main - entry point to our program
package main

import (
	"context"
	"fmt"

	"github.com/eugenshima/myapp/internal/handlers"
	"github.com/eugenshima/myapp/internal/repository"
	"github.com/eugenshima/myapp/internal/service"
	"github.com/labstack/echo/v4"
)

// Main struct
type Main struct {
	handler *handlers.Handler
}

// NewMain constructor function
func NewMain(handler *handlers.Handler) *Main {
	return &Main{handler: handler}
}

// Main - entry point
func main() {
	e := echo.New()
	// Initializing the Database Connector
	pool, err := repository.NewDatabasePsqlConnection()
	if err != nil {
		fmt.Println(err)
	}

	dbpool := repository.NewPsqlConnection(pool)
	Psqlservice := service.NewService(dbpool)
	handlr := handlers.NewHandler(Psqlservice)

	client, err := repository.CreateMongoConnect()
	if err != nil {
		fmt.Println(err)
	}
	db := repository.NewMongoDBConnection(client)
	mdb := service.NewMongoService(db)
	data, err := mdb.FindPersons(context.Background())
	if err != nil {
		fmt.Printf("Error in main: %v", err)
	}
	fmt.Println(data)
	handlr.HTTPHandler(e)

	e.Logger.Fatal(e.Start(":1323"))
}
