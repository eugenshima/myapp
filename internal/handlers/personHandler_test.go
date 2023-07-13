package handlers

import (
	"context"
	"os"
	"testing"

	mocks "github.com/eugenshima/myapp/internal/handlers/mocks"
	"github.com/eugenshima/myapp/internal/model"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

var (
	mockPersonService *mocks.PersonService
	mockPersonEntity  = model.Person{
		ID:        uuid.New(),
		Name:      "test",
		Age:       123,
		IsHealthy: true,
	}
)

// TestMain execute all tests
func TestMain(m *testing.M) {
	mockPersonService = new(mocks.PersonService)
	mockUserService = new(mocks.UserService)
	exitVal := m.Run()
	os.Exit(exitVal)
}

// TestCreate is a mocktest for Create method of interface Service
func TestCreate(t *testing.T) {
	mockPersonService.On("Create", mock.Anything, mock.AnythingOfType("*model.Person")).Return(uuid.UUID{}, nil).Once()

	id, err := mockPersonService.Create(context.Background(), &mockPersonEntity)
	require.NoError(t, err)
	require.NotNil(t, id)

	mockPersonService.AssertExpectations(t)
}

func TestDelete(t *testing.T) {
	mockPersonService.On("Delete", mock.Anything, mock.AnythingOfType("uuid.UUID")).Return(uuid.UUID{}, nil).Once()

	id, err := mockPersonService.Delete(context.Background(), mockPersonEntity.ID)
	require.NoError(t, err)
	require.NotNil(t, id)
}

func TestGetAll(t *testing.T) {
	mockPersonService.On("GetAll", mock.Anything).Return([]*model.Person{}, nil).Twice()
	handler := NewPersonHandler(mockPersonService, nil)
	res, err := mockPersonService.GetAll(context.Background())
	require.NoError(t, err)
	results, err := handler.srv.GetAll(context.Background())
	require.NoError(t, err)
	require.NotNil(t, results)
	require.Equal(t, len(res), len(results))
}

func TestUpdate(t *testing.T) {
	mockPersonService.On("Update", mock.Anything, mock.AnythingOfType("uuid.UUID"), mock.AnythingOfType("*model.Person")).Return(uuid.UUID{}, nil).Once()

	id, err := mockPersonService.Update(context.Background(), mockPersonEntity.ID, &mockPersonEntity)
	require.NoError(t, err)
	require.NotNil(t, id)
}

func TestGetByID(t *testing.T) {
	mockPersonService.On("GetByID", mock.Anything, mock.AnythingOfType("uuid.UUID")).Return(&model.Person{}, nil).Once()

	id, err := mockPersonService.GetByID(context.Background(), mockPersonEntity.ID)
	require.NoError(t, err)
	require.NotNil(t, id)
}
