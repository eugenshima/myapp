package handlers

import (
	"context"
	"fmt"
	"net/http"

	"github.com/eugenshima/myapp/internal/model"

	vld "github.com/go-playground/validator"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

// PersonHandler struct contains service service.PersonService
type PersonHandler struct {
	srv PersonService
	vl  *vld.Validate
}

// NewPersonHandler is a constructor
func NewPersonHandler(srv PersonService, vl *vld.Validate) *PersonHandler {
	return &PersonHandler{
		srv: srv,
		vl:  vl,
	}
}

//go:generate mockery --name=PersonService --case=underscore --output=./mocks

// PersonService interface, which contains Service methods
type PersonService interface {
	GetByID(ctx context.Context, id uuid.UUID) (*model.Person, error)
	GetAll(ctx context.Context) ([]*model.Person, error)
	Delete(ctx context.Context, uuidString uuid.UUID) (uuid.UUID, error)
	Create(ctx context.Context, entity *model.Person) (uuid.UUID, error)
	Update(ctx context.Context, uuidString uuid.UUID, entity *model.Person) (uuid.UUID, error)
}

// GetByID function receives Get request from client
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
func (handler *PersonHandler) GetByID(c echo.Context) error {
	ID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		logrus.Errorf("Parse: %v", err)
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Parse: %v", err))
	}
	result, err := handler.srv.GetByID(c.Request().Context(), ID)
	if err != nil {
		logrus.WithFields(logrus.Fields{"id": ID}).Errorf("GetByID: %v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("GetByID: %v", err))
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
func (handler *PersonHandler) GetAll(c echo.Context) error {
	results, err := handler.srv.GetAll(c.Request().Context())
	if err != nil {
		logrus.Errorf("GetAll: %v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("GetAll: %v", err))
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
func (handler *PersonHandler) Delete(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		logrus.Errorf("Parse: %v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("Parse: %v", err))
	}
	id, err = handler.srv.Delete(c.Request().Context(), id)
	if err != nil {
		logrus.WithFields(logrus.Fields{"id": &id}).Errorf("Delete: %v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("Delete: %v", err))
	}
	return c.String(http.StatusOK, fmt.Sprintf("Deleted ID: %v", id))
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
func (handler *PersonHandler) Create(c echo.Context) error {
	var person *model.Person
	err := c.Bind(&person)
	if err != nil {
		logrus.WithFields(logrus.Fields{"person": &person}).Errorf("Bind: %v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("Bind: %v", err))
	}
	person.ID = uuid.New()

	err = c.Validate(person)
	if err != nil {
		logrus.WithFields(logrus.Fields{"person": &person}).Errorf("Validate: %v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("Validate: %v", err))
	}
	id, err := handler.srv.Create(c.Request().Context(), person)
	if err != nil {
		logrus.WithFields(logrus.Fields{"person": &person}).Errorf("Create: %v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("Create: %v", err))
	}
	return c.String(http.StatusOK, fmt.Sprintf("inserted ID: %v", id))
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
func (handler *PersonHandler) Update(c echo.Context) error {
	var person *model.Person
	err := c.Bind(&person)
	if err != nil {
		logrus.WithFields(logrus.Fields{"person": &person}).Errorf("Bind: %v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("Bind: %v", err))
	}
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		logrus.Errorf("Parse: %v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("Parse: %v", err))
	}
	err = c.Validate(person)
	if err != nil {
		logrus.WithFields(logrus.Fields{"person": person}).Errorf("Validate: %v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("Validate: %v", err))
	}
	id, err = handler.srv.Update(c.Request().Context(), id, person)
	if err != nil {
		logrus.WithFields(logrus.Fields{"id": id, "person": person}).Errorf("Update: %v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("Update: %v", err))
	}
	return c.String(http.StatusOK, fmt.Sprintf("Updated id --> %v", id))
}
