package main

import (
	"context"
	"fmt"

	"github.com/eugenshima/myapp/internal/handlers"
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
	handlers.HttpHandler(e, conn, err)

	e.Logger.Fatal(e.Start(":1323"))

}
