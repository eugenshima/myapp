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

	vl "github.com/go-playground/validator"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

// UserHandler struct represents a user handler implementation
type UserHandler struct {
	srv       UserService
	validator *vl.Validate
}

// NewUserHandlerImpl creates a new Handler
func NewUserHandlerImpl(srv UserService, validator *vl.Validate) *UserHandler {
	return &UserHandler{
		srv:       srv,
		validator: validator,
	}
}

//go:generate mockgen -source=userHandler.go -destination=mocks/userMock.go

// UserService interface implementation
type UserService interface {
	GenerateTokens(ctx context.Context, login, password string) (string, string, error)
	Signup(ctx context.Context, entity *model.User) error
	RefreshTokenPair(ctx context.Context, accessToken string, refreshToken string, id uuid.UUID) (string, string, error)
	GetAll(ctx context.Context) ([]*model.User, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

// Login receives a GET request from client and returns a user(if exists)
// @Summary Login user
// @tags authentication methods
// @Description Logs in a user and returns access and refresh tokens
// @Accept json
// @Produce json
// @Param input body model.Login true "Login details"
// @Success 200 {object} map[string]interface{} " Generating access and refresh tokens"
// @Failure 404 {string} string "Error message"
// @Router /api/user/login [post]
func (handler *UserHandler) Login(c echo.Context) error {
	input := model.Login{}
	err := c.Bind(&input)
	if err != nil {
		logrus.WithFields(logrus.Fields{"input": &input}).Errorf("Bind: %v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("Bind: %v", err))
	}

	err = c.Validate(input)
	if err != nil {
		logrus.WithFields(logrus.Fields{"input": input}).Errorf("Validate: %v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("Validate: %v", err))
	}
	accessToken, refreshToken, err := handler.srv.GenerateTokens(c.Request().Context(), input.Login, input.Password)
	if err != nil {
		logrus.Errorf("GenerateTokens %v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("GenerateTokens: %v", err))
	}
	response := map[string]string{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	}
	return c.JSON(http.StatusOK, response)
}

// Signup receives a POST request from client to sign up a user
// @Summary Sign up user
// @tags authentication methods
// @Description Registers a new user
// @Accept json
// @Produce json
// @Param reqBody body model.Signup true "Signup details"
// @Success 200 {string} string "User created"
// @Failure 400 {string} string "Error message"
// @Failure 500 {string} string "Internal server error"
// @Router /api/user/signup [post]
func (handler *UserHandler) Signup(c echo.Context) error {
	reqBody := model.Signup{}
	err := c.Bind(&reqBody)
	if err != nil {
		logrus.Errorf("Bind: %v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("Bind: %v", err))
	}
	person := &model.User{
		ID:       uuid.New(),
		Login:    reqBody.Login,
		Password: []byte(reqBody.Password),
		Role:     reqBody.Role,
	}

	err = c.Validate(person)
	if err != nil {
		logrus.Errorf("Validate: %v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("Validate: %v", err))
	}
	err = handler.srv.Signup(c.Request().Context(), person)
	if err != nil {
		logrus.WithFields(logrus.Fields{"person": person}).Errorf("Signup: %v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("Signup: %v", err))
	}

	return c.JSON(http.StatusOK, "Created")
}

// GetAll receives a POST request from client for getting all entities from the server
func (handler *UserHandler) GetAll(c echo.Context) error {
	results, err := handler.srv.GetAll(c.Request().Context())
	if err != nil {
		logrus.Errorf("GetAll: %v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("GetAll: %v", err))
	}
	return c.JSON(http.StatusOK, results)
}

// RefreshTokenPair receives a POST request from client for refreshing a token pair
// @Summary Refresh token pair
// @tags authentication methods
// @Description Refreshes an access token and a refresh token
// @Accept json
// @Produce json
// @Param id path string true "ID of the user"
// @Param reqBody body model.Tokens true "Token pair details"
// @Success 200 {object} map[string]interface{} "Refreshed token pair"
// @Failure 400 {string} string "Bad request"
// @Failure 404 {string} string "Error message"
// @Router /api/user/refresh/{id} [post]
func (handler *UserHandler) RefreshTokenPair(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		logrus.Errorf("Parse: %v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("Parse: %v", err))
	}
	reqBody := model.Tokens{}
	err = c.Bind(&reqBody)
	if err != nil {
		logrus.Errorf("Bind: %v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("Bind: %v", err))
	}

	accessToken, refreshToken, err := handler.srv.RefreshTokenPair(c.Request().Context(), reqBody.AccessToken, reqBody.RefreshToken, id)
	if err != nil {
		logrus.WithFields(logrus.Fields{"reqBody.AccessToken": reqBody.AccessToken, "reqBody.RefreshToken": reqBody.RefreshToken, "id": id}).Errorf("RefreshTokenPair: %v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("RefreshTokenPair: %v", err))
	}
	response := map[string]string{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	}
	return c.JSON(http.StatusOK, response)
}

// Delete func receives a path variable abd return deleted id (if exists)
func (handler *UserHandler) Delete(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		logrus.Errorf("Parse: %v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("Parse: %v", err))
	}
	err = handler.srv.Delete(c.Request().Context(), id)
	if err != nil {
		logrus.WithFields(logrus.Fields{"id": id}).Errorf("Delete: %v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("Delete: %v", err))
	}
	return c.String(http.StatusOK, "OK")
}

// image struct for downlodaing images from internet
type image struct {
	Filename string `json:"filename"`
	URL      string `json:"url"`
}

// GetImage returns an image from database
// @Summary Get image by name
// @Security ApiKeyAuth
// @tags download/upload images
// @Description Retrieves an image by name
// @Produce octet-stream
// @Param name path string true "Name of the image"
// @Success 200 {file} file "Image file"
// @Failure 404 {string} string "Image not found"
// @Router /api/image/get/{name} [get]
func (handler *UserHandler) GetImage(c echo.Context) error {
	image := c.Param("name")
	filePath := "/home/yauhenishymanski/MyProject/myapp/internal/images/" + image
	return c.Inline(filePath, image)
}

// SetImage saves the image from the internet
// @Summary Set image
// @Security ApiKeyAuth
// @tags download/upload images
// @Description Sets an image from the provided URL
// @Accept json
// @Produce plain
// @Param img body image true "Image details"
// @Success 200 {string} string "Image has been set"
// @Failure 400 {string} string "Bad request"
// @Failure 404 {string} string "Error message"
// @Router /api/image/set [post]
func (handler *UserHandler) SetImage(c echo.Context) error {
	img := image{}
	err := c.Bind(&img)
	if err != nil {
		logrus.Errorf("Bind: %v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("Bind: %v", err))
	}
	response, err := http.Get(img.URL)
	if err != nil {
		logrus.Errorf("Get: %v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("Get: %v", err))
	}
	defer func() {
		if err = response.Body.Close(); err != nil {
			logrus.Errorf("Close: %v", err)
		}
	}()
	rootpath := "/home/yauhenishymanski/MyProject/myapp"
	basedir := "internal"
	subdir := "images"
	filename := filepath.Join(rootpath+string(filepath.Separator)+basedir+string(filepath.Separator)+subdir+string(filepath.Separator), img.Filename)
	file, err := os.Create(filename)
	if err != nil {
		logrus.Errorf("Create: %v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("Create: %v", err))
	}
	defer func() {
		if err = response.Body.Close(); err != nil {
			logrus.Errorf("Close: %v", err)
		}
	}()

	_, err = io.Copy(file, response.Body)
	if err != nil {
		logrus.Errorf("Copy: %v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("Copy: %v", err))
	}
	return c.String(http.StatusOK, "image has been set")
}
