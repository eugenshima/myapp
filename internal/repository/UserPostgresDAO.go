package repository

import (
	"context"
	"fmt"

	"github.com/eugenshima/myapp/internal/model"
)

func (db *PsqlConnection) GetAllUsers() ([]model.User, error) {
	rows, err := db.pool.Query(context.Background(), "SELECT id, login, password, role FROM goschema.user")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Create slice to store data from our SQL request
	var results []model.User

	// go;) through each line
	for rows.Next() {
		entity := model.User{}
		err := rows.Scan(&entity.ID, &entity.Login, &entity.Password, &entity.Role)
		if err != nil {
			return nil, err
		}
		results = append(results, entity)
	}
	fmt.Println(results)
	return results, rows.Err()
}
