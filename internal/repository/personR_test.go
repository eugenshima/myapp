package repository

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

var redisConnPerson *RedisConnection

func TestRedisSetByID(t *testing.T) {
	err := redisConnPerson.RedisSetByID(context.Background(), &entityEugen)
	require.NoError(t, err)
	err = redisConnPerson.RedisDeleteByID(context.Background(), entityEugen.ID)
	require.NoError(t, err)
}

func TestRedisGetByID(t *testing.T) {
	err := redisConnPerson.RedisSetByID(context.Background(), &entityEugen)
	require.NoError(t, err)
	entity, err := redisConnPerson.RedisGetByID(context.Background(), entityEugen.ID)
	require.NoError(t, err)
	require.Equal(t, entityEugen.Name, entity.Name)
	require.Equal(t, entityEugen.Age, entity.Age)
	require.Equal(t, entityEugen.IsHealthy, entity.IsHealthy)
}

func TestRedisGetByWrongID(t *testing.T) {
	entity, err := redisConnPerson.RedisGetByID(context.Background(), uuid.New())
	require.Error(t, err)
	require.Nil(t, entity)
}

func TestRedisDeleteByID(t *testing.T) {
	err := redisConnPerson.RedisSetByID(context.Background(), &entityEugen)
	require.NoError(t, err)
	err = redisConnPerson.RedisDeleteByID(context.Background(), entityEugen.ID)
	require.NoError(t, err)
}

func TestRedisDeleteNil(t *testing.T) {
	err := redisConnPerson.RedisDeleteByID(context.Background(), uuid.New())
	require.Error(t, err)
}
