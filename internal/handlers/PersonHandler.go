package handlers

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/eugenshima/myapp/internal/model"
	"github.com/eugenshima/myapp/internal/service"
	"github.com/labstack/echo/v4"
)

// type Handler struct {
// 	DB *repository.PsqlConnection
// }

type Handler struct {
	PersonService *service.Service
}

func NewHandler(PersonService *service.Service) *Handler {
	return &Handler{PersonService: PersonService}
}

// func NewHandler(DB *repository.PsqlConnection) *Handler {
// 	return &Handler{DB: DB}
// }

func (handler *Handler) HttpHandler(e *echo.Echo) {

	e.GET("/getByName/:name", handler.GetByName)

	e.GET("/getAll", handler.GetAll)

	e.DELETE("/delete/:id", handler.Delete)

	e.POST("/insert", handler.Insert)

	e.PATCH("/update/:id", handler.Update)

	e.GET("/user/getall", handler.UserGetAll)

}

func (handler *Handler) GetByName(c echo.Context) error {
	Name := c.Param("name")
	result, err := handler.PersonService.GetByName(Name)
	if err != nil {
		return c.String(http.StatusNotFound, "Bad")
	}
	return c.JSON(http.StatusOK, result)
}

func (handler *Handler) GetAll(c echo.Context) error {
	var results []model.Entity

	results, err := handler.PersonService.GetAll()
	if err != nil {
		return c.String(http.StatusNotFound, "Bad")
	}

	return c.JSON(http.StatusOK, results)
}

func (handler *Handler) Delete(c echo.Context) error {
	id := c.Param("id")
	handler.PersonService.Delete(id)
	return c.String(http.StatusOK, "delete request")
}

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
	handler.PersonService.Insert(entity)
	return c.String(http.StatusOK, "insert request")
}

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
	handler.PersonService.Update(id, entity)
	return c.String(http.StatusOK, "update request")
}
