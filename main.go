// Package main - entry point to our program
package main

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	_ "github.com/eugenshima/myapp/docs"
	cfgrtn "github.com/eugenshima/myapp/internal/config"
	"github.com/eugenshima/myapp/internal/handlers"
	middlwr "github.com/eugenshima/myapp/internal/middleware"
	"github.com/eugenshima/myapp/internal/repository"
	"github.com/eugenshima/myapp/internal/service"
	"github.com/go-playground/validator"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/labstack/echo/v4"
	swg "github.com/swaggo/echo-swagger"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type CustomValidator struct {
	validator *validator.Validate
}

func (cv *CustomValidator) Validate(i interface{}) error {
	if err := cv.validator.Struct(i); err != nil {
		logrus.Errorf("Validator: %v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("validator: %v", err))
	}
	return nil
}

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
	fmt.Println("Connected to PostgreSQL!")

	return pool, nil
}

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

	ch := "pgx"

	// Initializing the Database Connector (MongoDB)
	client, err := NewMongo(cfg.MongoDBAddr)
	if err != nil {
		echo.NewHTTPError(http.StatusBadRequest, fmt.Errorf("error creating database connection with MongoDB: %v", err))
		return
	}
	// Initializing the Database Connector (PostgreSQL)
	pool, err := NewDBPsql(cfg.PgxDBAddr)
	if err != nil {
		echo.NewHTTPError(http.StatusBadRequest, fmt.Errorf("error creating database connection with PostgreSQL: %v", err))
		return
	}
	rdbClient, err := NewDBRedis(cfg.RedisDBAddr)
	if err != nil {
		echo.NewHTTPError(http.StatusBadRequest, fmt.Errorf("error creating database connection with Redis: %v", err))
		return
	}
	var handlr *handlers.PersonHandler
	var uhandlr *handlers.UserHandler
	switch ch {
	case "mongo":

		// Person db mongodb
		rps := repository.NewMongoDBConnection(client)
		rdb := repository.NewRedisConnection(rdbClient)
		srv := service.NewPersonService(rps, rdb)
		handlr = handlers.NewPersonHandler(srv, validator.New())

	case "pgx":
		// Person db pgx
		rps := repository.NewPsqlConnection(pool)
		rdb := repository.NewRedisConnection(rdbClient)
		srv := service.NewPersonService(rps, rdb)
		handlr = handlers.NewPersonHandler(srv, validator.New())

		// User db pgx
		urps := repository.NewUserPsqlConnection(pool)
		usrv := service.NewUserServiceImpl(urps)
		uhandlr = handlers.NewUserHandlerImpl(usrv, validator.New())
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

		// Image requests
		image := api.Group("/image")
		image.Use(middlwr.AdminIdentity())
		image.GET("/get/:name", uhandlr.GetImage)
		image.POST("/set", uhandlr.SetImage)
	}
	e.GET("/swagger/*", swg.WrapHandler)

	//Redis Stream Producer
	go func() {
		identificator := 0
		for {
			id := strconv.FormatInt(time.Now().Unix(), 10)
			payload := map[string]interface{}{
				"timestamp": id,
				"content":   fmt.Sprintf("Redis streaming %d...", identificator),
			}

			identificator++
			id = id + "-" + strconv.Itoa(identificator)

			err := rdbClient.XAdd(context.Background(), &redis.XAddArgs{
				Stream: "testStream",
				MaxLen: 0,
				ID:     id,
				Values: payload,
			}).Err()
			if err != nil {
				fmt.Println("Error adding message to Redis Stream:", err)
			}

			time.Sleep(2 * time.Second)
		}
	}()

	// Redis Stream Consumer
	go func() {
		for {
			streams, err := rdbClient.XRead(context.Background(), &redis.XReadArgs{
				Streams: []string{"testStream", "0"},
				Count:   1,
				Block:   0,
			}).Result()
			if err != nil {
				fmt.Println("Error reading messages from Redis Stream:", err)
			}
			for _, stream := range streams {
				streamName := stream.Stream
				messages := stream.Messages

				for _, msg := range messages {
					messageID := msg.ID
					messageData := msg.Values

					fmt.Println("Received message from Redis Stream:", messageID, messageData)

					_, err := rdbClient.XDel(context.Background(), streamName, messageID).Result()
					if err != nil {
						fmt.Println("Error deleting message from Redis Stream:", err)
					}
				}
			}

			time.Sleep(2 * time.Second)
		}
	}()

	e.Logger.Fatal(e.Start(cfg.HTTPAddr))
}
