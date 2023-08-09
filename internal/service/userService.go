// Package service provides a set of functions, which include business-logic in it
package service

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/eugenshima/myapp/internal/config"
	mdlwr "github.com/eugenshima/myapp/internal/interceptors"
	"github.com/eugenshima/myapp/internal/model"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

const (
	accessTokenTTL  = 60 * time.Minute
	refreshTokenTTL = 72 * time.Hour
)

// tokenClaims struct contains information about the claims associated with the given token
type tokenClaims struct {
	Role string `json:"role"`
	jwt.StandardClaims
}

// UserService is a struct that contains a reference to the repository interface
type UserService struct {
	rps UserRepository
	rdb UserRepositoryRedis
}

// NewUserServiceImpl creates a new service
func NewUserServiceImpl(rps UserRepository, rdb UserRepositoryRedis) *UserService {
	return &UserService{
		rps: rps,
		rdb: rdb,
	}
}

// UserRepository interface, which contains psql/mongo repository methods
type UserRepository interface {
	GetUser(ctx context.Context, login string) (*model.User, error)
	Signup(context.Context, *model.User) error
	GetAll(context.Context) ([]*model.User, error)
	SaveRefreshToken(ctx context.Context, id uuid.UUID, token []byte) error
	GetRefreshToken(ctx context.Context, id uuid.UUID) ([]byte, error)
	GetRoleByID(ctx context.Context, id uuid.UUID) (string, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

// UserRepositoryRedis interface, which contains redis repository methods
type UserRepositoryRedis interface {
	Set(ctx context.Context, user *model.User) error
	Get(ctx context.Context, id uuid.UUID) (*model.User, error)
	Delete(ctx context.Context, id uuid.UUID) error
	GetRefreshToken(ctx context.Context, id uuid.UUID) ([]byte, error)
	SetRefreshToken(ctx context.Context, id uuid.UUID, token []byte) error
}

// GenerateTokens implements the UserServicePsql interface
func (db *UserService) GenerateTokens(ctx context.Context, login, password string) (accessToken, refreshToken string, err error) {
	cfg, err := config.NewConfig()
	if err != nil {
		return "", "", fmt.Errorf("NewConfig: %w", err)
	}

	// GetUser
	user, err := db.rps.GetUser(ctx, login)
	if err != nil {
		return "", "", fmt.Errorf("GetUser: %w", err)
	}

	// CompareHashAndPassword
	err = bcrypt.CompareHashAndPassword(user.Password, []byte(password))
	if err != nil {
		return "", "", fmt.Errorf("CompareHashAndPassword: %w", err)
	}
	// GenerateAccessToken
	accessToken, refreshToken, err = GenerateAccessAndRefreshTokens(cfg.SigningKey, user.Role, user.ID)
	if err != nil {
		return "", "", fmt.Errorf("GenerateAccessAndRefreshTokens: %w", err)
	}
	// HashRefreshToken
	hashedRefreshToken, err := HashRefreshToken(refreshToken)
	if err != nil {
		return "", "", fmt.Errorf("HashRefreshToken: %w", err)
	}
	// SaveRefreshToken
	err = db.rps.SaveRefreshToken(ctx, user.ID, hashedRefreshToken)
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
	user.RefreshToken = hashedRefreshToken
	user.Login = login
	err = db.rdb.Set(ctx, user)
	if err != nil {
		return "", "", fmt.Errorf("set: %w", err)
	}
	return accessToken, refreshToken, nil
}

// RefreshTokenPair func returns a pair of refresh tokens
func (db *UserService) RefreshTokenPair(ctx context.Context, accessToken, refreshToken string, id uuid.UUID) (access, refresh string, err error) {
	cfg, err := config.NewConfig()
	if err != nil {
		return "", "", fmt.Errorf("NewConfig: %w", err)
	}
	// Get RefreshToken
	savedRefreshToken, err := db.rdb.GetRefreshToken(ctx, id) // from cache
	if err != nil || savedRefreshToken == nil {
		savedRefreshToken, err = db.rps.GetRefreshToken(ctx, id) // from database
		if err != nil {
			return "", "", fmt.Errorf("GetRefreshToken: %w", err)
		}
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
	tokenID, role, err := mdlwr.GetPayloadFromToken(accessToken)
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
	access, refresh, err = GenerateAccessAndRefreshTokens(cfg.SigningKey, role, tokenID)
	if err != nil {
		return "", "", fmt.Errorf("GenerateAccessAndRefreshTokens: %w", err)
	}
	// HashRefreshToken
	hashedRefreshToken, err = HashRefreshToken(refresh)
	if err != nil {
		return "", "", fmt.Errorf("HashRefreshToken: %w", err)
	}
	// SaveRefreshToken
	err = db.rps.SaveRefreshToken(ctx, id, hashedRefreshToken)
	if err != nil {
		return "", "", fmt.Errorf("SaveRefreshToken: %w", err)
	}
	err = db.rdb.SetRefreshToken(ctx, id, hashedRefreshToken)
	if err != nil {
		return "", "", fmt.Errorf("SetRefreshToken: %w", err)
	}
	return access, refresh, nil
}

// Signup implements the UserServicePsql interface
func (db *UserService) Signup(ctx context.Context, user *model.User) error {
	hashedPassword := hashPassword(user.Password)
	user.Password = hashedPassword
	err := db.rdb.Set(ctx, user)
	if err != nil {
		return fmt.Errorf("set: %w", err)
	}
	return db.rps.Signup(ctx, user)
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
	hash := sha256.New()
	hash.Write([]byte(refreshToken))
	hashBytes := hash.Sum(nil)
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

// Delete calls delete method from repository level
func (db *UserService) Delete(ctx context.Context, id uuid.UUID) error {
	err := db.rdb.Delete(ctx, id)
	if err != nil {
		return db.rps.Delete(ctx, id)
	}
	return err
}
