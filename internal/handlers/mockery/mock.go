// mocks package for testing
package mocks

import (
	"context"
	"testing"

	"github.com/eugenshima/myapp/internal/model"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

type MockMongoDBConnection struct {
	mock.Mock
}

func (m *MockMongoDBConnection) GetAll(ctx context.Context) ([]model.Person, error) {
	args := m.Called(ctx)

	if args.Get(1) != nil {
		return nil, args.Get(1).(error)
	}

	return args.Get(0).([]model.Person), nil
}

func TestGetAll(t *testing.T) {
	mockConn := new(MockMongoDBConnection)

	// Setup expectations
	expected := []model.Person{
		{ID: uuid.New(), Name: "John"},
		{ID: uuid.New(), Name: "Jane"},
	}
	mockConn.On("GetAll", mock.Anything).Return(expected, nil)

	// Call the function being tested
	entities, err := mockConn.GetAll(context.Background())

	// Validate the results
	mockConn.AssertExpectations(t)
	if err != nil {
		t.Errorf("expected no error but got: %v", err)
	}
	if len(entities) != len(expected) {
		t.Errorf("expected %d entities but got %d", len(expected), len(entities))
	}
	for i, entity := range entities {
		if entity.ID != expected[i].ID || entity.Name != expected[i].Name {
			t.Errorf("expected %v but got %v", expected[i], entity)
		}
	}
}
