package repository

import (
	"context"
	"testing"

	"github.com/eugenshima/myapp/internal/model"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

var rps *PsqlConnection

var entityEugen = model.Person{
	ID:        uuid.New(),
	Name:      "Eugen",
	Age:       20,
	IsHealthy: true,
}

func TestPgxCreate(t *testing.T) {
	id, err := rps.Create(context.Background(), &entityEugen)
	require.NoError(t, err)
	// Step 1: Get byID
	testEntity, err := rps.GetByID(context.Background(), entityEugen.ID)
	require.NoError(t, err)
	// step 2: Data consistency check
	require.Equal(t, testEntity.ID, entityEugen.ID)
	require.Equal(t, testEntity.Name, entityEugen.Name)
	require.Equal(t, testEntity.Age, entityEugen.Age)
	require.Equal(t, testEntity.IsHealthy, entityEugen.IsHealthy)
	deletedID, err := rps.Delete(context.Background(), id)
	require.NotNil(t, deletedID)
	require.NoError(t, err)
}

func TestPgxDelete(t *testing.T) {
	id, err := rps.Create(context.Background(), &entityEugen)
	require.NoError(t, err)
	// step 1: Delete entity from database
	deletedID, err := rps.Delete(context.Background(), id)
	// step 2: Data consistency check
	require.NotNil(t, deletedID)
	require.NoError(t, err)
}

func TestPgxDeleteNil(t *testing.T) {
	id, err := rps.Create(context.Background(), &entityEugen)
	require.NoError(t, err)
	// step 1: Delete entity from database (wrong ID)
	deletedID, err := rps.Delete(context.Background(), uuid.New())
	require.Error(t, err)
	// step 2: Data consistency check
	require.NotNil(t, deletedID)
	require.NotEqual(t, id, deletedID)
	deletingTrash, err := rps.Delete(context.Background(), id)
	require.NotNil(t, deletingTrash)
	require.NoError(t, err)

}

func TestPgxGetAll(t *testing.T) {
	// Step 1: Get all entities
	res, err := rps.GetAll(context.Background())
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}
	require.NoError(t, err)
	// step 2: Data consistency check
	var count int
	err = rps.pool.QueryRow(context.Background(), "SELECT COUNT(*) FROM  goschema.person").Scan(&count)
	require.NoError(t, err)
	require.Equal(t, len(res), count)
}

func TestPgxUpdate(t *testing.T) {
	id, err := rps.Create(context.Background(), &entityEugen)
	require.NoError(t, err)
	// Step 1: Update entity with new data
	anotherID, err := rps.Update(context.Background(), entityEugen.ID, &entityEugen)
	require.NoError(t, err)
	// Step 2: Comparing IDs (Created with Updated)
	require.Equal(t, anotherID, id)
	deletedID, err := rps.Delete(context.Background(), id)
	require.NoError(t, err)
	require.NotNil(t, deletedID)
}

func TestPgxUpdateNil(t *testing.T) {
	id, err := rps.Create(context.Background(), &entityEugen)
	require.NoError(t, err)
	// step 1: Updating entity with Error check (Nil ID)
	anotherID, err := rps.Update(context.Background(), uuid.New(), &entityEugen)
	require.Error(t, err)
	// Step 2: Comparing IDs (Created with Updated)
	require.NotEqual(t, anotherID, id)
	deletedID, err := rps.Delete(context.Background(), id)
	require.NoError(t, err)
	require.NotNil(t, deletedID)
}

func TestGetByID(t *testing.T) {
	id, err := rps.Create(context.Background(), &entityEugen)
	require.NoError(t, err)
	// step 1: Getting ID with NoError check
	testEntity, err := rps.GetByID(context.Background(), id)
	require.NoError(t, err)
	// step 2: Data consistency check
	require.Equal(t, testEntity.ID, entityEugen.ID)
	require.Equal(t, testEntity.Name, entityEugen.Name)
	require.Equal(t, testEntity.Age, entityEugen.Age)
	require.Equal(t, testEntity.IsHealthy, entityEugen.IsHealthy)
	deletedID, err := rps.Delete(context.Background(), id)
	require.NoError(t, err)
	require.NotNil(t, deletedID)
}

func TestGetByWrongID(t *testing.T) {
	id, err := rps.Create(context.Background(), &entityEugen)
	require.NoError(t, err)
	entityEugen.ID = uuid.New()
	// step 1: Getting ID with Error check
	testEntity, err := rps.GetByID(context.Background(), entityEugen.ID)
	require.Error(t, err)
	// step 2: Data consistency check
	require.Nil(t, testEntity)
	deletedID, err := rps.Delete(context.Background(), id)
	require.NoError(t, err)
	require.NotNil(t, deletedID)
}
