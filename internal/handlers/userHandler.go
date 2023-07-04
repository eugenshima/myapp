// Package handlers provides HTTP request handler functions for a web service written in Go
package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/eugenshima/myapp/internal/model"
	"github.com/go-playground/validator"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	//"github.com/sirupsen/logrus"
)

// UserHandlerImpl struct represents a user handler implementation
type UserHandlerImpl struct {
	srv UserService
}

// NewHandler creates a new Handler
func NewUserHandlerImpl(srv UserService) *UserHandlerImpl {
	return &UserHandlerImpl{
		srv: srv,
	}
}

// UserService interface implementation
type UserService interface {
	GenerateToken(ctx context.Context, login, password string) (string, error)
	Signup(ctx context.Context, entity *model.User) error
	GetAll(ctx context.Context) ([]*model.User, error)
	ParseToken(ctx context.Context, accessToken string) (uuid.UUID, error)
}

// UserHandlerImpl represents
type signInInput struct {
	Login    string `db:"login"`
	Password string `db:"password"`
}

type RequestBody struct {
	Login    string `json:"login"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

// Login receives a GET request from client and returns a user(if exists)
func (handler *UserHandlerImpl) Login(c echo.Context) error {
	req := c.Request()

	// get request body
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return err
	}

	var input signInInput
	err = json.Unmarshal(body, &input)
	if err != nil {
		logrus.Errorf("Error unmarshalling...   %v", err)
		str := fmt.Sprintf("Error in userHandler: %v", err)
		return c.String(http.StatusNotFound, str)
	}

	token, err := handler.srv.GenerateToken(c.Request().Context(), input.Login, input.Password)
	fmt.Println("access token -->", token)
	if err != nil {
		logrus.Errorf("error Generating JWT token %v", err)
		str := fmt.Sprintf("Error in userHandler: %v", err)
		return c.String(http.StatusNotFound, str)
	}
	return c.JSON(http.StatusOK, token)
}

// Signup receives a POST request from client to sign up a user
func (handler *UserHandlerImpl) Signup(c echo.Context) error {
	var reqBody RequestBody

	req := c.Request()

	// get request body
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return err
	}

	err = json.Unmarshal(body, &reqBody)
	if err != nil {
		logrus.Errorf("Error unmarshalling...   %v", err)
		str := fmt.Sprintf("Error in userHandler: %v", err)
		return c.String(http.StatusNotFound, str)
	}

	entity := &model.User{
		ID:       uuid.New(),
		Login:    reqBody.Login,
		Password: []byte(reqBody.Password),
		Role:     reqBody.Role,
	}

	validate := validator.New()
	if err := validate.Struct(entity); err != nil {
		logrus.Errorf("error in handler: %v", err)
		str := fmt.Sprintf("Error in handler: %v", err)
		return c.String(http.StatusBadRequest, str)
	}

	err = handler.srv.Signup(c.Request().Context(), entity)
	if err != nil {
		logrus.Errorf("error calling Signup method: %v", err)
		str := fmt.Sprintf("Error in userHandler: %v", err)
		return c.String(http.StatusInternalServerError, str)
	}

	return c.JSON(http.StatusOK, "Created")
}

func (handler *UserHandlerImpl) GetAll(c echo.Context) error {
	results, err := handler.srv.GetAll(c.Request().Context())
	if err != nil {
		str := fmt.Sprintf("Error in userHandler: %v", err)
		return c.String(http.StatusNotFound, str)
	}
	return c.JSON(http.StatusOK, results)
}
