// Package service provides a set of functions, which include business-logic in it
package service

import (
	"context"
	"crypto/sha256"
	"fmt"
	"time"

	"github.com/eugenshima/myapp/internal/config"
	"github.com/eugenshima/myapp/internal/model"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

const (
	accessTokenTTL  = 15 * time.Minute
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

// GenerateTokens implements the UserServicePsql interface
func (db *UserServiceImpl) GenerateTokens(ctx context.Context, login, password string) (accessToken string, refreshToken string, err error) {
	cfg, err := config.NewConfig()
	if err != nil {
		return "", "", fmt.Errorf("error creating config: %v", err)
	}

	// GetUser
	id, pass, err := db.rps.GetUser(ctx, login)
	if err != nil {
		return "", "", fmt.Errorf("error in GenerateToken (GetUser): %v", err)
	}
	// CompareHashAndPassword
	err = bcrypt.CompareHashAndPassword(pass, []byte(password))
	if err != nil {
		return "", "", fmt.Errorf("error in GenerateToken (CompareHashAndPassword): %v", err)
	}
	// GenerateAccessToken
	accessToken, refreshToken, err = GenerateAccessAndRefreshTokens(cfg.SigningKey)
	if err != nil {
		return "", "", fmt.Errorf("error in GenerateToken (GenerateAccessAndRefreshTokens): %v", err)
	}
	// HashRefreshToken
	hashedRefreshToken, err := HashRefreshToken(refreshToken)
	if err != nil {
		return "", "", fmt.Errorf("error in refresh token: %v", err)
	}
	// SaveRefreshToken
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

// HashPassword func returns hashed password using bcrypt algorithm
func hashPassword(password []byte) []byte {
	fmt.Println(password)
	hashedPassword, err := bcrypt.GenerateFromPassword(password, 14)

	if err != nil {
		return nil
	}
	return hashedPassword
}

// HashRefreshToken func returns hashed refresh token using bcrypt algorithm
func HashRefreshToken(refreshToken string) ([]byte, error) {
	sum := sha256.Sum256([]byte(refreshToken))
	hashedRefreshToken := hashPassword(sum[:])

	return hashedRefreshToken, nil
}

// RefreshAccessToken func recreates access token
func RefreshAccessToken(refreshToken string, key string) (string, error) {
	// Refresh token validation
	valid, err := ValidateRefreshToken(refreshToken)
	if err != nil {
		return "", fmt.Errorf("error validating refresh token: %v", err)
	}
	if !valid {
		return "", fmt.Errorf("invalid refresh token : %v", err)
	}

	// Generating new access token
	accessToken, err := GenerateAccessToken(key)
	if err != nil {
		return "", fmt.Errorf("error generating access token: %v", err)
	}

	return accessToken, nil
}

// ValidateRefreshToken func validates refresh token
func ValidateRefreshToken(refreshToken string) (bool, error) {
	//TODO: validate refresh token
	return true, nil
}

// GenerateAccessAndRefreshTokens func returns access & refresh tokens
func GenerateAccessAndRefreshTokens(key string) (access string, refresh string, err error) {
	id, err := uuid.NewRandom()
	if err != nil {
		return "", "", fmt.Errorf("error in GenerateAccessToken (uuid.NewRandom): %v", err)
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
	access, err = accessToken.SignedString([]byte(key))
	if err != nil {
		return "", "", fmt.Errorf("error in GenerateAccessToken (accessToken.SignedString): %v", err)
	}
	refresh, err = refreshToken.SignedString([]byte(key))
	if err != nil {
		return "", "", fmt.Errorf("error in GenerateAccessToken (refreshToken.SignedString): %v", err)
	}
	return access, refresh, err
}

// GenerateAccessToken for signing requests
func GenerateAccessToken(key string) (string, error) {

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
	access, err := accessToken.SignedString([]byte(key))
	return access, err
}
