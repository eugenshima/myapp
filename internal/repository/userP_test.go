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
	// Step 1: Get all users
	res, err := urps.GetAll(context.Background())
	require.NotNil(t, &res)
	require.NoError(t, err)
	// step 2: Data consistency check
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
	// Step 1: Get user by login
	testGetUser, err := urps.GetUser(context.Background(), testLoginUser.Login)
	require.NotNil(t, testGetUser)
	require.NoError(t, err)
	// Step 2: Compare Hash And Password
	err = bcrypt.CompareHashAndPassword(testUser.Password, []byte(testSignupUser.Password))
	require.NoError(t, err)
	require.Equal(t, testUser.Role, testGetUser.Role)
}

func TestGetWrongUserName(t *testing.T) {
	// Step 1: Get User y wrong ID
	user, err := urps.GetUser(context.Background(), "wrong")
	require.Error(t, err)
	// step 2: Data consistency check
	require.Nil(t, user)
}

func TestSignUp(t *testing.T) {
	// Step 1: Sign Up using login&password "test"
	hashedPassword := hashPassword(testUser.Password)
	testUser.Password = hashedPassword
	err := urps.Signup(context.Background(), &testUser)
	require.NoError(t, err)
	del, err := rps.pool.Exec(context.Background(), "DELETE FROM goschema.user WHERE id=$1", testUser.ID)
	require.NoError(t, err)
	require.True(t, del.Delete())
}

func TestSaveRefreshToken(t *testing.T) {
	// Step 1: Sign Up using login&password "test"
	hashedPassword := hashPassword(testUser.Password)
	testUser.Password = hashedPassword
	err := urps.Signup(context.Background(), &testUser)
	require.NoError(t, err)
	// Step 2: Save Refresh Token
	err = urps.SaveRefreshToken(context.Background(), testUser.ID, []byte("testRefreshToken"))
	require.NoError(t, err)
	del, err := rps.pool.Exec(context.Background(), "DELETE FROM goschema.user WHERE id=$1", testUser.ID)
	require.NoError(t, err)
	require.True(t, del.Delete())
}

func TestSaveRefreshTokenToRandomID(t *testing.T) {
	// Step 1: Sign Up using login&password "test"
	hashedPassword := hashPassword(testUser.Password)
	testUser.Password = hashedPassword
	err := urps.Signup(context.Background(), &testUser)
	require.NoError(t, err)
	// Step 2: Save Refresh Token
	err = urps.SaveRefreshToken(context.Background(), uuid.New(), []byte("testRefreshToken"))
	require.Error(t, err)
	del, err := rps.pool.Exec(context.Background(), "DELETE FROM goschema.user WHERE id=$1", testUser.ID)
	require.NoError(t, err)
	require.True(t, del.Delete())
}

func TestGetRefreshToken(t *testing.T) {
	// Step 1: Sign Up using login&password "test"
	hashedPassword := hashPassword(testUser.Password)
	testUser.Password = hashedPassword
	err := urps.Signup(context.Background(), &testUser)
	require.NoError(t, err)
	// Step 2: Save Refresh Token
	err = urps.SaveRefreshToken(context.Background(), testUser.ID, []byte("testRefreshToken"))
	require.NoError(t, err)
	// Step 3: Get Refresh Token
	token, err := urps.GetRefreshToken(context.Background(), testUser.ID)
	require.NoError(t, err)
	require.Equal(t, []byte("testRefreshToken"), token)
}

func TestGetWrongRefreshToken(t *testing.T) {
	// Step 1: Sign Up using login&password "test"
	hashedPassword := hashPassword(testUser.Password)
	testUser.Password = hashedPassword
	err := urps.Signup(context.Background(), &testUser)
	require.NoError(t, err)
	// Step 2: Save Refresh Token
	err = urps.SaveRefreshToken(context.Background(), testUser.ID, []byte("testRefreshToken"))
	require.NoError(t, err)
	// Step 3: Get Refresh Token
	token, err := urps.GetRefreshToken(context.Background(), uuid.New())
	require.Error(t, err)
	require.Nil(t, token)
}

func TestGetRoleByID(t *testing.T) {
	//Step 1: Sign Up using login&password "test"
	hashedPassword := hashPassword(testUser.Password)
	testUser.Password = hashedPassword
	err := urps.Signup(context.Background(), &testUser)
	require.NoError(t, err)
	// Step 3: Get Role by ID
	role, err := urps.GetRoleByID(context.Background(), testUser.ID)
	require.NoError(t, err)
	require.Equal(t, "user", role)
}

func TestGetRoleByWrongID(t *testing.T) {
	//Step 1: Sign Up using login&password "test"
	hashedPassword := hashPassword(testUser.Password)
	testUser.Password = hashedPassword
	err := urps.Signup(context.Background(), &testUser)
	require.NoError(t, err)
	// Step 3: Get Role by ID
	role, err := urps.GetRoleByID(context.Background(), uuid.New())
	require.Error(t, err)
	require.Equal(t, role, "")
}

// HashPassword func returns hashed password using bcrypt algorithm
func hashPassword(password []byte) []byte {
	hashedPassword, err := bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)
	if err != nil {
		return nil
	}
	return hashedPassword
}
