package logger

import (
	"log/slog"
	"os"
)

const MessageKey = "MESSAGE"

type Logger struct {
	serviceName string
}

var logger *Logger

func NewAppLogger(serviceName string) *Logger {
	logger = &Logger{serviceName: serviceName}
	return logger
}

func (l *Logger) LogInfo(args ...any) {
	slog.Info(l.serviceName, args...)
}

func (l *Logger) LogDebug(args ...any) {
	slog.Debug(l.serviceName, args...)
}

func (l *Logger) LogError(args ...any) {
	slog.Error(l.serviceName, args...)
}

func (l *Logger) LogFatalError(args ...any) {
	slog.Error("FATAL "+l.serviceName, args...)
	os.Exit(1)
}
