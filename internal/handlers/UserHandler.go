package handlers

import (
	"net/http"

	"github.com/eugenshima/myapp/internal/model"
	"github.com/labstack/echo/v4"
)

func (handler *Handler) UserGetAll(c echo.Context) error {
	var results []model.User

	results, err := handler.PersonService.GetAllUsers()
	if err != nil {
		return c.String(http.StatusNotFound, "Something bad happened")
	}

	return c.JSON(http.StatusOK, results)
}
