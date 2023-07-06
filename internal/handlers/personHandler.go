package handlers

import (
	"context"
	"fmt"
	"net/http"
	"strings"

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
	Delete(ctx context.Context, uuidString uuid.UUID, accessToken string) error
	Create(ctx context.Context, entity *model.Person, accessToken string) error
	Update(ctx context.Context, uuidString uuid.UUID, entity *model.Person, accessToken string) error
}

// GetByID function receives GET request from client
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
func (handler *PersonHandlerImpl) GetAll(c echo.Context) error {
	results, err := handler.srv.GetAll(c.Request().Context())
	if err != nil {
		logrus.Errorf("Error in handler: %v", err)
		return c.String(http.StatusNotFound, fmt.Sprintf("Error in handler: %v", err))
	}
	return c.JSON(http.StatusOK, results)
}

// Delete function receives DELETE request from client
func (handler *PersonHandlerImpl) Delete(c echo.Context) error {
	authHeader := c.Request().Header.Get("Authorization")
	if authHeader == "" {
		return echo.NewHTTPError(http.StatusUnauthorized, "Missing authorization header")
	}

	// Разбиение значения заголовка Authorization на тип и токен
	authParts := strings.Split(authHeader, " ")
	if len(authParts) != 2 || strings.EqualFold(authParts[0], "bearer") {
		// Ошибка: некорректный формат заголовка Authorization
		return echo.NewHTTPError(http.StatusUnauthorized, "Invalid authorization header")
	}

	// Извлечение access токена
	accessToken := authParts[1]
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		logrus.Errorf("Error in handler: %v", err)
		return c.String(http.StatusNotFound, fmt.Sprintf("Error in handler: %v", err))
	}
	err = handler.srv.Delete(c.Request().Context(), id, accessToken)
	if err != nil {
		logrus.Errorf("Error in handler: %v", err)
		return c.String(http.StatusNotFound, fmt.Sprintf("Error in handler: %v", err))
	}
	return c.String(http.StatusOK, "delete request")
}

// Create function receives POST request from client
func (handler *PersonHandlerImpl) Create(c echo.Context) error {
	authHeader := c.Request().Header.Get("Authorization")
	if authHeader == "" {
		return echo.NewHTTPError(http.StatusUnauthorized, "Missing authorization header")
	}

	authParts := strings.Split(authHeader, " ")
	if len(authParts) != 2 || strings.EqualFold(authParts[0], "bearer") {
		return echo.NewHTTPError(http.StatusUnauthorized, "Invalid authorization header")
	}

	accessToken := authParts[1]
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
	err = handler.srv.Create(c.Request().Context(), entity, accessToken)
	if err != nil {
		logrus.Errorf("Error in handler: %v", err)
		return c.String(http.StatusNotFound, fmt.Sprintf("Error in handler: %v", err))
	}
	return c.String(http.StatusOK, "insert request")
}

// Update function receives PATCH request from client
func (handler *PersonHandlerImpl) Update(c echo.Context) error {
	authHeader := c.Request().Header.Get("Authorization")
	if authHeader == "" {
		return echo.NewHTTPError(http.StatusUnauthorized, "Missing authorization header")
	}

	authParts := strings.Split(authHeader, " ")
	if len(authParts) != 2 || !strings.EqualFold(authParts[0], "bearer") {
		return echo.NewHTTPError(http.StatusUnauthorized, "Invalid authorization header")
	}

	accessToken := authParts[1]
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
	err = handler.srv.Update(c.Request().Context(), id, entity, accessToken)
	if err != nil {
		logrus.Errorf("Error in handler: %v", err)
		return c.String(http.StatusNotFound, fmt.Sprintf("Error in handler: %v", err))
	}
	return c.String(http.StatusOK, "update request")
}
