// Package main - entry point to our program
package main

import (
	"context"
	"fmt"

	"github.com/eugenshima/myapp/internal/handlers"
	"github.com/eugenshima/myapp/internal/repository"
	"github.com/eugenshima/myapp/internal/service"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// PersonServiceImpl is a struct that contains a reference to the repository interface
type MainImpl struct {
	hndl PersonHandler
}

// NewPsqlService is a constructor for the PersonServiceImpl struct
func NewPsqlService(hndl PersonHandler) *MainImpl {
	return &MainImpl{hndl: hndl}
}

// PersonService interface, which contains repository methods
type PersonHandler interface {
	GetByName(c echo.Context) error
	GetAll(c echo.Context) error
	Delete(c echo.Context) error
	Create(c echo.Context) error
	Update(c echo.Context) error
}

// NewMongo creates a connection to MongoDB server
func NewMongo() (*mongo.Client, error) {
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")

	// Connect to MongoDB
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		return nil, err
	}

	// Check the connection
	err = client.Ping(context.Background(), nil)
	if err != nil {
		return nil, err
	}
	fmt.Println("Connected to MongoDB!")

	return client, nil
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

// Main - entry point
func main() {
	e := echo.New()

	ch := 2
	//TODO: switch case between databases(mongo or postgreSQL)

	//Initializing the Database Connector (MongoDB)
	client, err := NewMongo()
	if err != nil {
		fmt.Println(err)
	}
	// Initializing the Database Connector (PostgreSQL)
	pool, err := NewDBPsql()
	if err != nil {
		fmt.Println(err)
	}
	var handlr *handlers.PersonHandlerImpl
	switch ch {
	case 2:

		//TODO: mongoDB
		rps := repository.NewMongoDBConnection(client)
		srv := service.NewPersonService(rps)
		handlr = handlers.NewPersonHandler(srv)

	default:

		rps := repository.NewPsqlConnection(pool)
		srv := service.NewPersonService(rps)
		handlr = handlers.NewPersonHandler(srv)
	}

	e.GET("/person/getById/:id", handlr.GetById)
	e.GET("/person/getAll", handlr.GetAll)
	e.DELETE("/person/delete/:id", handlr.Delete)
	e.POST("/person/insert", handlr.Create)
	e.PATCH("/person/update/:id", handlr.Update)

	e.Logger.Fatal(e.Start(":1323"))

}
