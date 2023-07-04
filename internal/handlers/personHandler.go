package handlers

import (
	"fmt"
	"net/http"

	"github.com/eugenshima/myapp/internal/model"

	"github.com/go-playground/validator"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

// PersonHandlerImpl struct contains service service.PersonService
type PersonHandlerImpl struct {
	srv PersonService
}

// NewPersonHandler is a constructor
func NewPersonHandler(srv PersonService) *PersonHandlerImpl {
	return &PersonHandlerImpl{srv: srv}
}

// PersonService interface, which contains Service methods
type PersonService interface {
	GetByID(c echo.Context, id uuid.UUID) (*model.Person, error)
	GetAll(c echo.Context) ([]model.Person, error)
	Delete(c echo.Context, uuidString uuid.UUID) error
	Create(c echo.Context, entity *model.Person) error
	Update(c echo.Context, uuidString uuid.UUID, entity *model.Person) error
}

// GetByID function receives GET request from client
func (handler *PersonHandlerImpl) GetByID(c echo.Context) error {
	ID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		logrus.Errorf("Error parsing id: %v", err)
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	result, err := handler.srv.GetByID(c, ID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	if err != nil {
		logrus.Errorf("Error in handler: %v", err)
		return c.String(http.StatusNotFound, fmt.Sprintf("Error in handler: %v", err))
	}
	return c.JSON(http.StatusOK, result)
}

// GetAll function receives GET request from client
func (handler *PersonHandlerImpl) GetAll(c echo.Context) error {
	results, err := handler.srv.GetAll(c)
	if err != nil {
		logrus.Errorf("Error in handler: %v", err)
		return c.String(http.StatusNotFound, fmt.Sprintf("Error in handler: %v", err))
	}
	return c.JSON(http.StatusOK, results)
}

// Delete function receives DELETE request from client
func (handler *PersonHandlerImpl) Delete(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		logrus.Errorf("Error in handler: %v", err)
		return c.String(http.StatusNotFound, fmt.Sprintf("Error in handler: %v", err))
	}
	err = handler.srv.Delete(c, id)
	if err != nil {
		logrus.Errorf("Error in handler: %v", err)
		return c.String(http.StatusNotFound, fmt.Sprintf("Error in handler: %v", err))
	}
	return c.String(http.StatusOK, "delete request")
}

// Create function receives POST request from client
func (handler *PersonHandlerImpl) Create(c echo.Context) error {
	var entity *model.Person
	err := c.Bind(&entity)
	if err != nil {
		logrus.Errorf("Error in handler: %v", err)
		return c.String(http.StatusBadRequest, fmt.Sprintf("Error in handler: %v", err))
	}
	entity.ID = uuid.New()

	validate := validator.New()
	if err = validate.Struct(entity); err != nil {
		logrus.Errorf("error in handler: %v", err)
		return c.String(http.StatusNotFound, fmt.Sprintf("Error in handler: %v", err))
	}
	err = handler.srv.Create(c, entity)
	if err != nil {
		logrus.Errorf("Error in handler: %v", err)
		return c.String(http.StatusNotFound, fmt.Sprintf("Error in handler: %v", err))
	}
	return c.String(http.StatusOK, "insert request")
}

// Update function receives PATCH request from client
func (handler *PersonHandlerImpl) Update(c echo.Context) error {
	var entity *model.Person
	err := c.Bind(&entity)
	if err != nil {
		logrus.Errorf("Error in handler: %v", err)
		return c.String(http.StatusBadRequest, fmt.Sprintf("Error in handler: %v", err))
	}
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		logrus.Errorf("Error in handler: %v", err)
		return c.String(http.StatusNotFound, fmt.Sprintf("Error in handler: %v", err))
	}
	validate := validator.New()
	if err = validate.Struct(entity); err != nil {
		logrus.Errorf("error in handler: %v", err)
		return c.String(http.StatusNotFound, fmt.Sprintf("Error in handler: %v", err))
	}
	err = handler.srv.Update(c, id, entity)
	if err != nil {
		logrus.Errorf("Error in handler: %v", err)
		return c.String(http.StatusNotFound, fmt.Sprintf("Error in handler: %v", err))
	}
	return c.String(http.StatusOK, "update request")
}
