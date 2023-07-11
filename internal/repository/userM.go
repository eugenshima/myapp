package repository

import (
	"context"
	"fmt"

	"github.com/eugenshima/myapp/internal/model"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// import "go.mongodb.org/mongo-driver/mongo"

// MongoDBConnection is a struct, which contains *mongo.Client variable
type UserMongoDBConnection struct {
	client *mongo.Client
}

// NewMongoDBConnection func is a constructor of MongoDbConnection struct
func NewUserMongoDBConnection(client *mongo.Client) *UserMongoDBConnection {
	return &UserMongoDBConnection{client: client}
}

// GetUser function executes a query, which select all rows from user table
func (db *UserMongoDBConnection) GetUser(ctx context.Context, login string) (*model.GetUser, error) {
	var user *model.GetUser
	collection := db.client.Database("my_mongo_base").Collection("user")
	filter := bson.M{"login": login}
	err := collection.FindOne(ctx, filter).Decode(&user)
	if err != nil {
		return nil, fmt.Errorf("error: %v", err)
	}
	return user, nil
}

// Signup function executes a query, which insert a user to user table
func (db *UserMongoDBConnection) Signup(ctx context.Context, user *model.User) error {
	collection := db.client.Database("my_mongo_base").Collection("user")
	_, err := collection.InsertOne(ctx, user)
	if err != nil {
		return fmt.Errorf("InsertOne: %w", err)
	}
	return nil
}

// GetAll func executes a query, which returns all users
func (db *UserMongoDBConnection) GetAll(ctx context.Context) ([]*model.User, error) {
	collection := db.client.Database("my_mongo_base").Collection("user")
	filter := bson.M{}
	var all []*model.User
	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("Find(): %w", err)
	}

	for cursor.Next(ctx) {
		var pers *model.User
		err = cursor.Decode(&pers)

		if err != nil {
			return nil, fmt.Errorf("Decode(): %w", err)
		}
		all = append(all, pers)
	}
	return all, nil
}

// SaveRefreshToken func executes a query, which saves the refresh token to a specific user
func (db *UserMongoDBConnection) SaveRefreshToken(ctx context.Context, ID uuid.UUID, token []byte) error {
	collection := db.client.Database("my_mongo_base").Collection("user")
	filter := bson.M{"_id": ID}
	var user *model.User
	err := collection.FindOne(ctx, filter).Decode(&user)
	if err != nil {
		return fmt.Errorf("FindOne(): %w", err)
	}
	user.RefreshToken = token
	filter = bson.M{"_id": ID}
	update := bson.M{"$set": user}
	_, err = collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("InsertOne(): %w", err)
	}
	return nil
}

// GetRefreshToken returns a refresh token for the given user
func (db *UserMongoDBConnection) GetRefreshToken(ctx context.Context, ID uuid.UUID) ([]byte, error) {
	collection := db.client.Database("my_mongo_base").Collection("user")
	filter := bson.M{"_id": ID}
	var user *model.User
	err := collection.FindOne(ctx, filter).Decode(&user)
	if err != nil {
		return nil, fmt.Errorf("FindOne(): %w", err)
	}
	return user.RefreshToken, nil
}

// GetRoleByID returns a role for the given user ID
func (db *UserMongoDBConnection) GetRoleByID(ctx context.Context, ID uuid.UUID) (string, error) {
	collection := db.client.Database("my_mongo_base").Collection("user")
	filter := bson.M{"_id": ID}
	var user model.User
	err := collection.FindOne(ctx, filter).Decode(&user)
	if err != nil {
		return "", fmt.Errorf("FindOne(): %w", err)
	}
	return user.Role, nil
}

// Delete func deletes user from the database
func (db *UserMongoDBConnection) Delete(ctx context.Context, ID uuid.UUID) (uuid.UUID, error) {
	collection := db.client.Database("my_mongo_base").Collection("user")
	filter := bson.M{"_id": ID}
	res, err := collection.DeleteOne(ctx, filter)
	if err != nil {
		return uuid.Nil, fmt.Errorf("DeleteOne(): %v", err)
	}
	if res.DeletedCount == 0 {
		return uuid.Nil, fmt.Errorf("DeleteOne(): %v", err)
	}
	return ID, nil
}
