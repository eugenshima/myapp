// Package middleware will be called when the request is successful and the response is returned
package middleware

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
)

// JWTIdentity makes an authorization through access token
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
			if len(headerParts) != 2 || headerParts[0] != "Bearer" {
				return echo.NewHTTPError(http.StatusUnauthorized, "Invalid authorization header format")
			}

			// checking for valid access token
			token, err := ValidateToken(headerParts[1], "gyewgb2rf8r2b8437frb23")
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

// ValidateToken parses tokenString and checks if signing method is ok and return jwt token with filled Valid field
func ValidateToken(tokenString, signingKey string) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(signingKey), nil
	})
	if err != nil {
		return nil, err
	}
	return token, nil
}
