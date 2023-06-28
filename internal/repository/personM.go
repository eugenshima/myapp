// Package repository provides functions for interacting with a database
package repository

import (
	"fmt"

	"github.com/eugenshima/myapp/internal/model"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// MongoDBConnection is a struct, which contains *mongo.Client variable
type MongoDBConnection struct {
	client *mongo.Client
}

// NewMongoDBConnection func is a constructor of MongoDbConnection struct
func NewMongoDBConnection(client *mongo.Client) *MongoDBConnection {
	return &MongoDBConnection{client: client}
}

// Find executes the "db.person.find()" command
// func (db *MongoDBConnection) GetAll(ctx context.Context) (bson.M, error) {
// 	// Select the database and collection
// 	collection := db.client.Database("my_mongo_base").Collection("person")

// 	// Call the find function on the collection
// 	cursor, err := collection.Find(ctx, bson.M{})
// 	if err != nil {
// 		return nil, fmt.Errorf("error in Find method (Find): %v", err)
// 	}

// 	defer func() {
// 		if err := cursor.Close(ctx); err != nil {
// 			errMsg := fmt.Errorf("error defering func (Find): %v", err)
// 			log.Println(errMsg)
// 		}
// 	}()

// 	var result bson.M
// 	// Iterate over the cursor and print the result
// 	for cursor.Next(ctx) {
// 		err := cursor.Decode(&result)

// 		if err != nil {
// 			return nil, fmt.Errorf("error in decoding (Find): %v", err)
// 		}

// 		fmt.Println(result)
// 	}

// 	if err := cursor.Err(); err != nil {
// 		return nil, fmt.Errorf("error in cursor.Next cycle (Find): %v", err)
// 	}
// 	return result, nil
// }

// Create function executes "db.person.insertOne()" command
// func (db *MongoDBConnection) Create(ctx context.Context, entity model.Person) (interface{}, error) {
// 	// Select the database and collection
// 	collection := db.client.Database("my_mongo_base").Collection("person")

// 	// Call the InsertOne function on the collection
// 	cursor, err := collection.InsertOne(ctx, entity)
// 	if err != nil {
// 		return nil, fmt.Errorf("error in InsertOne method (Create): %v", err)
// 	}

// 	return cursor.InsertedID, nil
// }

// Delete is a func which executes MongoDB command db.person.deleteOne
// func (db *MongoDBConnection) Delete(ctx context.Context, entity model.Person) (interface{}, error) {
// 	// Select the database and collection
// 	collection := db.client.Database("my_mongo_base").Collection("person")

// 	// Call the felete function of the collection
// 	cursor, err := collection.DeleteOne(ctx, entity)
// 	if err != nil {
// 		return nil, fmt.Errorf("error in DeleteOne method (Create): %v", err)
// 	}
// 	return cursor.DeletedCount, nil
// }

// Update is a func which executes MongoDB command db.person.updateOne
// func (db *MongoDBConnection) Update(ctx context.Context, objectID string, updateData map[string]interface{}) error {
// 	// Define the filter
// 	filter := bson.M{"userID": objectID}

// 	// Define the update document using $set operator
// 	update := bson.M{"$set": updateData}
// 	// Select the database and collection
// 	collection := db.client.Database("my_mongo_base").Collection("person")

// 	// Call the update function of the collection
// 	cursor, err := collection.UpdateByID(ctx, filter, update)
// 	if err != nil {
// 		return fmt.Errorf("error in update method (personM): %v", err)
// 	}
// 	fmt.Println(cursor.UpsertedID)
// 	return nil
// }

// Update is a func which executes MongoDB command db.person.updateOne
func (db *MongoDBConnection) Update(c echo.Context, uuidString string, entity *model.Person) error {
	fmt.Println("Update is working")
	return nil
}

// Delete is a func which executes MongoDB command db.person.deleteOne
func (db *MongoDBConnection) Delete(c echo.Context, uuidString string) error {
	fmt.Println("Delete is working")
	return nil
}

// Create function executes "db.person.insertOne()" command
func (db *MongoDBConnection) Create(c echo.Context, entity *model.Person) error {
	fmt.Println("Create is working")
	return nil

}
func (db *MongoDBConnection) GetByName(c echo.Context, Name string) (*model.Person, error) {
	fmt.Println("GetByName is working")
	return nil, nil
}

// Find executes the "db.person.find()" command
func (db *MongoDBConnection) GetAll(c echo.Context) ([]model.Person, error) {

	collection := db.client.Database("my_mongo_base").Collection("person")

	filter := bson.M{}
	var allPers []model.Person
	cursor, err := collection.Find(c.Request().Context(), filter)

	if err != nil {
		return nil, fmt.Errorf("PersonMongo -> GetAll -> Find -> error: %w", err)
	}
	var pers model.Person
	for cursor.Next(c.Request().Context()) {
		err = cursor.Decode(&pers)
		if err != nil {
			return allPers, fmt.Errorf("PersonMongo -> GetAll -> Decode -> error: %w", err)
		}
		allPers = append(allPers, pers)
	}
	return allPers, nil
}

// // GetAll reads all documents from mongoDB collection
// func (rpsMongo *Mongo) GetAll(ctx context.Context) ([]model.Person, error) {
// 	coll := rpsMongo.client.Database("personMongoDB").Collection("persons")

// 	filter := bson.M{}
// 	var allPers []model.Person
// 	cursor, err := coll.Find(ctx, filter)

// 	if err != nil {
// 		return nil, fmt.Errorf("PersonMongo -> GetAll -> Find -> error: %w", err)
// 	}
// 	var pers model.Person
// 	for cursor.Next(ctx) {
// 		err = cursor.Decode(&pers)
// 		if err != nil {
// 			return allPers, fmt.Errorf("PersonMongo -> GetAll -> Decode -> error: %w", err)
// 		}
// 		allPers = append(allPers, pers)
// 	}
// 	return allPers, nil
// }
