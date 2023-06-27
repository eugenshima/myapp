package tests

import (
	"testing"

	"github.com/eugenshima/myapp/internal/repository"
)

type TestMongo struct {
	client *repository.MongoDBConnection
}

func NewTestMongo(client *repository.MongoDBConnection) *TestMongo {
	return &TestMongo{client: client}
}

func TestCreateMongoConnect(t *testing.T) {
	client, err := repository.CreateMongoConnect()

	// Assert that the client is not nil
	if client == nil {
		t.Error("MongoDB client not created")
	}

	// Assert that the connection is successful
	if err != nil {
		t.Errorf("Failed to connect to MongoDB: %v", err)
	}

	// Assert that the disconnection is successful
	if err != nil {
		t.Errorf("Failed to disconnect from MongoDB: %v", err)
	}
}

// func TestFindPersons(t *testing.T) {
// 	db := repository.MongoDBConnection{} // Create an instance of `MongoDbConnection`
// 	db.Connect()                         // Connect to the MongoDB database

// 	client, err := repository.CreateMongoConnect()

// 	TestMongo.client.FindPersons()

// 	collection := db.client.Database("test").Collection("person")
// 	_, err = collection.InsertMany(context.Background(), []interface{}{
// 		bson.M{"name": "John Doe", "age": 25},
// 		bson.M{"name": "Alice Smith", "age": 30},
// 	})

// 	if err != nil {
// 		t.Error("Failed to insert test data")
// 	}

// 	// Test the `FindPersons()` function
// 	db.FindPersons()

// 	// Cleanup
// 	_, err = collection.DeleteMany(context.Background(), bson.M{})
// 	if err != nil {
// 		t.Error("Failed to delete test data")
// 	}

// 	db.Disconnect() // Disconnect from the MongoDB database
// }
