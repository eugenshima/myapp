package repository

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/eugenshima/myapp/internal/model"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/ory/dockertest"
	"github.com/stretchr/testify/require"
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
	resource, err := pool.Run("mongo", "latest", []string{
		"MONGO_INITDB_ROOT_USERNAME=eugenshima",
		"MONGO_INITDB_ROOT_PASSWORD=ur2qly1ini",
		"MONGO_INITDB_DATABASE=my_mongo_base"})
	if err != nil {
		return nil, nil, fmt.Errorf("could not start resource: %w", err)
	}
	port := resource.GetPort("27017/tcp")
	mongoURL := fmt.Sprintf("mongodb://eugenshima:ur2qly1ini@localhost:%s", port)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to parse dbURL: %w", err)
	}
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(mongoURL))
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
	client, cleanupMongo, err := SetupTestMongoDB()
	if err != nil {
		fmt.Println(err)
		cleanupMongo()
		os.Exit(1)
	}
	rpsM = NewMongoDBConnection(client)
	exitVal := m.Run()
	cleanupPgx()
	cleanupMongo()
	os.Exit(exitVal)
}

var rps *PsqlConnection

var rpsM *MongoDBConnection

var entityEugen = model.Person{
	ID:        uuid.New(),
	Name:      "Eugen",
	Age:       20,
	IsHealthy: true,
}

// func TestPgxCreate(t *testing.T) {
// 	err := rps.Create(context.Background(), &entityEugen)
// 	require.NoError(t, err)
// 	testEntity, err := rps.GetByID(context.Background(), entityEugen.ID)
// 	require.NoError(t, err)
// 	require.Equal(t, testEntity.ID, entityEugen.ID)
// 	require.Equal(t, testEntity.Name, entityEugen.Name)
// 	require.Equal(t, testEntity.Age, entityEugen.Age)
// 	require.Equal(t, testEntity.IsHealthy, entityEugen.IsHealthy)
// }

func TestPgxDelete(t *testing.T) {
	err := rps.Delete(context.Background(), uuid.Nil)
	require.NoError(t, err)
	require.True(t, true, "not deleting entity")
}

func TestPgxDeleteNil(t *testing.T) {
	err := rps.Delete(context.Background(), entityEugen.ID)
	require.NoError(t, err)
	require.True(t, true)
}

func TestPgxGetAll(t *testing.T) {
	allPers, err := rps.GetAll(context.Background())
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}
	require.NoError(t, err)

	var numberPersons int
	err = rps.pool.QueryRow(context.Background(), "SELECT COUNT(*) FROM  goschema.person").Scan(&numberPersons)
	require.NoError(t, err)
	require.Equal(t, len(allPers), numberPersons)
}

func TestPgxUpdate(t *testing.T) {
	// Test case 1: Valid update
	err := rps.Update(context.Background(), entityEugen.ID, &entityEugen)
	require.NoError(t, err)
	// Test case 2: Invalid uuidString
	err = rps.Update(context.Background(), uuid.Nil, &entityEugen)
	require.NoError(t, err)
}

// // Фиктивный тест пока что =================
// func TestPgxCreateWithNegativeAge(t *testing.T) {
// 	entityEugen.Age = -1
// 	validate := validator.New()
// 	err := validate.Struct(entityEugen)
// 	require.Error(t, err)
// 	if err != nil {
// 		err = rps.Create(context.Background(), &entityEugen)
// 		require.NoError(t, err)
// 		require.True(t, true, "not creating entity")
// 	}
// }
