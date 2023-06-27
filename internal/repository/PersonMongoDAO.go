// Package repository provides functions for interacting with a database
// or other persistent storage system in a web service.
// It includes functions for creating, reading, updating, and deleting data from the storage system.
package repository

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// MongoDBConnection is a struct, which contains *mongo.Client variable
type MongoDBConnection struct {
	client *mongo.Client
}

// NewMongoDBConnection func is a constructor of MongoDbConnection struct
func NewMongoDBConnection(client *mongo.Client) *MongoDBConnection {
	return &MongoDBConnection{client: client}
}

// CreateMongoConnect creates a connection to MongoDB server
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

	return client, nil
}

// FindPersons executes the "db.person.find()" command
func (db *MongoDBConnection) FindPersons(ctx context.Context) (bson.M, error) {
	// Select the database and collection
	collection := db.client.Database("my_mongo_base").Collection("person")

	// Call the find function on the collection
	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}

	defer cursor.Close(ctx)
	if cursor.Close(ctx) != nil {
		return nil, err
	}

	var result bson.M
	// Iterate over the cursor and print the result
	for cursor.Next(ctx) {
		err := cursor.Decode(&result)

		if err != nil {
			return nil, err
		}

		fmt.Println(result)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}
	return result, nil
}
