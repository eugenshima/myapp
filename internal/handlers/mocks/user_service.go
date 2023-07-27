// Code generated by mockery v2.18.0. DO NOT EDIT.

package mocks

import (
	context "context"

	mock "github.com/stretchr/testify/mock"

	model "github.com/eugenshima/myapp/internal/model"

	uuid "github.com/google/uuid"
)

// UserService is an autogenerated mock type for the UserService type
type UserService struct {
	mock.Mock
}

// Delete provides a mock function with given fields: ctx, id
func (_m *UserService) Delete(ctx context.Context, id uuid.UUID) error {
	ret := _m.Called(ctx, id)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, uuid.UUID) error); ok {
		r0 = rf(ctx, id)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// GenerateTokens provides a mock function with given fields: ctx, login, password
func (_m *UserService) GenerateTokens(ctx context.Context, login string, password string) (string, string, error) {
	ret := _m.Called(ctx, login, password)

	var r0 string
	if rf, ok := ret.Get(0).(func(context.Context, string, string) string); ok {
		r0 = rf(ctx, login, password)
	} else {
		r0 = ret.Get(0).(string)
	}

	var r1 string
	if rf, ok := ret.Get(1).(func(context.Context, string, string) string); ok {
		r1 = rf(ctx, login, password)
	} else {
		r1 = ret.Get(1).(string)
	}

	var r2 error
	if rf, ok := ret.Get(2).(func(context.Context, string, string) error); ok {
		r2 = rf(ctx, login, password)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}

// GetAll provides a mock function with given fields: ctx
func (_m *UserService) GetAll(ctx context.Context) ([]*model.User, error) {
	ret := _m.Called(ctx)

	var r0 []*model.User
	if rf, ok := ret.Get(0).(func(context.Context) []*model.User); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*model.User)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = rf(ctx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// RefreshTokenPair provides a mock function with given fields: ctx, accessToken, refreshToken, id
func (_m *UserService) RefreshTokenPair(ctx context.Context, accessToken string, refreshToken string, id uuid.UUID) (string, string, error) {
	ret := _m.Called(ctx, accessToken, refreshToken, id)

	var r0 string
	if rf, ok := ret.Get(0).(func(context.Context, string, string, uuid.UUID) string); ok {
		r0 = rf(ctx, accessToken, refreshToken, id)
	} else {
		r0 = ret.Get(0).(string)
	}

	var r1 string
	if rf, ok := ret.Get(1).(func(context.Context, string, string, uuid.UUID) string); ok {
		r1 = rf(ctx, accessToken, refreshToken, id)
	} else {
		r1 = ret.Get(1).(string)
	}

	var r2 error
	if rf, ok := ret.Get(2).(func(context.Context, string, string, uuid.UUID) error); ok {
		r2 = rf(ctx, accessToken, refreshToken, id)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}

// Signup provides a mock function with given fields: ctx, entity
func (_m *UserService) Signup(ctx context.Context, entity *model.User) error {
	ret := _m.Called(ctx, entity)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *model.User) error); ok {
		r0 = rf(ctx, entity)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

type mockConstructorTestingTNewUserService interface {
	mock.TestingT
	Cleanup(func())
}

// NewUserService creates a new instance of UserService. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewUserService(t mockConstructorTestingTNewUserService) *UserService {
	mock := &UserService{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
