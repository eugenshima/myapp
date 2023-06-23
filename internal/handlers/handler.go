package handlers

import (
	"io/ioutil"
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
		// Get object from the context
		req := c.Request()

		// get request body
		body, err := ioutil.ReadAll(req.Body)
		if err != nil {
			return err
		}
		// Call repository function
		repository.CreatePerson()
		return c.JSON(http.StatusOK, body)
	})

	e.GET("/all", func(c echo.Context) error {
		repository.GetAll()
		return c.String(http.StatusOK, "Get Request(GetAll)")
	})

	e.DELETE("/delete/:id", func(c echo.Context) error {
		Uuid := c.Param("id")
		repository.Delete(Uuid)
		return c.String(http.StatusOK, "Delete Request(delete)")
	})
	e.PATCH("/update", func(c echo.Context) error {
		return c.String(http.StatusOK, "Put Request(update)")
	})
}
