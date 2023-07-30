package handlers

// import (
// 	"context"
// 	"testing"

// 	mocks "github.com/eugenshima/myapp/internal/handlers/mocks"
// 	"github.com/eugenshima/myapp/internal/model"
// 	"github.com/google/uuid"
// 	"github.com/stretchr/testify/mock"
// 	"github.com/stretchr/testify/require"
// )

// var (
// 	mockUserService *mocks.UserService
// 	mockUserEntity  = model.User{
// 		ID:           uuid.New(),
// 		Login:        "test",
// 		Password:     []byte("test"),
// 		Role:         "user",
// 		RefreshToken: nil,
// 	}
// 	str string
// )

// func TestUserhandlerGetAll(t *testing.T) {
// 	mockUserService.On("GetAll", mock.Anything).Return([]*model.User{}, nil).Twice()
// 	handler := NewUserHandler(mockUserService, nil)
// 	res, err := mockUserService.GetAll(context.Background())
// 	require.NoError(t, err)
// 	results, err := handler.srv.GetAll(context.Background())
// 	require.NoError(t, err)
// 	require.NotNil(t, results)
// 	require.Equal(t, len(res), len(results))
// }

// func TestUserHandlerSignUp(t *testing.T) {
// 	mockUserService.On("Signup", mock.Anything, mock.AnythingOfType("*model.User")).Return(nil).Once()

// 	err := mockUserService.Signup(context.Background(), &mockUserEntity)
// 	require.NoError(t, err)
// }

// func TestUserhandlerLogin(t *testing.T) {
// 	mockUserService.On("GenerateTokens", mock.Anything, mock.AnythingOfType("string"), mock.AnythingOfType("string")).Return(str, str, nil).Once()
// 	access, refresh, err := mockUserService.GenerateTokens(context.Background(), mockUserEntity.Login, string(mockUserEntity.Password))
// 	require.NoError(t, err)
// 	require.IsType(t, "string", access)
// 	require.IsType(t, "string", refresh)
// }

// func TestUserHandlerRefreshTokenPair(t *testing.T) {
// 	mockUserService.On("RefreshTokenPair", mock.Anything, mock.AnythingOfType("string"), mock.AnythingOfType("string"), mock.AnythingOfType("uuid.UUID")).Return(str, str, nil).Once()
// 	access, refresh, err := mockUserService.RefreshTokenPair(context.Background(), mockUserEntity.Login, string(mockUserEntity.Password), mockUserEntity.ID)
// 	require.NoError(t, err)
// 	require.IsType(t, "string", access)
// 	require.IsType(t, "string", refresh)
// }

// func TestUserHandlerDelete(t *testing.T) {
// 	mockUserService.On("Delete", mock.Anything, mock.AnythingOfType("uuid.UUID")).Return(nil).Once()
// 	err := mockUserService.Delete(context.Background(), mockUserEntity.ID)
// 	require.NoError(t, err)
// }
