// Package repository provides functions for interacting with a database
// or other persistent storage system in a web service.
// It includes functions for creating, reading, updating, and deleting data from the storage system.
package repository

import (
	"context"
	"fmt"

	"github.com/eugenshima/myapp/internal/model"
	"github.com/jackc/pgx/v4/pgxpool"
)

// UserPsqlConnection struct represents a connection to a database
type UserPsqlConnection struct {
	pool *pgxpool.Pool
}

// NewPsqlConnection constructor for PsqlConnection
func NewUserPsqlConnection(pool *pgxpool.Pool) *UserPsqlConnection {
	return &UserPsqlConnection{pool: pool}
}

// Login function executes a query, which select all rows from user table
func (db *UserPsqlConnection) GetUser(ctx context.Context, login, password string) (*model.User, error) {
	fmt.Println(password)
	var user model.User
	err := db.pool.QueryRow(ctx, "SELECT id FROM goschema.user WHERE login = $1 AND password = $2", login, password).Scan(&user.ID)
	if err != nil {
		return nil, fmt.Errorf("error executing query: %v", err)
	}
	return &user, nil
}

// signup function executes a query, which insert a user to user table
func (db *UserPsqlConnection) Signup(ctx context.Context, entity *model.User) error {
	bd, err := db.pool.Exec(ctx,
		`INSERT INTO goschema.user (id, login, password, email) 
		 values ($1, $2, $3, $4)`,
		entity.ID, entity.Login, entity.Password, entity.Email)
	if err != nil && !bd.Insert() {
		return fmt.Errorf("error inserting user: %v", err)
	}
	return nil
}

func (db *UserPsqlConnection) GetAll(ctx context.Context) ([]*model.User, error) {
	rows, err := db.pool.Query(ctx, `SELECT id, login, password, email FROM goschema.user`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var users []*model.User
	for rows.Next() {
		var user model.User
		err := rows.Scan(&user.ID, &user.Login, &user.Password, &user.Email)
		if err != nil {
			return nil, err
		}
		users = append(users, &user)
	}
	return users, nil
}
