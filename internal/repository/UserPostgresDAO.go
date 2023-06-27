package repository

import (
	"context"
	"fmt"

	"github.com/eugenshima/myapp/internal/model"
)

func (db *PsqlConnection) GetAllUsers() ([]model.Entity, error) {
	rows, err := db.pool.Query(context.Background(), "SELECT id, login, password FROM goschema.user")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Create slice to store data from our SQL request
	var results []model.Entity

	// go;) through each line
	for rows.Next() {
		entity := model.Entity{}
		err := rows.Scan(&entity.ID, &entity.Name, &entity.Age, &entity.IsHealthy)
		if err != nil {
			return nil, err
		}
		results = append(results, entity)
	}
	fmt.Println(results)
	return results, rows.Err()
}
