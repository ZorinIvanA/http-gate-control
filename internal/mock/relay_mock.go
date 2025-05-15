package mock

import (
	"context"

	"github.com/stretchr/testify/mock"
)

// RelayMock - мок для RelayClient
type RelayMock struct {
	mock.Mock
}

// NewRelayMock создает новый мок для RelayClient
func NewRelayMock() *RelayMock {
	return &RelayMock{}
}

// OpenGate реализует метод интерфейса RelayClient
func (m *RelayMock) OpenGate(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}
