package relay

import (
	"github.com/stretchr/testify/mock"
)

type MockRelayClient struct {
	mock.Mock
}

func NewMockRelayClient() *MockRelayClient {
	return &MockRelayClient{}
}

func (m *MockRelayClient) OpenGate() error {
	args := m.Called()
	return args.Error(0)
}
