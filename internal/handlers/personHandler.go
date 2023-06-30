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
	ID := c.Param("id")
	var entity model.Person
	var err error
	entity.ID, err = uuid.Parse(ID)
	if err != nil {
		logrus.Errorf("Error parsing id: %v", err)
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	result, err := handler.srv.GetByID(c, entity.ID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	if err != nil {
		logrus.Errorf("Error in handler: %v", err)
		str := fmt.Sprintf("Error in handler: %v", err)
		return c.String(http.StatusNotFound, str)
	}
	return c.JSON(http.StatusOK, result)
}

// GetAll function receives GET request from client
func (handler *PersonHandlerImpl) GetAll(c echo.Context) error {
	results, err := handler.srv.GetAll(c)
	if err != nil {
		logrus.Errorf("Error in handler: %v", err)
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
		logrus.Errorf("Error in handler: %v", err)
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	err = handler.srv.Delete(c, entity.ID)
	if err != nil {
		logrus.Errorf("Error in handler: %v", err)
		str := fmt.Sprintf("Error in handler: %v", err)
		return c.String(http.StatusNotFound, str)
	}
	return c.String(http.StatusOK, "delete request")
}

// Create function receives POST request from client
func (handler *PersonHandlerImpl) Create(c echo.Context) error {
	req := c.Request()

	// get request body
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		logrus.Errorf("Error in handler: %v", err)

		return fmt.Errorf("error reading request body: %v", err)
	}

	var entity *model.Person
	err = json.Unmarshal(body, &entity)
	if err != nil {
		logrus.Errorf("error in handler: %v", err)
		str := fmt.Sprintf("Error in handler: %v", err)
		return c.String(http.StatusNotFound, str)
	}
	entity.ID = uuid.New()
	err = handler.srv.Create(c, entity)
	if err != nil {
		logrus.Errorf("Error in handler: %v", err)
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
		logrus.Errorf("Error in handler: %v", err)
		return c.String(http.StatusNotFound, "error unmarshalling request body")
	}

	err = handler.srv.Update(c, entity.ID, &entity)
	if err != nil {
		logrus.Errorf("Error in handler: %v", err)
		return c.String(http.StatusNotFound, "error updating (handler error)")
	}
	return c.String(http.StatusOK, "update request")
}
