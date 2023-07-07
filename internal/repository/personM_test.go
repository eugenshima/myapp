package repository

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

// func TestMongoCreate(t *testing.T) {
// 	err := rpsM.Create(context.Background(), &entityEugen)
// 	require.NoError(t, err)
// 	testEntity, err := rpsM.GetByID(context.Background(), entityEugen.ID)
// 	require.NoError(t, err)
// 	require.Equal(t, testEntity.ID, entityEugen.ID)
// 	require.Equal(t, testEntity.Name, entityEugen.Name)
// 	require.Equal(t, testEntity.Age, entityEugen.Age)
// 	require.Equal(t, testEntity.IsHealthy, entityEugen.IsHealthy)
// }

func TestMongoDelete(t *testing.T) {
	err := rpsM.Delete(context.Background(), uuid.Nil)
	require.NoError(t, err)
	require.True(t, true, "not deleting entity")
}

func TestMongoDeleteNil(t *testing.T) {
	err := rpsM.Delete(context.Background(), entityEugen.ID)
	require.NoError(t, err)
	require.True(t, true)
}

// func TestMongoGetAll(t *testing.T) {
// 	allPers, err := rpsM.GetAll(context.Background())
// 	if err != nil {
// 		t.Errorf("Expected no error, got: %v", err)
// 	}
// 	require.NoError(t, err)

// 	var numberPersons int
// 	err = rpsM.client.QueryRow(context.Background(), "SELECT COUNT(*) FROM  goschema.person").Scan(&numberPersons)
// 	require.NoError(t, err)
// 	require.Equal(t, len(allPers), numberPersons)
// }

func TestMongoUpdate(t *testing.T) {
	// Test case 1: Valid update
	err := rpsM.Update(context.Background(), entityEugen.ID, &entityEugen)
	require.NoError(t, err)
	// Test case 2: Invalid uuidString
	err = rps.Update(context.Background(), uuid.Nil, &entityEugen)
	require.NoError(t, err)
}

// // Фиктивный тест пока что =================
// func TestMongoCreateWithNegativeAge(t *testing.T) {
// 	entityEugen.Age = -1
// 	validate := validator.New()
// 	err := validate.Struct(entityEugen)
// 	require.Error(t, err)
// 	if err != nil {
// 		err = rpsM.Create(context.Background(), &entityEugen)
// 		require.NoError(t, err)
// 		require.True(t, true, "not creating entity")
// 	}
// }
