package repository_test

// // MongoDBConnection is a struct, which contains *mongo.Client variable
// type MongoDBConnection struct {
// 	person *mongo.Client
// }

// // NewMongoDBConnection func is a constructor of MongoDbConnection struct
// func NewMongoDBConnection(person *mongo.Client) *MongoDBConnection {
// 	return &MongoDBConnection{person: person}
// }

// type TestMongo struct {
// 	client *MongoDBConnection
// }

// func NewTestMongo(client *MongoDBConnection) *TestMongo {
// 	return &TestMongo{client: client}
// }

// // NewMongo creates a connection to MongoDB server
// func NewMongo() (*mongo.Client, error) {
// 	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")

// 	// Connect to MongoDB
// 	client, err := mongo.Connect(context.Background(), clientOptions)
// 	if err != nil {
// 		return nil, err
// 	}

// 	// Check the connection
// 	err = client.Ping(context.Background(), nil)
// 	if err != nil {
// 		return nil, err
// 	}
// 	fmt.Println("Connected to MongoDB!")

// 	return client, nil
// }

// func TestCreateMongoConnect(t *testing.T) {
// 	client, err := NewMongo()
// 	// Assert that the client is not nil
// 	if client == nil {
// 		t.Error("MongoDB client not created")
// 	}

// 	// Assert that the connection is successful
// 	if err != nil {
// 		t.Errorf("Failed to connect to MongoDB: %v", err)
// 	}

// 	// Assert that the disconnection is successful
// 	if err != nil {
// 		t.Errorf("Failed to disconnect from MongoDB: %v", err)
// 	}
// }

// func TestFindPersons(t *testing.T) {
// 	db := MongoDBConnection{} // Create an instance of `MongoDbConnection`
// 	db.Connect()                         // Connect to the MongoDB database

// 	client, err := NewMongo()

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
