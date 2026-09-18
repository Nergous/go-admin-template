package config

import (
	"fmt"
	"io"
	"log/slog"
	"strings"
)

type SlogHandlerFunc func(io.Writer, *slog.HandlerOptions) slog.Handler

type LogFormat string

const (
	LogFormatUnknown LogFormat = ""
	LogFormatText    LogFormat = "text"
	LogFormatJSON    LogFormat = "json"
)

func LogFormats() []LogFormat {
	return []LogFormat{
		LogFormatText,
		LogFormatJSON,
	}
}

const (
	logFormatEnv     = "LOG_FORMAT"
	defaultLogFormat = LogFormatText
)

func LoadLogFormat() (LogFormat, error) {
	format := stringEnv(defaultLogFormat.String(), logFormatEnv)
	return NewLogFormat(format)
}

func NewLogFormat(l string) (LogFormat, error) {
	l = strings.ToLower(strings.TrimSpace(l))
	lf := LogFormat(l)
	if err := lf.Validate(); err != nil {
		return defaultLogFormat, err
	}

	return lf, nil
}

func (l LogFormat) String() string {
	return string(l)
}

func (l LogFormat) Validate() error {
	switch l {
	case LogFormatText, LogFormatJSON:
		return nil
	default:
		return fmt.Errorf("invalid log format: %q", l)
	}
}

func (l LogFormat) SlogHandlerFunc() SlogHandlerFunc {
	if l == LogFormatJSON {
		return func(writer io.Writer, options *slog.HandlerOptions) slog.Handler {
			return slog.NewJSONHandler(writer, options)
		}
	}

	return func(writer io.Writer, options *slog.HandlerOptions) slog.Handler {
		return slog.NewTextHandler(writer, options)
	}
}
