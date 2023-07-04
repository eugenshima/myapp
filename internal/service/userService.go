// Package service provides a set of functions, which include business-logic in it
package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/caarlos0/env/v9"
	"github.com/eugenshima/myapp/internal/config"
	"github.com/eugenshima/myapp/internal/model"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

const (
	salt            = "hjgrhjgw124617ajfhajs"
	accessTokenTTL  = 3 * time.Minute
	refreshTokenTTL = 24 * time.Hour
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
	SaveRefreshToken(ctx context.Context, id uuid.UUID, token string) error
}

// Login implements the UserServicePsql interface
func (db *UserServiceImpl) GenerateToken(ctx context.Context, login, password string) (string, error) {

	cfg := config.Config{}
	err := env.Parse(&cfg)
	if err != nil {
		return "", fmt.Errorf("error parsing environment variable: %v", err)
	}
	// hashedPassword := hashPassword(password)
	id, pass, err := db.rps.GetUser(ctx, login)
	if err != nil {
		return "", fmt.Errorf("error in GenerateToken (GetUser): %v", err)
	}
	err = bcrypt.CompareHashAndPassword(pass, []byte(password))
	if err != nil {
		return "", fmt.Errorf("error in GenerateToken (CompareHashAndPassword): %v", err)
	}

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, &tokenClaims{
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(accessTokenTTL).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
		id,
	})
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, &tokenClaims{
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(refreshTokenTTL).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
		id,
	})
	refreshToDB, err := refreshToken.SignedString([]byte(cfg.SigningKey))
	if err != nil {
		return "", fmt.Errorf("error in refresh token: %v", err)
	}
	err = db.rps.SaveRefreshToken(ctx, id, refreshToDB)
	if err != nil {
		return "", fmt.Errorf("error in GenerateToken(SaveRefreshToken): %v", err)
	}
	return accessToken.SignedString([]byte(cfg.SigningKey))
}

// Signup implements the UserServicePsql interface
func (db *UserServiceImpl) Signup(ctx context.Context, entity *model.User) error {
	hashedPassword := hashPassword(entity.Password)
	entity.Password = hashedPassword
	return db.rps.Signup(ctx, entity)
}

func (db *UserServiceImpl) GetAll(ctx context.Context) ([]*model.User, error) {
	return db.rps.GetAll(ctx)
}

func (db *UserServiceImpl) ParseToken(ctx context.Context, accessToken string) (uuid.UUID, error) {
	cfg := config.Config{}
	err := env.Parse(&cfg)
	if err != nil {
		return uuid.Nil, fmt.Errorf("error parsing environment variable: %v", err)
	}
	token, err := jwt.ParseWithClaims(accessToken, &tokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}

		return []byte(cfg.SigningKey), nil
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
