package logger

import (
	"github.com/akrindev/gonely/internal/infrastructure/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Logger wraps zap.Logger
type Logger struct {
	*zap.Logger
}

// New creates a new logger instance based on configuration
func New(cfg *config.Config) (*Logger, error) {
	var zapConfig zap.Config

	if cfg.App.IsDevelopment() {
		zapConfig = zap.NewDevelopmentConfig()
		zapConfig.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	} else {
		zapConfig = zap.NewProductionConfig()
	}

	// Set log level
	level, err := zapcore.ParseLevel(cfg.Log.Level)
	if err != nil {
		level = zapcore.InfoLevel
	}
	zapConfig.Level = zap.NewAtomicLevelAt(level)

	// Set encoding format
	if cfg.Log.Format == "console" {
		zapConfig.Encoding = "console"
	} else {
		zapConfig.Encoding = "json"
	}

	zapLogger, err := zapConfig.Build(
		zap.AddCallerSkip(1),
		zap.AddStacktrace(zapcore.ErrorLevel),
	)
	if err != nil {
		return nil, err
	}

	return &Logger{zapLogger}, nil
}

// With creates a child logger with additional fields
func (l *Logger) With(fields ...zap.Field) *Logger {
	return &Logger{l.Logger.With(fields...)}
}

// Named creates a named logger
func (l *Logger) Named(name string) *Logger {
	return &Logger{l.Logger.Named(name)}
}

// Debug logs a debug level message
func (l *Logger) Debug(msg string, fields ...zap.Field) {
	l.Logger.Debug(msg, fields...)
}

// Info logs an info level message
func (l *Logger) Info(msg string, fields ...zap.Field) {
	l.Logger.Info(msg, fields...)
}

// Warn logs a warn level message
func (l *Logger) Warn(msg string, fields ...zap.Field) {
	l.Logger.Warn(msg, fields...)
}

// Error logs an error level message
func (l *Logger) Error(msg string, fields ...zap.Field) {
	l.Logger.Error(msg, fields...)
}

// Fatal logs a fatal level message and exits the application
func (l *Logger) Fatal(msg string, fields ...zap.Field) {
	l.Logger.Fatal(msg, fields...)
}

// Sync flushes any buffered log entries
func (l *Logger) Sync() error {
	return l.Logger.Sync()
}
