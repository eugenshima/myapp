package handlers

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/eugenshima/myapp/internal/model"
	"github.com/eugenshima/myapp/internal/service"
	"github.com/labstack/echo/v4"
)

// Handler struct contains *service.Service variable
type Handler struct {
	PersonService *service.Service
}

// NewHandler function is a constructor for Handler
func NewHandler(PersonService *service.Service) *Handler {
	return &Handler{PersonService: PersonService}
}

// HTTPHandler routing function
func (handler *Handler) HTTPHandler(e *echo.Echo) {
	e.GET("/getByName/:name", handler.GetByName)

	e.GET("/getAll", handler.GetAll)

	e.DELETE("/delete/:id", handler.Delete)

	e.POST("/insert", handler.Insert)

	e.PATCH("/update/:id", handler.Update)

	e.GET("/user/getall", handler.UserGetAll)
}

// GetByName function receives GET request from client
func (handler *Handler) GetByName(c echo.Context) error {
	Name := c.Param("name")
	result, err := handler.PersonService.GetByName(c, Name)
	if err != nil {
		return c.String(http.StatusNotFound, "Bad")
	}
	return c.JSON(http.StatusOK, result)
}

// GetAll function receives GET request from client
func (handler *Handler) GetAll(c echo.Context) error {
	var results []model.Entity

	results, err := handler.PersonService.GetAll()
	if err != nil {
		return c.String(http.StatusNotFound, "Bad")
	}

	return c.JSON(http.StatusOK, results)
}

// Delete function receives DELETE request from client
func (handler *Handler) Delete(c echo.Context) error {
	id := c.Param("id")
	err := handler.PersonService.Delete(id)
	if err != nil {
		return err
	}
	return c.String(http.StatusOK, "delete request")
}

// Insert function receives POST request from client
func (handler *Handler) Insert(c echo.Context) error {
	req := c.Request()

	// get request body
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return err
	}

	var entity *model.Entity
	err = json.Unmarshal(body, &entity)
	if err != nil {
		return c.String(http.StatusNotFound, "Bad")
	}
	err = handler.PersonService.Insert(entity)
	if err != nil {
		return err
	}
	return c.String(http.StatusOK, "insert request")
}

// Update function receives PATCH request from client
func (handler *Handler) Update(c echo.Context) error {
	id := c.Param("id")
	req := c.Request()

	// get request body
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return err
	}
	var entity *model.Entity
	err = json.Unmarshal(body, &entity)
	if err != nil {
		return c.String(http.StatusNotFound, "Bad")
	}
	err = handler.PersonService.Update(id, entity)
	if err != nil {
		return err
	}
	return c.String(http.StatusOK, "update request")
}
