package repository

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoDbConnection struct {
	client *mongo.Client
}

func NewMongoDbConnection(client *mongo.Client) *MongoDbConnection {
	return &MongoDbConnection{client: client}
}

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

func (db *MongoDbConnection) FindPersons(ctx context.Context) (bson.M, error) {
	// Select the database and collection
	collection := db.client.Database("my_mongo_base").Collection("person")

	// Call the find() function on the collection
	cursor, err := collection.Find(ctx, bson.M{})

	if err != nil {
		return nil, err
	}

	defer cursor.Close(ctx)

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
