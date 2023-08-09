// Package main - entry point to our program
package main

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"time"

	_ "github.com/eugenshima/myapp/docs"
	cfgrtn "github.com/eugenshima/myapp/internal/config"
	"github.com/eugenshima/myapp/internal/consumer"
	"github.com/eugenshima/myapp/internal/handlers"
	"github.com/eugenshima/myapp/internal/interceptors"
	"github.com/eugenshima/myapp/internal/producer"
	"github.com/eugenshima/myapp/internal/repository"
	"github.com/eugenshima/myapp/internal/service"
	"google.golang.org/grpc"

	proto "github.com/eugenshima/myapp/proto_services"
	"github.com/go-playground/validator"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/labstack/echo/v4"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
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
	timeOut = 6
	pgx     = "pgx"
	mongod  = "mongo"
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

// NewDBRedis function provides Connection with Redis database
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
//
// nolint:funlen // required
func main() {
	cfg, err := cfgrtn.NewConfig()
	if err != nil {
		fmt.Printf("Error extracting env variables: %v", err)
		return
	}

	ch := pgx
	// Initializing the Database Connector (MongoDB)
	client, err := NewMongo(cfg.MongoDBAddr)
	if err != nil {
		err := echo.NewHTTPError(http.StatusBadRequest, fmt.Errorf("error creating database connection with MongoDB: %w", err))
		logrus.Fatal(err)
	}
	// Initializing the Database Connector (PostgreSQL)
	pool, err := NewDBPsql(cfg.PgxDBAddr)
	if err != nil {
		err := echo.NewHTTPError(http.StatusBadRequest, fmt.Errorf("error creating database connection with PostgreSQL: %w", err))
		logrus.Fatal(err)
	}
	rdbClient, err := NewDBRedis(cfg.RedisDBAddr)
	if err != nil {
		err := echo.NewHTTPError(http.StatusBadRequest, fmt.Errorf("error creating database connection with Redis: %w", err))
		logrus.Fatal(err)
	}
	var handlr *handlers.GRPCPersonHandler
	var uhandlr *handlers.GRPCUserHandler
	switch ch {
	case mongod:

		// Person db mongodb
		rps := repository.NewMongoDBConnection(client)
		rdb := repository.NewRedisConnection(rdbClient)
		srv := service.NewPersonService(rps, rdb)
		handlr = handlers.NewGRPCPersonHandler(srv)

		// User db mongodb
		urps := repository.NewUserMongoDBConnection(client)
		urdb := repository.NewUserRedisConnection(rdbClient)
		usrv := service.NewUserServiceImpl(urps, urdb)
		uhandlr = handlers.NewGRPCUserHandler(usrv)

	case pgx:
		// Person db pgx
		rps := repository.NewPsqlConnection(pool)
		rdb := repository.NewRedisConnection(rdbClient)
		srv := service.NewPersonService(rps, rdb)
		handlr = handlers.NewGRPCPersonHandler(srv)

		// User db pgx
		urps := repository.NewUserPsqlConnection(pool)
		urdb := repository.NewUserRedisConnection(rdbClient)
		usrv := service.NewUserServiceImpl(urps, urdb)
		uhandlr = handlers.NewGRPCUserHandler(usrv)
	}
	lis, err := net.Listen("tcp", "127.0.0.1:8080")
	if err != nil {
		logrus.Fatalf("cannot create listener: %s", err)
	}
	serverRegistrar := grpc.NewServer(
		grpc.UnaryInterceptor(interceptors.AdminUnaryInterceptor),
		grpc.StreamInterceptor(interceptors.ServerStreamInterceptor),
	)
	proto.RegisterPersonHandlerServer(serverRegistrar, handlr)
	proto.RegisterUserHandlerServer(serverRegistrar, uhandlr)
	err = serverRegistrar.Serve(lis)
	fmt.Println("service started successfully")
	if err != nil {
		logrus.Fatalf("cannot start server: %s", err)
	}

	redisProd := producer.NewProducer(rdbClient)
	redisCons := consumer.NewConsumer(rdbClient)
	ctx, cancel := context.WithTimeout(context.Background(), timeOut*time.Second)
	defer cancel()
	// Redis Stream
	go redisProd.RedisProducer(ctx)
	go redisCons.RedisConsumer(ctx)
}
