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
	testEntity, err := rps.GetByID(context.Background(), entityEugen.ID)
	require.NoError(t, err)
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
	deletedID, err := rps.Delete(context.Background(), id)
	require.NotNil(t, deletedID)
	require.NoError(t, err)
}

func TestPgxDeleteNil(t *testing.T) {
	id, err := rps.Create(context.Background(), &entityEugen)
	require.NoError(t, err)
	deletedID, err := rps.Delete(context.Background(), uuid.New())
	require.Error(t, err)
	require.NotNil(t, deletedID)
	require.NotEqual(t, id, deletedID)
	deletingTrash, err := rps.Delete(context.Background(), id)
	require.NotNil(t, deletingTrash)
	require.NoError(t, err)

}

func TestPgxGetAll(t *testing.T) {
	res, err := rps.GetAll(context.Background())
	require.NoError(t, err)
	var count int
	err = rps.pool.QueryRow(context.Background(), "SELECT COUNT(*) FROM  goschema.person").Scan(&count)
	require.NoError(t, err)
	require.Equal(t, len(res), count)
}

func TestPgxUpdate(t *testing.T) {
	id, err := rps.Create(context.Background(), &entityEugen)
	require.NoError(t, err)
	anotherID, err := rps.Update(context.Background(), entityEugen.ID, &entityEugen)
	require.NoError(t, err)
	require.Equal(t, anotherID, id)
	deletedID, err := rps.Delete(context.Background(), id)
	require.NoError(t, err)
	require.NotNil(t, deletedID)
}

func TestPgxUpdateNil(t *testing.T) {
	id, err := rps.Create(context.Background(), &entityEugen)
	require.NoError(t, err)
	anotherID, err := rps.Update(context.Background(), uuid.New(), &entityEugen)
	require.Error(t, err)
	require.NotEqual(t, anotherID, id)
	deletedID, err := rps.Delete(context.Background(), id)
	require.NoError(t, err)
	require.NotNil(t, deletedID)
}

func TestGetByID(t *testing.T) {
	id, err := rps.Create(context.Background(), &entityEugen)
	require.NoError(t, err)
	testEntity, err := rps.GetByID(context.Background(), id)
	require.NoError(t, err)
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
	testEntity, err := rps.GetByID(context.Background(), entityEugen.ID)
	require.Error(t, err)
	require.Nil(t, testEntity)
	deletedID, err := rps.Delete(context.Background(), id)
	require.NoError(t, err)
	require.NotNil(t, deletedID)
}
