package repository

//here will be a model of our entity from database
import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type Entity struct {
	ID        uuid.UUID `db:"id"`
	Name      string    `db:"name"`
	Age       int       `db:"age"`
	IsHealthy bool      `db:"is_healthy"`
}

func NewEntity(name string, age int, isHealthy bool) *Entity {
	return &Entity{
		ID:        uuid.New(),
		Name:      name,
		Age:       age,
		IsHealthy: isHealthy,
	}
}

func get() {
	// Подключаемся к базе данных
	conn, err := pgx.Connect(context.Background(), "postgres://user:password@localhost:5432/your-db")
	if err != nil {
		panic(err)
	}
	defer conn.Close(context.Background())

	// Получаем результаты SELECT запроса
	rows, err := conn.Query(context.Background(), "SELECT * FROM your_table")
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	// Выбираем и обрабатываем данные
	for rows.Next() {
		var id int
		var name string
		if err := rows.Scan(&id, &name); err != nil {
			panic(err)
		}
		fmt.Println("id:", id, "name:", name)
	}
	if err := rows.Err(); err != nil {
		panic(err)
	}
}
