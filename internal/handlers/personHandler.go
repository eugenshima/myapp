package handlers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/eugenshima/myapp/internal/model"
	"github.com/google/uuid"
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
	GetById(c echo.Context, id uuid.UUID) (*model.Person, error)
	GetAll(c echo.Context) ([]model.Person, error)
	Delete(c echo.Context, uuidString uuid.UUID) error
	Create(c echo.Context, entity *model.Person) error
	Update(c echo.Context, uuidString uuid.UUID, entity *model.Person) error
}

// GetById function receives GET request from client
func (handler *PersonHandlerImpl) GetById(c echo.Context) error {
	Id := c.Param("id")
	var entity model.Person
	var err error
	entity.ID, err = uuid.Parse(Id)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	result, err := handler.srv.GetById(c, entity.ID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	if err != nil {
		str := fmt.Sprintf("Error in handler: %v", err)
		return c.String(http.StatusNotFound, str)
	}
	return c.JSON(http.StatusOK, result)
}

// GetAll function receives GET request from client
func (handler *PersonHandlerImpl) GetAll(c echo.Context) error {
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
	var entity model.Person
	var err error
	entity.ID, err = uuid.Parse(id)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	err = handler.srv.Delete(c, entity.ID)
	if err != nil {
		str := fmt.Sprintf("Error in handler: %v", err)
		return c.String(http.StatusNotFound, str)
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
		str := fmt.Sprintf("Error in handler: %v", err)
		return c.String(http.StatusNotFound, str)
	}
	entity.ID = uuid.New()
	err = handler.srv.Create(c, entity)
	if err != nil {
		str := fmt.Sprintf("Error in handler: %v", err)
		return c.String(http.StatusNotFound, str)
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
	var entity model.Person
	entity.ID, err = uuid.Parse(id)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	err = json.Unmarshal(body, &entity)
	if err != nil {
		return c.String(http.StatusNotFound, "Bad")
	}

	err = handler.srv.Update(c, entity.ID, &entity)
	if err != nil {
		return err
	}
	return c.String(http.StatusOK, "update request")
}
