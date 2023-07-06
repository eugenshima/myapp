// Package handlers provides HTTP request handler functions for a web service written in Go
package handlers

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/eugenshima/myapp/internal/model"

	"github.com/go-playground/validator"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

// UserHandlerImpl struct represents a user handler implementation
type UserHandlerImpl struct {
	srv UserService
}

// NewUserHandlerImpl creates a new Handler
func NewUserHandlerImpl(srv UserService) *UserHandlerImpl {
	return &UserHandlerImpl{
		srv: srv,
	}
}

// UserService interface implementation
type UserService interface {
	GenerateTokens(ctx context.Context, login, password string) (string, string, error)
	Signup(ctx context.Context, entity *model.User) error
	RefreshTokenPair(ctx context.Context, accessToken string, refreshToken string, id uuid.UUID) (string, string, error)
	GetAll(ctx context.Context) ([]*model.User, error)
}

// Login receives a GET request from client and returns a user(if exists)
func (handler *UserHandlerImpl) Login(c echo.Context) error {
	input := model.Login{}
	err := c.Bind(&input)
	if err != nil {
		logrus.Errorf("Error in userHandler: %v", err)
		return c.String(http.StatusNotFound, fmt.Sprintf("Error in userHandler: %v", err))
	}

	accessToken, refreshToken, err := handler.srv.GenerateTokens(c.Request().Context(), input.Login, input.Password)
	if err != nil {
		logrus.Errorf("error Generating JWT token %v", err)
		return c.String(http.StatusNotFound, fmt.Sprintf("Error in userHandler: %v", err))
	}
	response := map[string]interface{}{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	}
	return c.JSON(http.StatusOK, response)
}

// Signup receives a POST request from client to sign up a user
func (handler *UserHandlerImpl) Signup(c echo.Context) error {
	reqBody := model.Signup{}
	err := c.Bind(&reqBody)
	if err != nil {
		logrus.Errorf("Error in userHandler: %v", err)
		return c.String(http.StatusNotFound, fmt.Sprintf("Error in userHandler: %v", err))
	}
	entity := &model.User{
		ID:       uuid.New(),
		Login:    reqBody.Login,
		Password: []byte(reqBody.Password),
		Role:     reqBody.Role,
	}

	validate := validator.New()
	if err = validate.Struct(entity); err != nil {
		logrus.Errorf("error in handler: %v", err)
		str := fmt.Sprintf("Error in handler: %v", err)
		return c.String(http.StatusBadRequest, str)
	}

	err = handler.srv.Signup(c.Request().Context(), entity)
	if err != nil {
		logrus.Errorf("error calling Signup method: %v", err)
		return c.String(http.StatusInternalServerError, fmt.Sprintf("Error in userHandler: %v", err))
	}

	return c.JSON(http.StatusOK, "Created")
}

// GetAll receives a POST request from client for getting all entities from the server
func (handler *UserHandlerImpl) GetAll(c echo.Context) error {
	results, err := handler.srv.GetAll(c.Request().Context())
	if err != nil {
		logrus.Errorf("error calling GetAll method: %v", err)
		return c.String(http.StatusNotFound, fmt.Sprintf("Error in userHandler: %v", err))
	}
	return c.JSON(http.StatusOK, results)
}

// RefreshTokenPair receives a POST request from client for refreshing a token pair
func (handler *UserHandlerImpl) RefreshTokenPair(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		logrus.Errorf("error calling RefreshTokenPair method: %v", err)
		return c.String(http.StatusBadRequest, fmt.Sprintf("Error in userHandler: %v", err))
	}
	reqBody := model.Tokens{}
	err = c.Bind(&reqBody)
	if err != nil {
		logrus.Errorf("error calling RefreshTokenPair method: %v", err)
		return c.String(http.StatusNotFound, fmt.Sprintf("Error in userHandler: %v", err))
	}

	accessToken, refreshToken, err := handler.srv.RefreshTokenPair(c.Request().Context(), reqBody.AccessToken, reqBody.RefreshToken, id)
	if err != nil {
		logrus.Errorf("error calling RefreshTokenPair method: %v", err)
		return c.String(http.StatusNotFound, fmt.Sprintf("Error in userHandler: %v", err))
	}
	response := map[string]interface{}{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	}
	return c.JSON(http.StatusOK, response)
}

type image struct {
	Filename string `json:"filename"`
	Url      string `json:"url"`
}

func (handler *UserHandlerImpl) GetImage(c echo.Context) error {
	image := c.Param("name")
	filePath := "/home/yauhenishymanski/MyProject/myapp/internal/images/" + image
	return c.Inline(filePath, image)
}

// SetImage saves the image from the internet
func (handler *UserHandlerImpl) SetImage(c echo.Context) error {
	img := image{}
	err := c.Bind(&img)
	if err != nil {
		logrus.Errorf("error in handler: %v", err)
		return c.String(http.StatusBadRequest, fmt.Sprintf("Error in handler: %v", err))
	}
	response, err := http.Get(img.Url)
	if err != nil {
		logrus.Errorf("error in handler: %v", err)
		return c.String(http.StatusBadRequest, fmt.Sprintf("Error in handler: %v", err))
	}
	defer response.Body.Close()

	filename := filepath.Join("/home/yauhenishymanski/MyProject/myapp/internal/images", img.Filename)
	file, err := os.Create(filename)
	if err != nil {
		logrus.Errorf("error in handler: %v", err)
		return c.String(http.StatusBadRequest, fmt.Sprintf("Error in handler: %v", err))
	}
	defer file.Close()

	_, err = io.Copy(file, response.Body)
	if err != nil {
		logrus.Errorf("error in handler: %v", err)
		return c.String(http.StatusBadRequest, fmt.Sprintf("Error in handler: %v", err))
	}
	return c.String(http.StatusOK, "image has been set")
}
