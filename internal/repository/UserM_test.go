package repository

import (
	"context"
	"testing"

	"github.com/eugenshima/myapp/internal/model"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var urpsM *UserMongoDBConnection

var mongotestUser = model.User{
	ID:           uuid.New(),
	Login:        "test",
	Password:     []byte("test"),
	Role:         "user",
	RefreshToken: nil,
}

func TestMongoUserGetAll(t *testing.T) {
	users, err := urpsM.GetAll(context.Background())
	require.NoError(t, err)
	require.NotEmpty(t, users)
}

func TestMongoUserGetUser(t *testing.T) {
	err := urpsM.Signup(context.Background(), &testUser)
	require.NoError(t, err)
	user, err := urpsM.GetUser(context.Background(), "test")
	assert.NoError(t, err)
	assert.NotEmpty(t, user)
	err = urpsM.Delete(context.Background(), testUser.ID)
	require.NoError(t, err)
}

func TestMongoUserGetWrongUser(t *testing.T) {
	_, err := urpsM.GetUser(context.Background(), "wrong")
	require.Error(t, err)
}

func TestMongoUserDelete(t *testing.T) {
	err := urpsM.Signup(context.Background(), &testUser)
	require.NoError(t, err)
	err = urpsM.Delete(context.Background(), testUser.ID)
	require.NoError(t, err)
}

func TestMongouserDeleteNil(t *testing.T) {
	err := urpsM.Signup(context.Background(), &testUser)
	require.NoError(t, err)
	err = urpsM.Delete(context.Background(), uuid.New())
	assert.Error(t, err)
	err = urpsM.Delete(context.Background(), testUser.ID)
	require.NoError(t, err)
}
func TestMongoUserSignUp(t *testing.T) {
	err := urpsM.Signup(context.Background(), &testUser)
	require.NoError(t, err)
	test, err := urpsM.GetUser(context.Background(), "test")
	assert.NoError(t, err)
	assert.Equal(t, testUser.ID, test.ID)
	assert.Equal(t, testUser.Login, test.Login)
	assert.Equal(t, testUser.Password, test.Password)
	assert.Equal(t, testUser.Role, test.Role)
	err = urpsM.Delete(context.Background(), testUser.ID)
	require.NoError(t, err)
}

func TestMongoUserSaveRefreshToken(t *testing.T) {
	err := urpsM.Signup(context.Background(), &mongotestUser)
	require.NoError(t, err)
	err = urpsM.SaveRefreshToken(context.Background(), mongotestUser.ID, []byte("testRefreshToken"))
	require.NoError(t, err)
	test, err := urpsM.GetUser(context.Background(), "test")
	assert.NoError(t, err)
	assert.Equal(t, mongotestUser.ID, test.ID)
	assert.Equal(t, mongotestUser.Login, test.Login)
	assert.Equal(t, mongotestUser.Password, test.Password)
	assert.Equal(t, mongotestUser.Role, test.Role)
	assert.Equal(t, []byte("testRefreshToken"), test.RefreshToken)
	err = urpsM.Delete(context.Background(), mongotestUser.ID)
	require.NoError(t, err)
}

func TestMongoUserGetRefreshToken(t *testing.T) {
	err := urpsM.Signup(context.Background(), &mongotestUser)
	require.NoError(t, err)
	err = urpsM.SaveRefreshToken(context.Background(), mongotestUser.ID, []byte("testRefreshToken"))
	require.NoError(t, err)
	test, err := urpsM.GetUser(context.Background(), "test")
	assert.NotEmpty(t, test)
	assert.NoError(t, err)
	err = urpsM.Delete(context.Background(), mongotestUser.ID)
	require.NoError(t, err)
}

func TestMongoUserGetWrongToken(t *testing.T) {
	err := urpsM.Signup(context.Background(), &mongotestUser)
	require.NoError(t, err)
	err = urpsM.SaveRefreshToken(context.Background(), mongotestUser.ID, []byte("testRefreshToken"))
	require.NoError(t, err)
	test, err := urpsM.GetUser(context.Background(), "test")
	assert.NotEmpty(t, test)
	assert.NoError(t, err)
	testToken, err := urpsM.GetRefreshToken(context.Background(), uuid.New())
	assert.Error(t, err)
	assert.Nil(t, testToken)
	err = urpsM.Delete(context.Background(), mongotestUser.ID)
	require.NoError(t, err)
}

func TestMongoUserGetRoleByID(t *testing.T) {
	err := urpsM.Signup(context.Background(), &mongotestUser)
	require.NoError(t, err)
	testRole, err := urpsM.GetRoleByID(context.Background(), mongotestUser.ID)
	assert.NotEmpty(t, testRole)
	assert.NoError(t, err)
	err = urpsM.Delete(context.Background(), mongotestUser.ID)
	require.NoError(t, err)
}

func TestMongoUserGetRoleByWrongID(t *testing.T) {
	err := urpsM.Signup(context.Background(), &mongotestUser)
	require.NoError(t, err)
	testRole, err := urpsM.GetRoleByID(context.Background(), uuid.New())
	assert.Error(t, err)
	assert.Equal(t, testRole, "")
	err = urpsM.Delete(context.Background(), mongotestUser.ID)
	require.NoError(t, err)
}
