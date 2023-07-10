// Package service provides a set of functions, which include business-logic in it
package service

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/eugenshima/myapp/internal/config"
	mdlwr "github.com/eugenshima/myapp/internal/middleware"
	"github.com/eugenshima/myapp/internal/model"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

const (
	accessTokenTTL  = 1 * time.Hour
	refreshTokenTTL = 72 * time.Hour
)

// tokenClaims struct contains information about the claims associated with the given token
type tokenClaims struct {
	Role string `json:"role"`
	jwt.StandardClaims
}

// UserService is a struct that contains a reference to the repository interface
type UserService struct {
	rps UserRepositoryPsql
}

// NewUserServiceImpl creates a new service
func NewUserServiceImpl(rps UserRepositoryPsql) *UserService {
	return &UserService{
		rps: rps,
	}
}

// UserRepositoryPsql interface, which contains repository methods
type UserRepositoryPsql interface {
	GetUser(ctx context.Context, login string) (uuid.UUID, []byte, string, error)
	Signup(context.Context, *model.User) error
	GetAll(context.Context) ([]*model.User, error)
	SaveRefreshToken(ctx context.Context, id uuid.UUID, token []byte) error
	GetRefreshToken(ctx context.Context, id uuid.UUID) ([]byte, error)
	GetRoleByID(ctx context.Context, id uuid.UUID) (string, error)
}

// GenerateTokens implements the UserServicePsql interface
func (db *UserService) GenerateTokens(ctx context.Context, login, password string) (accessToken, refreshToken string, err error) {
	cfg, err := config.NewConfig()
	if err != nil {
		return "", "", fmt.Errorf("NewConfig: %w", err)
	}

	// GetUser
	id, pass, role, err := db.rps.GetUser(ctx, login)
	if err != nil {
		return "", "", fmt.Errorf("GetUser: %w", err)
	}
	// CompareHashAndPassword
	err = bcrypt.CompareHashAndPassword(pass, []byte(password))
	if err != nil {
		return "", "", fmt.Errorf("CompareHashAndPassword: %w", err)
	}
	// GenerateAccessToken
	accessToken, refreshToken, err = GenerateAccessAndRefreshTokens(cfg.SigningKey, role, id)
	if err != nil {
		return "", "", fmt.Errorf("GenerateAccessAndRefreshTokens: %w", err)
	}
	// HashRefreshToken
	hashedRefreshToken, err := HashRefreshToken(refreshToken)
	if err != nil {
		return "", "", fmt.Errorf("HashRefreshToken: %w", err)
	}
	// SaveRefreshToken
	err = db.rps.SaveRefreshToken(ctx, id, hashedRefreshToken)
	if err != nil {
		return "", "", fmt.Errorf("SaveRefreshToken: %w", err)
	}
	// CompareTokenIDs
	compID, err := CompareTokenIDs(accessToken, refreshToken, cfg.SigningKey)
	if err != nil {
		return "", "", fmt.Errorf("CompareTokenIDs: %w", err)
	}
	if !compID {
		return "", "", fmt.Errorf("invalid token(campare error): %w", err)
	}

	return accessToken, refreshToken, nil
}

func (db *UserService) RefreshTokenPair(ctx context.Context, accessToken, refreshToken string, id uuid.UUID) (access, refresh string, err error) {
	cfg, err := config.NewConfig()
	if err != nil {
		return "", "", fmt.Errorf("NewConfig: %w", err)
	}
	// Get RefreshToken
	savedRefreshToken, err := db.rps.GetRefreshToken(ctx, id)
	if err != nil {
		return "", "", fmt.Errorf("GetRefreshToken: %w", err)
	}
	// HashRefreshToken
	hashedRefreshToken, err := HashRefreshToken(refreshToken)
	if err != nil {
		return "", "", fmt.Errorf("HashRefreshToken: %w", err)
	}
	// CompareHashedTokens
	isEqual := CompareHashedTokens(savedRefreshToken, hashedRefreshToken)
	if !isEqual {
		return "", "", fmt.Errorf("CompareHashedTokens: %w", err)
	}
	id, role, err := mdlwr.GetPayloadFromToken(accessToken)
	if err != nil {
		return "", "", fmt.Errorf("GetPayloadFromToken: %w", err)
	}
	// CompareTokenIDs
	compID, err := CompareTokenIDs(accessToken, refreshToken, cfg.SigningKey)
	if err != nil {
		return "", "", fmt.Errorf("CompareTokenIDs: %w", err)
	}
	if !compID {
		return "", "", fmt.Errorf("invalid token(campare error): %w", err)
	}
	// GenerateAccessAndRefreshTokens
	access, refresh, err = GenerateAccessAndRefreshTokens(cfg.SigningKey, role, id)
	if err != nil {
		return "", "", fmt.Errorf("GenerateAccessAndRefreshTokens: %w", err)
	}
	return access, refresh, nil
}

