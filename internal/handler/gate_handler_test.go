package handler

import (
	"net/http"
	"net/http/httptest"
	"testing"

	internalMock "github.com/ZorinIvanA/http-gate-control/internal/mock"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// FakeGateService — фейковая реализация GateServiceInterface для тестов
type FakeGateService struct {
	checkAccess bool
}

func (f *FakeGateService) CheckAccess() bool {
	return f.checkAccess
}

func (f *FakeGateService) GetBlockedCount() int {
	return 0
}

func TestHandleOpen_ShouldOpenGate(t *testing.T) {
	// Создаем моки
	relayMock := new(internalMock.RelayMock)
	loggerMock := new(internalMock.LoggerMock)

	// Настраиваем ожидаемое поведение
	relayMock.On("OpenGate", mock.Anything).Return(nil)
	loggerMock.On("LogEvent", mock.Anything, "Gate opened", mock.Anything).Return()

	// Создаем фейковый сервис, который разрешает доступ
	fakeService := &FakeGateService{checkAccess: true}

	// Создаем хендлер
	handler := NewGateHandler(fakeService, relayMock, loggerMock)

	// Создаем тестовый запрос
	req, _ := http.NewRequest("GET", "/open", nil)
	rec := httptest.NewRecorder()

	// Вызываем хендлер
	handler.HandleOpen(rec, req)

	// Проверяем результат
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Contains(t, rec.Body.String(), "Gate opened")

	// Проверяем, что моки были вызваны
	relayMock.AssertExpectations(t)
	loggerMock.AssertExpectations(t)
}
