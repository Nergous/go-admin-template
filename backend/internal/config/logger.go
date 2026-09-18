package config

import (
	"errors"
	"fmt"
)

const (
	logAddSourceEnv = "LOG_ADD_SOURCE"
)

type Logger struct {
	Level     LogLevel
	Format    LogFormat
	Sink      LogSink
	AddSource bool
}

func LoadLogger() (Logger, error) {
	var errs []error

	addSource, err := boolEnv(false, logAddSourceEnv)
	if err != nil {
		errs = append(errs, fmt.Errorf("addSource: %w", err))
	}

	logLevel, levelErr := LoadLogLevel()
	if levelErr != nil {
		errs = append(errs, fmt.Errorf("level: %w", levelErr))
	}

	logFormat, formatErr := LoadLogFormat()
	if formatErr != nil {
		errs = append(errs, fmt.Errorf("format: %w", formatErr))
	}

	logSink, sinkErr := LoadLogSink()
	if sinkErr != nil {
		errs = append(errs, fmt.Errorf("sink: %w", sinkErr))
	}

	if err := errors.Join(errs...); err != nil {
		return Logger{}, err
	}

	return NewLogger(
		logLevel,
		logFormat,
		logSink,
		addSource,
	)
}

func NewLogger(logLevel LogLevel, logFormat LogFormat, logSink LogSink, addSource bool) (Logger, error) {
	logger := Logger{
		Level:     logLevel,
		Format:    logFormat,
		Sink:      logSink,
		AddSource: addSource,
	}

	if err := logger.Validate(); err != nil {
		return Logger{}, err
	}

	return logger, nil
}

func (l *Logger) Validate() error {
	if l == nil {
		return errors.New("logger is nil")
	}

	var errs []error

	if err := l.Level.Validate(); err != nil {
		errs = append(errs, fmt.Errorf("level: %w", err))
	}

	if err := l.Format.Validate(); err != nil {
		errs = append(errs, fmt.Errorf("format: %w", err))
	}

	if err := l.Sink.Validate(); err != nil {
		errs = append(errs, fmt.Errorf("sink: %w", err))
	}

	return errors.Join(errs...)
}
