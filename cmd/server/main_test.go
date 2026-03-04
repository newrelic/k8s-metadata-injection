package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap/zapcore"
)

func TestSetupLogger(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name          string
		level         string
		expectedLevel zapcore.Level
		shouldWarn    bool
	}{
		{
			name:          "debug level",
			level:         "debug",
			expectedLevel: zapcore.DebugLevel,
			shouldWarn:    false,
		},
		{
			name:          "info level",
			level:         "info",
			expectedLevel: zapcore.InfoLevel,
			shouldWarn:    false,
		},
		{
			name:          "warn level",
			level:         "warn",
			expectedLevel: zapcore.WarnLevel,
			shouldWarn:    false,
		},
		{
			name:          "error level",
			level:         "error",
			expectedLevel: zapcore.ErrorLevel,
			shouldWarn:    false,
		},
		{
			name:          "uppercase level",
			level:         "INFO",
			expectedLevel: zapcore.InfoLevel,
			shouldWarn:    false,
		},
		{
			name:          "level with whitespace",
			level:         "  warn  ",
			expectedLevel: zapcore.WarnLevel,
			shouldWarn:    false,
		},
		{
			name:          "invalid level defaults to info",
			level:         "invalid",
			expectedLevel: zapcore.InfoLevel,
			shouldWarn:    true,
		},
		{
			name:          "empty level defaults to info",
			level:         "",
			expectedLevel: zapcore.InfoLevel,
			shouldWarn:    true,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			t.Parallel()

			logger := setupLogger(c.level)
			assert.NotNil(t, logger)

			// Verify the logger's level
			core := logger.Desugar().Core()
			assert.True(t, core.Enabled(c.expectedLevel))

			// For invalid levels, verify it still works at info level
			if c.shouldWarn {
				assert.True(t, core.Enabled(zapcore.InfoLevel))
			}
		})
	}
}
