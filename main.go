package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/jackc/pgx/v5"
	"github.com/labstack/echo/v4"
)

func createConn() (*pgx.Conn, error) {
	connConfig, err := pgx.ParseConfig("postgres://eugen:ur2qly1ini@localhost:5432/eugen")
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
	str := "Greetings, traveller"
	e.GET("/bd", func(c echo.Context) error {
		return c.String(http.StatusOK, str)
	})
	e.GET("/", func(c echo.Context) error {
		createConn()
		repository.get()
		return c.String(http.StatusOK, "Interract with PostgreSQL")
	})
	e.Logger.Fatal(e.Start(":1323"))

}
