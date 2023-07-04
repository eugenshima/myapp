// Package service provides a set of functions, which include business-logic in it
package service

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
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

// tokenClaims struct contains information about the claims associated with the given token
type tokenClaims struct {
	jwt.StandardClaims
	UserID uuid.UUID `db:"id"`
}

// UserServiceImpl is a struct that contains a reference to the repository interface
type UserServiceImpl struct {
	rps UserRepositoryPsql
}

// NewUserServiceImpl creates a new service
func NewUserServiceImpl(rps UserRepositoryPsql) *UserServiceImpl {
	return &UserServiceImpl{
		rps: rps,
	}
}

// UserRepositoryPsql interface, which contains repository methods
type UserRepositoryPsql interface {
	GetUser(ctx context.Context, login string) (uuid.UUID, []byte, error)
	Signup(context.Context, *model.User) error
	GetAll(context.Context) ([]*model.User, error)
	SaveRefreshToken(ctx context.Context, id uuid.UUID, token []byte) error
}

// GenerateToken implements the UserServicePsql interface
func (db *UserServiceImpl) GenerateToken(ctx context.Context, login, password string) (string, string, error) {
	cfg := config.Config{}
	err := env.Parse(&cfg)
	if err != nil {
		return "", "", fmt.Errorf("error parsing environment variable: %v", err)
	}
	id, pass, err := db.rps.GetUser(ctx, login)
	if err != nil {
		return "", "", fmt.Errorf("error in GenerateToken (GetUser): %v", err)
	}
	err = bcrypt.CompareHashAndPassword(pass, []byte(password))
	if err != nil {
		return "", "", fmt.Errorf("error in GenerateToken (CompareHashAndPassword): %v", err)
	}

	accessToken, err := GenerateAccessToken()
	if err != nil {
		return "", "", fmt.Errorf("error in GenerateToken (GenerateAccessToken): %v", err)
	}

	refreshToken, err := GenerateRefreshToken()
	if err != nil {
		return "", "", fmt.Errorf("error in GenerateToken (GenerateRefreshToken): %v", err)
	}

	hashedRefreshToken, err := HashRefreshToken(refreshToken)
	if err != nil {
		return "", "", fmt.Errorf("error in refresh token: %v", err)
	}

	err = db.rps.SaveRefreshToken(ctx, id, hashedRefreshToken)
	if err != nil {
		return "", "", fmt.Errorf("error in GenerateToken(SaveRefreshToken): %v", err)
	}
	return accessToken, refreshToken, nil
}

// Signup implements the UserServicePsql interface
func (db *UserServiceImpl) Signup(ctx context.Context, entity *model.User) error {
	hashedPassword := hashPassword(entity.Password)
	entity.Password = hashedPassword
	return db.rps.Signup(ctx, entity)
}

// GetAll implements the UserServicePsql interface
func (db *UserServiceImpl) GetAll(ctx context.Context) ([]*model.User, error) {
	return db.rps.GetAll(ctx)
}

func hashPassword(password []byte) []byte {
	fmt.Println(password)
	hashedPassword, err := bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)

	if err != nil {
		return nil
	}
	return hashedPassword
}

func HashRefreshToken(refreshToken string) ([]byte, error) {
	hash := sha256.New()
	_, err := hash.Write([]byte(refreshToken))
	if err != nil {
		return nil, fmt.Errorf("error hashing refresh token: %v", err)
	}
	hashedToken := hex.EncodeToString(hash.Sum(nil))

	return []byte(hashedToken), nil
}

func RefreshAccessToken(refreshToken string) (string, error) {
	// Проверяем валидность refresh токена
	valid, err := ValidateRefreshToken(refreshToken)
	if err != nil {
		return "", fmt.Errorf("error validating refresh token: %v", err)
	}
	if !valid {
		return "", fmt.Errorf("invalid refresh token : %v", err)
	}

	// Генерируем новый access токен
	accessToken, err := GenerateAccessToken()
	if err != nil {
		return "", fmt.Errorf("error generating access token: %v", err)
	}

	return accessToken, nil
}

func ValidateRefreshToken(refreshToken string) (bool, error)

func GenerateAccessToken() (string, error) {
	cfg := config.Config{}
	err := env.Parse(&cfg)
	if err != nil {
		return "", fmt.Errorf("error parsing environment variable: %v", err)
	}
	id, err := uuid.NewRandom()
	if err != nil {
		return "", fmt.Errorf("error in GenerateAccessToken (uuid.NewRandom): %v", err)
	}
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, &tokenClaims{
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(accessTokenTTL).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
		id,
	})
	access, err := accessToken.SignedString([]byte(cfg.SigningKey))
	return access, err
}

func GenerateRefreshToken() (string, error) {
	cfg := config.Config{}
	err := env.Parse(&cfg)
	if err != nil {
		return "", fmt.Errorf("error parsing environment variable: %v", err)
	}
	id, err := uuid.NewRandom()
	if err != nil {
		return "", fmt.Errorf("error in GenerateAccessToken (uuid.NewRandom): %v", err)
	}
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, &tokenClaims{
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(refreshTokenTTL).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
		id,
	})
	refresh, err := refreshToken.SignedString([]byte(cfg.SigningKey))
	return refresh, err
}
