package config

import (
	"fmt"
	"log/slog"
	"strings"
)

type LogLevel string

const (
	LogLevelUnknown LogLevel = ""
	LogLevelDebug   LogLevel = "debug"
	LogLevelInfo    LogLevel = "info"
	LogLevelWarn    LogLevel = "warn"
	LogLevelError   LogLevel = "error"
)

func LogLevels() []LogLevel {
	return []LogLevel{
		LogLevelDebug,
		LogLevelInfo,
		LogLevelWarn,
		LogLevelError,
	}
}

const (
	logLevelEnv     = "LOG_LEVEL"
	defaultLogLevel = LogLevelInfo
)

func LoadLogLevel() (LogLevel, error) {
	level := stringEnv(defaultLogLevel.String(), logLevelEnv)
	return NewLogLevel(level)
}

func NewLogLevel(l string) (LogLevel, error) {
	l = strings.ToLower(strings.TrimSpace(l))
	ll := LogLevel(l)
	if err := ll.Validate(); err != nil {
		return defaultLogLevel, err
	}

	return ll, nil
}

func (l LogLevel) String() string {
	return string(l)
}

func (l LogLevel) Validate() error {
	switch l {
	case LogLevelDebug, LogLevelInfo, LogLevelWarn, LogLevelError:
		return nil
	default:
		return fmt.Errorf("invalid log level: %q", l)
	}
}

func (l LogLevel) ToSlogLevel() slog.Level {
	switch l {
	case LogLevelDebug:
		return slog.LevelDebug
	case LogLevelInfo:
		return slog.LevelInfo
	case LogLevelWarn:
		return slog.LevelWarn
	case LogLevelError:
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}
