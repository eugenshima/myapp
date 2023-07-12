package repository

import (
	"context"
	"testing"

	"github.com/eugenshima/myapp/internal/model"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

var redisConnUser *UserRedisConnection

var testUserRedis = model.User{
	ID:       uuid.New(),
	Login:    "test",
	Password: []byte("test"),
	Role:     "user",
}

func TestUserRedisGet(t *testing.T) {
	err := redisConnUser.Set(context.Background(), &testUserRedis)
	require.NoError(t, err)
	user, err := redisConnUser.Get(context.Background(), testUserRedis.ID)
	require.NoError(t, err)
	require.Equal(t, testUserRedis.Login, user.Login)
	require.Equal(t, testUserRedis.Password, user.Password)
	require.Equal(t, testUserRedis.Role, user.Role)
	require.Nil(t, user.RefreshToken)
	err = redisConnUser.Delete(context.Background(), testUserRedis.ID)
	require.NoError(t, err)
}

func TestUserRedisGetWrongID(t *testing.T) {
	err := redisConnUser.Set(context.Background(), &testUserRedis)
	require.NoError(t, err)
	user, err := redisConnUser.Get(context.Background(), uuid.New())
	require.Error(t, err)
	require.Nil(t, user)
	err = redisConnUser.Delete(context.Background(), testUserRedis.ID)
	require.NoError(t, err)
}

func TestUserRedisSet(t *testing.T) {
	err := redisConnUser.Set(context.Background(), &testUserRedis)
	require.NoError(t, err)
	user, err := redisConnUser.Get(context.Background(), testUserRedis.ID)
	require.NoError(t, err)
	require.Equal(t, testUserRedis.Login, user.Login)
	require.Equal(t, testUserRedis.Password, user.Password)
	require.Equal(t, testUserRedis.Role, user.Role)
	require.Nil(t, user.RefreshToken)
	err = redisConnUser.Delete(context.Background(), testUserRedis.ID)
	require.NoError(t, err)
}

func TestUserRedisSetRefreshToken(t *testing.T) {
	err := redisConnUser.Set(context.Background(), &testUserRedis)
	require.NoError(t, err)
	err = redisConnUser.SetRefreshToken(context.Background(), testUserRedis.ID, []byte("testRefreshToken"))
	require.NoError(t, err)
	user, err := redisConnUser.Get(context.Background(), testUserRedis.ID)
	require.NoError(t, err)
	require.Equal(t, testUserRedis.Login, user.Login)
	require.Equal(t, testUserRedis.Password, user.Password)
	require.Equal(t, testUserRedis.Role, user.Role)
	require.Equal(t, user.RefreshToken, []byte("testRefreshToken"))
	err = redisConnUser.Delete(context.Background(), testUserRedis.ID)
	require.NoError(t, err)
}
func TestUserRedisGetRefreshToken(t *testing.T) {
	err := redisConnUser.Set(context.Background(), &testUserRedis)
	require.NoError(t, err)
	err = redisConnUser.SetRefreshToken(context.Background(), testUserRedis.ID, []byte("testRefreshToken"))
	require.NoError(t, err)
	user, err := redisConnUser.Get(context.Background(), testUserRedis.ID)
	require.NoError(t, err)
	require.Equal(t, user.RefreshToken, []byte("testRefreshToken"))
	err = redisConnUser.Delete(context.Background(), testUserRedis.ID)
	require.NoError(t, err)
}

func TestUserRedisGetWrongRefreshToken(t *testing.T) {
	err := redisConnUser.Set(context.Background(), &testUserRedis)
	require.NoError(t, err)
	err = redisConnUser.SetRefreshToken(context.Background(), testUserRedis.ID, []byte("testRefreshToken"))
	require.NoError(t, err)
	user, err := redisConnUser.Get(context.Background(), uuid.New())
	require.Error(t, err)
	require.Nil(t, user)
	err = redisConnUser.Delete(context.Background(), testUserRedis.ID)
	require.NoError(t, err)
}

func TestUserRedisDelete(t *testing.T) {
	err := redisConnUser.Set(context.Background(), &testUserRedis)
	require.NoError(t, err)
	err = redisConnUser.Delete(context.Background(), testUserRedis.ID)
	require.NoError(t, err)
	user, err := redisConnUser.Get(context.Background(), testUserRedis.ID)
	require.Error(t, err)
	require.Nil(t, user)
}

func TestUserRedisDeleteNil(t *testing.T) {
	err := redisConnUser.Set(context.Background(), &testUserRedis)
	require.NoError(t, err)
	err = redisConnUser.Delete(context.Background(), uuid.New())
	require.Error(t, err)
	err = redisConnUser.Delete(context.Background(), testUserRedis.ID)
	require.NoError(t, err)
}
