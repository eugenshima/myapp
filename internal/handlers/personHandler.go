package handlers

import (
	"context"
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
	GetByID(ctx context.Context, id uuid.UUID) (*model.Person, error)
	GetAll(ctx context.Context) ([]model.Person, error)
	Delete(ctx context.Context, uuidString uuid.UUID) error
	Create(ctx context.Context, entity *model.Person) (uuid.UUID, error)
	Update(ctx context.Context, uuidString uuid.UUID, entity *model.Person) error
}

// Create function receives Get request from client
// @Summary Get person by ID
// @Security ApiKeyAuth
// @Tags Person CRUD
// @Description Retrieves a person by ID
// @Produce json
// @Param id path string true "ID of the person"
// @Success 200 {object} model.Person "Person object"
// @Failure 400 {string} string "Bad request"
// @Failure 404 {string} string "Person not found"
// @Router /api/person/getById/{id} [get]
func (handler *PersonHandlerImpl) GetByID(c echo.Context) error {
	ID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		logrus.Errorf("Error parsing id: %v", err)
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	result, err := handler.srv.GetByID(c.Request().Context(), ID)
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
// @Summary Get All
// @Security ApiKeyAuth
// @Tags Person CRUD
// @Description Get All ^ :)
// @Produce json
// @Failure 400 {string} string "Bad request"
// @Failure 404 {string} string "Person not found"
// @Router /api/person/getAll [get]
func (handler *PersonHandlerImpl) GetAll(c echo.Context) error {
	results, err := handler.srv.GetAll(c.Request().Context())
	if err != nil {
		logrus.Errorf("Error in handler: %v", err)
		return c.String(http.StatusNotFound, fmt.Sprintf("Error in handler: %v", err))
	}
	return c.JSON(http.StatusOK, results)
}

// Delete function receives DELETE request from client
// @Summary Delete Person
// @Security ApiKeyAuth
// @Tags Person CRUD
// @Description Delete person from database by ID
// @Accept json
// @Produce json
// @Param id path string true "ID of the person"
// @Success 200 {object} model.Person "Person object"
// @Failure 400 {string} string "Bad request"
// @Failure 404 {string} string "Person not found"
// @Router /api/person/delete/{id} [delete]
func (handler *PersonHandlerImpl) Delete(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		logrus.Errorf("Error in handler: %v", err)
		return c.String(http.StatusNotFound, fmt.Sprintf("Error in handler: %v", err))
	}
	err = handler.srv.Delete(c.Request().Context(), id)
	if err != nil {
		logrus.Errorf("Error in handler: %v", err)
		return c.String(http.StatusNotFound, fmt.Sprintf("Error in handler: %v", err))
	}
	return c.String(http.StatusOK, "delete request")
}

// Create function receives POST request from client
// @Summary Create person
// @Security ApiKeyAuth
// @Tags Person CRUD
// @Description Creates a new person
// @Accept json
// @Produce json
// @Param entity body model.Person true "Person object to be created"
// @Success 200 {string} string "ID of the created person"
// @Failure 400 {string} string "Bad request"
// @Failure 404 {string} string "Error message"
// @Router /api/person/insert [post]
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
	id, err := handler.srv.Create(c.Request().Context(), entity)
	if err != nil {
		logrus.Errorf("Error in handler: %v", err)
		return c.String(http.StatusNotFound, fmt.Sprintf("Error in handler: %v", err))
	}
	return c.String(http.StatusOK, fmt.Sprintf("inserted this ID: %v", id))
}

// Update function receives PATCH request from client
// @Summary Update person's information
// @Security ApiKeyAuth
// @Tags Person CRUD
// @Description updates person
// @Accept json
// @Produce json
// @Param id path string true "ID of the person"
// @Param entity body model.Person true "Updated person data"
// @Success 200 {string} string "ID of the created person"
// @Failure 400 {string} string "Bad request"
// @Failure 404 {string} string "Error message"
// @Router /api/person/update/{id} [patch]
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
	err = handler.srv.Update(c.Request().Context(), id, entity)
	if err != nil {
		logrus.Errorf("Error in handler: %v", err)
		return c.String(http.StatusNotFound, fmt.Sprintf("Error in handler: %v", err))
	}
	return c.String(http.StatusOK, "update request")
}
