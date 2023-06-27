package repository

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func CreateMongoConnect() (*mongo.Client, error) {
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")

	// Connect to MongoDB
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		return nil, err
	}

	// Check the connection
	err = client.Ping(context.Background(), nil)
	if err != nil {
		return nil, err
	}
	fmt.Println("Connected to MongoDB!")

	// // Disconnect from MongoDB
	// ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	// err = client.Disconnect(ctx)
	// if err != nil {
	// 	return nil, err
	// }

	// fmt.Println("Disconnected from MongoDB!")
	return client, nil
}
