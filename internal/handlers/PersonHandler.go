package handlers

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/eugenshima/myapp/internal/model"
	"github.com/eugenshima/myapp/internal/repository"
	"github.com/labstack/echo/v4"
)

type Handler struct {
	DB *repository.PsqlConnection
	//Service *service.Service
}

// func NewHandler(Service *service.Service) *Handler {
// 	return &Handler{Service: Service}
// }

func NewHandler(DB *repository.PsqlConnection) *Handler {
	return &Handler{DB: DB}
}

func (handler *Handler) HttpHandler(e *echo.Echo) {

	e.GET("/getByName/:name", func(c echo.Context) error {
		Name := c.Param("name")
		result, err := handler.DB.GetByName(Name)
		if err != nil {
			return c.String(http.StatusNotFound, "Bad")
		}
		return c.JSON(http.StatusOK, result)
	})
	e.GET("/getAll", func(c echo.Context) error {
		var results []model.Entity

		results, err := handler.DB.GetAll()
		if err != nil {
			return c.String(http.StatusNotFound, "Bad")
		}

		return c.JSON(http.StatusOK, results)
	})
	e.DELETE("/delete/:id", func(c echo.Context) error {
		id := c.Param("id")
		handler.DB.Delete(id)
		return c.String(http.StatusOK, "delete request")
	})
	e.POST("/insert", func(c echo.Context) error {
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
		handler.DB.Insert(entity)
		return c.String(http.StatusOK, "insert request")
	})

	e.PATCH("/update/:id", func(c echo.Context) error {
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
		handler.DB.Update(id, entity)
		return c.String(http.StatusOK, "update request")
	})

}