// Signup implements the UserServicePsql interface
func (db *UserService) Signup(ctx context.Context, entity *model.User) error {
	hashedPassword := hashPassword(entity.Password)
	entity.Password = hashedPassword
	return db.rps.Signup(ctx, entity)
}

// GetAll implements the UserServicePsql interface
func (db *UserService) GetAll(ctx context.Context) ([]*model.User, error) {
	return db.rps.GetAll(ctx)
}

// HashPassword func returns hashed password using bcrypt algorithm
func hashPassword(password []byte) []byte {
	hashedPassword, err := bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)
	if err != nil {
		return nil
	}
	return hashedPassword
}

// HashRefreshToken func returns hashed refresh token using bcrypt algorithm
func HashRefreshToken(refreshToken string) ([]byte, error) {
	// Creating a new SHA-256 hash
	hash := sha256.New()

	// Writing a refresh token to a hash
	hash.Write([]byte(refreshToken))

	// Getting the hashed value as a slice of bytes
	hashBytes := hash.Sum(nil)

	// Convert slice of bytes to hexadecimal string
	hashString := hex.EncodeToString(hashBytes)

	return []byte(hashString), nil
}

// CompareHashedTokens func compairs hashed tokens from database and request
func CompareHashedTokens(token1, token2 []byte) bool {
	return sha256.Sum256(token1) == sha256.Sum256(token2)
}

// CompareTokenIDs func compares token ids
func CompareTokenIDs(accessToken, refreshToken, key string) (bool, error) {
	accessID, err := ExtractIDFromToken(accessToken, key)
	if err != nil {
		return false, fmt.Errorf("ExtractIDFromToken: %w", err)
	}

	refreshID, err := ExtractIDFromToken(refreshToken, key)
	if err != nil {
		return false, fmt.Errorf("ExtractIDFromToken: %w", err)
	}
	return accessID == refreshID, nil
}

// ExtractIDFromToken extracts the identifier (ID) from the payload (claims) of the token.
func ExtractIDFromToken(tokenString, key string) (string, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(key), nil
	})
	if err != nil {
		return "", fmt.Errorf("Parse(): %w", err)
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if id, ok := claims["jti"].(string); ok {
			return id, nil
		}
	}

	return "", fmt.Errorf("error extracting ID from token: %v", token)
}

// GenerateAccessAndRefreshTokens func returns access & refresh tokens
func GenerateAccessAndRefreshTokens(key, role string, id uuid.UUID) (access, refresh string, err error) {
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, &tokenClaims{
		Role: role,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(accessTokenTTL).Unix(),
			IssuedAt:  time.Now().Unix(),
			Id:        id.String(),
		},
	})

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, &tokenClaims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(refreshTokenTTL).Unix(),
			IssuedAt:  time.Now().Unix(),
			Id:        id.String(),
		},
	})
	access, err = accessToken.SignedString([]byte(key))
	if err != nil {
		return "", "", fmt.Errorf("SignedString(access): %w", err)
	}
	refresh, err = refreshToken.SignedString([]byte(key))
	if err != nil {
		return "", "", fmt.Errorf("SignedString(refresh): %w", err)
	}
	return access, refresh, err
}
