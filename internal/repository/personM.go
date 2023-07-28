// Package repository provides functions for interacting with a database
package repository

import (
	"context"
	"fmt"

	"github.com/eugenshima/myapp/internal/model"

	"github.com/google/uuid"
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

// Update is a func which executes MongoDB command db.person.updateOne
func (db *MongoDBConnection) Update(ctx context.Context, uuidString uuid.UUID, person *model.Person) (uuid.UUID, error) {
	collection := db.client.Database("my_mongo_base").Collection("person")
	filter := bson.M{"_id": uuidString}
	update := bson.M{"$set": person}
	_, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return uuid.Nil, fmt.Errorf("UpdateOne: %w", err)
	}
	return person.ID, nil
}

// Delete is a func which executes MongoDB command db.person.deleteOne
func (db *MongoDBConnection) Delete(ctx context.Context, uuidString uuid.UUID) (uuid.UUID, error) {
	collection := db.client.Database("my_mongo_base").Collection("person")
	filter := bson.M{"_id": uuidString}
	_, err := collection.DeleteOne(ctx, filter)
	if err != nil {
		return uuid.Nil, fmt.Errorf("error in DeleteOne : %v", err)
	}
	return uuidString, nil
}

// Create function executes "db.person.insertOne()" command
func (db *MongoDBConnection) Create(ctx context.Context, person *model.Person) (uuid.UUID, error) {
	collection := db.client.Database("my_mongo_base").Collection("person")
	_, err := collection.InsertOne(ctx, person)
	if err != nil {
		return uuid.Nil, fmt.Errorf("InsertOne: %w", err)
	}
	return person.ID, nil
}

// GetByID function executes "db.person.FindOne()" command
func (db *MongoDBConnection) GetByID(ctx context.Context, ID uuid.UUID) (*model.Person, error) {
	collection := db.client.Database("my_mongo_base").Collection("person")
	filter := bson.M{"_id": ID}
	var person model.Person
	err := collection.FindOne(ctx, filter).Decode(&person)

	if err != nil {
		return nil, fmt.Errorf("Decode(): %w", err)
	}
	return &person, nil
}

// GetAll function executes "db.person.FindOne()" command
//nolint:dupl
func (db *MongoDBConnection) GetAll(ctx context.Context) ([]*model.Person, error) {
	collection := db.client.Database("my_mongo_base").Collection("person")
	filter := bson.M{}
	var all []*model.Person
	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("Find(): %w", err)
	}

	for cursor.Next(ctx) {
		var pers *model.Person
		err = cursor.Decode(&pers)

		if err != nil {
			return nil, fmt.Errorf("Decode(): %w", err)
		}
		all = append(all, pers)
	}
	return all, nil
}
