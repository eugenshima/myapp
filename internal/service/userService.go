// Package service provides a set of functions, which include business-logic in it
package service

import (
	"context"
	"time"

	"github.com/eugenshima/myapp/internal/model"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

const (
	salt       = "hjgrhjgw124617ajfhajs"
	signingKey = "gyewgb2rf8r2b8437frb23"
	tokenTTL   = 1 * time.Hour
)

type tokenClaims struct {
	jwt.StandardClaims
	id uuid.UUID `db:"id"`
}

// UserServiceImpl is a struct that contains a reference to the repository interface
type UserServiceImpl struct {
	rps UserRepositoryPsql
}

// NewService creates a new service
func NewUserServiceImpl(rps UserRepositoryPsql) *UserServiceImpl {
	return &UserServiceImpl{
		rps: rps,
	}
}

// PersonRepositoryPsql interface, which contains repository methods
type UserRepositoryPsql interface {
	GetUser(ctx context.Context, login, password string) (*model.User, error)
	Signup(context.Context, *model.User) error
	GetAll(context.Context) ([]*model.User, error)
}

// Login implements the UserServicePsql interface
func (db *UserServiceImpl) GenerateToken(c echo.Context, login, password string) (string, error) {
	user, err := db.rps.GetUser(c.Request().Context(), login, password)
	if err != nil {
		return "", err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodES256, &tokenClaims{
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(tokenTTL).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
		user.ID,
	})
	return token.SignedString([]byte(signingKey))
	//return db.rps.GetUser(c.Request().Context(), login, password)
}

// Signup implements the UserServicePsql interface
func (db *UserServiceImpl) Signup(c echo.Context, entity *model.User) error {
	return db.rps.Signup(c.Request().Context(), entity)
}

func (db *UserServiceImpl) GetAll(c echo.Context) ([]*model.User, error) {
	return db.rps.GetAll(c.Request().Context())
}

// func hashPassword(password string) string {
// 	fmt.Println(password)
// 	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
// 	if err != nil {
// 		return ""
// 	}
// 	return string(hashedPassword)
// }
