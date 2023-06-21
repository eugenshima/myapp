package repository

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
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
