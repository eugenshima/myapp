package handlers

import (
	"net/http"

	"github.com/eugenshima/myapp/internal/repository"
	"github.com/labstack/echo/v4"
)

func HttpHandler(e *echo.Echo) {

	e.GET("/getByName/:name", func(c echo.Context) error {
		Name := c.Param("name")
		repository.GetByName(Name)
		return c.String(http.StatusOK, "Get Request(get by name)")
	})

	e.POST("/add", func(c echo.Context) error {
		repository.CreatePerson()
		return c.String(http.StatusOK, "create")
	})

	e.GET("/all", func(c echo.Context) error {
		repository.GetAll()
		return c.String(http.StatusOK, "Get Request(GetAll)")
	})

	e.DELETE("/delete/:name", func(c echo.Context) error {
		Name := c.Param("name")
		repository.Delete(Name)
		return c.String(http.StatusOK, "Delete Request(delete)")
	})
	e.PUT("/update", func(c echo.Context) error {
		return c.String(http.StatusOK, "Put Request(update)")
	})
}
