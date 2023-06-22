package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/eugenshima/myapp/internal/model"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

func CreatePerson(conn *pgx.Conn, err error) {
	var entity model.Entity
	entity.ID = uuid.New()
	entity.Name = "eugen"
	entity.Age = 20
	entity.IsHealthy = true
	bd, err := conn.Exec(context.Background(), "INSERT INTO yauhenishymanski.newtable (id, name, age, ishealthy) VALUES ($1, $2, $3, $4)", entity.ID, entity.Name, entity.Age, entity.IsHealthy)
	if err != nil {
		fmt.Println("Error inserting data into table:", err)
		return
	}
	fmt.Println(bd, " <-- result of the request")
	fmt.Println("Data successfully inserted into table yauhenishymanski.newtable")
}

func GetById(conn *pgx.Conn, err error) {
	var entity model.Entity
	err = conn.QueryRow(context.Background(), "select * from yauhenishymanski.newtable where name='eugen'").Scan(&entity.ID, &entity.Name, &entity.Age, &entity.IsHealthy)
	if err != nil {
		log.Fatal(err)
	}

	jsonString, err := json.Marshal(entity)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(string(jsonString))
}

//func GetAll(conn *pgx.Conn, err error)
//func Delete(conn *pgx.Conn, err error)
//func Update(conn *pgx.Conn, err error)
