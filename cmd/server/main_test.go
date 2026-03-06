package main

import (
	"testing"

	"go.uber.org/zap/zapcore"
)

func TestSetupLogger_ValidLogLevels(test *testing.T) {
	test.Parallel()
	tests := []struct {
		name          string
		logLevel      string
		expectedLevel zapcore.Level
	}{
		{
			name:          "debug level",
			logLevel:      "debug",
			expectedLevel: zapcore.DebugLevel,
		},
		{
			name:          "info level",
			logLevel:      "info",
			expectedLevel: zapcore.InfoLevel,
		},
		{
			name:          "warn level",
			logLevel:      "warn",
			expectedLevel: zapcore.WarnLevel,
		},
		{
			name:          "error level",
			logLevel:      "error",
			expectedLevel: zapcore.ErrorLevel,
		},
		{
			name:          "dpanic level",
			logLevel:      "dpanic",
			expectedLevel: zapcore.DPanicLevel,
		},
		{
			name:          "panic level",
			logLevel:      "panic",
			expectedLevel: zapcore.PanicLevel,
		},
		{
			name:          "fatal level",
			logLevel:      "fatal",
			expectedLevel: zapcore.FatalLevel,
		},
		{
			name:          "uppercase DEBUG",
			logLevel:      "DEBUG",
			expectedLevel: zapcore.DebugLevel,
		},
		{
			name:          "uppercase INFO",
			logLevel:      "INFO",
			expectedLevel: zapcore.InfoLevel,
		},
		{
			name:          "mixed case WaRn",
			logLevel:      "WaRn",
			expectedLevel: zapcore.WarnLevel,
		},
	}

	for _, tt := range tests {
		test.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			logger := setupLogger(tt.logLevel)

			if logger == nil {
				t.Fatal("setupLogger returned nil")
			}

			// Get the underlying zap logger
			zapLogger := logger.Desugar()
			if zapLogger == nil {
				t.Fatal("logger.Desugar() returned nil")
			}

			// Check that the logger's core is enabled at the expected level
			if !zapLogger.Core().Enabled(tt.expectedLevel) {
				t.Errorf("logger is not enabled at level %v", tt.expectedLevel)
			}

		})
	}
}

func TestSetupLogger_InvalidLogLevels(test *testing.T) {
	test.Parallel()
	tests := []struct {
		name     string
		logLevel string
	}{
		{
			name:     "invalid level",
			logLevel: "invalid",
		},
		{
			name:     "empty string",
			logLevel: "",
		},
		{
			name:     "random string",
			logLevel: "foobar",
		},
		{
			name:     "trace level (not supported by zap)",
			logLevel: "trace",
		},
	}

	for _, tt := range tests {
		test.Run(tt.name, func(test *testing.T) {
			test.Parallel()

			logger := setupLogger(tt.logLevel)

			if logger == nil {
				test.Fatal("setupLogger returned nil even for invalid level")
			}

			// Invalid levels should default to info
			zapLogger := logger.Desugar()
			if !zapLogger.Core().Enabled(zapcore.InfoLevel) {
				test.Error("logger with invalid level should default to info level")
			}
		})
	}
}
