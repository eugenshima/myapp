package repository

import (
	"context"
	"testing"

	"github.com/eugenshima/myapp/internal/model"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

var urpsM *UserMongoDBConnection

var mongotestLoginUser = model.Login{
	Login:    "test",
	Password: "test",
}

var mongotestSignupUser = model.Signup{
	Login:    "test",
	Password: "test",
	Role:     "user",
}

var mongotestUser = model.User{
	ID:       uuid.New(),
	Login:    "test",
	Password: []byte("test"),
	Role:     "user",
}

//TODO: dopilitb tests
func TestMongoUserGetAll(t *testing.T) {
	users, err := urpsM.GetAll(context.Background())
	require.NoError(t, err)
	require.NotEmpty(t, users)
}

func TestMongoUserGetUser(t *testing.T) {
	// creating...
	user, err := urpsM.GetUser(context.Background(), "test")
	require.NoError(t, err)
	require.NotEmpty(t, user)
	// deleting...
}

func TestMongoUserGetWrongUser(t *testing.T) {
	_, err := urpsM.GetUser(context.Background(), "wrong")
	require.Error(t, err)
}

func TestMongouserDeleteNil(t *testing.T) {
	err := urpsM.Signup(context.Background(), &testUser)
	require.NoError(t, err)
	id, err := urpsM.Delete(context.Background(), uuid.New())
	require.Error(t, err)
	require.NotEqual(t, testUser.ID, id)
	id, err = urpsM.Delete(context.Background(), testUser.ID)
	require.NoError(t, err)
	require.Equal(t, testUser.ID, id)
}
func TestMongoUserSignUp(t *testing.T) {
	// creating...
	err := urpsM.Signup(context.Background(), &testUser)
	require.NoError(t, err)
	test, err := urpsM.GetUser(context.Background(), "test")
	require.NoError(t, err)
	require.Equal(t, testUser.ID, test.ID)
	require.Equal(t, testUser.Login, test.Login)
	require.Equal(t, testUser.Password, test.Password)
	require.Equal(t, testUser.Role, test.Role)
	// deleting...
	id, err := urpsM.Delete(context.Background(), testUser.ID)
	require.NoError(t, err)
	require.Equal(t, testUser.ID, id)
}
