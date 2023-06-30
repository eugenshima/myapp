// Package handlers provides HTTP request handler functions for a web service written in Go
package handlers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/eugenshima/myapp/internal/model"
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
	GenerateToken(c echo.Context, login, password string) (string, error)
	Signup(c echo.Context, entity *model.User) error
	GetAll(c echo.Context) ([]*model.User, error)
}

// UserHandlerImpl represents
type signInInput struct {
	Login    string `db:"login"`
	Password string `db:"password"`
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

	token, err := handler.srv.GenerateToken(c, input.Login, input.Password)
	if err != nil {
		logrus.Errorf("error Generating JWT token %v", err)
		str := fmt.Sprintf("Error in userHandler: %v", err)
		return c.String(http.StatusNotFound, str)
	}
	return c.JSON(http.StatusOK, token)
}

// Signup receives a POST request from client to sign up a user
func (handler *UserHandlerImpl) Signup(c echo.Context) error {
	req := c.Request()

	// get request body
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return err
	}

	var entity *model.User
	err = json.Unmarshal(body, &entity)
	if err != nil {
		logrus.Errorf("error Unmarshalling...    %v", err)
		str := fmt.Sprintf("error in handler: %v", err)
		return c.String(http.StatusNotFound, str)
	}
	entity.ID = uuid.New()
	err = handler.srv.Signup(c, entity)
	if err != nil {
		logrus.Errorf("error calling Signup method %v", err)
		str := fmt.Sprintf("Error in userHandler: %v", err)
		return c.String(http.StatusNotFound, str)
	}
	return c.JSON(http.StatusOK, "Creared")
}

func (handler *UserHandlerImpl) GetAll(c echo.Context) error {
	results, err := handler.srv.GetAll(c)
	if err != nil {
		str := fmt.Sprintf("Error in userHandler: %v", err)
		return c.String(http.StatusNotFound, str)
	}
	return c.JSON(http.StatusOK, results)
}
