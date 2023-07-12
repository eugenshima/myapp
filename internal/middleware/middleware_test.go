package middleware

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	cfgration "github.com/eugenshima/myapp/internal/config"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/require"
)

func TestUserIdentity(t *testing.T) {
	e := echo.New()
	cfg, err := cfgration.NewConfig()
	if err != nil {
		fmt.Printf("Error extracting env variables: %v", err)
		return
	}
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, &tokenClaims{
		Role: "admin",
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(1 * time.Minute).Unix(),
			IssuedAt:  time.Now().Unix(),
			Id:        uuid.New().String(),
		},
	})
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	tokenString, err := accessToken.SignedString([]byte(cfg.SigningKey))
	require.NoError(t, err)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", tokenString))

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	h := UserIdentity()(func(c echo.Context) error {
		return c.String(http.StatusOK, "Authorized")
	})

	err = h(c)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, rec.Code)
	require.Equal(t, "Authorized", rec.Body.String())
}
