package logger

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// NewLogger создает новый логгер с выводом в stdout и указанный файл
func NewLogger(logPath string) (*zap.Logger, error) {
	config := zap.NewProductionConfig()
	config.OutputPaths = []string{"stdout", logPath}
	config.ErrorOutputPaths = []string{"stderr", logPath}
	config.Level = zap.NewAtomicLevelAt(zapcore.InfoLevel)

	// Проверяем, доступен ли путь для логов
	if err := ensureLogPath(logPath); err != nil {
		return nil, err
	}

	return config.Build()
}

// ensureLogPath проверяет доступность пути для логов и создает необходимые директории
func ensureLogPath(path string) error {
	// Извлекаем директорию из пути к файлу
	dir := path[:len(path)-len("/gate-service.log")]

	// Создаем директорию, если она не существует
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	return nil
}
