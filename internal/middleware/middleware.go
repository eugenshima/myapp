// Package middleware will be called when the request is successful and the response is returned
package middleware

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/eugenshima/myapp/internal/config"
	"github.com/google/uuid"

	"github.com/caarlos0/env/v9"
	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
)

// tokenClaims struct consists od JWT claims
type tokenClaims struct {
	Role string `json:"role"`
	jwt.StandardClaims
}

// const for middlware
const (
	Bearer = "Bearer"
	Admin  = "admin"
)

// UserIdentity makes an authorization through access token
func UserIdentity() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Chtcking for auth header
			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" {
				return echo.NewHTTPError(http.StatusUnauthorized, "Missing authorization header")
			}
			// checking for auth header format
			headerParts := strings.Split(authHeader, " ")
			if len(headerParts) != 2 || headerParts[0] != Bearer {
				return echo.NewHTTPError(http.StatusUnauthorized, "Invalid authorization header format")
			}
			// getting environment variable
			cfg := config.Config{}
			err := env.Parse(&cfg)
			if err != nil {
				return echo.NewHTTPError(http.StatusUnauthorized, "Invalid env variable")
			}
			// checking for valid access token
			token, err := ValidateToken(headerParts[1], cfg.SigningKey)
			if err != nil || !token.Valid {
				return echo.NewHTTPError(http.StatusUnauthorized, "Invalid token")
			}
			// checking for token expiration
			if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
				exp := claims["exp"].(float64)
				if exp < float64(time.Now().Unix()) {
					return echo.NewHTTPError(http.StatusUnauthorized, "Token is expired")
				}
			}

			return next(c)
		}
	}
}

// AdminIdentity makes an authorization through access token for admin only
func AdminIdentity() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Chtcking for auth header
			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" {
				return echo.NewHTTPError(http.StatusUnauthorized, "Missing authorization header")
			}
			// checking for auth header format
			headerParts := strings.Split(authHeader, " ")
			if len(headerParts) != 2 || headerParts[0] != Bearer {
				return echo.NewHTTPError(http.StatusUnauthorized, "Invalid authorization header format")
			}
			// getting environment variable
			cfg := config.Config{}
			err := env.Parse(&cfg)
			if err != nil {
				return echo.NewHTTPError(http.StatusUnauthorized, "Invalid env variable")
			}
			// checking for valid access token
			token, err := ValidateToken(headerParts[1], cfg.SigningKey)
			if err != nil || !token.Valid {
				return echo.NewHTTPError(http.StatusUnauthorized, "Invalid token")
			}
			id, role, err := GetPayloadFromToken(headerParts[1])
			fmt.Println("middleware id working --> ", id)
			if err != nil {
				return err
			}
			if role != Admin {
				return echo.NewHTTPError(http.StatusUnauthorized, "Invalid role")
			}
			// checking for token expiration
			if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
				exp := claims["exp"].(float64)
				if exp < float64(time.Now().Unix()) {
					return echo.NewHTTPError(http.StatusUnauthorized, "Token is expired")
				}
			}
			return next(c)
		}
	}
}

// RoleValidation is used to validate the role
func RoleValidation(tokenString string) (bool, error) {
	parts := strings.Split(tokenString, ".")
	payload := parts[1]

	payloadBytes, err := base64.RawURLEncoding.DecodeString(payload)
	if err != nil {
		return false, fmt.Errorf("DecodeString: %w", err)
	}

	var claims tokenClaims
	err = json.Unmarshal(payloadBytes, &claims)
	if err != nil {
		return false, fmt.Errorf("Unmarshal(): %w", err)
	}

	role := claims.Role
	if role != Admin {
		return false, fmt.Errorf("invalid role: %w", err)
	}
	return true, nil
}

// ValidateToken parses tokenString and returns valid jwt token string
func ValidateToken(tokenString, signingKey string) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Parse(): %v", token.Header["alg"])
		}
		return []byte(signingKey), nil
	})
	if err != nil {
		return nil, fmt.Errorf("Parse(): %w", err)
	}
	return token, nil
}

// GetPayloadFromToken returns a payload from the given token
func GetPayloadFromToken(token string) (uuid.UUID, string, error) {
	parts := strings.Split(token, ".")
	payload := parts[1]

	// Декодирование Base64url полезной нагрузки в формат JSON
	payloadBytes, err := base64.RawURLEncoding.DecodeString(payload)
	if err != nil {
		return uuid.Nil, "", fmt.Errorf("DecodeString: %w", err)
	}

	// Распаковка полезной нагрузки в структуру CustomClaims
	var claims tokenClaims
	err = json.Unmarshal(payloadBytes, &claims)
	if err != nil {
		return uuid.Nil, "", fmt.Errorf("Unmarshal(): %w", err)
	}

	// Получение значения ролей
	role := claims.Role
	id, err := uuid.Parse(claims.Id)
	if err != nil {
		return uuid.Nil, "", fmt.Errorf("Parse(): %w", err)
	}
	return id, role, nil
}
