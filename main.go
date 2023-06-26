package main

import (
	"fmt"

	"github.com/eugenshima/myapp/internal/handlers"
	"github.com/eugenshima/myapp/internal/repository"
	"github.com/labstack/echo/v4"
)

type Main struct {
	handler *handlers.Handler
}

func NewMain(handler *handlers.Handler) *Main {
	return &Main{handler: handler}
}

func main() {
	e := echo.New()
	// Initializing the DATABASE CONNECTOR
	pool, err := repository.NewDatabasePsqlConnection()
	if err != nil {
		fmt.Println(err)
	}

	dbpool := repository.NewPsqlConnection(pool)
	handlr := handlers.NewHandler(dbpool)
	handlr.HttpHandler(e)

	e.Logger.Fatal(e.Start(":1323"))

}
