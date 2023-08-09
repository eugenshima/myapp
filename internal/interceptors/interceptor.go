// Package interceptors will be called when the request is successful and the response is returned
package interceptors

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/caarlos0/env/v9"
	"github.com/eugenshima/myapp/internal/config"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

// tokenClaims struct consists od JWT claims
type tokenClaims struct {
	Role string `json:"role"`
	jwt.StandardClaims
}

// const for interceptor
const (
	Bearer = "Bearer"
	Admin  = "admin"
)

// Admin Unary Interceptor defines the admin intercept interface
func AdminUnaryInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	cfg := config.Config{}
	err := env.Parse(&cfg)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "Invalid env variable: %v", info.FullMethod)
	}
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, "missing metadata")
	}
	if info.FullMethod == "/userHandler/Login" || info.FullMethod == "/userHandler/SignUp" {
		return handler(ctx, req)
	}
	authHeader, ok := md["authorization"]
	if !ok || len(authHeader) == 0 {
		return nil, status.Errorf(codes.Unauthenticated, "missing authorization header")
	}

	headerParts := strings.Split(authHeader[0], " ")
	if len(headerParts) != 2 || headerParts[0] != Bearer {
		return nil, status.Errorf(codes.Unauthenticated, "invalid authorization header format")
	}

	token, err := ValidateToken(headerParts[1], cfg.SigningKey)
	if err != nil || !token.Valid {
		return nil, status.Errorf(codes.Unauthenticated, "invalid token")
	}

	id, role, err := GetPayloadFromToken(headerParts[1])
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "invalid role")
	}
	if role != Admin && (info.FullMethod == "/userHandler/RefreshTokenPair" || info.FullMethod == "/userHandler/Delete") {
		return nil, status.Errorf(codes.PermissionDenied, "invalid role")
	}
	logrus.Infof("interceptor works as expected, returning Unary Handler\nid of token pair: %v", id)
	return handler(ctx, req)
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

type WrappedStream struct {
	grpc.ServerStream
}

func (w *WrappedStream) RecvMsg(m interface{}) error {
	logrus.Printf("received %T - %v\n", m, m)

	return w.ServerStream.RecvMsg(m)
}

func (w *WrappedStream) SendMsg(m interface{}) error {
	logrus.Printf("sent %T - %v\n", m, m)

	return w.ServerStream.SendMsg(m)
}

func newWrappedStream(s grpc.ServerStream) grpc.ServerStream {
	return &WrappedStream{s}
}

func ServerStreamInterceptor(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
	logrus.Println("Stream Interceptor PRE: ", info.FullMethod)
	logrus.Println("Stream Interceptor POST: ", info.FullMethod)
	return handler(srv, newWrappedStream(ss))
}
