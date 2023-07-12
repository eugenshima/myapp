package repository

import (
	"context"
	"testing"

	"github.com/eugenshima/myapp/internal/model"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"github.com/stretchr/testify/require"
)

var urps *UserPsqlConnection

var testLoginUser = model.Login{
	Login:    "test",
	Password: "test",
}

var testSignupUser = model.Signup{
	Login:    "test",
	Password: "test",
	Role:     "user",
}

var testUser = model.User{
	ID:       uuid.New(),
	Login:    "test",
	Password: []byte("test"),
	Role:     "user",
}

func TestGetAll(t *testing.T) {
	res, err := urps.GetAll(context.Background())
	require.NotNil(t, &res)
	require.NoError(t, err)
	var count int
	err = rps.pool.QueryRow(context.Background(), "SELECT COUNT(*) FROM  goschema.user").Scan(&count)
	require.NoError(t, err)
	require.Equal(t, len(res), count)
}

func TestGetUser(t *testing.T) {
	hashedPassword := hashPassword(testUser.Password)
	testUser.Password = hashedPassword
	err := urps.Signup(context.Background(), &testUser)
	require.NoError(t, err)
	testGetUser, err := urps.GetUser(context.Background(), testLoginUser.Login)
	require.NotNil(t, testGetUser)
	require.NoError(t, err)
	err = bcrypt.CompareHashAndPassword(testUser.Password, []byte(testSignupUser.Password))
	require.NoError(t, err)
	require.Equal(t, testUser.Role, testGetUser.Role)
	err = urps.Delete(context.Background(), testUser.ID)
	require.NoError(t, err)
}

func TestGetWrongUserName(t *testing.T) {
	hashedPassword := hashPassword(testUser.Password)
	testUser.Password = hashedPassword
	err := urps.Signup(context.Background(), &testUser)
	require.NoError(t, err)
	user, err := urps.GetUser(context.Background(), "wrong")
	require.Error(t, err)
	require.Nil(t, user)
	err = urps.Delete(context.Background(), testUser.ID)
	require.NoError(t, err)
}

func TestSignUp(t *testing.T) {
	hashedPassword := hashPassword(testUser.Password)
	testUser.Password = hashedPassword
	err := urps.Signup(context.Background(), &testUser)
	require.NoError(t, err)
	err = urps.Delete(context.Background(), testUser.ID)
	require.NoError(t, err)
}

func TestSaveRefreshToken(t *testing.T) {
	hashedPassword := hashPassword(testUser.Password)
	testUser.Password = hashedPassword
	err := urps.Signup(context.Background(), &testUser)
	require.NoError(t, err)
	err = urps.SaveRefreshToken(context.Background(), testUser.ID, []byte("testRefreshToken"))
	require.NoError(t, err)
	err = urps.Delete(context.Background(), testUser.ID)
	require.NoError(t, err)
}

func TestSaveRefreshTokenToRandomID(t *testing.T) {
	hashedPassword := hashPassword(testUser.Password)
	testUser.Password = hashedPassword
	err := urps.Signup(context.Background(), &testUser)
	require.NoError(t, err)
	err = urps.SaveRefreshToken(context.Background(), uuid.New(), []byte("testRefreshToken"))
	require.Error(t, err)
	err = urps.Delete(context.Background(), testUser.ID)
	require.NoError(t, err)
}

func TestGetRefreshToken(t *testing.T) {
	hashedPassword := hashPassword(testUser.Password)
	testUser.Password = hashedPassword
	err := urps.Signup(context.Background(), &testUser)
	require.NoError(t, err)
	err = urps.SaveRefreshToken(context.Background(), testUser.ID, []byte("testRefreshToken"))
	require.NoError(t, err)
	token, err := urps.GetRefreshToken(context.Background(), testUser.ID)
	require.NoError(t, err)
	require.Equal(t, []byte("testRefreshToken"), token)
	err = urps.Delete(context.Background(), testUser.ID)
	require.NoError(t, err)
}

func TestGetWrongRefreshToken(t *testing.T) {
	hashedPassword := hashPassword(testUser.Password)
	testUser.Password = hashedPassword
	err := urps.Signup(context.Background(), &testUser)
	require.NoError(t, err)
	err = urps.SaveRefreshToken(context.Background(), testUser.ID, []byte("testRefreshToken"))
	require.NoError(t, err)
	token, err := urps.GetRefreshToken(context.Background(), uuid.New())
	require.Error(t, err)
	require.Nil(t, token)
	err = urps.Delete(context.Background(), testUser.ID)
	require.NoError(t, err)
}

func TestGetRoleByID(t *testing.T) {
	hashedPassword := hashPassword(testUser.Password)
	testUser.Password = hashedPassword
	err := urps.Signup(context.Background(), &testUser)
	require.NoError(t, err)
	role, err := urps.GetRoleByID(context.Background(), testUser.ID)
	require.NoError(t, err)
	require.Equal(t, "user", role)
	err = urps.Delete(context.Background(), testUser.ID)
	require.NoError(t, err)
}

func TestGetRoleByWrongID(t *testing.T) {
	hashedPassword := hashPassword(testUser.Password)
	testUser.Password = hashedPassword
	err := urps.Signup(context.Background(), &testUser)
	require.NoError(t, err)
	role, err := urps.GetRoleByID(context.Background(), uuid.New())
	require.Error(t, err)
	require.Equal(t, role, "")
	err = urps.Delete(context.Background(), testUser.ID)
	require.NoError(t, err)
}

func TestDeleteUserByID(t *testing.T) {
	hashedPassword := hashPassword(testUser.Password)
	testUser.Password = hashedPassword
	err := urps.Signup(context.Background(), &testUser)
	require.NoError(t, err)
	err = urps.Delete(context.Background(), testUser.ID)
	require.NoError(t, err)
}

func TestDeleteWrongID(t *testing.T) {
	hashedPassword := hashPassword(testUser.Password)
	testUser.Password = hashedPassword
	err := urps.Signup(context.Background(), &testUser)
	require.NoError(t, err)
	err = urps.Delete(context.Background(), uuid.New())
	require.Error(t, err)
	err = urps.Delete(context.Background(), testUser.ID)
	require.NoError(t, err)
}
