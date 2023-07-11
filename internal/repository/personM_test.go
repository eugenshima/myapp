package repository

import (
	"context"
	"testing"

	"github.com/eugenshima/myapp/internal/model"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

var rpsM *MongoDBConnection

var entityMongoEugen = model.Person{
	ID:        uuid.New(),
	Name:      "Eugen",
	Age:       20,
	IsHealthy: true,
}

func TestMongoCreate(t *testing.T) {
	// Step 1: Create test entity
	id, err := rpsM.Create(context.Background(), &entityMongoEugen)
	require.NoError(t, err)
	require.NotEmpty(t, id)
	// Step 2: Data consistency check
	testEntity, err := rpsM.GetByID(context.Background(), entityMongoEugen.ID)
	require.NoError(t, err)
	require.Equal(t, testEntity.ID, entityMongoEugen.ID)
	require.Equal(t, testEntity.Name, entityMongoEugen.Name)
	require.Equal(t, testEntity.Age, entityMongoEugen.Age)
	require.Equal(t, testEntity.IsHealthy, entityMongoEugen.IsHealthy)
	// Step 3: Delete test entity
	deletedID, err := rpsM.Delete(context.Background(), entityMongoEugen.ID)
	require.NoError(t, err)
	require.NotEmpty(t, deletedID)
}

func TestMongoDelete(t *testing.T) {
	// Step 1: Create test entity
	id, err := rpsM.Create(context.Background(), &entityMongoEugen)
	require.NoError(t, err)
	require.NotEmpty(t, id)
	// Step 2: Delete test entity
	deletedID, err := rpsM.Delete(context.Background(), id)
	require.NoError(t, err)
	require.NotEmpty(t, deletedID)
}

func TestMongoDeleteNil(t *testing.T) {
	// Step 1: Create test entity
	id, err := rpsM.Create(context.Background(), &entityMongoEugen)
	require.NoError(t, err)
	require.NotEmpty(t, id)
	// Step 2: Try to delete with wrong ID
	wrongID, err := rpsM.Delete(context.Background(), uuid.New())
	require.NoError(t, err)
	require.NotEqual(t, wrongID, id)
	// Step 3: Delete test entity
	deletedID, err := rpsM.Delete(context.Background(), id)
	require.NoError(t, err)
	require.NotEmpty(t, deletedID)
}

func TestMongoGetAll(t *testing.T) {
	// Step 1: Get all entities
	all, err := rpsM.GetAll(context.Background())
	require.NoError(t, err)
	require.NotEmpty(t, all)
}

func TestMongoUpdate(t *testing.T) {
	// Step 1: Create test entity
	id, err := rpsM.Create(context.Background(), &entityMongoEugen)
	require.NoError(t, err)
	require.NotEmpty(t, id)
	// Step 2: Update test entity
	id, err = rpsM.Update(context.Background(), entityEugen.ID, &entityEugen)
	require.NoError(t, err)
	require.NotEmpty(t, id)
	require.Equal(t, id, entityEugen.ID)
	// Step 3: Delete test entity
	deletedID, err := rpsM.Delete(context.Background(), entityMongoEugen.ID)
	require.NoError(t, err)
	require.NotEmpty(t, deletedID)
}

func TestMongoUpdateWrongID(t *testing.T) {
	// Step 1: Create test entity
	id, err := rpsM.Create(context.Background(), &entityMongoEugen)
	require.NoError(t, err)
	require.NotEmpty(t, id)
	// Step 2: Try to update test entity with wrong ID
	wrongID, err := rpsM.Update(context.Background(), uuid.New(), &entityEugen)
	require.NoError(t, err)
	require.NotEqual(t, wrongID, id)
	// Step 3: Delete test entity
	deletedID, err := rpsM.Delete(context.Background(), entityMongoEugen.ID)
	require.NoError(t, err)
	require.NotEmpty(t, deletedID)
}

func TestMongoGetByID(t *testing.T) {
	// Step 1: Create test entity
	id, err := rpsM.Create(context.Background(), &entityMongoEugen)
	require.NoError(t, err)
	require.NotEmpty(t, id)
	// Step 2: Get test entity
	testEntity, err := rpsM.GetByID(context.Background(), entityMongoEugen.ID)
	require.NoError(t, err)
	require.Equal(t, testEntity.ID, entityMongoEugen.ID)
	require.Equal(t, testEntity.Name, entityMongoEugen.Name)
	require.Equal(t, testEntity.Age, entityMongoEugen.Age)
	require.Equal(t, testEntity.IsHealthy, entityMongoEugen.IsHealthy)
	// Step 3: Delete test entity
	deletedID, err := rpsM.Delete(context.Background(), entityMongoEugen.ID)
	require.NoError(t, err)
	require.NotEmpty(t, deletedID)
}

func TestMongoGetByIDWrongID(t *testing.T) {
	// Step 1: Create test entity
	id, err := rpsM.Create(context.Background(), &entityMongoEugen)
	require.NoError(t, err)
	require.NotEmpty(t, id)
	// Step 2: Try to get test entity with wrong ID
	wrongID, err := rpsM.GetByID(context.Background(), uuid.New())
	require.Error(t, err)
	require.NotEqual(t, wrongID, id)
	// Step 3: Delete test entity
	deletedID, err := rpsM.Delete(context.Background(), entityMongoEugen.ID)
	require.NoError(t, err)
	require.NotEmpty(t, deletedID)
}
