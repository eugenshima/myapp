package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5"
)

func GetById(conn *pgx.Conn, err error) {
	var greeting Entity
	err = conn.QueryRow(context.Background(), "select * from yauhenishymanski.db_table where id=2").Scan(&greeting.ID, &greeting.Name, &greeting.Age, &greeting.IsHealthy)
	if err != nil {
		log.Fatal(err)
	}

	jsonString, err := json.Marshal(greeting)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(string(jsonString))
}

//func GetAll(conn *pgx.Conn, err error)
