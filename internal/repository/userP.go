// Package repository provides functions for interacting with a database
// or other persistent storage system in a web service.
// It includes functions for creating, reading, updating, and deleting data from the storage system.
package repository

import (
	"context"
	"fmt"

	"github.com/eugenshima/myapp/internal/model"
	"github.com/google/uuid"

	"github.com/jackc/pgx/v4/pgxpool"
)

// UserPsqlConnection struct represents a connection to a database
type UserPsqlConnection struct {
	pool *pgxpool.Pool
}

// NewUserPsqlConnection constructor for PsqlConnection
func NewUserPsqlConnection(pool *pgxpool.Pool) *UserPsqlConnection {
	return &UserPsqlConnection{pool: pool}
}

// GetUser function executes a query, which select all rows from user table
func (db *UserPsqlConnection) GetUser(ctx context.Context, login string) (uuid.UUID, []byte, error) {
	var user model.User
	err := db.pool.QueryRow(ctx, "SELECT id, password FROM goschema.user WHERE login = $1", login).Scan(&user.ID, &user.Password)
	if err != nil {
		return uuid.Nil, nil, fmt.Errorf("error executing query: %v", err)
	}
	return user.ID, user.Password, nil
}

// Signup function executes a query, which insert a user to user table
func (db *UserPsqlConnection) Signup(ctx context.Context, entity *model.User) error {
	bd, err := db.pool.Exec(ctx,
		`INSERT INTO goschema.user (id, login, password, role) 
		 values ($1, $2, $3, $4)`,
		entity.ID, entity.Login, entity.Password, entity.Role)
	if err != nil && !bd.Insert() {
		return fmt.Errorf("error inserting user: %v", err)
	}
	return nil
}

// GetAll func executes a query, which returns all users
func (db *UserPsqlConnection) GetAll(ctx context.Context) ([]*model.User, error) {
	rows, err := db.pool.Query(ctx, `SELECT id, login, password, role, refreshtoken FROM goschema.user`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var users []*model.User
	for rows.Next() {
		var user model.User
		err := rows.Scan(&user.ID, &user.Login, &user.Password, &user.Role, &user.RefreshToken)
		if err != nil {
			return nil, err
		}
		users = append(users, &user)
	}
	return users, nil
}

// SaveRefreshToken func executes a query, which saves the refresh token to a specific user
func (db *UserPsqlConnection) SaveRefreshToken(ctx context.Context, ID uuid.UUID, token []byte) error {
	var user model.User
	err := db.pool.QueryRow(ctx, "SELECT id, login, password, role FROM goschema.user WHERE id=$1", ID).Scan(&user.ID, &user.Login, &user.Password, &user.Role)
	if err != nil {
		return fmt.Errorf("error in SaveRefreshToken: %v ", err)
	}
	bd, err := db.pool.Exec(ctx, "UPDATE goschema.user SET refreshtoken=$1 WHERE id=$2", token, user.ID)
	if err != nil && !bd.Update() {
		return fmt.Errorf("error updating user: %v", err)
	}
	return nil
}
