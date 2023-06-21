package main

import (
	//"myapp/internal/repository"
	"net/http"

	"github.com/eugenshima/myapp/internal/repository"
	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()
	str := repository.Greet()
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, str)
	})
	e.Logger.Fatal(e.Start(":1323"))
}
