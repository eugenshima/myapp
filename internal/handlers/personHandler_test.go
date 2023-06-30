package handlers

import (
	"github.com/stretchr/testify/mock"
)

type MockMongoDBConnection struct {
	mock.Mock
}

func (m *MockMongoDBConnection) Connect() error {
	args := m.Called()
	return args.Error(0)
}
