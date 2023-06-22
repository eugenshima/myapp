package handlers

import (
	"net/http"

	"github.com/eugenshima/myapp/internal/repository"
	"github.com/jackc/pgx/v5"
	"github.com/labstack/echo/v4"
)

func HttpHandler(e *echo.Echo, conn *pgx.Conn, err error) {
	e.GET("/getById", func(c echo.Context) error {
		repository.GetById(conn, err)
		return c.String(http.StatusOK, "getbyid")
	})
	e.GET("/add", func(c echo.Context) error {
		repository.CreatePerson(conn, err)
		return c.String(http.StatusOK, "post request")
	})
}
