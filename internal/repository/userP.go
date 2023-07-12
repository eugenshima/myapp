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
func (db *UserPsqlConnection) GetUser(ctx context.Context, login string) (*model.User, error) {
	var user model.User
	err := db.pool.QueryRow(ctx, "SELECT id, password, role FROM goschema.user WHERE login = $1", login).Scan(&user.ID, &user.Password, &user.Role)
	if err != nil {
		return nil, fmt.Errorf("QueryRow: %w", err)
	}
	return &user, nil
}

// Signup function executes a query, which insert a user to user table
func (db *UserPsqlConnection) Signup(ctx context.Context, entity *model.User) error {
	bd, err := db.pool.Exec(ctx,
		`INSERT INTO goschema.user (id, login, password, role) 
		 values ($1, $2, $3, $4)`,
		entity.ID, entity.Login, entity.Password, entity.Role)
	if err != nil && !bd.Insert() {
		return fmt.Errorf("Exec(): %w", err)
	}
	return nil
}

// GetAll func executes a query, which returns all users
func (db *UserPsqlConnection) GetAll(ctx context.Context) ([]*model.User, error) {
	rows, err := db.pool.Query(ctx, "SELECT id, login, password, role, refresh_token FROM goschema.user")
	if err != nil {
		return nil, fmt.Errorf("Query(): %w", err)
	}
	defer rows.Close()

	var users []*model.User
	for rows.Next() {
		var user model.User
		err := rows.Scan(&user.ID, &user.Login, &user.Password, &user.Role, &user.RefreshToken)
		if err != nil {
			return nil, fmt.Errorf("Scan(): %w", err)
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
		return fmt.Errorf("QueryRow: %w", err)
	}
	bd, err := db.pool.Exec(ctx, "UPDATE goschema.user SET refresh_token=$1 WHERE id=$2", token, user.ID)
	if err != nil && !bd.Update() {
		return fmt.Errorf("Exec(): %w", err)
	}
	return nil
}

// GetRefreshToken returns a refresh token for the given user
func (db *UserPsqlConnection) GetRefreshToken(ctx context.Context, ID uuid.UUID) ([]byte, error) {
	var user model.User
	err := db.pool.QueryRow(ctx, "SELECT refresh_token FROM goschema.user WHERE id=$1", ID).Scan(&user.RefreshToken)
	if err != nil {
		return nil, fmt.Errorf("QueryRow: %w ", err)
	}
	return user.RefreshToken, nil
}

// GetRoleByID returns a role for the given user ID
func (db *UserPsqlConnection) GetRoleByID(ctx context.Context, ID uuid.UUID) (string, error) {
	var user model.User
	err := db.pool.QueryRow(ctx, "SELECT role FROM goschema.user WHERE id=$1", ID).Scan(&user.Role)
	if err != nil {
		return "", fmt.Errorf("QueryRow: %w ", err)
	}
	return user.Role, nil
}

func (db *UserPsqlConnection) Delete(ctx context.Context, ID uuid.UUID) error {
	bd, err := db.pool.Exec(ctx, "DELETE FROM goschema.user WHERE id=$1", ID)
	if err != nil || bd.String() == "DELETE 0" {
		return fmt.Errorf("Exec(): %w", err)
	}
	return nil
}
