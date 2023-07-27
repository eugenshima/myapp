// Package main - entry point to our program
package main

import (
	"context"
	"fmt"
	"net/http"
	"time"

	_ "github.com/eugenshima/myapp/docs"
	cfgrtn "github.com/eugenshima/myapp/internal/config"
	"github.com/eugenshima/myapp/internal/consumer"
	"github.com/eugenshima/myapp/internal/handlers"
	middlwr "github.com/eugenshima/myapp/internal/middleware"
	"github.com/eugenshima/myapp/internal/producer"
	"github.com/eugenshima/myapp/internal/repository"
	"github.com/eugenshima/myapp/internal/service"

	"github.com/go-playground/validator"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/labstack/echo/v4"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	swg "github.com/swaggo/echo-swagger"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// CustomValidator is a custom validator struct
type CustomValidator struct {
	validator *validator.Validate
}

// Validate func validates your model
func (cv *CustomValidator) Validate(i interface{}) error {
	if err := cv.validator.Struct(i); err != nil {
		logrus.Errorf("Validator: %v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("validator: %v", err))
	}
	return nil
}

const (
	pgx    = "pgx"
	mongod = "mongo"
)

// NewMongo creates a connection to MongoDB server
func NewMongo(env string) (*mongo.Client, error) {
	clientOptions := options.Client().ApplyURI(env)

	// Connect to MongoDB
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		return nil, fmt.Errorf("Connect(): %w", err)
	}

	// Check the connection
	err = client.Ping(context.Background(), nil)
	if err != nil {
		return nil, fmt.Errorf("Ping(): %w", err)
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
	fmt.Println("Connected to PostgreSQL!")

	return pool, nil
}

//NewDBRedis function provides Connection with Redis database
func NewDBRedis(env string) (*redis.Client, error) {
	opt, err := redis.ParseURL(env)
	if err != nil {
		return nil, fmt.Errorf("error parsing redis: %v", err)
	}

	fmt.Println("Connected to redis!")
	rdb := redis.NewClient(opt)
	return rdb, nil
}

// @title Golang Web Service
// @version 1.0
// @description This is my golang server.

// @host localhost:8080
// @BasePath /

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
func main() {
	e := echo.New()
	e.Validator = &CustomValidator{validator: validator.New()}
	cfg, err := cfgrtn.NewConfig()
	if err != nil {
		fmt.Printf("Error extracting env variables: %v", err)
		return
	}

	ch := mongod
	// Initializing the Database Connector (MongoDB)
	client, err := NewMongo(cfg.MongoDBAddr)
	if err != nil {
		err := echo.NewHTTPError(http.StatusBadRequest, fmt.Errorf("error creating database connection with MongoDB: %w", err))
		e.Logger.Fatal(err)
	}
	// Initializing the Database Connector (PostgreSQL)
	pool, err := NewDBPsql(cfg.PgxDBAddr)
	if err != nil {
		err := echo.NewHTTPError(http.StatusBadRequest, fmt.Errorf("error creating database connection with PostgreSQL: %w", err))
		e.Logger.Fatal(err)
	}
	rdbClient, err := NewDBRedis(cfg.RedisDBAddr)
	if err != nil {
		err := echo.NewHTTPError(http.StatusBadRequest, fmt.Errorf("error creating database connection with Redis: %w", err))
		e.Logger.Fatal(err)
	}

	var handlr *handlers.PersonHandler
	var uhandlr *handlers.UserHandler
	switch ch {
	case mongod:

		// Person db mongodb
		rps := repository.NewMongoDBConnection(client)
		rdb := repository.NewRedisConnection(rdbClient)
		srv := service.NewPersonService(rps, rdb)
		handlr = handlers.NewPersonHandler(srv, validator.New())

		// User db mongodb
		urps := repository.NewUserMongoDBConnection(client)
		urdb := repository.NewUserRedisConnection(rdbClient)
		usrv := service.NewUserServiceImpl(urps, urdb)
		uhandlr = handlers.NewUserHandler(usrv, validator.New())

	case pgx:
		// Person db pgx
		rps := repository.NewPsqlConnection(pool)
		rdb := repository.NewRedisConnection(rdbClient)
		srv := service.NewPersonService(rps, rdb)
		handlr = handlers.NewPersonHandler(srv, validator.New())

		// User db pgx
		urps := repository.NewUserPsqlConnection(pool)
		urdb := repository.NewUserRedisConnection(rdbClient)
		usrv := service.NewUserServiceImpl(urps, urdb)
		uhandlr = handlers.NewUserHandler(usrv, validator.New())
	}

	api := e.Group("/api")
	{
		// Person Api
		person := api.Group("/person")
		person.POST("/insert", handlr.Create, middlwr.AdminIdentity())
		person.GET("/getAll", handlr.GetAll, middlwr.UserIdentity())
		person.GET("/getById/:id", handlr.GetByID, middlwr.UserIdentity())
		person.PATCH("/update/:id", handlr.Update, middlwr.AdminIdentity())
		person.DELETE("/delete/:id", handlr.Delete, middlwr.AdminIdentity())

		// User Api
		user := api.Group("/user")
		user.POST("/login", uhandlr.Login)
		user.POST("/signup", uhandlr.Signup)
		user.GET("/getAll", uhandlr.GetAll)
		user.POST("/refresh/:id", uhandlr.RefreshTokenPair)
		user.DELETE("/delete/:id", uhandlr.Delete)

		// Image requests
		image := api.Group("/image")
		image.Use(middlwr.AdminIdentity())
		image.GET("/get/:name", uhandlr.GetImage)
		image.POST("/set", uhandlr.SetImage)
	}
	e.GET("/swagger/*", swg.WrapHandler)

	redisProd := producer.NewProducer(rdbClient)
	redisCons := consumer.NewConsumer(rdbClient)
	ctx, cancel := context.WithTimeout(context.Background(), 6*time.Second)
	defer cancel()
	// Redis Stream
	go redisProd.RedisProducer(ctx)
	go redisCons.RedisConsumer(ctx)

	e.Logger.Fatal(e.Start(cfg.HTTPAddr))
}
