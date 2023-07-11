package repository

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/ory/dockertest"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func SetupTestPgx() (*pgxpool.Pool, func(), error) {
	pool, err := dockertest.NewPool("")
	if err != nil {
		return nil, nil, fmt.Errorf("could not construct pool: %w", err)
	}
	resource, err := pool.Run("postgres", "latest", []string{
		"POSTGRES_USER=eugen",
		"POSTGRESQL_PASSWORD=ur2qly1ini",
		"POSTGRES_DB=eugen"})
	if err != nil {
		return nil, nil, fmt.Errorf("could not start resource: %w", err)
	}
	dbURL := "postgres://eugen:ur2qly1ini@localhost:5432/eugen"
	cfg, err := pgxpool.ParseConfig(dbURL)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to parse dbURL: %w", err)
	}
	dbpool, err := pgxpool.ConnectConfig(context.Background(), cfg)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to connect pgxpool: %w", err)
	}
	cleanup := func() {
		dbpool.Close()
		pool.Purge(resource)
	}
	return dbpool, cleanup, nil
}

func SetupTestMongoDB() (*mongo.Client, func(), error) {
	pool, err := dockertest.NewPool("")
	if err != nil {
		return nil, nil, fmt.Errorf("could not construct pool: %w", err)
	}
	resource, err := pool.Run("mongo", "6.0.6", []string{
		"MONGO_INITDB_ROOT_USERNAME=eugenshima",
		"MONGO_INITDB_ROOT_PASSWORD=ur2qly1ini",
		"MONGO_INITDB_DATABASE=my_mongo_base"})
	if err != nil {
		return nil, nil, fmt.Errorf("could not start resource: %w", err)
	}
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI("mongodb://localhost:27017/"))
	if err != nil {
		return nil, nil, fmt.Errorf("failed to connect mongoDB: %w", err)
	}
	cleanup := func() {
		client.Disconnect(context.Background())
		pool.Purge(resource)
	}
	return client, cleanup, nil
}

func TestMain(m *testing.M) {
	dbpool, cleanupPgx, err := SetupTestPgx()
	if err != nil {
		fmt.Println("Could not construct the pool: ", err)
		cleanupPgx()
		os.Exit(1)
	}
	rps = NewPsqlConnection(dbpool)
	urps = NewUserPsqlConnection(dbpool)
	client, cleanupMongo, err := SetupTestMongoDB()
	if err != nil {
		fmt.Println(err)
		cleanupMongo()
		os.Exit(1)
	}
	rpsM = NewMongoDBConnection(client)
	urpsM = NewUserMongoDBConnection(client)
	//TODO write redis connection + tests
	//redisCli, cleanupRedis, err := SetupTestRedis()
	// if err != nil {
	// 	fmt.Println(err)
	// 	cleanupRedis()
	// 	os.Exit(1)
	// }
	// redisP = NewRedisConnection(redisCli)
	exitVal := m.Run()
	cleanupPgx()
	cleanupMongo()
	// cleanupRedis()
	os.Exit(exitVal)
}
