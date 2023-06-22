package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/eugenshima/myapp/internal/repository"
	"github.com/jackc/pgx/v5"
	"github.com/labstack/echo/v4"
)

func createConn() (*pgx.Conn, error) {
	connConfig, err := pgx.ParseConfig("postgres://eugen:ur2qly1ini@localhost:5432/eugene")
	if err != nil {
		PsqlErr := fmt.Errorf("failed to create connection to PostgreSQL: %v", err)
		return nil, PsqlErr
	}

	conn, err := pgx.ConnectConfig(context.Background(), connConfig)
	if err != nil {
		PsqlErr := fmt.Errorf("failed to create connection to PostgreSQL: %v", err)
		return nil, PsqlErr
	}

	fmt.Println("Connection to PostgreSQL successful")
	return conn, nil
}

func main() {
	e := echo.New()

	conn, err := createConn()

	e.GET("/getById", func(c echo.Context) error {
		repository.GetById(conn, err)
		return c.String(http.StatusOK, "getbyid")
	})
	e.GET("/add", func(c echo.Context) error {
		repository.CreatePerson(conn, err)
		return c.String(http.StatusOK, "post request")
	})

	e.Logger.Fatal(e.Start(":1323"))

}
