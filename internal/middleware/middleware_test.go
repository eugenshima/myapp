package middleware

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	cfgration "github.com/eugenshima/myapp/internal/config"
	"github.com/stretchr/testify/require"
	"gotest.tools/assert"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

func setupMiddlware() (e *echo.Echo, cfg *cfgration.Config, tokenString string, invalidTokenString string, err error) {
	e = echo.New()
	cfg, err = cfgration.NewConfig()
	if err != nil {
		return nil, nil, "", "", fmt.Errorf("Error extracting env variables: %w", err)
	}
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, &tokenClaims{
		Role: "admin",
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(50 * time.Millisecond).Unix(),
			IssuedAt:  time.Now().Unix(),
			Id:        uuid.New().String(),
		},
	})
	tokenString, err = accessToken.SignedString([]byte(cfg.SigningKey))
	if err != nil {
		return nil, nil, "", "", fmt.Errorf("Error creating token string: %w", err)
	}
	invalidTokenString, err = accessToken.SignedString([]byte("invalidSigningKey"))
	if err != nil {
		return nil, nil, "", "", fmt.Errorf("Error creating token string: %w", err)
	}
	return e, cfg, tokenString, invalidTokenString, nil
}

func TestMain(m *testing.M) {
	e, cfg, tokenString, invalidTokenString, err = setupMiddlware()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	exitVal := m.Run()
	os.Exit(exitVal)
}

var (
	invalidTokenString string
	err                error
	tokenString        string
	cfg                *cfgration.Config
	e                  *echo.Echo
)

func TestUserIdentity(t *testing.T) {
	e.Use(UserIdentity())

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "OK")
	})

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Authorization", "Bearer "+tokenString)
	rec := httptest.NewRecorder()

	e.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestAdminIdentity(t *testing.T) {
	e.Use(AdminIdentity())

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "OK")
	})

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Authorization", "Bearer "+tokenString)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	require.Equal(t, http.StatusOK, rec.Code)
}

func TestRoleValidation(t *testing.T) {
	validRole, err := RoleValidation(tokenString)
	require.True(t, validRole)
	require.NoError(t, err)
}

func TestInvalidRoleValidation(t *testing.T) {
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, &tokenClaims{
		Role: "invalid",
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(50 * time.Millisecond).Unix(),
			IssuedAt:  time.Now().Unix(),
			Id:        uuid.New().String(),
		},
	})
	invalidRoleToken, err := accessToken.SignedString([]byte(cfg.SigningKey))
	require.NoError(t, err)
	invalidRole, err := RoleValidation(invalidRoleToken)
	require.False(t, invalidRole)
	require.Error(t, err)
}

func TestValidateToken(t *testing.T) {
	token, err := ValidateToken(tokenString, cfg.SigningKey)
	require.NotNil(t, token)
	require.NoError(t, err)
}

func TestValidateWrongToken(t *testing.T) {
	invalidToken, err := ValidateToken(invalidTokenString, cfg.SigningKey)
	require.Nil(t, invalidToken)
	require.Error(t, err)
}

func TestGetPayloadFromToken(t *testing.T) {
	id, role, err := GetPayloadFromToken(tokenString)
	require.NotEqual(t, "", id)
	require.NotEqual(t, "", role)
	require.NoError(t, err)
}

func TestMiddlewareWithoutAuthHeader(t *testing.T) {
	e.Use(UserIdentity())

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "OK")
	})

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()

	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusUnauthorized, rec.Code)
}

func TestMiddlewareInvalidTokenFormat(t *testing.T) {
	e.Use(UserIdentity())
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "OK")
	})

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Authorization", "notBearer "+tokenString)
	rec := httptest.NewRecorder()

	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusUnauthorized, rec.Code)
}

func TestMiddlewareInvalidToken(t *testing.T) {
	e.Use(UserIdentity())

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "OK")
	})

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Authorization", "Bearer "+invalidTokenString)
	rec := httptest.NewRecorder()

	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusUnauthorized, rec.Code)
}

func TestMiddlewareExpiredToken(t *testing.T) {
	e.Use(UserIdentity())

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "OK")
	})

	time.Sleep(1000 * time.Millisecond)
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Authorization", "Bearer "+tokenString)

	rec := httptest.NewRecorder()

	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusUnauthorized, rec.Code)
}
