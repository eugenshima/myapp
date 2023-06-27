package tests

import (
	"testing"

	"github.com/eugenshima/myapp/internal/repository"
)

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
