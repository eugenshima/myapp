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
	connConfig, err := pgx.ParseConfig("postgres://eugen:ur2qly1ini@localhost:5432/eugen")
	if err != nil {
		return nil, fmt.Errorf("failed to parse connection config: %v", err)
	}

	conn, err := pgx.ConnectConfig(context.Background(), connConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create connection to PostgreSQL: %v", err)
	}

	fmt.Println("Connection to PostgreSQL successful")
	return conn, nil
}

func main() {
	e := echo.New()
	str := repository.Greet()
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, str)
	})
	e.GET("/bd", func(c echo.Context) error {
		createConn()
		return c.String(http.StatusOK, str)
	})
	e.Logger.Fatal(e.Start(":1323"))

}
