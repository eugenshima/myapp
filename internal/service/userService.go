// Package service provides a set of functions, which include business-logic in it
package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/eugenshima/myapp/internal/model"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

const (
	salt       = "hjgrhjgw124617ajfhajs"
	signingKey = "gyewgb2rf8r2b8437frb23"
	tokenTTL   = 1 * time.Hour
)

type tokenClaims struct {
	jwt.StandardClaims
	UserId uuid.UUID `db:"id"`
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
	GetUser(ctx context.Context, login string) (uuid.UUID, []byte, error)
	Signup(context.Context, *model.User) error
	GetAll(context.Context) ([]*model.User, error)
}

// Login implements the UserServicePsql interface
func (db *UserServiceImpl) GenerateToken(c echo.Context, login, password string) (string, error) {
	// hashedPassword := hashPassword(password)
	id, pass, err := db.rps.GetUser(c.Request().Context(), login)
	if err != nil {
		return "", err
	}
	err = bcrypt.CompareHashAndPassword(pass, []byte(password))
	if err != nil {
		return "", err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &tokenClaims{
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(tokenTTL).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
		id,
	})
	return token.SignedString([]byte(signingKey))
	//return db.rps.GetUser(c.Request().Context(), login, password)
}

// Signup implements the UserServicePsql interface
func (db *UserServiceImpl) Signup(c echo.Context, entity *model.User) error {
	hashedPassword := hashPassword(entity.Password)
	entity.Password = hashedPassword
	return db.rps.Signup(c.Request().Context(), entity)
}

func (db *UserServiceImpl) GetAll(c echo.Context) ([]*model.User, error) {
	return db.rps.GetAll(c.Request().Context())
}

func (db *UserServiceImpl) ParseToken(accessToken string) (uuid.UUID, error) {
	token, err := jwt.ParseWithClaims(accessToken, &tokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}

		return []byte(signingKey), nil
	})
	if err != nil {
		return uuid.Nil, err
	}

	claims, ok := token.Claims.(*tokenClaims)
	if !ok {
		return uuid.Nil, errors.New("token claims are not of type *tokenClaims")
	}
	return claims.UserId, nil
}

func hashPassword(password []byte) []byte {
	fmt.Println(password)
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	if err != nil {
		return nil
	}
	return hashedPassword
}
