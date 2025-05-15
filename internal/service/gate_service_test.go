package service

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestGateService_CheckAccess(t *testing.T) {
	// Создаем сервис
	service := NewGateService(5 * time.Second)

	// Первый запрос - должен открыть
	assert.True(t, service.CheckAccess())

	// Второй запрос через 1 секунду - должен заблокировать
	time.Sleep(1 * time.Second)
	assert.False(t, service.CheckAccess())

	// Третий запрос через 5 секунд после первого - должен открыть
	time.Sleep(4 * time.Second)
	assert.True(t, service.CheckAccess())
}
