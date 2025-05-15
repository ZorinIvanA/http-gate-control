package mock

import (
	"context"

	"github.com/stretchr/testify/mock"
)

// LoggerMock - мок для Logger
type LoggerMock struct {
	mock.Mock
}

// NewLoggerMock создает новый мок для Logger
func NewLoggerMock() *LoggerMock {
	return &LoggerMock{}
}

// LogEvent реализует метод интерфейса Logger
func (m *LoggerMock) LogEvent(ctx context.Context, event string, metadata map[string]interface{}) {
	m.Called(ctx, event, metadata)
}
