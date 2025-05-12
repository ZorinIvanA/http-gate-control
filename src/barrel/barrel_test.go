package barrel

import (
	"testing"
	"time"

	"github.com/ZorinIvanA/http-gate-control/metrics"
	"github.com/ZorinIvanA/http-gate-control/relay"

	"github.com/stretchr/testify/assert"
)

func Test_BarrelController(t *testing.T) {
	// Создаем мок метрик
	metrics := metrics.NewMetrics(3600, "/tmp/test.log")

	// Создаем мок реле
	relayMock := new(relay.MockRelayClient)

	// Создаем контроллер с блокировкой на 1 секунду
	controller := NewController(1, metrics)

	t.Run("First call should open gate", func(t *testing.T) {
		relayMock.On("OpenGate").Return(nil).Once()

		assert.True(t, controller.ShouldOpen())
		relayMock.AssertExpectations(t)
	})

	t.Run("Second call within block duration should not open gate", func(t *testing.T) {
		relayMock.On("OpenGate").Return(nil).Maybe()

		assert.False(t, controller.ShouldOpen())
		relayMock.AssertNumberOfCalls(t, "OpenGate", 0)
		assert.Equal(t, uint64(1), metrics.DailyBlocked())
	})

	t.Run("Call after block duration should open gate again", func(t *testing.T) {
		time.Sleep(1 * time.Second)
		relayMock.On("OpenGate").Return(nil).Once()

		assert.True(t, controller.ShouldOpen())
		relayMock.AssertExpectations(t)
		assert.Equal(t, uint64(1), metrics.DailyBlocked()) // счетчик блокировок не должен сброситься
	})
}
