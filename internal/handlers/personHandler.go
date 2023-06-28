package handlers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/eugenshima/myapp/internal/model"
	"github.com/labstack/echo/v4"
)

// PersonHandler struct contains service service.PersonService
type PersonHandlerImpl struct {
	srv PersonService
}

// NewPersonHandler is a constructor
func NewPersonHandler(srv PersonService) *PersonHandlerImpl {
	return &PersonHandlerImpl{srv: srv}
}

type PersonService interface {
	GetByName(c echo.Context, Name string) (*model.Person, error)
	GetAll(c echo.Context) ([]model.Person, error)
	Delete(c echo.Context, uuidString string) error
	Create(c echo.Context, entity *model.Person) error
	Update(c echo.Context, uuidString string, entity *model.Person) error
}

// GetByName function receives GET request from client
func (handler *PersonHandlerImpl) GetByName(c echo.Context) error {
	Name := c.Param("name")

	result, err := handler.srv.GetByName(c, Name)
	if err != nil {
		str := fmt.Sprintf("Error in handler: %v", err)
		return c.String(http.StatusNotFound, str)
	}
	return c.JSON(http.StatusOK, result)
}

// GetAll function receives GET request from client
func (handler *PersonHandlerImpl) GetAll(c echo.Context) error {
	var results []model.Person

	results, err := handler.srv.GetAll(c)
	if err != nil {
		str := fmt.Sprintf("Error in handler: %v", err)
		return c.String(http.StatusNotFound, str)
	}

	return c.JSON(http.StatusOK, results)
}

// Delete function receives DELETE request from client
func (handler *PersonHandlerImpl) Delete(c echo.Context) error {
	id := c.Param("id")
	err := handler.srv.Delete(c, id)
	if err != nil {
		return err
	}
	return c.String(http.StatusOK, "delete request")
}

// Insert function receives POST request from client
func (handler *PersonHandlerImpl) Create(c echo.Context) error {
	req := c.Request()

	// get request body
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return err
	}

	var entity *model.Person
	err = json.Unmarshal(body, &entity)
	if err != nil {
		return c.String(http.StatusNotFound, "Bad")
	}
	err = handler.srv.Create(c, entity)
	if err != nil {
		return err
	}
	return c.String(http.StatusOK, "insert request")
}

// Update function receives PATCH request from client
func (handler *PersonHandlerImpl) Update(c echo.Context) error {
	id := c.Param("id")
	req := c.Request()

	// get request body
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return err
	}
	var entity *model.Person
	err = json.Unmarshal(body, &entity)
	if err != nil {
		return c.String(http.StatusNotFound, "Bad")
	}
	err = handler.srv.Update(c, id, entity)
	if err != nil {
		return err
	}
	return c.String(http.StatusOK, "update request")
}
