package config

import (
	"errors"
	"testing"
)

func TestNewEnumValuesReturnUnwrappedValidationErrors(t *testing.T) {
	tests := []struct {
		name     string
		newValue func(string) error
	}{
		{
			name: "environment",
			newValue: func(value string) error {
				_, err := NewEnvironment(value)
				return err
			},
		},
		{
			name: "log format",
			newValue: func(value string) error {
				_, err := NewLogFormat(value)
				return err
			},
		},
		{
			name: "log level",
			newValue: func(value string) error {
				_, err := NewLogLevel(value)
				return err
			},
		},
		{
			name: "log sink",
			newValue: func(value string) error {
				_, err := NewLogSink(value)
				return err
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			err := tt.newValue("invalid")
			if err == nil {
				t.Errorf("New%s(%q) error = nil, want validation error", tt.name, "invalid")
			}
			if errors.Unwrap(err) != nil {
				t.Errorf("New%s(%q) error unwrap = %v, want nil", tt.name, "invalid", errors.Unwrap(err))
			}
		})
	}
}

func TestNewLogValuesReturnDefaultsOnInvalidInput(t *testing.T) {
	tests := []struct {
		name     string
		newValue func(string) (any, error)
		want     any
	}{
		{
			name: "format",
			newValue: func(value string) (any, error) {
				return NewLogFormat(value)
			},
			want: defaultLogFormat,
		},
		{
			name: "level",
			newValue: func(value string) (any, error) {
				return NewLogLevel(value)
			},
			want: defaultLogLevel,
		},
		{
			name: "sink",
			newValue: func(value string) (any, error) {
				return NewLogSink(value)
			},
			want: defaultLogSink,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := tt.newValue("invalid")
			if err == nil {
				t.Errorf("NewLog%s(%q) error = nil, want validation error", tt.name, "invalid")
			}
			if got != tt.want {
				t.Errorf("NewLog%s(%q) = %v, want %v", tt.name, "invalid", got, tt.want)
			}
		})
	}
}

func TestLoggerValidateJoinsInvalidComponentErrors(t *testing.T) {
	logger := Logger{
		Level:  "invalid",
		Format: "invalid",
		Sink:   "relative-path",
	}

	err := logger.Validate()
	if err == nil {
		t.Fatal("Logger.Validate() error = nil, want three validation errors")
	}

	joined, ok := err.(interface{ Unwrap() []error })
	if !ok {
		t.Fatalf("Logger.Validate() error type = %T, want joined error", err)
	}
	if got := len(joined.Unwrap()); got != 3 {
		t.Errorf("Logger.Validate() error count = %d, want 3", got)
	}
}
