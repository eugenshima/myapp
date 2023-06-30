// Package main - entry point to our program
package main

import (
	"context"
	"fmt"

	"github.com/caarlos0/env/v9"
	"github.com/eugenshima/myapp/internal/config"
	"github.com/eugenshima/myapp/internal/handlers"
	"github.com/eugenshima/myapp/internal/repository"
	"github.com/eugenshima/myapp/internal/service"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// NewMongo creates a connection to MongoDB server
func NewMongo(env string) (*mongo.Client, error) {

	clientOptions := options.Client().ApplyURI(env)

	// Connect to MongoDB
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		return nil, fmt.Errorf("error connecting to MongoDB: %v", err)
	}

	// Check the connection
	err = client.Ping(context.Background(), nil)
	if err != nil {
		return nil, fmt.Errorf("error connecting to MongoDB: %v", err)
	}
	fmt.Println("Connected to MongoDB!")

	return client, nil
}

// NewDBPsql function provides Connection with PostgreSQL database
func NewDBPsql(env string) (*pgxpool.Pool, error) {
	// Initialization a connect configuration for a PostgreSQL using pgx driver
	config, err := pgxpool.ParseConfig(env)
	if err != nil {
		return nil, fmt.Errorf("error connection to PostgreSQL: %v", err)
	}

	// Establishing a new connection to a PostgreSQL database using the pgx driver
	pool, err := pgxpool.ConnectConfig(context.Background(), config)
	if err != nil {
		return nil, fmt.Errorf("error connection to PostgreSQL: %v", err)
	}
	// Output to console
	fmt.Println("Connection to PostgreSQL successful")

	return pool, nil
}

// Main - entry point
func main() {
	e := echo.New()

	cfg := config.Config{}

	// Используйте env.Parse() для разбора существующих переменных окружения
	err := env.Parse(&cfg)
	if err != nil {
		fmt.Println("Ошибка при разборе переменных окружения:", err)
		return
	}

	// Доступ к системной переменной окружения

	ch := 1

	// Initializing the Database Connector (MongoDB)
	client, err := NewMongo(cfg.MongoDBADDR)
	if err != nil {
		fmt.Printf("Error creating database connection with PostgreSQL: %v", err)
	}
	// Initializing the Database Connector (PostgreSQL)
	pool, err := NewDBPsql(cfg.PgxDBAddr)
	if err != nil {
		fmt.Printf("Error creating database connection with PostgreSQL: %v", err)
	}
	var handlr *handlers.PersonHandlerImpl
	var uhandlr *handlers.UserHandlerImpl
	switch ch {
	case 2:

		// Person db mongodb
		rps := repository.NewMongoDBConnection(client)
		srv := service.NewPersonService(rps)
		handlr = handlers.NewPersonHandler(srv)

	default:
		// Person db pgx
		rps := repository.NewPsqlConnection(pool)
		srv := service.NewPersonService(rps)
		handlr = handlers.NewPersonHandler(srv)

		// User db pgx
		urps := repository.NewUserPsqlConnection(pool)
		usrv := service.NewUserServiceImpl(urps)
		uhandlr = handlers.NewUserHandlerImpl(usrv)
	}

	// Person requests
	e.GET("/person/getById/:id", handlr.GetByID)
	e.GET("/person/getAll", handlr.GetAll)
	e.DELETE("/person/delete/:id", handlr.Delete)
	e.POST("/person/insert", handlr.Create)
	e.PATCH("/person/update/:id", handlr.Update)

	// User requests
	e.GET("/user/login", uhandlr.Login)
	e.POST("/user/signup", uhandlr.Signup)
	e.GET("/user/getAll", uhandlr.GetAll)

	e.Logger.Fatal(e.Start(":1323"))
}
